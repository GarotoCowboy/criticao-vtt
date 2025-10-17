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
	s.Logger.InfoF("Starting GPRC to SendMessage")

	userIDFromCtx := ctx.Value("user_id")

	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	if err := Validate(req); err != nil {
		s.Logger.ErrorF("Validation error: %v", err)
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	var tableUserModel models.TableUser

	if err := s.Db.WithContext(ctx).Preload("User").Where("user_id = ? AND table_id = ?", userID, req.TableId).First(&tableUserModel).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("user is not in a table ")
			return nil, status.Errorf(codes.NotFound, "user is not in a table ")
		}
		s.Logger.ErrorF("database error: %v", err.Error())
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var attachmentsJSON datatypes.JSON
	if len(req.Attachments) > 0 {
		jsonBytes, err := json.Marshal(req.GetAttachments())
		if err != nil {
			s.Logger.ErrorF("json marshal error: %v", err.Error())
			return nil, status.Errorf(codes.Internal, "failed to process attachments")
		}
		attachmentsJSON = jsonBytes
	}

	chatMessageModel := models.ChatMessage{
		TableUserID:   tableUserModel.ID,
		TableID:       tableUserModel.TableID,
		Message:       req.GetMessageText(),
		MessageType:   consts.MessageType(req.GetMessageType()),
		MessageStatus: consts.MessageStatus(chat.MessageStatus_SENT),
		Attachments:   attachmentsJSON,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	if req.MediaUrl != nil {
		chatMessageModel.MediaURL = req.MediaUrl
	}
	if req.ReplyToMessageId != nil {
		chatMessageModel.ReplyToMessageId = req.ReplyToMessageId
	}

	if err := s.Db.WithContext(ctx).Create(&chatMessageModel).Error; err != nil {
		s.Logger.ErrorF("Database error creating message: %v", err)
		return nil, status.Errorf(codes.Internal, "could not save message")
	}
	s.Logger.InfoF("Message %s saved to database for table %d", chatMessageModel.ID, req.TableId)

	respProto := &chat.ChatMessageResponse{
		MessageUuid:      chatMessageModel.ID.String(),
		TableId:          req.GetTableId(),
		SenderId:         uint64(tableUserModel.ID),
		SenderUsername:   tableUserModel.User.Username,
		MessageText:      chatMessageModel.Message,
		MessageType:      req.GetMessageType(),
		MessageStatus:    chat.MessageStatus_SENT,
		SentAt:           timestamppb.New(chatMessageModel.CreatedAt),
		MediaUrl:         req.MediaUrl,
		Attachments:      req.Attachments,
		ReplyToMessageId: req.ReplyToMessageId,
	}

	syncMsg := &syncBroker.SyncResponse{
		TableId: req.GetTableId(),
		Action: &syncBroker.SyncResponse_MessageSended{
			MessageSended: &chat.ChatMessageSent{
				Message: respProto,
			},
		},
	}

	s.Broker.Publish(pubSubSyncConst.TableSync, req.GetTableId(), syncMsg)
	s.Logger.InfoF("Published ChatMessageSent event to broker for table %d", req.TableId)

	return &emptypb.Empty{}, nil
}
func (s *ChatService) ListMessages(ctx context.Context, req *chat.ListMessagesRequest) (*chat.ListChatMessageResponse, error) {
	s.Logger.InfoF("gRPC ChatService: ListMessages initiated for table %d", req.TableId)

	if req.GetTableId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "table_id is required")
	}

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	var count int64
	s.Db.WithContext(ctx).Model(&models.TableUser{}).Where("user_id = ? AND table_id = ?", userID, req.TableId).Count(&count)
	if count == 0 {
		s.Logger.ErrorF("Authorization error: user %d is not a member of table %d", userID, req.TableId)
		return nil, status.Errorf(codes.PermissionDenied, "user is not a member of the specified table")
	}

	pageSize := req.GetPageSize()
	if pageSize <= 0 || pageSize > 100 {
		pageSize = 1
	}

	var messages []models.ChatMessage

	query := s.Db.WithContext(ctx).Preload("TableUser.User").Where("table_id = ?", req.GetTableId()).Order("created_at DESC").Limit(int(pageSize))

	if req.GetLastMessageId() != "" {
		var cursorMessage models.ChatMessage
		if err := s.Db.First(&cursorMessage, "id = ?", req.GetLastMessageId()).Error; err == nil {
			query = query.Where("created_at < ?", cursorMessage.CreatedAt)
		}
	}

	if err := query.Find(&messages).Error; err != nil {
		s.Logger.ErrorF("Database error listing messages for table %d: %v", req.TableId, err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve messages")
	}

	respMessages := make([]*chat.ChatMessageResponse, 0, len(messages))

	for _, msg := range messages {

		var attachments []string
		if msg.Attachments != nil {
			if err := json.Unmarshal(msg.Attachments, &attachments); err != nil {
				s.Logger.ErrorF("Failed to unmarshal attachments for message %s: %v", msg.ID, err)

				continue
			}
		}

		respMsg := &chat.ChatMessageResponse{
			MessageUuid:      msg.ID.String(),
			TableId:          uint64(msg.TableID),
			SenderId:         uint64(msg.TableUserID),
			SenderUsername:   msg.TableUser.User.Username,
			MessageText:      msg.Message,
			MessageType:      chat.MessageType(msg.MessageType),
			MessageStatus:    chat.MessageStatus(msg.MessageStatus),
			SentAt:           timestamppb.New(msg.CreatedAt),
			MediaUrl:         msg.MediaURL,
			Attachments:      attachments,
			ReplyToMessageId: msg.ReplyToMessageId,
			IsDeleted:        !msg.DeletedAt.Time.IsZero(),
		}
		if msg.UpdatedAt.After(msg.CreatedAt) {
			respMsg.UpdatedAt = timestamppb.New(msg.UpdatedAt)
		}

		respMessages = append(respMessages, respMsg)
	}

	var nextCursor string
	if len(messages) > 0 && len(messages) == int(pageSize) {
		nextCursor = messages[len(messages)-1].ID.String()
	}

	return &chat.ListChatMessageResponse{
		Messages:   respMessages,
		NextCursor: &nextCursor,
	}, nil

}
func (s *ChatService) UpdateMessage(ctx context.Context, req *chat.UpdateMessageRequest) (*emptypb.Empty, error) {

	s.Logger.InfoF("gRPC ChatService: UpdateMessage initiated for message %s", req.MessageUuid)

	// --- 1. Validação e Autenticação ---
	if req.GetMessageUuid() == "" || strings.TrimSpace(req.GetNewMessage()) == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message_uuid and new_message_text are required")
	}

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 2. Encontrar a Mensagem e Autorizar a Ação ---
	var messageToUpdate models.ChatMessage
	// Usamos Preload para já trazer os dados do autor
	if err := s.Db.Preload("TableUser.User").First(&messageToUpdate, "id = ?", req.GetMessageUuid()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Message with ID %s not found", req.MessageUuid)
			return nil, status.Errorf(codes.NotFound, "message not found")
		}
		s.Logger.ErrorF("Database error finding message: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	// Autorização: O user_id do autor da mensagem (TableUser.UserID) deve ser o mesmo do usuário no token.
	if messageToUpdate.TableUser.UserID != userID {
		s.Logger.WarningF("Authorization failed: User %d attempted to edit message %s owned by user %d", userID, req.MessageUuid, messageToUpdate.TableUser.UserID)
		return nil, status.Errorf(codes.PermissionDenied, "you can only edit your own messages")
	}

	// Impede a edição de mensagens que não são de texto
	if messageToUpdate.MessageType != consts.MessageType(chat.MessageType_TEXT) {
		return nil, status.Errorf(codes.InvalidArgument, "only text messages can be edited")
	}

	// --- 3. Atualizar, Salvar e Publicar ---
	messageToUpdate.Message = req.GetNewMessage()
	messageToUpdate.UpdatedAt = time.Now()

	if err := s.Db.Save(&messageToUpdate).Error; err != nil {
		s.Logger.ErrorF("Database error updating message %s: %v", req.MessageUuid, err)
		return nil, status.Errorf(codes.Internal, "could not update message")
	}

	s.Logger.InfoF("Message %s updated successfully in database", req.MessageUuid)

	// Publicar o evento no Broker

	// (Reutilizamos a lógica de converter anexos de JSON para []string do ListMessages)
	var attachments []string
	if messageToUpdate.Attachments != nil {
		err := json.Unmarshal(messageToUpdate.Attachments, &attachments)
		if err != nil {
			return nil, err
		}
	}

	respProto := &chat.ChatMessageResponse{
		MessageUuid:      messageToUpdate.ID.String(),
		TableId:          uint64(messageToUpdate.TableID),
		SenderId:         uint64(messageToUpdate.TableUserID),
		SenderUsername:   messageToUpdate.TableUser.User.Username,
		MessageText:      messageToUpdate.Message,
		MessageType:      chat.MessageType(messageToUpdate.MessageType),
		MessageStatus:    chat.MessageStatus(messageToUpdate.MessageStatus),
		SentAt:           timestamppb.New(messageToUpdate.CreatedAt),
		UpdatedAt:        timestamppb.New(messageToUpdate.UpdatedAt),
		Attachments:      attachments,
		MediaUrl:         messageToUpdate.MediaURL,
		ReplyToMessageId: messageToUpdate.ReplyToMessageId,
	}

	syncMsg := &syncBroker.SyncResponse{
		TableId: respProto.TableId,
		Action: &syncBroker.SyncResponse_MessageUpdated{ // Usamos o evento ChatMessageUpdated
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

	// --- 1. Validação e Autenticação ---
	if req.GetMessageUuid() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "message_uuid is required")
	}

	userIDFromCtx := ctx.Value("user_id")
	userID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing user identity")
	}

	// --- 2. Encontrar a Mensagem e Autorizar a Ação ---
	var messageToDelete models.ChatMessage
	if err := s.Db.Preload("TableUser").First(&messageToDelete, "id = ?", req.GetMessageUuid()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Message with ID %s not found for deletion", req.MessageUuid)
			// Retornamos sucesso mesmo se a mensagem não existir para evitar que um cliente saiba se uma mensagem existiu ou não.
			return &emptypb.Empty{}, nil
		}
		s.Logger.ErrorF("Database error finding message for deletion: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	// Autorização: Apenas o autor original pode deletar a mensagem.
	// (Futuramente, você pode adicionar uma lógica para permitir que o "Mestre" da mesa também delete)
	if messageToDelete.TableUser.UserID != userID {
		s.Logger.WarningF("Authorization failed: User %d attempted to delete message %s owned by user %d", userID, req.MessageUuid, messageToDelete.TableUser.UserID)
		return nil, status.Errorf(codes.PermissionDenied, "you can only delete your own messages")
	}

	// --- 3. Deletar (Soft Delete), Salvar e Publicar ---
	if err := s.Db.Delete(&messageToDelete).Error; err != nil {
		s.Logger.ErrorF("Database error deleting message %s: %v", req.MessageUuid, err)
		return nil, status.Errorf(codes.Internal, "could not delete message")
	}

	s.Logger.InfoF("Message %s soft-deleted successfully from database", req.MessageUuid)

	// Publicar o evento no Broker para que a UI de todos os clientes possa atualizar
	syncMsg := &syncBroker.SyncResponse{
		TableId: uint64(messageToDelete.TableUser.TableID),
		Action: &syncBroker.SyncResponse_MessageDeleted{ // Usamos o evento ChatMessageDeleted
			MessageDeleted: &chat.ChatMessageDeleted{
				MessageUuid: req.GetMessageUuid(),
				TableId:     uint64(messageToDelete.TableUser.TableID),
			},
		},
	}

	s.Broker.Publish(pubSubSyncConst.TableSync, uint64(messageToDelete.TableUser.TableID), syncMsg)
	s.Logger.InfoF("Published ChatMessageDeleted event to broker for table %d", messageToDelete.TableUser.TableID)

	return &emptypb.Empty{}, nil
}
func (s *ChatService) SendPrivateMessage(ctx context.Context, req *chat.SendPrivateMessageRequest) (*emptypb.Empty, error) {
	s.Logger.InfoF("gRPC ChatService: SendPrivateMessage initiated")

	// --- 1. Validação e Autenticação ---
	if req.GetToTableUserId() == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "recipient user ID (to_table_user_id) is required")
	}

	userIDFromCtx := ctx.Value("user_id")
	senderUserID, ok := userIDFromCtx.(uint)
	if !ok {
		s.Logger.ErrorF("Sender UserID not found in context")
		return nil, status.Errorf(codes.Internal, "error processing sender identity")
	}

	// --- 2. Encontrar Remetente e Destinatário ---
	// NOTA: Esta busca assume que um usuário só pode estar em uma mesa por vez,
	// ou que o contexto gRPC já foi filtrado para uma mesa específica.
	// Uma solução mais robusta seria passar o table_id na requisição.
	var senderTableUser models.TableUser
	if err := s.Db.Preload("User").First(&senderTableUser, "user_id = ?", senderUserID).Error; err != nil {
		s.Logger.ErrorF("Could not find sender's table profile: %v", err)
		return nil, status.Errorf(codes.FailedPrecondition, "sender profile not found")
	}

	var recipientTableUser models.TableUser
	if err := s.Db.Preload("User").First(&recipientTableUser, "id = ?", req.GetToTableUserId()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.ErrorF("Recipient table user with ID %d not found", req.GetToTableUserId())
			return nil, status.Errorf(codes.NotFound, "recipient not found")
		}
		s.Logger.ErrorF("Database error finding recipient: %v", err)
		return nil, status.Errorf(codes.Internal, "database error")
	}

	if senderTableUser.ID == recipientTableUser.ID {
		return nil, status.Errorf(codes.InvalidArgument, "cannot send a private message to yourself")
	}

	// --- 3. Autorização Crucial ---
	if senderTableUser.TableID != recipientTableUser.TableID {
		s.Logger.WarningF("Authorization failed: User %d (table %d) attempted to send PM to user %d (table %d)",
			senderUserID, senderTableUser.TableID, recipientTableUser.UserID, recipientTableUser.TableID)
		return nil, status.Errorf(codes.PermissionDenied, "sender and recipient are not in the same table")
	}

	// --- 4. Montar e Salvar a Mensagem ---
	var attachmentsJSON datatypes.JSON
	if len(req.GetAttachments()) > 0 {
		jsonBytes, err := json.Marshal(req.GetAttachments())
		if err != nil {
			s.Logger.ErrorF("Failed to marshal attachments for PM: %v", err)
			return nil, status.Errorf(codes.Internal, "failed to process attachments")
		}
		attachmentsJSON = jsonBytes
	}

	recipientTableUserID := recipientTableUser.ID // Para usar no ponteiro
	privateMessageModel := models.ChatMessage{
		TableUserID:   senderTableUser.ID,    // O autor
		ToTableUserId: &recipientTableUserID, // O destinatário
		Message:       req.GetMessage(),
		MessageType:   consts.MessageType(req.GetMessageType()),
		MessageStatus: consts.MessageStatus(chat.MessageStatus_SENT),
		Attachments:   attachmentsJSON,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	// Adiciona campos opcionais
	if req.MediaUrl != nil {
		privateMessageModel.MediaURL = req.MediaUrl
	}
	if req.ReplyToMessageId != nil {
		privateMessageModel.ReplyToMessageId = req.ReplyToMessageId
	}

	if err := s.Db.Create(&privateMessageModel).Error; err != nil {
		s.Logger.ErrorF("Database error creating private message: %v", err)
		return nil, status.Errorf(codes.Internal, "could not save private message")
	}
	s.Logger.InfoF("Private message %s from %d to %d saved to database", privateMessageModel.ID, senderTableUser.ID, recipientTableUser.ID)

	// --- 5. Publicar no Broker ---

	recipientID := uint64(recipientTableUser.ID)

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
		PrivateRecipientId: &recipientID,
	}

	syncMsg := &syncBroker.SyncResponse{
		TableId: respProto.TableId,
		Action: &syncBroker.SyncResponse_MessageSended{
			MessageSended: &chat.ChatMessageSent{
				Message: respProto,
			},
		},
	}

	s.Broker.Publish(pubSubSyncConst.TableSync, respProto.TableId, syncMsg)
	s.Logger.InfoF("Published private ChatMessageSent event to broker for table %d", respProto.TableId)

	return &emptypb.Empty{}, nil
}

//// SendMessage sends a Message at table using pub/sub
//func (s *ChatService) SendMessage(stream grpc.BidiStreamingServer[chat.ChatMessageRequest, chat.ChatMessageResponse]) error {
//
//	req, err := stream.Recv()
//	if err == io.EOF {
//		return nil
//	}
//	if err != nil {
//		return status.Errorf(codes.Internal, "initial recv error: %v", err)
//	}
//
//	if err := Validate(req); err != nil {
//		return status.Errorf(codes.InvalidArgument, err.Error())
//	}
//
//	var tableUser models.TableUser
//
//	if err := s.Db.Preload("User").Where("id = ?", req.GetTableUserId()).First(&tableUser).Error; err != nil {
//		return status.Errorf(codes.NotFound, "table user not found: %v", err)
//	}
//	tableID := tableUser.TableID
//
//	subID := s.publicChatSubscribe(tableID, stream)
//	defer s.publicChatUnsubscribe(tableID, subID)
//
//	// Enters the loop after checking the connected users
//	for {
//		sendedMessage := time.Now()
//
//		// Prepara ponteiros opcionais
//		var attachmentsPtr, mediaURLPtr, replyToPtr *string
//		if req.Attachments != "" {
//			attachmentsPtr = &req.Attachments
//		}
//		if req.MediaUrl != "" {
//			mediaURLPtr = &req.MediaUrl
//		}
//		if req.ReplyToMessageId != "" {
//			replyToPtr = &req.ReplyToMessageId
//		}
//
//		//Send Audio
//		if req.MessageType == chat.MessageType_AUDIO {
//			//Cannot send audio and text at the same time
//			if req.Message != "" {
//				return status.Errorf(codes.InvalidArgument, "cannot send audio and message at the same time")
//			}
//			//pick audio URL
//			mediaURLPtr = &req.MediaUrl
//
//		}
//
//		// assemble and save the message
//		chatMessage := models.ChatMessage{
//			TableUserID:      tableUser.ID,
//			Message:          req.Message,
//			Username:         tableUser.User.Username,
//			MessageStatus:    consts.MessageStatus(chat.MessageStatus_SENT),
//			MessageType:      consts.MessageType(req.MessageType),
//			Attachments:      attachmentsPtr,
//			MediaURL:         mediaURLPtr,
//			ReplyToMessageId: replyToPtr,
//			CreatedAt:        sendedMessage,
//			UpdatedAt:        sendedMessage,
//		}
//
//		if err := s.Db.Create(&chatMessage).Error; err != nil {
//			return status.Errorf(codes.Internal, "error creating message: %v", err)
//		}
//
//		s.Logger.InfoF("Message sent: %v", chatMessage)
//
//		// creating response
//		resp := &chat.ChatMessageResponse{
//			Message:          chatMessage.Message,
//			Username:         chatMessage.Username,
//			MediaUrl:         req.MediaUrl,
//			Attachments:      req.Attachments,
//			MessageType:      req.MessageType,
//			MessageStatus:    chat.MessageStatus_SENT,
//			SendAt:           timestamppb.New(sendedMessage),
//			MessageId:        chatMessage.ID.String(),
//			IsDeleted:        false,
//			ReplyToMessageId: req.ReplyToMessageId,
//		}
//
//		// send the message to all users on the table
//		s.mu.RLock()
//		for id, sub := range s.publicSubscribers[tableID] {
//			if err := sub.Send(resp); err != nil {
//				s.Logger.ErrorF("error sending to subscriber %s: %v", id, err)
//				s.publicChatUnsubscribe(tableID, id)
//			}
//		}
//		s.mu.RUnlock()
//		s.Logger.InfoF("broadcast update for table_user: %v", tableID)
//
//		// wait next message
//		req, err = stream.Recv()
//		if err == io.EOF {
//			return nil
//		}
//		if err != nil {
//			return status.Errorf(codes.Internal, "recv error: %v", err)
//		}
//		if err := Validate(req); err != nil {
//			return status.Errorf(codes.InvalidArgument, err.Error())
//		}
//	}
//}
//
//func (s *ChatService) ListMessages(req *chat.ListMessagesRequest, stream grpc.ServerStreamingServer[chat.ChatMessageResponse]) error {
//	var tableUser models.TableUser
//
//	// Fetches the TableUser and gets the associated table_id
//	if err := s.Db.Where("id = ?", req.GetTableId()).First(&tableUser).Error; err != nil {
//		return status.Errorf(codes.NotFound, "TableUser not found: %v", err)
//	}
//
//	var allTableUsers []models.TableUser
//	if err := s.Db.Where("table_id = ?", tableUser.TableID).Find(&allTableUsers).Error; err != nil {
//		return status.Errorf(codes.Internal, "error finding table users: %v", err)
//	}
//
//	if len(allTableUsers) == 0 {
//		return status.Errorf(codes.NotFound, "no users found in this table")
//	}
//
//	// Extracts all user IDs participating in the table
//	var tableUserIDs []uint
//	for _, tu := range allTableUsers {
//		tableUserIDs = append(tableUserIDs, tu.ID)
//	}
//
//	// Search all messages from these users
//	var listMessages []models.ChatMessage
//	if err := s.Db.
//		Where("table_user_id IN ?", tableUserIDs).
//		Order("created_at asc").
//		Find(&listMessages).Error; err != nil {
//		return status.Errorf(codes.Internal, "error finding messages: %v", err)
//	}
//
//	if len(listMessages) == 0 {
//		return status.Errorf(codes.NotFound, "no messages were found in this table")
//	}
//
//	for _, msg := range listMessages {
//		resp := &chat.ChatMessageResponse{
//			Message:       msg.Message,
//			Username:      msg.Username,
//			MessageType:   chat.MessageType(msg.MessageType),
//			MessageStatus: chat.MessageStatus(msg.MessageStatus),
//			SendAt:        timestamppb.New(msg.CreatedAt),
//			MessageId:     msg.ID.String(),
//			IsDeleted:     false,
//		}
//
//		if msg.MediaURL != nil {
//			resp.MediaUrl = *msg.MediaURL
//		}
//		if msg.Attachments != nil {
//			resp.Attachments = *msg.Attachments
//		}
//		if msg.ReplyToMessageId != nil {
//			resp.ReplyToMessageId = *msg.ReplyToMessageId
//		}
//
//		if err := stream.Send(resp); err != nil {
//			return status.Errorf(codes.Internal, "error sending message: %v", err)
//		}
//	}
//
//	return nil
//}
//
//func (s *ChatService) SendPrivateMessage(stream grpc.BidiStreamingServer[chat.ChatMessagePrivateRequest, chat.ChatMessageResponse]) error {
//
//	req, err := stream.Recv()
//
//	if err == io.EOF {
//		return nil
//	}
//	if err != nil {
//		return status.Errorf(codes.Internal, "initial recv error: %v", err)
//	}
//
//	if req.TableUserId <= 0 || req.ToTableUserId <= 0 {
//		return status.Errorf(codes.InvalidArgument, "destin or origin endress is invalid:")
//	}
//
//	if req.TableUserId == req.ToTableUserId {
//		return status.Errorf(codes.InvalidArgument, "cannot send private message to yourself")
//	}
//
//	var fromUser models.TableUser
//	var toUser models.TableUser
//
//	if err := s.Db.Preload("User").Where("id = ?", req.GetTableUserId()).First(&fromUser).Error; err != nil {
//		return status.Errorf(codes.NotFound, "table user not found: %v", err)
//	}
//
//	if err := s.Db.Preload("User").Where("id = ?", req.GetToTableUserId()).First(&toUser).Error; err != nil {
//		return status.Errorf(codes.NotFound, "table user not found: %v", err)
//	}
//
//	if fromUser.TableID != toUser.TableID {
//		return status.Errorf(codes.FailedPrecondition, "users are not on the same table")
//	}
//
//	subID := s.privateChatSubscribe(fromUser.ID, toUser.ID, stream)
//	defer s.privateChatUnsubscribe(fromUser.ID, toUser.ID, subID)
//
//	for {
//		sendedMessage := time.Now()
//
//		req, err = stream.Recv()
//		if err == io.EOF {
//			return nil
//		}
//		if err != nil {
//			return status.Errorf(codes.Internal, "recv error: %v", err)
//		}
//
//		var attachmentsPtr, mediaURLPtr, replyToPtr *string
//		if req.Attachments != "" {
//			attachmentsPtr = &req.Attachments
//		}
//		if req.MediaUrl != "" {
//			mediaURLPtr = &req.MediaUrl
//		}
//		if req.ReplyToMessageId != "" {
//			replyToPtr = &req.ReplyToMessageId
//		}
//
//		chatMessage := &models.ChatMessage{
//			TableUserID:      fromUser.ID,
//			ToTableUserId:    &toUser.ID,
//			Username:         fromUser.User.Username,
//			Message:          req.GetMessage(),
//			Attachments:      attachmentsPtr,
//			MediaURL:         mediaURLPtr,
//			ReplyToMessageId: replyToPtr,
//			MessageType:      consts.MessageType(req.MessageType),
//			MessageStatus:    consts.MessageStatus(req.MessageStatus),
//			CreatedAt:        sendedMessage,
//			UpdatedAt:        sendedMessage,
//		}
//
//		if err := s.Db.Create(&chatMessage).Error; err != nil {
//			return status.Errorf(codes.Internal, "error creating message: %v", err)
//		}
//		s.Logger.InfoF("Message sent: %v", chatMessage)
//
//		// creating response
//		resp := &chat.ChatMessageResponse{
//			Message:          chatMessage.Message,
//			Username:         chatMessage.Username,
//			MediaUrl:         req.MediaUrl,
//			Attachments:      req.Attachments,
//			MessageType:      req.MessageType,
//			MessageStatus:    chat.MessageStatus_SENT,
//			SendAt:           timestamppb.New(sendedMessage),
//			MessageId:        chatMessage.ID.String(),
//			IsDeleted:        false,
//			ReplyToMessageId: req.ReplyToMessageId,
//		}
//
//		// send the message to all users on the table
//		s.mu.RLock()
//		if subsMap, ok := s.privateSubscribers[toUser.ID]; ok {
//			if subs, ok := subsMap[fromUser.ID]; ok {
//				for id, sub := range subs {
//					if err := sub.Send(resp); err != nil {
//						s.Logger.ErrorF("error sending to subscriber %s: %v", id, err)
//						s.privateChatUnsubscribe(toUser.ID, fromUser.ID, id)
//					}
//				}
//			}
//		}
//		s.mu.RUnlock()
//		s.Logger.InfoF("broadcast update for table_user: %v", toUser.UserID)
//
//	}
//}
