package chat

import (
	"fmt"
	"strings"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/chat"
)

func ErrParamIsRequired(name, typ string) error {
	return fmt.Errorf("param %s (type: %s) is required", name, typ)
}

func Validate(req *chat.SendMessageRequest) error {

	if req.GetTableId() == 0 {
		return ErrParamIsRequired("table_id", "uint64")
	}

	switch req.GetMessageType() {
	case chat.MessageType_TEXT:
		if strings.TrimSpace(req.GetMessageText()) == "" {
			return fmt.Errorf("message_text cannot be empty for message type TEXT")
		}
		if req.GetMediaUrl() != "" || len(req.GetAttachments()) > 0 {
			return fmt.Errorf("media_url and attachments are not allowed for message type TEXT")
		}
	case chat.MessageType_IMAGE, chat.MessageType_VIDEO, chat.MessageType_AUDIO:
		if strings.TrimSpace(req.GetMediaUrl()) == "" {
			return fmt.Errorf("media_url cannot be empty for media message types")
		}
		if req.GetMessageText() != "" || len(req.GetAttachments()) > 0 {
			return fmt.Errorf("message_text and attachments are not allowed for media message types")
		}
	case chat.MessageType_DOCUMENT:
		// Para documentos, deve haver pelo menos um anexo.
		if len(req.GetAttachments()) == 0 {
			return fmt.Errorf("at least one attachment is required for message type DOCUMENT")
		}
		// E não deve haver texto ou URL de mídia.
		if req.GetMessageText() != "" || req.GetMediaUrl() != "" {
			return fmt.Errorf("message_text and media_url are not allowed for message type DOCUMENT")
		}

	case chat.MessageType_SYSTEM:
		// Mensagens de sistema não devem ser enviadas por usuários.
		return fmt.Errorf("message type SYSTEM cannot be sent by a user")

	default:
		// Garante que um tipo de mensagem válido foi fornecido.
		return fmt.Errorf("invalid message type provided")
	}

	return nil
}
