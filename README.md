# Educ Retro - Sistema de Retrospectivas

Sistema completo para gerenciar retrospectivas de equipes ágeis, desenvolvido em Go (backend) e React (frontend).

## 🚀 Features

### Core Features
- ✅ **Gestão de Usuários** - Registro, login e perfis
- ✅ **Gestão de Times** - Criar times, adicionar membros, gerenciar permissões
- ✅ **Retrospectivas** - Criar e gerenciar sessões de retrospectiva
- ✅ **Templates** - Diferentes formatos (Start/Stop/Continue, 4Ls, Mad/Sad/Glad, Sailboat)
- ✅ **Participação Colaborativa** - Múltiplos participantes contribuindo
- 🔄 **Tempo Real** - Atualizações em tempo real via WebSocket
- 🔄 **Sistema de Votação** - Votar em itens importantes
- 🔄 **Action Items** - Tracking de ações e follow-up
- 🔄 **Relatórios** - Analytics e métricas do time

### Features de UX
- 🔄 **Interface Responsiva** - Funciona em mobile e desktop
- 🔄 **Temas** - Dark/Light mode
- 🔄 **Timer** - Cronômetro para sessões
- 🔄 **Notificações** - Lembretes e updates
- 🔄 **Export** - Exportar retrospectivas (PDF, CSV)

## 🏗️ Arquitetura

```
educ-retro/
├── cmd/server/           # Ponto de entrada da aplicação
├── internal/
│   ├── models/          # Modelos de dados
│   ├── repositories/    # Camada de acesso a dados
│   ├── services/        # Lógica de negócio
│   ├── handlers/        # Controllers da API
│   ├── auth/           # Autenticação e JWT
│   ├── database/       # Conexão com banco
│   └── utils/          # Utilitários
├── migrations/         # Migrations do banco
├── frontend/          # Aplicação React (em desenvolvimento)
└── docs/             # Documentação
```

## 🛠️ Setup e Instalação

### Pré-requisitos
- Go 1.21+
- PostgreSQL 12+
- Node.js 18+ (para frontend)

### 1. Clone o repositório
```bash
git clone <repository-url>
cd educ-retro
```

### 2. Configure o banco de dados
```bash
# Crie um banco PostgreSQL
createdb educ_retro

# Configure as variáveis de ambiente
cp env.example .env
# Edite o arquivo .env com suas configurações
```

### 3. Execute as migrations
```bash
# Instale migrate CLI (se não tiver)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Execute as migrations
make migrate-up
```

### 4. Instale as dependências e execute
```bash
# Instalar dependências do backend
go mod tidy

# Executar o servidor
make run
```

O servidor estará rodando em `http://localhost:8080`

### 5. Configurar e executar o frontend
```bash
# Navegar para o diretório do frontend
cd frontend

# Instalar dependências do Node.js
npm install

# Configurar variáveis de ambiente
cp env.example .env

# Executar o frontend
npm start
```

O frontend estará rodando em `http://localhost:3000`

## 📚 API Endpoints

### Autenticação
- `POST /api/v1/auth/register` - Registrar usuário
- `POST /api/v1/auth/login` - Login

### Usuários
- `GET /api/v1/users/profile` - Perfil do usuário
- `PUT /api/v1/users/profile` - Atualizar perfil

### Times
- `GET /api/v1/teams` - Listar times do usuário
- `POST /api/v1/teams` - Criar time
- `GET /api/v1/teams/:id` - Detalhes do time
- `PUT /api/v1/teams/:id` - Atualizar time
- `DELETE /api/v1/teams/:id` - Deletar time
- `POST /api/v1/teams/:id/members` - Adicionar membro
- `DELETE /api/v1/teams/:id/members/:userId` - Remover membro

### Retrospectivas (Em desenvolvimento)
- `GET /api/v1/retrospectives` - Listar retrospectivas
- `POST /api/v1/retrospectives` - Criar retrospectiva
- `GET /api/v1/retrospectives/:id` - Detalhes da retrospectiva
- `POST /api/v1/retrospectives/:id/items` - Adicionar item
- `POST /api/v1/retrospectives/:id/vote` - Votar em item

## 🧪 Testando a API

### Registrar um usuário
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "name": "João Silva",
    "password": "123456"
  }'
```

### Fazer login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "usuario@exemplo.com",
    "password": "123456"
  }'
```

### Criar um time (use o token retornado no login)
```bash
curl -X POST http://localhost:8080/api/v1/teams \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SEU_TOKEN_AQUI" \
  -d '{
    "name": "Time de Desenvolvimento",
    "description": "Time responsável pelo desenvolvimento do produto"
  }'
```

## 🔧 Comandos Úteis

```bash
# Executar o servidor
make run

# Executar migrations
make migrate-up
make migrate-down

# Limpar e reinstalar dependências
make clean
make deps

# Executar testes
make test

# Build para produção
make build
```

## 🚧 Próximos Passos

1. ✅ **Implementar repositórios e serviços para retrospectivas**
2. ✅ **Criar sistema de WebSocket para tempo real**
3. ✅ **Desenvolver frontend React**
4. ✅ **Implementar templates de retrospectiva**
5. ✅ **Criar relatórios e analytics**
6. 🔄 **Adicionar sistema de votação**
7. 🔄 **Implementar action items**
8. 🔄 **Adicionar testes unitários e de integração**
9. 🔄 **Implementar funcionalidades de retrospectiva no backend**
10. 🔄 **Conectar frontend com WebSocket**

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo LICENSE para mais detalhes.
