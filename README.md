<div align = "center"> <h1>Criticão (Projeto de Estudo)</h1></div>
<div align="center"><img src="https://github.com/user-attachments/assets/a478f526-e66e-41de-a6e9-1379f93c5f88" width="250px">
  <p><i>A nossa mascote Lili mordendo um d20</i></p>
</div>

---

## 🧠 Objetivo
Este projeto visa a criação de uma plataforma robusta para jogadores de RPG de mesa, servindo como um estudo prático e aprofundado nas seguintes tecnologias:
- **Go (Golang)**: Desenvolvimento de APIs RESTful, concorrência, gRPC.
- **PostgreSQL**: Modelagem de dados e interações com banco de dados relacional.
- **React**: Desenvolvimento da interface do usuário (UI) da plataforma.
- **Arquitetura de Software**: Aplicação de conceitos como arquitetura em camadas (Services, Handlers, DTOs).
- **Ferramentas e ORMs**: Utilização de GORM para mapeamento objeto-relacional e Swagger para documentação de API.
- **Protocol Buffers (gRPC)**: Definição de contratos de serviço para comunicação em tempo real.

O projeto busca ser uma alternativa às plataformas existentes no mercado para RPG de mesa.

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

## 🛠 Tecnologias
### Backend
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Gin Gonic](https://img.shields.io/badge/Gin%20Gonic-009485?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-C42B9F?style=for-the-badge&logo=gorm&logoColor=white)](https://gorm.io/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)
[![gRPC](https://img.shields.io/badge/gRPC-4283F3?style=for-the-badge&logo=grpc&logoColor=white)](https://grpc.io/)

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

# AUTH / JWT

JWT_SECRET=minha_chave_super_super_secreta


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


### 3. 📚 Documentação da API
A documentação da API está no link: https://vttproject.postman.co/workspace/golangapi~d97bdf1e-aada-4788-86b2-8949b8d429bb/collection/24061336-6431ac82-57f0-4799-ae4f-61b9c5be2dac?action=share&creator=24061336

## 📋 Etapas do Projeto
- Desenvolvimento dos diagramas de caso de uso, diagrama de classe e diagrama entidade relacionamento.
- Desenvolvimento das classes (models).
- Desenvolvimento do banco de dados.
- Implementar funcionalidades tais como por exemplo GORM e SWAGGER.
- Desenvolvimento das regras de negócio (services).
- Desenvolvimento das funcionalidades que utilizarão gRPC.
- Desenvolver a UI da plataforma.
- Realizar testes unitários.
- Realizar testes de performance.
- Corrigir bugs encontrados após os testes.
- Lançar a plataforma.

### 📊 Diagramas
#### Diagrama de Casos de Uso (Inicial)
![projeto vtt-Caso de Uso drawio](https://github.com/user-attachments/assets/4ecb1797-9342-4c5a-aa71-516118f249bd)
*O projeto está ainda em desenvolvimento e poderá haver alterações dos diagramas*.



Os diagramas podem evoluir conforme o projeto avança.

📌 Status do Projeto

✅ v1.0 – Backend concluído

- REST + gRPC
- Arquitetura em camadas
- Tempo real funcional
- Banco de dados integrado

🚧 v2.0 – Frontend
- Desenvolvimento das telas em Flutter
- Interface multi-plataforma
- Consumo completo dos serviços backend

👨‍💻 Autor

Pedro Henrique Marques Rocha
Aluno de Sistemas de Informação
Instituto Federal Goiano – Campus Urutaí


