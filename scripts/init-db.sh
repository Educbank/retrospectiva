#!/bin/bash

# Script para inicializar o banco de dados PostgreSQL
# Este script cria o banco de dados se não existir e executa as migrations

set -e

# Verificar se é modo de verificação apenas
CHECK_ONLY=false
if [ "$1" = "--check-only" ]; then
    CHECK_ONLY=true
fi

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

# Carregar variáveis de ambiente
if [ -f .env ]; then
    print_message "Carregando variáveis do arquivo .env"
    export $(cat .env | grep -v '^#' | xargs)
else
    print_warning "Arquivo .env não encontrado, usando variáveis padrão"
fi

# Configurações padrão
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}
DB_NAME=${DB_NAME:-educ_retro}
DB_SSLMODE=${DB_SSLMODE:-disable}

print_message "Configurações do banco:"
print_message "  Host: $DB_HOST"
print_message "  Port: $DB_PORT"
print_message "  User: $DB_USER"
print_message "  Database: $DB_NAME"

# Verificar se o PostgreSQL está rodando
print_message "Verificando se o PostgreSQL está rodando..."
export PGPASSWORD=$DB_PASSWORD

# Verificar se psql está disponível localmente
if command -v psql &> /dev/null; then
    if ! psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
        print_error "PostgreSQL não está rodando ou não está acessível"
        print_message "Por favor, inicie o PostgreSQL e tente novamente"
        exit 1
    fi
else
    # Usar Docker se psql não estiver disponível
    if ! docker exec educ-retro-postgres psql -U $DB_USER -d postgres -c "SELECT 1;" > /dev/null 2>&1; then
        print_error "PostgreSQL no Docker não está rodando ou não está acessível"
        print_message "Por favor, inicie o PostgreSQL e tente novamente"
        exit 1
    fi
fi
print_success "PostgreSQL está rodando"

# Se for modo de verificação apenas, verificar se o banco existe e tem tabelas
if [ "$CHECK_ONLY" = true ]; then
    export PGPASSWORD=$DB_PASSWORD
    
    # Verificar se o banco existe
    if command -v psql &> /dev/null; then
        DB_EXISTS=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" 2>/dev/null || echo "0")
        TABLE_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'" 2>/dev/null || echo "0")
    else
        DB_EXISTS=$(docker exec educ-retro-postgres psql -U $DB_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'" 2>/dev/null || echo "0")
        TABLE_COUNT=$(docker exec educ-retro-postgres psql -U $DB_USER -d $DB_NAME -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'" 2>/dev/null || echo "0")
    fi
    
    if [ "$DB_EXISTS" != "1" ]; then
        print_error "Banco '$DB_NAME' não existe"
        exit 1
    fi
    
    if [ "$TABLE_COUNT" -eq 0 ]; then
        print_error "Banco existe mas não tem tabelas. Execute 'make setup-db' para configurar"
        exit 1
    fi
    
    print_success "Banco de dados está configurado corretamente ($TABLE_COUNT tabelas)"
    exit 0
fi

# Conectar ao PostgreSQL e criar o banco se não existir
print_message "Verificando se o banco '$DB_NAME' existe..."

# Usar PGPASSWORD para evitar prompt de senha
export PGPASSWORD=$DB_PASSWORD

# Verificar se o banco existe
if command -v psql &> /dev/null; then
    DB_EXISTS=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")
else
    DB_EXISTS=$(docker exec educ-retro-postgres psql -U $DB_USER -d postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$DB_NAME'")
fi

if [ "$DB_EXISTS" != "1" ]; then
    print_message "Banco '$DB_NAME' não existe, criando..."
    if command -v createdb &> /dev/null; then
        createdb -h $DB_HOST -p $DB_PORT -U $DB_USER $DB_NAME
    else
        docker exec educ-retro-postgres createdb -U $DB_USER $DB_NAME
    fi
    print_success "Banco '$DB_NAME' criado com sucesso"
else
    print_success "Banco '$DB_NAME' já existe"
fi

# Verificar se a ferramenta migrate está instalada
if ! command -v migrate &> /dev/null; then
    print_error "Ferramenta 'migrate' não encontrada"
    print_message "Por favor, instale a ferramenta migrate:"
    print_message "  macOS: brew install golang-migrate"
    print_message "  Linux: https://github.com/golang-migrate/migrate/releases"
    print_message "  Ou baixe de: https://github.com/golang-migrate/migrate/releases/latest"
    exit 1
fi

# Executar migrations
print_message "Executando migrations..."

DB_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE"

# Verificar se há migrations para executar
MIGRATION_COUNT=$(migrate -path migrations -database "$DB_URL" version 2>/dev/null || echo "0")

if [ "$MIGRATION_COUNT" = "0" ]; then
    print_message "Executando todas as migrations..."
    migrate -path migrations -database "$DB_URL" up
    print_success "Todas as migrations executadas com sucesso"
else
    print_message "Verificando migrations pendentes..."
    migrate -path migrations -database "$DB_URL" up
    print_success "Migrations atualizadas com sucesso"
fi

# Verificar se as tabelas foram criadas
print_message "Verificando se as tabelas foram criadas..."
if command -v psql &> /dev/null; then
    TABLE_COUNT=$(psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
else
    TABLE_COUNT=$(docker exec educ-retro-postgres psql -U $DB_USER -d $DB_NAME -tAc "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public'")
fi

if [ "$TABLE_COUNT" -gt 0 ]; then
    print_success "Banco de dados inicializado com sucesso!"
    print_message "Tabelas criadas: $TABLE_COUNT"
    
    # Listar tabelas criadas
    print_message "Tabelas no banco:"
    if command -v psql &> /dev/null; then
        psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -c "\dt" | grep -E "^\s+[a-z_]+" | awk '{print "  - " $3}'
    else
        docker exec educ-retro-postgres psql -U $DB_USER -d $DB_NAME -c "\dt" | grep -E "^\s+[a-z_]+" | awk '{print "  - " $3}'
    fi
else
    print_error "Nenhuma tabela foi criada. Verifique os logs acima para erros."
    exit 1
fi

print_success "Setup do banco de dados concluído!"
print_message "Você pode agora executar: make run"
