#!/bin/bash

# Script de setup completo do projeto Educ Retro
# Este script configura o ambiente de desenvolvimento completo

set -e

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para imprimir mensagens coloridas
print_message() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_header() {
    echo -e "${BLUE}================================${NC}"
    echo -e "${BLUE}  Setup do Educ Retro${NC}"
    echo -e "${BLUE}================================${NC}"
}

# Verificar se estamos no diretório correto
if [ ! -f "go.mod" ] || [ ! -f "package.json" ]; then
    print_error "Execute este script no diretório raiz do projeto"
    exit 1
fi

print_header

# 1. Verificar dependências do sistema
print_message "Verificando dependências do sistema..."

# Verificar Go
if ! command -v go &> /dev/null; then
    print_error "Go não está instalado. Por favor, instale Go: https://golang.org/dl/"
    exit 1
fi
print_success "Go está instalado: $(go version)"

# Verificar Node.js
if ! command -v node &> /dev/null; then
    print_error "Node.js não está instalado. Por favor, instale Node.js: https://nodejs.org/"
    exit 1
fi
print_success "Node.js está instalado: $(node --version)"

# Verificar npm
if ! command -v npm &> /dev/null; then
    print_error "npm não está instalado. Por favor, instale npm"
    exit 1
fi
print_success "npm está instalado: $(npm --version)"

# Verificar PostgreSQL
if ! command -v psql &> /dev/null; then
    print_error "PostgreSQL não está instalado. Por favor, instale PostgreSQL:"
    print_message "  macOS: brew install postgresql"
    print_message "  Ubuntu: sudo apt-get install postgresql postgresql-contrib"
    print_message "  CentOS: sudo yum install postgresql postgresql-server"
    exit 1
fi
print_success "PostgreSQL está instalado: $(psql --version)"

# Verificar se PostgreSQL está rodando
if ! pg_isready -h localhost -p 5432 -U postgres > /dev/null 2>&1; then
    print_warning "PostgreSQL não está rodando. Iniciando PostgreSQL..."
    if command -v brew &> /dev/null; then
        brew services start postgresql
    elif command -v systemctl &> /dev/null; then
        sudo systemctl start postgresql
    else
        print_error "Não foi possível iniciar o PostgreSQL automaticamente"
        print_message "Por favor, inicie o PostgreSQL manualmente e execute o script novamente"
        exit 1
    fi
    sleep 3
fi
print_success "PostgreSQL está rodando"

# 2. Instalar dependências do Go
print_message "Instalando dependências do Go..."
go mod tidy
go mod download
print_success "Dependências do Go instaladas"

# 3. Instalar dependências do frontend
print_message "Instalando dependências do frontend..."
cd frontend
npm install
cd ..
print_success "Dependências do frontend instaladas"

# 4. Verificar se migrate está instalado
if ! command -v migrate &> /dev/null; then
    print_message "Instalando ferramenta migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
    print_success "Ferramenta migrate instalada"
else
    print_success "Ferramenta migrate já está instalada"
fi

# 5. Configurar arquivo .env
print_message "Configurando arquivo .env..."
if [ ! -f .env ]; then
    cp env.example .env
    print_success "Arquivo .env criado a partir do exemplo"
    print_warning "Por favor, edite o arquivo .env com suas configurações se necessário"
else
    print_success "Arquivo .env já existe"
fi

# 6. Configurar banco de dados
print_message "Configurando banco de dados..."
./scripts/init-db.sh
print_success "Banco de dados configurado"

# 7. Verificar se tudo está funcionando
print_message "Verificando configuração..."
if ./scripts/init-db.sh --check-only; then
    print_success "Banco de dados está configurado corretamente"
else
    print_error "Erro na configuração do banco de dados"
    exit 1
fi

print_header
print_success "Setup completo do Educ Retro concluído!"
print_message "Próximos passos:"
print_message "  1. Edite o arquivo .env se necessário"
print_message "  2. Execute 'make run' para iniciar o servidor backend"
print_message "  3. Execute 'cd frontend && npm start' para iniciar o frontend"
print_message "  4. Ou execute 'make run-all' para iniciar ambos simultaneamente"
print_message ""
print_message "URLs:"
print_message "  Backend: http://localhost:8080"
print_message "  Frontend: http://localhost:3000"
print_message "  Health Check: http://localhost:8080/health"
