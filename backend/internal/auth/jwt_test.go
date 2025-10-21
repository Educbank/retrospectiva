package auth

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	tests := []struct {
		name     string
		userID   uuid.UUID
		email    string
		userName string
		secret   string
		wantErr  bool
	}{
		{
			name:     "valid token generation with custom secret",
			userID:   uuid.New(),
			email:    "test@example.com",
			userName: "Test User",
			secret:   "test-secret-key",
			wantErr:  false,
		},
		{
			name:     "valid token generation with default secret",
			userID:   uuid.New(),
			email:    "user@test.com",
			userName: "User Test",
			secret:   "",
			wantErr:  false,
		},
		{
			name:     "valid token with special characters in name",
			userID:   uuid.New(),
			email:    "joão@example.com",
			userName: "João da Silva",
			secret:   "test-secret",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment
			if tt.secret != "" {
				os.Setenv("JWT_SECRET", tt.secret)
			} else {
				os.Unsetenv("JWT_SECRET")
			}
			jwtSecret = []byte(os.Getenv("JWT_SECRET"))

			// Execute
			token, err := GenerateToken(tt.userID, tt.email, tt.userName)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				require.NoError(t, err)
				assert.NotEmpty(t, token)
				assert.IsType(t, "", token)

				// Verify token structure (basic check)
				assert.Contains(t, token, ".")
				parts := len([]rune(token))
				assert.Greater(t, parts, 100) // JWT tokens are typically long
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Test data
	testUserID := uuid.New()
	testEmail := "test@example.com"
	testName := "Test User"
	testSecret := "test-secret-key"

	// Setup environment
	os.Setenv("JWT_SECRET", testSecret)
	jwtSecret = []byte(testSecret)

	// Generate a valid token for testing
	validToken, err := GenerateToken(testUserID, testEmail, testName)
	require.NoError(t, err)

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
		wantClaims  *Claims
	}{
		{
			name:        "valid token",
			tokenString: validToken,
			wantErr:     false,
			wantClaims: &Claims{
				UserID: testUserID,
				Email:  testEmail,
				Name:   testName,
			},
		},
		{
			name:        "invalid token format",
			tokenString: "invalid.token.format",
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "empty token",
			tokenString: "",
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "malformed token",
			tokenString: "not.a.jwt.token",
			wantErr:     true,
			wantClaims:  nil,
		},
		{
			name:        "token with wrong signature",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			wantErr:     true,
			wantClaims:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute
			claims, err := ValidateToken(tt.tokenString)

			// Assert
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				require.NoError(t, err)
				require.NotNil(t, claims)
				assert.Equal(t, tt.wantClaims.UserID, claims.UserID)
				assert.Equal(t, tt.wantClaims.Email, claims.Email)
				assert.Equal(t, tt.wantClaims.Name, claims.Name)
				assert.NotNil(t, claims.ExpiresAt)
				assert.NotNil(t, claims.IssuedAt)
			}
		})
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	testSecret := "test-secret-key"
	os.Setenv("JWT_SECRET", testSecret)
	jwtSecret = []byte(testSecret)

	// Create an expired token manually
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// Test expired token
	claims, err = ValidateToken(expiredToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Contains(t, err.Error(), "token is expired")
}

func TestValidateToken_WrongSecret(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Generate token with one secret
	os.Setenv("JWT_SECRET", "secret1")
	jwtSecret = []byte("secret1")

	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	// Try to validate with different secret
	os.Setenv("JWT_SECRET", "secret2")
	jwtSecret = []byte("secret2")

	claims, err := ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestClaims_Structure(t *testing.T) {
	// Test Claims struct fields
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Verify struct fields
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
}

func TestGenerateToken_DefaultSecret(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Unset JWT_SECRET to test default behavior
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")

	// Should not error and should use default secret
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Token should be valid
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
}

func TestValidateToken_DefaultSecret(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Unset JWT_SECRET to test default behavior
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	// Validate token
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, "test@example.com", claims.Email)
	assert.Equal(t, "Test User", claims.Name)
}

// Benchmark tests
func BenchmarkGenerateToken(b *testing.B) {
	// Setup
	os.Setenv("JWT_SECRET", "benchmark-secret")
	jwtSecret = []byte("benchmark-secret")

	userID := uuid.New()
	email := "benchmark@example.com"
	name := "Benchmark User"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := GenerateToken(userID, email, name)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkValidateToken(b *testing.B) {
	// Setup
	os.Setenv("JWT_SECRET", "benchmark-secret")
	jwtSecret = []byte("benchmark-secret")

	userID := uuid.New()
	token, err := GenerateToken(userID, "benchmark@example.com", "Benchmark User")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ValidateToken(token)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Test helper functions
func TestTokenRoundTrip(t *testing.T) {
	// Test that a generated token can be validated and returns the same claims
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "roundtrip-secret")
	jwtSecret = []byte("roundtrip-secret")

	userID := uuid.New()
	email := "roundtrip@example.com"
	name := "Roundtrip User"

	// Generate token
	token, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	// Validate token
	claims, err := ValidateToken(token)
	require.NoError(t, err)

	// Verify claims match original data
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, name, claims.Name)
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()))
	assert.True(t, claims.IssuedAt.Time.Before(time.Now().Add(time.Second)))
}

// TestValidateToken_DefaultSecretFallback - Testa se o secret padrão é usado quando JWT_SECRET não está definido
func TestValidateToken_DefaultSecretFallback(t *testing.T) {
	// Setup - garantir que JWT_SECRET não está definido
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Unset JWT_SECRET
	os.Unsetenv("JWT_SECRET")
	jwtSecret = []byte(os.Getenv("JWT_SECRET"))

	// Gerar token com secret padrão
	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	// Validar token - deve usar secret padrão
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
}

// TestValidateToken_InvalidToken - Testa o caso onde token.Valid é false
func TestValidateToken_InvalidToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	// Criar um token válido
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "test@example.com",
		Name:   "Test User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// Simular token inválido alterando o secret após a geração
	os.Setenv("JWT_SECRET", "different-secret")
	jwtSecret = []byte("different-secret")

	// Tentar validar com secret diferente - deve falhar
	claims, err = ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestValidateToken_ManipulatedToken - Testa token manipulado
func TestValidateToken_ManipulatedToken(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	// Gerar token válido
	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	// Manipular o token (alterar último caractere)
	manipulatedToken := token[:len(token)-1] + "X"

	// Validar token manipulado - deve falhar
	claims, err := ValidateToken(manipulatedToken)
	assert.Error(t, err)
	assert.Nil(t, claims)
}

// TestValidateToken_EmptySecret - Testa comportamento com secret vazio
func TestValidateToken_EmptySecret(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Definir JWT_SECRET como string vazia
	os.Setenv("JWT_SECRET", "")
	jwtSecret = []byte("")

	// Gerar token com secret vazio (deve usar padrão)
	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	// Validar token
	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
}

// TestGenerateToken_EmptyInputs - Testa com inputs vazios
func TestGenerateToken_EmptyInputs(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	tests := []struct {
		name     string
		userID   uuid.UUID
		email    string
		userName string
		wantErr  bool
	}{
		{
			name:     "empty email",
			userID:   uuid.New(),
			email:    "",
			userName: "Test User",
			wantErr:  false, // JWT permite email vazio
		},
		{
			name:     "empty name",
			userID:   uuid.New(),
			email:    "test@example.com",
			userName: "",
			wantErr:  false, // JWT permite nome vazio
		},
		{
			name:     "zero UUID",
			userID:   uuid.Nil,
			email:    "test@example.com",
			userName: "Test User",
			wantErr:  false, // JWT permite UUID zero
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email, tt.userName)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}

// TestValidateToken_EdgeCases - Testa casos extremos
func TestValidateToken_EdgeCases(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	tests := []struct {
		name        string
		tokenString string
		wantErr     bool
	}{
		{
			name:        "token with only dots",
			tokenString: "...",
			wantErr:     true,
		},
		{
			name:        "token with spaces",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c ",
			wantErr:     true,
		},
		{
			name:        "token with newlines",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c\n",
			wantErr:     true,
		},
		{
			name:        "token with tabs",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c\t",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.tokenString)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

// TestSecretConsistency - Testa consistência do secret entre geração e validação
func TestSecretConsistency(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	// Teste 1: Secret customizado
	os.Setenv("JWT_SECRET", "custom-secret")
	jwtSecret = []byte("custom-secret")

	userID := uuid.New()
	token, err := GenerateToken(userID, "test@example.com", "Test User")
	require.NoError(t, err)

	claims, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)

	// Teste 2: Mudança de secret após geração
	os.Setenv("JWT_SECRET", "different-secret")
	jwtSecret = []byte("different-secret")

	claims, err = ValidateToken(token)
	assert.Error(t, err)
	assert.Nil(t, claims)

	// Teste 3: Volta ao secret original
	os.Setenv("JWT_SECRET", "custom-secret")
	jwtSecret = []byte("custom-secret")

	claims, err = ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
}

// TestClaimsIntegrity - Testa integridade dos claims
func TestClaimsIntegrity(t *testing.T) {
	// Setup
	originalSecret := os.Getenv("JWT_SECRET")
	defer func() {
		if originalSecret != "" {
			os.Setenv("JWT_SECRET", originalSecret)
		} else {
			os.Unsetenv("JWT_SECRET")
		}
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	}()

	os.Setenv("JWT_SECRET", "test-secret")
	jwtSecret = []byte("test-secret")

	// Dados de teste com caracteres especiais
	userID := uuid.New()
	email := "test+special@example.com"
	name := "José da Silva & Cia."

	// Gerar token
	token, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	// Validar e verificar integridade dos claims
	claims, err := ValidateToken(token)
	require.NoError(t, err)

	// Verificar se todos os dados foram preservados corretamente
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, name, claims.Name)
	assert.NotNil(t, claims.ExpiresAt)
	assert.NotNil(t, claims.IssuedAt)
	assert.True(t, claims.ExpiresAt.Time.After(claims.IssuedAt.Time))
}
