package chat

import (
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/chat/pb"
	"github.com/GarotoCowboy/vttProject/api/models"
	"github.com/GarotoCowboy/vttProject/api/models/consts"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"time"
)

// SendMessage sends a Message at table using pub/sub
func (s *ChatService) SendMessage(stream grpc.BidiStreamingServer[pb.ChatMessageRequest, pb.ChatMessageResponse]) error {

	req, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return status.Errorf(codes.Internal, "initial recv error: %v", err)
	}

	if err := Validate(req); err != nil {
		return status.Errorf(codes.InvalidArgument, err.Error())
	}

	var tableUser models.TableUser

	if err := s.Db.Preload("User").Where("id = ?", req.GetTableUserId()).First(&tableUser).Error; err != nil {
		return status.Errorf(codes.NotFound, "table user not found: %v", err)
	}
	tableID := tableUser.TableID

	subID := s.publicChatSubscribe(tableID, stream)
	defer s.publicChatUnsubscribe(tableID, subID)

	// Enters the loop after checking the connected users
	for {
		sendedMessage := time.Now()

		// Prepara ponteiros opcionais
		var attachmentsPtr, mediaURLPtr, replyToPtr *string
		if req.Attachments != "" {
			attachmentsPtr = &req.Attachments
		}
		if req.MediaUrl != "" {
			mediaURLPtr = &req.MediaUrl
		}
		if req.ReplyToMessageId != "" {
			replyToPtr = &req.ReplyToMessageId
		}

		//Send Audio
		if req.MessageType == pb.MessageType_AUDIO {
			//Cannot send audio and text at the same time
			if req.Message != "" {
				return status.Errorf(codes.InvalidArgument, "cannot send audio and message at the same time")
			}
			//pick audio URL
			mediaURLPtr = &req.MediaUrl

		}

		// assemble and save the message
		chatMessage := models.ChatMessage{
			TableUserID:      tableUser.ID,
			Message:          req.Message,
			Username:         tableUser.User.Username,
			MessageStatus:    consts.MessageStatus(pb.MessageStatus_SENT),
			MessageType:      consts.MessageType(req.MessageType),
			Attachments:      attachmentsPtr,
			MediaURL:         mediaURLPtr,
			ReplyToMessageId: replyToPtr,
			CreatedAt:        sendedMessage,
			UpdatedAt:        sendedMessage,
		}

		if err := s.Db.Create(&chatMessage).Error; err != nil {
			return status.Errorf(codes.Internal, "error creating message: %v", err)
		}

		s.Logger.InfoF("Message sent: %v", chatMessage)

		// creating response
		resp := &pb.ChatMessageResponse{
			Message:          chatMessage.Message,
			Username:         chatMessage.Username,
			MediaUrl:         req.MediaUrl,
			Attachments:      req.Attachments,
			MessageType:      req.MessageType,
			MessageStatus:    pb.MessageStatus_SENT,
			SendAt:           timestamppb.New(sendedMessage),
			MessageId:        chatMessage.ID.String(),
			IsDeleted:        false,
			ReplyToMessageId: req.ReplyToMessageId,
		}

		// send the message to all users on the table
		s.mu.RLock()
		for id, sub := range s.publicSubscribers[tableID] {
			if err := sub.Send(resp); err != nil {
				s.Logger.ErrorF("error sending to subscriber %s: %v", id, err)
				s.publicChatUnsubscribe(tableID, id)
			}
		}
		s.mu.RUnlock()
		s.Logger.InfoF("broadcast update for table_user: %v", tableID)

		// wait next message
		req, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}
		if err := Validate(req); err != nil {
			return status.Errorf(codes.InvalidArgument, err.Error())
		}
	}
}

func (s *ChatService) ListMessages(req *pb.ListMessagesRequest, stream grpc.ServerStreamingServer[pb.ChatMessageResponse]) error {
	var tableUser models.TableUser

	// Fetches the TableUser and gets the associated table_id
	if err := s.Db.Where("id = ?", req.GetTableId()).First(&tableUser).Error; err != nil {
		return status.Errorf(codes.NotFound, "TableUser not found: %v", err)
	}

	var allTableUsers []models.TableUser
	if err := s.Db.Where("table_id = ?", tableUser.TableID).Find(&allTableUsers).Error; err != nil {
		return status.Errorf(codes.Internal, "error finding table users: %v", err)
	}

	if len(allTableUsers) == 0 {
		return status.Errorf(codes.NotFound, "no users found in this table")
	}

	// Extracts all user IDs participating in the table
	var tableUserIDs []uint
	for _, tu := range allTableUsers {
		tableUserIDs = append(tableUserIDs, tu.ID)
	}

	// Search all messages from these users
	var listMessages []models.ChatMessage
	if err := s.Db.
		Where("table_user_id IN ?", tableUserIDs).
		Order("created_at asc").
		Find(&listMessages).Error; err != nil {
		return status.Errorf(codes.Internal, "error finding messages: %v", err)
	}

	if len(listMessages) == 0 {
		return status.Errorf(codes.NotFound, "no messages were found in this table")
	}

	for _, msg := range listMessages {
		resp := &pb.ChatMessageResponse{
			Message:       msg.Message,
			Username:      msg.Username,
			MessageType:   pb.MessageType(msg.MessageType),
			MessageStatus: pb.MessageStatus(msg.MessageStatus),
			SendAt:        timestamppb.New(msg.CreatedAt),
			MessageId:     msg.ID.String(),
			IsDeleted:     false,
		}

		if msg.MediaURL != nil {
			resp.MediaUrl = *msg.MediaURL
		}
		if msg.Attachments != nil {
			resp.Attachments = *msg.Attachments
		}
		if msg.ReplyToMessageId != nil {
			resp.ReplyToMessageId = *msg.ReplyToMessageId
		}

		if err := stream.Send(resp); err != nil {
			return status.Errorf(codes.Internal, "error sending message: %v", err)
		}
	}

	return nil
}

func (s *ChatService) SendPrivateMessage(stream grpc.BidiStreamingServer[pb.ChatMessagePrivateRequest, pb.ChatMessageResponse]) error {

	req, err := stream.Recv()

	if err == io.EOF {
		return nil
	}
	if err != nil {
		return status.Errorf(codes.Internal, "initial recv error: %v", err)
	}

	if req.TableUserId <= 0 || req.ToTableUserId <= 0 {
		return status.Errorf(codes.InvalidArgument, "destin or origin endress is invalid:")
	}

	if req.TableUserId == req.ToTableUserId {
		return status.Errorf(codes.InvalidArgument, "cannot send private message to yourself")
	}

	var fromUser models.TableUser
	var toUser models.TableUser

	if err := s.Db.Preload("User").Where("id = ?", req.GetTableUserId()).First(&fromUser).Error; err != nil {
		return status.Errorf(codes.NotFound, "table user not found: %v", err)
	}

	if err := s.Db.Preload("User").Where("id = ?", req.GetToTableUserId()).First(&toUser).Error; err != nil {
		return status.Errorf(codes.NotFound, "table user not found: %v", err)
	}

	if fromUser.TableID != toUser.TableID {
		return status.Errorf(codes.FailedPrecondition, "users are not on the same table")
	}

	subID := s.privateChatSubscribe(fromUser.ID, toUser.ID, stream)
	defer s.privateChatUnsubscribe(fromUser.ID, toUser.ID, subID)

	for {
		sendedMessage := time.Now()

		req, err = stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return status.Errorf(codes.Internal, "recv error: %v", err)
		}

		var attachmentsPtr, mediaURLPtr, replyToPtr *string
		if req.Attachments != "" {
			attachmentsPtr = &req.Attachments
		}
		if req.MediaUrl != "" {
			mediaURLPtr = &req.MediaUrl
		}
		if req.ReplyToMessageId != "" {
			replyToPtr = &req.ReplyToMessageId
		}

		chatMessage := &models.ChatMessage{
			TableUserID:      fromUser.ID,
			ToTableUserId:    &toUser.ID,
			Username:         fromUser.User.Username,
			Message:          req.GetMessage(),
			Attachments:      attachmentsPtr,
			MediaURL:         mediaURLPtr,
			ReplyToMessageId: replyToPtr,
			MessageType:      consts.MessageType(req.MessageType),
			MessageStatus:    consts.MessageStatus(req.MessageStatus),
			CreatedAt:        sendedMessage,
			UpdatedAt:        sendedMessage,
		}

		if err := s.Db.Create(&chatMessage).Error; err != nil {
			return status.Errorf(codes.Internal, "error creating message: %v", err)
		}
		s.Logger.InfoF("Message sent: %v", chatMessage)

		// creating response
		resp := &pb.ChatMessageResponse{
			Message:          chatMessage.Message,
			Username:         chatMessage.Username,
			MediaUrl:         req.MediaUrl,
			Attachments:      req.Attachments,
			MessageType:      req.MessageType,
			MessageStatus:    pb.MessageStatus_SENT,
			SendAt:           timestamppb.New(sendedMessage),
			MessageId:        chatMessage.ID.String(),
			IsDeleted:        false,
			ReplyToMessageId: req.ReplyToMessageId,
		}

		// send the message to all users on the table
		s.mu.RLock()
		if subsMap, ok := s.privateSubscribers[toUser.ID]; ok {
			if subs, ok := subsMap[fromUser.ID]; ok {
				for id, sub := range subs {
					if err := sub.Send(resp); err != nil {
						s.Logger.ErrorF("error sending to subscriber %s: %v", id, err)
						s.privateChatUnsubscribe(toUser.ID, fromUser.ID, id)
					}
				}
			}
		}
		s.mu.RUnlock()
		s.Logger.InfoF("broadcast update for table_user: %v", toUser.UserID)

	}
}
