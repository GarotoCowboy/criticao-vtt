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
	testTableID = 1 // Use IDs que existem no seu banco de dados de teste
	testSceneID = 3
	testTokenID = 4
)

// TestTokenSync é o nosso teste de integração completo.
func TestTokenSync(t *testing.T) {
	// --- 1. CONFIGURAÇÃO DA CONEXÃO COM O SERVIDOR ---
	// Criamos uma única conexão gRPC que será usada por ambos os "jogadores".
	conn, err := grpc.Dial(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Falha ao conectar ao servidor gRPC: %v", err)
	}
	defer conn.Close()

	// Criamos os clientes para os serviços que vamos usar.
	syncClient := syncBroker.NewSyncServiceClient(conn)
	placedTokenClient := placedToken.NewPlacedTokenServiceClient(conn)

	// --- 2. CONFIGURAÇÃO DA COMUNICAÇÃO ENTRE GOROUTINES ---
	// Usamos um canal para que o "Ouvinte" possa enviar os eventos que ele recebe
	// de volta para a goroutine principal do teste para verificação.
	receivedEvents := make(chan *syncBroker.SyncResponse, 10)                // Canal com buffer
	var wg sync.WaitGroup                                                    // WaitGroup para garantir que a goroutine do ouvinte termine
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Timeout para o teste
	defer cancel()

	// --- 3. INICIANDO O JOGADOR A (O OUVINTE) ---
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[Ouvinte] Conectando ao SyncService...")

		// Estabelece o stream bidirecional com o SyncService
		stream, err := syncClient.SyncScene(ctx)
		if err != nil {
			log.Printf("[Ouvinte] Erro ao conectar ao Sync: %v", err)
			return
		}

		// Envia a requisição inicial para se inscrever nos tópicos corretos
		initialReq := &syncBroker.SyncRequest{
			TableId: testTableID,
			SceneId: testSceneID,
		}
		if err := stream.Send(initialReq); err != nil {
			log.Printf("[Ouvinte] Erro ao enviar requisição inicial: %v", err)
			return
		}
		log.Println("[Ouvinte] Inscrito com sucesso! Aguardando eventos...")

		// Loop para receber eventos
		for {
			event, err := stream.Recv()
			if err != nil {
				// Se o contexto for cancelado, é o fim normal do teste.
				if ctx.Err() != nil {
					log.Println("[Ouvinte] Desconectando devido ao fim do teste.")
					return
				}
				log.Printf("[Ouvinte] Erro ao receber evento: %v", err)
				return
			}
			log.Printf("[Ouvinte] Evento recebido: %T", event.GetAction())
			receivedEvents <- event // Envia o evento para o canal para ser verificado
		}
	}()

	// --- 4. EXECUTANDO AS AÇÕES DO JOGADOR B (O ATOR) ---

	// Damos um pequeno tempo para o Ouvinte se conectar e se inscrever.
	time.Sleep(1 * time.Second)

	// AÇÃO 1: Criar um Placed Token
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

		// Verificação: Esperamos receber o evento de criação do Ouvinte
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

	// AÇÃO 2: Mover o Placed Token
	t.Run("Sincronização do MoveToken", func(t *testing.T) {
		log.Println("[Ator] Chamando MoveToken...")
		moveReq := &placedToken.MoveTokenRequest{
			SceneId:       testSceneID,
			PlacedTokenId: 1, // Assumindo que o ID do PlacedToken é 1. Ajuste se necessário.
			PosX:          250,
			PosY:          300,
		}
		_, err := placedTokenClient.MoveToken(ctx, moveReq)
		assert.NoError(t, err, "A chamada MoveToken não deve retornar erro")

		// Verificação: Esperamos receber o evento de movimento do Ouvinte
		select {
		case event := <-receivedEvents:
			log.Println("[Teste] Verificando evento PlacedTokenMoved...")
			movedEvent, ok := event.GetAction().(*syncBroker.SyncResponse_PlacedTokenMoved)
			assert.True(t, ok, "O tipo do evento deve ser PlacedTokenMoved")
			assert.Equal(t, uint64(testSceneID), movedEvent.PlacedTokenMoved.GetSceneId())
			// assert.Equal(t, moveReq.GetPlacedTokenId(), movedEvent.PlacedTokenMoved.GetTokenId()) // O ID do token movido
			assert.Equal(t, int32(250), movedEvent.PlacedTokenMoved.GetPosX())
		case <-time.After(5 * time.Second):
			t.Fatal("Timeout: Nenhum evento PlacedTokenMoved recebido do Ouvinte")
		}
	})

	// --- 5. FINALIZAÇÃO ---
	cancel()  // Sinaliza para a goroutine do ouvinte parar
	wg.Wait() // Espera a goroutine do ouvinte terminar de forma limpa
	close(receivedEvents)
}
