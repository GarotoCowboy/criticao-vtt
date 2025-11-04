<div align = "center"> <h1>CriticãoVTT: Uma Plataforma Gratuita para RPG de Mesa</h1></div>
<div align="center"><img src="https://github.com/user-attachments/assets/a478f526-e66e-41de-a6e9-1379f93c5f88" width="250px">
  <p><i>A nossa mascote Lili mordendo um d20</i></p>
</div>
<div align="center">
  <h3>Plataforma RESTful com gRPC para RPG de Mesa</h3>
  <p><i>Projeto desenvolvido para aperfeiçoar conhecimentos em Go (Golang), gRPC, Flutter e arquitetura de software.</i></p>
</div>

---

## 🧠 Objetivo
Este projeto visa a criação de uma plataforma robusta para jogadores de RPG de mesa. A ideia surgiu para preencher uma lacuna no mercado brasileiro, onde as plataformas existentes costumam ser muito caras (muitas vezes cobrando em dólar), oferecendo uma solução acessível e de alta performance para a comunidade.

Servindo como um estudo prático e aprofundado nas seguintes áreas:
- **Go (Golang)**: Desenvolvimento de APIs RESTful e serviços gRPC concorrentes.
- **PostgreSQL**: Modelagem de dados e interações com banco de dados relacional.
- **Flutter**: Desenvolvimento da interface do usuário (UI) multi-plataforma.
- **Arquitetura de Software**: Aplicação de arquitetura em camadas (Services, Handlers, DTOs) e sistemas orientados a eventos (Pub/Sub).
- **Ferramentas e ORMs**: Utilização de GORM e documentação com Swagger.

---

## 🚀 Funcionalidades (Backend v1.0)
A primeira versão do backend está quase concluida, implementando a lógica de negócio principal da plataforma.

### Arquitetura Híbrida: REST e gRPC
A aplicação utiliza uma abordagem híbrida para máxima eficiência:
- **REST API**: Usada para operações de gerenciamento de estado, como CRUD de usuários e mesas de RPG.
- **gRPC**: Usado para comunicação de alta performance e baixa latência, ideal para:
  - Gerenciamento de sessões de jogo.
  - Chat em tempo real (bidirecional).
  - Criação e atualização de fichas de personagem.
  - Manipulação de tokens e imagens em cena.

### Funcionalidades Implementadas
- **Autenticação Segura**: Sistema de autenticação JWT utilizando Bearer Tokens para garantir a segurança nas interações e acessos de usuários.
- **Gerenciamento de Mesas**: CRUD completo para criação de mesas de RPG, com geração de links de convite únicos, listagem e associação de participantes.
- **Gerenciamento de Cenas**: CRUD completo para criação de cenas em uma mesa de RPG, sendo possível inserir tokens e imagens para que sirva de tabuleiro para os jogadores.
- **Gerenciamento de Usuários**: CRUD completo para contas de usuário.
- **Motor de Fichas de Personagem**: Sistema que possibilita a criação de fichas de personagens para diferentes sistemas de RPG (Sistema Tormenta 20 implementado; D&D e GURPS planejados).
- **Chat em Tempo Real**: Implementação de um chat bidirecional (via gRPC) utilizando um broker Pub/Sub para interação entre os jogadores de forma orientada a eventos.
- **Token de personagens e suas barras**: Implementação de tokens e barras utilizando um broker Pub/Sub para interação entre os elementos e jogadores de forma orientada a eventos.
-  **Tokens em cenas e imagens em cenas**: Implementação de tokens e imagens inseridos em uma cena sendo possível movimentar e alterar a camada desses objetos, utiliza um broker Pub/Sub para interação entre os objetos e jogadores de forma orientada a eventos.
- **Atualização em Tempo Real das Fichas**: Fichas de personagem são atualizadas em tempo real, propagando as mudanças instantaneamente para todos os clientes conectados na sessão.

---

## 🛠 Tecnologias
### Backend (Concluído)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
[![gRPC](https://img.shields.io/badge/gRPC-4283F3?style=for-the-badge&logo=grpc&logoColor=white)](https://grpc.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Gin Gonic](https://img.shields.io/badge/Gin%20Gonic-009485?style=for-the-badge&logo=gin&logoColor=white)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-C42B9F?style=for-the-badge&logo=gorm&logoColor=white)](https://gorm.io/)
[![Swagger](https://img.shields.io/badge/Swagger-85EA2D?style=for-the-badge&logo=swagger&logoColor=black)](https://swagger.io/)

### Frontend (Planejado)
[![Flutter](https://img.shields.io/badge/Flutter-02569B?style=for-the-badge&logo=flutter&logoColor=white)](https://flutter.dev/)

---

## 📋 Status do Projeto
- **v1.0 - Backend:** O desenvolvimento da API REST e dos serviços gRPC em Go (Golang) está finalizado. A arquitetura em camadas (Services, Handlers, DTOs), a integração com banco de dados (PostgreSQL + GORM) e os sistemas de tempo real (gRPC + Pub/Sub) estão implementados e funcionais.
- **v2.0 - Frontend (Próximos Passos):** O foco agora será no desenvolvimento das telas e da interface do usuário (UI) utilizando Flutter, para criar uma interface amigável, fluida e multi-plataforma que consumirá os serviços do backend.

---

## 📊 Diagramas
### Diagrama de Casos de Uso (Inicial)
![projeto vtt-Caso de Uso drawio](https://github.com/user-attachments/assets/4ecb1797-9342-4c5a-aa71-516118f249bd)
*O projeto está em desenvolvimento e poderá haver alterações dos diagramas conforme a implementação do frontend avança.*

---

## 🧑‍💻 Autor
Pedro Henrique Marques Rocha - Aluno de Sistemas de Informação do Instituto Federal Goiano Campus Urutaí.

---
*Este projeto está em fase de desenvolvimento.*
