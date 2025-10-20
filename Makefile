.PHONY: run build clean deps test migrate-up migrate-down run-all setup-db setup check-db docker-db docker-db-start docker-db-stop docker-db-restart help

# Variáveis
BINARY_NAME=educ-retro
DB_URL=postgres://postgres:password@localhost:5432/educ_retro?sslmode=disable

# Executar o servidor
run:
	@echo "🚀 Iniciando servidor..."
	go run cmd/server/main.go

# Executar backend e frontend simultaneamente
run-all:
	@echo "🚀 Iniciando backend e frontend..."
	@echo "📡 Backend rodando em: http://localhost:8080"
	@echo "🌐 Frontend rodando em: http://localhost:3000"
	@echo "⏹️  Pressione Ctrl+C para parar ambos"
	@trap 'kill %1 %2' INT; \
	(cd frontend && npm start) & \
	(sleep 5 && go run cmd/server/main.go) & \
	wait

# Build para produção
build:
	@echo "🔨 Building binary..."
	go build -o bin/$(BINARY_NAME) cmd/server/main.go

# Limpar arquivos gerados
clean:
	@echo "🧹 Limpando arquivos..."
	go clean
	rm -rf bin/

# Instalar dependências
deps:
	@echo "📦 Instalando dependências..."
	go mod tidy
	go mod download

# Executar testes
test:
	@echo "🧪 Executando testes..."
	go test ./...

# Executar testes com coverage
test-coverage:
	@echo "🧪 Executando testes com coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Executar migrations
migrate-up:
	@echo "📈 Executando migrations..."
	migrate -path migrations -database "$(DB_URL)" up

# Reverter migrations
migrate-down:
	@echo "📉 Revertendo migrations..."
	migrate -path migrations -database "$(DB_URL)" down

# Resetar banco (cuidado!)
migrate-reset:
	@echo "⚠️  Resetando banco de dados..."
	migrate -path migrations -database "$(DB_URL)" down
	migrate -path migrations -database "$(DB_URL)" up

# Criar nova migration
migrate-create:
	@echo "📝 Criando nova migration..."
	@read -p "Nome da migration: " name; \
	migrate create -ext sql -dir migrations $$name

# Executar linter
lint:
	@echo "🔍 Executando linter..."
	golangci-lint run

# Formatar código
fmt:
	@echo "💅 Formatando código..."
	go fmt ./...
	goimports -w .

# Gerar documentação da API
docs:
	@echo "📚 Gerando documentação..."
	swag init -g cmd/server/main.go

# Setup inicial do banco de dados
setup-db:
	@echo "🗄️  Configurando banco de dados..."
	@if [ ! -f .env ]; then \
		echo "📋 Criando arquivo .env a partir do exemplo..."; \
		cp env.example .env; \
		echo "⚠️  Por favor, edite o arquivo .env com suas configurações"; \
	fi
	@echo "🚀 Executando script de inicialização do banco..."
	@./scripts/init-db.sh

# Setup completo do projeto
setup: deps setup-db
	@echo "✅ Setup completo do projeto concluído!"
	@echo "🚀 Execute 'make run' para iniciar o servidor"

# Verificar status do banco
check-db:
	@echo "🔍 Verificando status do banco de dados..."
	@./scripts/init-db.sh --check-only || echo "❌ Banco não está configurado. Execute 'make setup-db'"

# Docker PostgreSQL - Iniciar container
docker-db-start:
	@echo "🐳 Iniciando PostgreSQL no Docker..."
	@if docker ps -a --format "table {{.Names}}" | grep -q "educ-retro-postgres"; then \
		echo "📦 Container já existe, iniciando..."; \
		docker start educ-retro-postgres; \
	else \
		echo "📦 Criando e iniciando novo container..."; \
		docker run --name educ-retro-postgres \
			-e POSTGRES_PASSWORD=password \
			-e POSTGRES_DB=educ_retro \
			-p 5432:5432 \
			-d postgres:14; \
	fi
	@echo "⏳ Aguardando PostgreSQL inicializar..."
	@sleep 5
	@echo "✅ PostgreSQL no Docker está rodando!"

# Docker PostgreSQL - Parar container
docker-db-stop:
	@echo "🛑 Parando PostgreSQL no Docker..."
	@docker stop educ-retro-postgres 2>/dev/null || echo "Container não estava rodando"
	@echo "✅ PostgreSQL parado!"

# Docker PostgreSQL - Reiniciar container
docker-db-restart: docker-db-stop docker-db-start
	@echo "🔄 PostgreSQL reiniciado!"

# Docker PostgreSQL - Remover container
docker-db-remove:
	@echo "🗑️  Removendo container PostgreSQL..."
	@docker stop educ-retro-postgres 2>/dev/null || true
	@docker rm educ-retro-postgres 2>/dev/null || true
	@echo "✅ Container removido!"

# Setup completo com Docker
docker-setup: docker-db-start setup-db
	@echo "✅ Setup completo com Docker concluído!"
	@echo "🚀 Execute 'make run' para iniciar o servidor"

# Help
help:
	@echo "📋 Comandos disponíveis:"
	@echo ""
	@echo "🚀 Execução:"
	@echo "  run           - Executar o servidor"
	@echo "  run-all       - Executar backend e frontend simultaneamente"
	@echo ""
	@echo "🐳 Docker PostgreSQL:"
	@echo "  docker-db-start    - Iniciar PostgreSQL no Docker"
	@echo "  docker-db-stop     - Parar PostgreSQL no Docker"
	@echo "  docker-db-restart  - Reiniciar PostgreSQL no Docker"
	@echo "  docker-db-remove   - Remover container PostgreSQL"
	@echo "  docker-setup       - Setup completo com Docker"
	@echo ""
	@echo "🗄️  Banco de Dados:"
	@echo "  setup-db      - Configurar banco de dados"
	@echo "  check-db      - Verificar status do banco"
	@echo "  migrate-up    - Executar migrations"
	@echo "  migrate-down  - Reverter migrations"
	@echo "  migrate-reset - Resetar banco (cuidado!)"
	@echo "  migrate-create- Criar nova migration"
	@echo ""
	@echo "🔧 Desenvolvimento:"
	@echo "  setup         - Setup completo do projeto"
	@echo "  deps          - Instalar dependências"
	@echo "  test          - Executar testes"
	@echo "  test-coverage - Executar testes com coverage"
	@echo "  lint          - Executar linter"
	@echo "  fmt           - Formatar código"
	@echo ""
	@echo "📦 Produção:"
	@echo "  build         - Build para produção"
	@echo "  clean         - Limpar arquivos gerados"
	@echo "  docs          - Gerar documentação da API"
	@echo ""
	@echo "❓ Ajuda:"
	@echo "  help          - Mostrar esta ajuda"