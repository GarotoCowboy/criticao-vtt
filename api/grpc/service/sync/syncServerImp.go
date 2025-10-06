package sync

import (
	"io"

	"github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"
	"github.com/GarotoCowboy/vttProject/api/models/consts/pubSubSyncConst"
	"google.golang.org/grpc"
)

func (s *SyncServer) Sync(stream grpc.BidiStreamingServer[sync.SyncRequest, sync.SyncResponse]) error {
	s.Logger.InfoF("GRPC Sync Service - Sync - start")

	req, err := stream.Recv()

	if err == io.EOF {
		s.Logger.InfoF("connection closed from client")
		return nil
	}

	if err != nil {
		s.Logger.ErrorF("error to receive first syncRequest: %v", err)
		return nil
	}

	sceneID := req.GetSceneId()
	tableId := req.GetTableId()

	s.Logger.InfoF("Client connected for table: %v and scene: %v", tableId, sceneID)

	msgChan := make(chan *sync.SyncResponse, 100)

	s.Broker.SubscribeToTopic(pubSubSyncConst.TableSync, tableId, msgChan)
	s.Broker.SubscribeToTopic(pubSubSyncConst.SceneSync, sceneID, msgChan)

	defer func() {
		s.Logger.InfoF("Cleaning up subscriptions for clients from table: %v and scene:", tableId, sceneID)
		s.Broker.UnsubscribeToTopic(pubSubSyncConst.TableSync, tableId, msgChan)
		s.Broker.UnsubscribeToTopic(pubSubSyncConst.SceneSync, sceneID, msgChan)

		close(msgChan)
	}()

	ctx := stream.Context()

	go func() {
		<-ctx.Done()
		s.Logger.InfoF("client disconnected from scene %v", sceneID)
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case msg := <-msgChan:
			if err := stream.Send(msg); err != nil {
				s.Logger.ErrorF("error to send message from client scene: %d %v", sceneID, err)
				return err
			}
		}
	}
}
