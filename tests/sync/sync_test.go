package sync

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"

	// Importe os pacotes de proto gerados
	"github.com/GarotoCowboy/vttProject/api/grpc/pb/placedToken"
	syncBroker "github.com/GarotoCowboy/vttProject/api/grpc/pb/sync"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	grpcAddress = "localhost:50051"
	testTableID = 1 
	testSceneID = 3
	testTokenID = 4
)


func TestTokenSync(t *testing.T) {

	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Falha ao conectar ao servidor gRPC: %v", err)
	}
	defer conn.Close()

	
	syncClient := syncBroker.NewSyncServiceClient(conn)
	placedTokenClient := placedToken.NewPlacedTokenServiceClient(conn)

	
	receivedEvents := make(chan *syncBroker.SyncResponse, 10)                
	var wg sync.WaitGroup                                                    
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) 
	defer cancel()

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[Ouvinte] Conectando ao SyncService...")

		
		stream, err := syncClient.SyncScene(ctx)
		if err != nil {
			log.Printf("[Ouvinte] Erro ao conectar ao Sync: %v", err)
			return
		}

		initialReq := &syncBroker.SyncRequest{
			TableId: testTableID,
			SceneId: testSceneID,
		}
		if err := stream.Send(initialReq); err != nil {
			log.Printf("[Ouvinte] Erro ao enviar requisição inicial: %v", err)
			return
		}
		log.Println("[Ouvinte] Inscrito com sucesso! Aguardando eventos...")

		for {
			event, err := stream.Recv()
			if err != nil {
				if ctx.Err() != nil {
					log.Println("[Ouvinte] Desconectando devido ao fim do teste.")
					return
				}
				log.Printf("[Ouvinte] Erro ao receber evento: %v", err)
				return
			}
			log.Printf("[Ouvinte] Evento recebido: %T", event.GetAction())
			receivedEvents <- event
		}
	}()


	time.Sleep(1 * time.Second)

	t.Run("Sincronização do CreatePlacedToken", func(t *testing.T) {
		log.Println("[Ator] Chamando CreatePlacedToken...")
		createReq := &placedToken.CreatePlacedTokenRequest{
			SceneId: testSceneID,
			TokenId: testTokenID,
			PosX:    100,
			PosY:    150,
		}
		_, err := placedTokenClient.CreatePlacedToken(ctx, createReq)
		assert.NoError(t, err, "A chamada CreatePlacedToken não deve retornar erro")

		select {
		case event := <-receivedEvents:
			log.Println("[Teste] Verificando evento PlacedTokenCreated...")
			createdEvent, ok := event.GetAction().(*syncBroker.SyncResponse_PlacedTokenCreated)
			assert.True(t, ok, "O tipo do evento deve ser PlacedTokenCreated")
			assert.Equal(t, uint64(testSceneID), createdEvent.PlacedTokenCreated.GetPlacedToken().GetSceneId())
			assert.Equal(t, uint64(testTokenID), createdEvent.PlacedTokenCreated.GetPlacedToken().GetTokenId())
			assert.Equal(t, int32(100), createdEvent.PlacedTokenCreated.GetPlacedToken().GetPosX())
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout: Nenhum evento PlacedTokenCreated recebido do Ouvinte")
		}
	})

	t.Run("Sincronização do MoveToken", func(t *testing.T) {
		log.Println("[Ator] Chamando MoveToken...")
		moveReq := &placedToken.MoveTokenRequest{
			SceneId:       testSceneID,
			PlacedTokenId: 1, 
			PosX:          250,
			PosY:          300,
		}
		_, err := placedTokenClient.MoveToken(ctx, moveReq)
		assert.NoError(t, err, "A chamada MoveToken não deve retornar erro")

		select {
		case event := <-receivedEvents:
			log.Println("[Teste] Verificando evento PlacedTokenMoved...")
			movedEvent, ok := event.GetAction().(*syncBroker.SyncResponse_PlacedTokenMoved)
			assert.True(t, ok, "O tipo do evento deve ser PlacedTokenMoved")
			assert.Equal(t, uint64(testSceneID), movedEvent.PlacedTokenMoved.GetSceneId())
			assert.Equal(t, int32(250), movedEvent.PlacedTokenMoved.GetPosX())
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout: Nenhum evento PlacedTokenMoved recebido do Ouvinte")
		}
	})

	cancel()  
	wg.Wait() 
	close(receivedEvents)
}
