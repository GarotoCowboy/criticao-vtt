package chat

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/chat"
	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func (s *ChatService) SendMessage(ctx context.Context, req *chat.SendMessageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("Starting GPRC to SendMessage for table %d", req.TableId)

	// --- 1. Authentication: Get UserID from context ---
	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 2. Validation: Check request payload ---
	if err := Validate(req); err != nil {
		s.Logger.ErrorF("Validation error for SendMessage request: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// --- 3. Authorization: Verify user is a member of the table ---
	s.Logger.InfoF("Authorizing: checking if user %d is a member of table %d", userID, req.TableId)
	var tableUserModel models.TableUser
	// Preload User to get the username for the event payload later
	if err := s.Db.WithContext(ctx).Preload("User").Where("user_id = ? AND table_id = ?", userID, req.TableId).First(&tableUserModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Authorization failed: User %d is not in table %d", userID, req.TableId)
			return nil, status.Errorf(codes.NotFound, "user is not in this table")
		}
		s.Logger.ErrorF("Database error checking user membership: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	// --- 4. Process Attachments ---
	var attachmentsJSON datatypes.JSON
	if len(req.Attachments) > 0 {
		s.Logger.InfoF("Marshaling %d attachments to JSON", len(req.Attachments))
		jsonBytes, err := json.Marshal(req.GetAttachments())
		if err != nil {
			s.Logger.ErrorF("JSON marshal error for attachments: %v", err.Error())
			return nil, status.Errorf(codes.Internal, "failed to process attachments")
		}
		attachmentsJSON = jsonBytes
	}

	// --- 5. Create Message Model ---
	chatMessageModel := models.ChatMessage{
		TableUserID:   tableUserModel.ID,      // ID of the TableUser link
		TableID:       tableUserModel.TableID, // Table ID for partitioning
		Message:       req.GetMessageText(),
		MessageType:   consts.MessageType(req.GetMessageType()),
		MessageStatus: consts.MessageStatus(chat.MessageStatus_SENT), // Set initial status
		Attachments:   attachmentsJSON,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Add optional fields
	if req.MediaUrl != nil {
		chatMessageModel.MediaURL = req.MediaUrl
	}
	if req.ReplyToMessageId != nil {
		chatMessageModel.ReplyToMessageId = req.ReplyToMessageId
	}

	// --- 6. Save to Database ---
	s.Logger.InfoF("Saving new message to database...")
	if err := s.Db.WithContext(ctx).Create(&chatMessageModel).Error; err != nil {
		s.Logger.ErrorF("Database error creating message: %v", err)
		return nil, status.Errorf(codes.Internal, "could not save message")
	}
	s.Logger.InfoF("Message %s saved to database for table %d", chatMessageModel.ID, req.TableId)

	// --- 7. Prepare and Publish Event ---
	respProto := &chat.ChatMessageResponse{
		MessageUuid:      chatMessageModel.ID.String(),
		TableId:          req.GetTableId(),
		SenderId:         uint64(tableUserModel.ID),
		SenderUsername:   tableUserModel.User.Username, // Got this from Preload
		MessageText:      chatMessageModel.Message,
		MessageType:      req.GetMessageType(),
		MessageStatus:    chat.MessageStatus_SENT,
		SentAt:           timestamppb.New(chatMessageModel.CreatedAt),
		MediaUrl:         req.MediaUrl,
		Attachments:      req.Attachments,
		ReplyToMessageId: req.ReplyToMessageId,
	}

	// Create the sync event payload
	syncMsg := &syncBroker.SyncResponse{
		TableId: req.GetTableId(),
		Action: &syncBroker.SyncResponse_MessageSended{
			MessageSended: &chat.ChatMessageSent{
				Message: respProto,
			},
		},
	}

	// Publish to the table-specific channel
	s.Broker.Publish(pubSubSyncConst.TableSync, req.GetTableId(), syncMsg)
	s.Logger.InfoF("Published ChatMessageSent event to broker for table %d", req.TableId)

	return &emptypb.Empty{}, nil
}

func (s *ChatService) ListMessages(ctx context.Context, req *chat.ListMessagesRequest) (*chat.ListChatMessageResponse, error) {
	s.Logger.InfoF("gRPC ChatService: ListMessages initiated for table %d", req.TableId)

	// --- 1. Validation ---
	if req.GetTableId() == 0 {
		s.Logger.ErrorF("Invalid request: table_id is required")
		return nil, status.Errorf(codes.InvalidArgument, "table_id is required")
	}

	// --- 2. Authentication ---
	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 3. Authorization: Verify user is a member ---
	s.Logger.InfoF("Authorizing: checking if user %d is a member of table %d", userID, req.TableId)
	var count int64
	// Efficient check without fetching the full model
	s.Db.WithContext(ctx).Model(&models.TableUser{}).Where("user_id = ? AND table_id = ?", userID, req.TableId).Count(&count)
	if count == 0 {
		s.Logger.ErrorF("Authorization error: User %d is not a member of table %d", userID, req.TableId)
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of the specified table")
	}

	// --- 4. Pagination Setup ---
	pageSize := req.GetPageSize()
	if pageSize <= 0 || pageSize > 100 {
		s.Logger.WarningF("Invalid page size %d, defaulting to 25", pageSize)
		pageSize = 25 // Default page size
	}

	// --- 5. Database Query ---
	s.Logger.InfoF("Fetching messages for table %d, page size %d", req.TableId, pageSize)
	var messages []models.ChatMessage
	// Base query: get messages for the table, preload sender's user info, order by newest first
	query := s.Db.WithContext(ctx).
		Preload("TableUser.User"). // Preload user info through TableUser
		Where("table_id = ?", req.GetTableId()).
		Order("created_at DESC").
		Limit(int(pageSize))

	// Cursor-based pagination: if a cursor is provided
	if req.GetLastMessageId() != "" {
		s.Logger.InfoF("Using cursor: fetching messages created before message %s", req.GetLastMessageId())
		var cursorMessage models.ChatMessage
		// Find the cursor message to get its creation time
		if err := s.Db.First(&cursorMessage, "id = ?", req.GetLastMessageId()).Error; err == nil {
			// Add to query: only get messages *older* than the cursor
			query = query.Where("created_at < ?", cursorMessage.CreatedAt)
		} else {
			s.Logger.WarningF("Could not find cursor message %s, ignoring cursor", req.GetLastMessageId())
		}
	}

	if err := query.Find(&messages).Error; err != nil {
		s.Logger.ErrorF("Database error listing messages for table %d: %v", req.TableId, err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve messages")
	}

	// --- 6. Process Results ---
	s.Logger.InfoF("Found %d messages, processing response...", len(messages))
	respMessages := make([]*chat.ChatMessageResponse, 0, len(messages))

	for _, msg := range messages {
		// Process attachments from JSON
		var attachments []string
		if msg.Attachments != nil {
			if err := json.Unmarshal(msg.Attachments, &attachments); err != nil {
				s.Logger.ErrorF("Failed to unmarshal attachments for message %s: %v", msg.ID, err)
				// Don't fail the whole request, just skip attachments for this message
			}
		}

		// Build the protobuf response message
		respMsg := &chat.ChatMessageResponse{
			MessageUuid:      msg.ID.String(),
			TableId:          uint64(msg.TableID),
			SenderId:         uint64(msg.TableUserID),
			SenderUsername:   msg.TableUser.User.Username, // From preload
			MessageText:      msg.Message,
			MessageType:      chat.MessageType(msg.MessageType),
			MessageStatus:    chat.MessageStatus(msg.MessageStatus),
			SentAt:           timestamppb.New(msg.CreatedAt),
			MediaUrl:         msg.MediaURL,
			Attachments:      attachments,
			ReplyToMessageId: msg.ReplyToMessageId,
			IsDeleted:        !msg.DeletedAt.Time.IsZero(), // Check if soft-deleted
		}
		// Only add UpdatedAt if it's different from CreatedAt
		if msg.UpdatedAt.After(msg.CreatedAt) {
			respMsg.UpdatedAt = timestamppb.New(msg.UpdatedAt)
		}

		respMessages = append(respMessages, respMsg)
	}

	// --- 7. Determine Next Cursor ---
	var nextCursor string
	// If we fetched a full page, the last message ID is the cursor for the next page
	if len(messages) > 0 && len(messages) == int(pageSize) {
		nextCursor = messages[len(messages)-1].ID.String()
		s.Logger.InfoF("Next cursor for pagination: %s", nextCursor)
	}

	return &chat.ListChatMessageResponse{
		Messages:   respMessages,
		NextCursor: &nextCursor,
	}, nil
}

func (s *ChatService) UpdateMessage(ctx context.Context, req *chat.UpdateMessageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("gRPC ChatService: UpdateMessage initiated for message %s", req.MessageUuid)

	// --- 1. Validation and Authentication ---
	if req.GetMessageUuid() == "" || strings.TrimSpace(req.GetNewMessage()) == "" {
		s.Logger.ErrorF("Invalid request: message_uuid and new_message_text are required")
		return nil, status.Errorf(codes.InvalidArgument, "message_uuid and new_message_text are required")
	}

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 2. Find Message and Authorize Action ---
	s.Logger.InfoF("Fetching message %s for update...", req.MessageUuid)
	var messageToUpdate models.ChatMessage
	// Preload TableUser to get the UserID of the author
	if err := s.Db.Preload("TableUser.User").First(&messageToUpdate, "id = ?", req.GetMessageUuid()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Message with ID %s not found", req.MessageUuid)
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		s.Logger.ErrorF("Database error finding message: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	// Authorization: The user_id from the token must match the user_id of the message's author
	if messageToUpdate.TableUser.UserID != userID {
		s.Logger.WarningF("Authorization failed: User %d attempted to edit message %s owned by user %d (via TableUser %d)",
			userID, req.MessageUuid, messageToUpdate.TableUser.UserID, messageToUpdate.TableUser.ID)
		return nil, status.Errorf(codes.PermissionDenied, "you can only edit your own messages")
	}

	// Business Logic: Only allow editing of text messages
	if messageToUpdate.MessageType != consts.MessageType(chat.MessageType_TEXT) {
		s.Logger.WarningF("User %d attempted to edit non-text message %s (type: %d)", userID, req.MessageUuid, messageToUpdate.MessageType)
		return nil, status.Errorf(codes.InvalidArgument, "only text messages can be edited")
	}

	// --- 3. Update, Save, and Publish ---
	s.Logger.InfoF("Authorization successful. Updating message %s...", req.MessageUuid)
	messageToUpdate.Message = req.GetNewMessage()
	messageToUpdate.UpdatedAt = time.Now() // GORM will update this automatically, but being explicit is good

	if err := s.Db.Save(&messageToUpdate).Error; err != nil {
		s.Logger.ErrorF("Database error updating message %s: %v", req.MessageUuid, err)
		return nil, status.Errorf(codes.Internal, "could not update message")
	}

	s.Logger.InfoF("Message %s updated successfully in database", req.MessageUuid)

	// Publish the event to the Broker
	s.Logger.InfoF("Publishing ChatMessageUpdated event for table %d", messageToUpdate.TableID)

	// Re-process attachments for the event payload
	var attachments []string
	if messageToUpdate.Attachments != nil {
		if err := json.Unmarshal(messageToUpdate.Attachments, &attachments); err != nil {
			s.Logger.ErrorF("Failed to unmarshal attachments for updated message event %s: %v", messageToUpdate.ID, err)
			// Don't fail the request, just send an empty list
		}
	}

	// Build the response payload for the event
	respProto := &chat.ChatMessageResponse{
		MessageUuid:      messageToUpdate.ID.String(),
		TableId:          uint64(messageToUpdate.TableID),
		SenderId:         uint64(messageToUpdate.TableUserID),
		SenderUsername:   messageToUpdate.TableUser.User.Username,
		MessageText:      messageToUpdate.Message,
		MessageType:      chat.MessageType(messageToUpdate.MessageType),
		MessageStatus:    chat.MessageStatus(messageToUpdate.MessageStatus),
		SentAt:           timestamppb.New(messageToUpdate.CreatedAt),
		UpdatedAt:        timestamppb.New(messageToUpdate.UpdatedAt), // The new updated time
		Attachments:      attachments,
		MediaUrl:         messageToUpdate.MediaURL,
		ReplyToMessageId: messageToUpdate.ReplyToMessageId,
	}

	// Create the sync event
	syncMsg := &syncBroker.SyncResponse{
		TableId: respProto.TableId,
		Action: &syncBroker.SyncResponse_MessageUpdated{ // Use the ChatMessageUpdated event
			MessageUpdated: &chat.ChatMessageUpdated{
				Message: respProto,
			},
		},
	}

	s.Broker.Publish(pubSubSyncConst.TableSync, respProto.TableId, syncMsg)
	s.Logger.InfoF("Published ChatMessageUpdated event to broker for table %d", respProto.TableId)

	return &emptypb.Empty{}, nil
}

func (s *ChatService) DeleteMessage(ctx context.Context, req *chat.DeleteMessageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("gRPC ChatService: DeleteMessage initiated for message %s", req.MessageUuid)

	// --- 1. Validation and Authentication ---
	if req.GetMessageUuid() == "" {
		s.Logger.ErrorF("Invalid request: message_uuid is required")
		return nil, status.Errorf(codes.InvalidArgument, "message_uuid is required")
	}

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 2. Find Message and Authorize Action ---
	s.Logger.InfoF("Fetching message %s for deletion...", req.MessageUuid)
	var messageToDelete models.ChatMessage
	if err := s.Db.Preload("TableUser").First(&messageToDelete, "id = ?", req.GetMessageUuid()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.WarningF("Message with ID %s not found for deletion (idempotent success)", req.MessageUuid)
			// Return success (idempotency) if message is already gone
			return &emptypb.Empty{}, nil
		}
		s.Logger.ErrorF("Database error finding message for deletion: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	// Authorization: Only the original author can delete.
	if messageToDelete.TableUser.UserID != userID {
		s.Logger.WarningF("Authorization failed: User %d attempted to delete message %s owned by user %d (via TableUser %d)",
			userID, req.MessageUuid, messageToDelete.TableUser.UserID, messageToDelete.TableUser.ID)
		return nil, status.Errorf(codes.PermissionDenied, "you can only delete your own messages")
	}

	// --- 3. Delete (Soft Delete), Save, and Publish ---
	s.Logger.InfoF("Authorization successful. Soft-deleting message %s...", req.MessageUuid)
	if err := s.Db.Delete(&messageToDelete).Error; err != nil { // GORM performs a soft delete
		s.Logger.ErrorF("Database error deleting message %s: %v", req.MessageUuid, err)
		return nil, status.Errorf(codes.Internal, "could not delete message")
	}

	s.Logger.InfoF("Message %s soft-deleted successfully from database", req.MessageUuid)

	// Publish the delete event to the Broker
	syncMsg := &syncBroker.SyncResponse{
		TableId: uint64(messageToDelete.TableID), // Use TableID from the message
		Action: &syncBroker.SyncResponse_MessageDeleted{ // Use the ChatMessageDeleted event
			MessageDeleted: &chat.ChatMessageDeleted{
				MessageUuid: req.GetMessageUuid(),
				TableId:     uint64(messageToDelete.TableID),
			},
		},
	}

	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(messageToDelete.TableID), syncMsg)
	s.Logger.InfoF("Published ChatMessageDeleted event to broker for table %d", messageToDelete.TableID)

	return &emptypb.Empty{}, nil
}

func (s *ChatService) SendPrivateMessage(ctx context.Context, req *chat.SendPrivateMessageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("gRPC ChatService: SendPrivateMessage initiated")

	// --- 1. Validation and Authentication ---
	if req.GetToTableUserId() == 0 {
		s.Logger.ErrorF("Invalid request: recipient tableUser ID (to_table_user_id) is required")
		return nil, status.Errorf(codes.InvalidArgument, "recipient tableUser ID (to_table_user_id) is required")
	}

	userIDFromCtx := ctx.Value("user_id")
	senderUserID, ok := userIDFromCtx.(uint) // This is the models.User ID
	if !ok {
		s.Logger.ErrorF("Sender UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing sender identity")
	}

	// --- 2. Find Sender and Recipient TableUser profiles ---
	// NOTE: This logic assumes a user might be in multiple tables, so we must
	// find the *correct* sender TableUser profile. A simple First() is not enough.
	// A robust solution would require the TableID in the request.

	s.Logger.InfoF("Fetching recipient TableUser profile for ID %d", req.GetToTableUserId())
	var recipientTableUser models.TableUser
	if err := s.Db.Preload("User").First(&recipientTableUser, "id = ?", req.GetToTableUserId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Recipient table user with ID %d not found", req.GetToTableUserId())
			return nil, status.Errorf(codes.NotFound, "recipient not found")
		}
		s.Logger.ErrorF("Database error finding recipient: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	// Now find the sender's TableUser profile *for the same table*
	s.Logger.InfoF("Fetching sender TableUser profile for user %d in table %d", senderUserID, recipientTableUser.TableID)
	var senderTableUser models.TableUser
	if err := s.Db.Preload("User").First(&senderTableUser, "user_id = ? AND table_id = ?", senderUserID, recipientTableUser.TableID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.WarningF("Authorization failed: Sender %d is not in the same table as recipient %d (Table %d)",
				senderUserID, recipientTableUser.ID, recipientTableUser.TableID)
			return nil, status.Errorf(codes.PermissionDenied, "sender and recipient are not in the same table")
		}
		s.Logger.ErrorF("Could not find sender's table profile: %v", err)
		return nil, status.Errorf(codes.FailedPrecondition, "sender profile not found in recipient's table")
	}

	// Check for self-message
	if senderTableUser.ID == recipientTableUser.ID {
		s.Logger.WarningF("User %d (TableUser %d) attempted to send PM to themselves", senderUserID, senderTableUser.ID)
		return nil, status.Errorf(codes.InvalidArgument, "cannot send a private message to yourself")
	}

	// --- 3. (Authorization was handled in step 2) ---

	// --- 4. Assemble and Save Message ---
	s.Logger.InfoF("All checks passed. Assembling private message...")
	var attachmentsJSON datatypes.JSON
	if len(req.GetAttachments()) > 0 {
		jsonBytes, err := json.Marshal(req.GetAttachments())
		if err != nil {
			s.Logger.ErrorF("Failed to marshal attachments for PM: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to process attachments")
		}
		attachmentsJSON = jsonBytes
	}

	recipientTableUserID := recipientTableUser.ID // Get ID to use as a pointer
	privateMessageModel := models.ChatMessage{
		TableID:       senderTableUser.TableID, // Set the TableID for partitioning
		TableUserID:   senderTableUser.ID,      // The author
		ToTableUserId: &recipientTableUserID,   // The recipient (pointer)
		Message:       req.GetMessage(),
		MessageType:   consts.MessageType(req.GetMessageType()),
		MessageStatus: consts.MessageStatus(chat.MessageStatus_SENT),
		Attachments:   attachmentsJSON,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Add optional fields
	if req.MediaUrl != nil {
		privateMessageModel.MediaURL = req.MediaUrl
	}
	if req.ReplyToMessageId != nil {
		privateMessageModel.ReplyToMessageId = req.ReplyToMessageId
	}

	s.Logger.InfoF("Saving private message from TableUser %d to %d in table %d...",
		senderTableUser.ID, recipientTableUser.ID, senderTableUser.TableID)
	if err := s.Db.Create(&privateMessageModel).Error; err != nil {
		s.Logger.ErrorF("Database error creating private message: %v", err)
		return nil, status.Errorf(codes.Internal, "could not save private message")
	}
	s.Logger.InfoF("Private message %s saved to database", privateMessageModel.ID)

	// --- 5. Publish to Broker ---
	recipientID := uint64(recipientTableUser.ID)

	// Build the event payload
	respProto := &chat.ChatMessageResponse{
		MessageUuid:        privateMessageModel.ID.String(),
		TableId:            uint64(senderTableUser.TableID),
		SenderId:           uint64(senderTableUser.ID),
		SenderUsername:     senderTableUser.User.Username,
		MessageText:        privateMessageModel.Message,
		MessageType:        req.GetMessageType(),
		MessageStatus:      chat.MessageStatus_SENT,
		SentAt:             timestamppb.New(privateMessageModel.CreatedAt),
		MediaUrl:           req.MediaUrl,
		Attachments:        req.Attachments,
		ReplyToMessageId:   req.ReplyToMessageId,
		PrivateRecipientId: &recipientID, // Mark this as a PM
	}

	syncMsg := &syncBroker.SyncResponse{
		TableId: respProto.TableId,
		Action: &syncBroker.SyncResponse_MessageSended{
			MessageSended: &chat.ChatMessageSent{
				Message: respProto,
			},
		},
	}

	// Publish to the main table sync
	// The clients will be responsible for filtering messages
	s.Broker.Publish(pubSubSyncConst.TableSync, respProto.TableId, syncMsg)
	s.Logger.InfoF("Published private ChatMessageSent event to broker for table %d", respProto.TableId)

	return &emptypb.Empty{}, nil
}
