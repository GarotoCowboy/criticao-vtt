<div align="center">
  <h1>🎲 CriticãoVTT</h1>
  <p><strong>Uma plataforma gratuita e brasileira para RPG de Mesa</strong></p>
</div>

<div align="center">
  <img src="https://github.com/user-attachments/assets/a478f526-e66e-41de-a6e9-1379f93c5f88" width="250px"/>
  <p><i>Lili, a mascote oficial do CriticãoVTT, mordendo um d20</i></p>
</div>

---

## 📌 Visão Geral

O **CriticãoVTT** é uma plataforma de **Virtual Tabletop (VTT)** desenvolvida para jogadores de RPG de mesa, com foco no **mercado brasileiro**, oferecendo uma alternativa **gratuita**, **local** e **sem custos em dólar**.

O projeto foi criado como um estudo prático e aprofundado em **Go (Golang)**, **gRPC**, **arquitetura de software**, **sistemas em tempo real** e **desenvolvimento backend moderno**, servindo também como base para evolução futura em frontend e mobile.

---

## 🎯 Objetivo do Projeto

- Criar uma plataforma robusta para RPG de mesa
- Evitar dependência de soluções caras e estrangeiras
- Explorar arquitetura em camadas e sistemas orientados a eventos
- Estudar comunicação em tempo real com gRPC
- Desenvolver uma base extensível para múltiplos sistemas de RPG

---

## 🚀 Funcionalidades

### 👤 Gerenciamento de Usuários
- CRUD completo de usuários
- Upload de imagem de perfil
- Autenticação com **JWT**

### 🎲 Gerenciamento de Mesas de RPG
- CRUD de mesas
- Geração de link de convite
- Definição de proprietário da mesa (Mestre)

### 👥 Participantes da Mesa (TableUser)
- Associação usuário ↔ mesa
- Definição de papéis (Jogador, Mestre)
- Listagem de participantes por mesa

---

### 💬 Chat em Tempo Real (gRPC)
- Envio de mensagens via **Pub/Sub**
- Listagem de mensagens com **Server Streaming**
- Mensagens privadas entre usuários da mesa

---

### 🗺️ Tabuleiro em Tempo Real
- Criação de cenas
- Movimentação de tokens em tempo real
- Envio de imagens para o tabuleiro pelo mestre
- Sincronização via eventos gRPC

---

### 🧙 Personagens
- Criação e gerenciamento de fichas
- Atualização em tempo real (streams bidirecionais)
- Sistema de regras implementado para **Tormenta 20**
- Estrutura genérica para suportar futuramente:
  - D&D
  - GURPS
  - Outros sistemas

> Funcionalidades futuras planejadas:
> - Chat por vídeo
> - Loja de plugins e sistemas

---

## 🧱 Arquitetura

- Arquitetura em camadas:
  - **Handlers**
  - **Services**
  - **DTOs**
  - **Models**
- Backend orientado a eventos
- Comunicação REST + gRPC
- Autenticação via JWT
- Pub/Sub para tempo real

---

## 🛠 Tecnologias Utilizadas

### Backend (Concluído)
- **Go (Golang)**
- **Gin Gonic**
- **gRPC**
- **PostgreSQL**
- **GORM**
- **JWT**
- **Swagger**

### Frontend (Planejado)
- **React** (alternativo)
- **HTML / CSS**

---

## 📋 Pré-requisitos

- Go **1.25.0** ou superior
- PostgreSQL **17.5**
- Protobuf Compiler (`protoc`)

---

## ⚙️ Configuração do Ambiente

### 1️⃣ Clone o repositório

```bash 
git clone https://github.com/GarotoCowboy/criticao-vtt
cd criticao-vtt
```

### 2️⃣ Configure o arquivo .env

#### Crie um arquivo .env na raiz do projeto:

```env
# DATABASE
DB_HOST=localhost
DB_USERNAME=postgres
DB_PASSWORD=senha_database
DB_URL=postgres://usuario:senha@host:porta/database

# REST
REST_HOST=localhost
PORT_REST=8080

# GRPC
GRPC_HOST=localhost
PORT_GRPC=50051

```

### ▶️ Executando a Aplicação
```
#Desenvolvimento
go run main.go

#Produção
go build
./criticao-vtt
```

### 🌐 Endpoints

```
#REST API:
http://{REST_HOST}:{PORT_REST}

#gRPC:
{GRPC_HOST}:{PORT_GRPC}
```


### 📚 Documentação da API

Postman Collection:
https://vttproject.postman.co/workspace/golangapi~d97bdf1e-aada-4788-86b2-8949b8d429bb/collection/24061336-6431ac82-57f0-4799-ae4f-61b9c5be2dac?action=share&creator=24061336

### 📊 Diagramas
#### Diagrama de Casos de Uso (Inicial)
![projeto vtt-Caso de Uso drawio](https://github.com/user-attachments/assets/4ecb1797-9342-4c5a-aa71-516118f249bd)
*O projeto está em desenvolvimento e poderá haver alterações dos diagramas conforme a implementação do sistema avança.*



Os diagramas podem evoluir conforme o projeto avança.

📌 Status do Projeto

✅ v1.0 – Backend concluído

REST + gRPC

Arquitetura em camadas

Tempo real funcional

Banco de dados integrado

🚧 v2.0 – Frontend

Desenvolvimento das telas em Flutter

Interface multi-plataforma

Consumo completo dos serviços backend

👨‍💻 Autor

Pedro Henrique Marques Rocha
Aluno de Sistemas de Informação
Instituto Federal Goiano – Campus Urutaí


