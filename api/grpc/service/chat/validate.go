package chat

import (
	"fmt"
	"github.com/GarotoCowboy/vttProject/api/grpc/proto/chat/pb"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func Validate(req *pb.ChatMessageRequest) error {
	if req.TableUserId == 0 /*&& req.Username == "" && req.Message == ""*/ {
		if req.TableUserId == 0 /*&& req.Username == "" && req.Message == ""*/ {
			return fmt.Errorf("request body is empty")
		}
		if req.TableUserId == 0 {
			return ErrParamIsRequired("tableUserId", "string")
		}

		if _, ok := pb.MessageStatus_name[int32(req.MessageStatus)]; !ok {
			if _, ok := pb.MessageStatus_name[int32(req.MessageStatus)]; !ok {
				return ErrParamIsRequired("messageStatus", "pb.MessageStatus enum")
			}
			if _, ok := pb.MessageType_name[int32(req.MessageType)]; !ok {
				if _, ok := pb.MessageType_name[int32(req.MessageType)]; !ok {
					return ErrParamIsRequired("messageType", "pb.MessageType enum")
				}

				return nil

			}
		}
	}
	return nil
}
