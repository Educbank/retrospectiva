.PHONY: run build clean deps test migrate-up migrate-down run-all

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

# Executar em modo de desenvolvimento com hot reload
dev:
	@echo "🔥 Iniciando modo de desenvolvimento..."
	@command -v air >/dev/null 2>&1 || { echo "Instalando air..."; go install github.com/cosmtrek/air@latest; }
	air

# Setup inicial do projeto
setup:
	@echo "⚙️  Configurando projeto..."
	@echo "1. Instalando dependências..."
	$(MAKE) deps
	@echo "2. Criando arquivo .env..."
	@if [ ! -f .env ]; then cp env.example .env; echo "Arquivo .env criado. Configure as variáveis."; fi
	@echo "3. Instalando migrate CLI..."
	@command -v migrate >/dev/null 2>&1 || { echo "Instalando migrate..."; go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; }
	@echo "✅ Setup concluído!"
	@echo "📝 Configure o arquivo .env e execute 'make migrate-up' para configurar o banco"

# Help
help:
	@echo "📋 Comandos disponíveis:"
	@echo "  run           - Executar o servidor"
	@echo "  run-all       - Executar backend e frontend simultaneamente"
	@echo "  build         - Build para produção"
	@echo "  clean         - Limpar arquivos gerados"
	@echo "  deps          - Instalar dependências"
	@echo "  test          - Executar testes"
	@echo "  test-coverage - Executar testes com coverage"
	@echo "  migrate-up    - Executar migrations"
	@echo "  migrate-down  - Reverter migrations"
	@echo "  migrate-reset - Resetar banco (cuidado!)"
	@echo "  migrate-create- Criar nova migration"
	@echo "  lint          - Executar linter"
	@echo "  fmt           - Formatar código"
	@echo "  docs          - Gerar documentação da API"
	@echo "  dev           - Modo de desenvolvimento com hot reload"
	@echo "  setup         - Setup inicial do projeto"
	@echo "  help          - Mostrar esta ajuda"
