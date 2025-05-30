<div align = "center"> <h1>Criticão (Projeto de Estudo)</h1></div>
<div align="center"><img src="https://github.com/user-attachments/assets/a478f526-e66e-41de-a6e9-1379f93c5f88" width="250px">
  <p><i>A nossa mascote Lili mordendo um d20</i></p>
</div>
<div align="center">
  <h3>Plataforma RESTful com gRPC para RPG de Mesa</h3>
  <p><i>Projeto desenvolvido para aperfeiçoar conhecimentos em Go (Golang), tecnologias de backend e frontend.</i></p>
</div>

---

## 🧠 Objetivo
Este projeto visa a criação de uma plataforma robusta para jogadores de RPG de mesa, servindo como um estudo prático e aprofundado nas seguintes áreas:
- **Go (Golang)**: Desenvolvimento de APIs RESTful, concorrência, gRPC.
- **PostgreSQL**: Modelagem de dados e interações com banco de dados relacional.
- **React**: Desenvolvimento da interface do usuário (UI) da plataforma.
- **Arquitetura de Software**: Aplicação de conceitos como arquitetura em camadas (Services, Handlers, DTOs).
- **Ferramentas e ORMs**: Utilização de GORM para mapeamento objeto-relacional e Swagger para documentação de API.

O projeto busca ser uma alternativa às plataformas existentes no mercado para RPG de mesa.

---

## 🚀 Funcionalidades Principais
O sistema permitirá o gerenciamento de usuários, mesas de RPG e a relação entre eles, com as seguintes funcionalidades:

### Gerenciamento de Usuários
- **CRUD de Usuários**: Criação, visualização, listagem, atualização e exclusão de contas de usuário.
- **Upload de Imagem de Usuário**: Permitir que usuários adicionem imagens aos seus perfis.

### Gerenciamento de Mesas de RPG
- **CRUD de Mesas**: Criação (com geração de link de convite), visualização, listagem, atualização e exclusão de mesas de RPG.
- **Propriedade de Mesas**: Cada mesa possui um usuário proprietário (Mestre do Jogo).

### Gerenciamento de Participantes da Mesa (TableUser)
- **Associação Usuário-Mesa**: Adicionar e remover usuários de mesas, definindo seus papéis (ex: Jogador, Mestre).
- **Listagem de Participantes**: Visualizar os usuários associados a uma mesa específica.

*(Funcionalidades adicionais como chat em tempo real, rolagem de dados, fichas de personagem, e outras interações via gRPC estão planejadas para fases futuras do desenvolvimento)*.

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
[![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)](https://reactjs.org/)
[![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/HTML5)
[![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/CSS)

*Ferramentas complementares:*
- Testes Unitários (planejado/em desenvolvimento inicial).

---

## 📋 Etapas do Projeto (Conforme README Original)
- Desenvolvimento dos diagramas de caso de uso, diagrama de classe e diagrama entidade relacionamento.
- Desenvolvimento das classes (models).
- Desenvolvimento do banco de dados.
- Implementar funcionalidades tais como GORM e SWAGGER.
- Desenvolvimento das regras de negócio (services).
- Desenvolvimento das funcionalidades que utilizarão gRPC.
- Desenvolver a UI da plataforma.
- Realizar testes unitários.
- Realizar testes de performance.
- Corrigir bugs encontrados após os testes.
- Lançar a plataforma.

---

## 📊 Diagramas
### Diagrama de Casos de Uso (Inicial)
![projeto vtt-Caso de Uso drawio](https://github.com/user-attachments/assets/4ecb1797-9342-4c5a-aa71-516118f249bd)
*O projeto está ainda em desenvolvimento e poderá haver alterações dos diagramas*.

---

## 🧑‍💻 Autor
Pedro Henrique Marques Rocha - Aluno de Sistemas de Informação do Instituto Federal Goiano Campus Urutaí.

---
*Este projeto está em fase de desenvolvimento.*
