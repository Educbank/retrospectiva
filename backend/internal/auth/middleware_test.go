package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthMiddleware(t *testing.T) {
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

	// Test data
	userID := uuid.New()
	email := "test@example.com"
	name := "Test User"

	// Generate valid token
	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
		shouldSetUser  bool
		expectedUserID uuid.UUID
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			expectedError:  "",
			shouldSetUser:  true,
			expectedUserID: userID,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Authorization header required",
			shouldSetUser:  false,
		},
		{
			name:           "invalid bearer format",
			authHeader:     "Invalid " + validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
			shouldSetUser:  false,
		},
		{
			name:           "no bearer prefix",
			authHeader:     validToken,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
			shouldSetUser:  false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
			shouldSetUser:  false,
		},
		{
			name:           "empty bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
			shouldSetUser:  false,
		},
		{
			name:           "malformed token",
			authHeader:     "Bearer not.a.jwt.token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
			shouldSetUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin in test mode
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// Setup route with middleware
			r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute request
			c.Request = req
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			// Check if user context was set
			if tt.shouldSetUser {
				// We need to check the context in the actual handler
				// Since we can't access the context directly, we'll test it differently
				assert.Equal(t, http.StatusOK, w.Code)
			}
		})
	}
}

func TestAuthMiddleware_ContextValues(t *testing.T) {
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

	// Test data
	userID := uuid.New()
	email := "context@example.com"
	name := "Context User"

	// Generate valid token
	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	// Setup route with middleware that checks context
	r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
		// Check if user context was set correctly
		userIDFromContext, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, userID, userIDFromContext)

		userEmailFromContext, exists := c.Get("user_email")
		assert.True(t, exists)
		assert.Equal(t, email, userEmailFromContext)

		userNameFromContext, exists := c.Get("user_name")
		assert.True(t, exists)
		assert.Equal(t, name, userNameFromContext)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with valid token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	// Execute request
	c.Request = req
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOptionalAuthMiddleware(t *testing.T) {
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

	// Test data
	userID := uuid.New()
	email := "optional@example.com"
	name := "Optional User"

	// Generate valid token
	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		shouldSetUser  bool
		expectedUserID uuid.UUID
	}{
		{
			name:           "valid token",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
			shouldSetUser:  true,
			expectedUserID: userID,
		},
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedStatus: http.StatusOK,
			shouldSetUser:  false,
		},
		{
			name:           "invalid bearer format",
			authHeader:     "Invalid " + validToken,
			expectedStatus: http.StatusOK,
			shouldSetUser:  false,
		},
		{
			name:           "no bearer prefix",
			authHeader:     validToken,
			expectedStatus: http.StatusOK,
			shouldSetUser:  false,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedStatus: http.StatusOK,
			shouldSetUser:  false,
		},
		{
			name:           "empty bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusOK,
			shouldSetUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin in test mode
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			// Setup route with optional middleware
			r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			// Create request
			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Execute request
			c.Request = req
			r.ServeHTTP(w, req)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			// Optional middleware should never return error status
			assert.Equal(t, http.StatusOK, w.Code)
		})
	}
}

func TestOptionalAuthMiddleware_ContextValues(t *testing.T) {
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

	// Test data
	userID := uuid.New()
	email := "optional-context@example.com"
	name := "Optional Context User"

	// Generate valid token
	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	// Test with valid token
	t.Run("with valid token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
			// Check if user context was set correctly
			userIDFromContext, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, userID, userIDFromContext)

			userEmailFromContext, exists := c.Get("user_email")
			assert.True(t, exists)
			assert.Equal(t, email, userEmailFromContext)

			userNameFromContext, exists := c.Get("user_name")
			assert.True(t, exists)
			assert.Equal(t, name, userNameFromContext)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Test without token
	t.Run("without token", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
			// Check that no user context was set
			_, exists := c.Get("user_id")
			assert.False(t, exists)

			_, exists = c.Get("user_email")
			assert.False(t, exists)

			_, exists = c.Get("user_name")
			assert.False(t, exists)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)

		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
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

	// Create an expired token manually
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "expired@example.com",
		Name:   "Expired User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with expired token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	// Execute request
	c.Request = req
	r.ServeHTTP(w, req)

	// Assert response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "Invalid token", response["error"])
}

func TestOptionalAuthMiddleware_ExpiredToken(t *testing.T) {
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

	// Create an expired token manually
	userID := uuid.New()
	claims := &Claims{
		UserID: userID,
		Email:  "expired-optional@example.com",
		Name:   "Expired Optional User",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)), // Expired 1 hour ago
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredToken, err := token.SignedString(jwtSecret)
	require.NoError(t, err)

	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
		// Should not have user context set
		_, exists := c.Get("user_id")
		assert.False(t, exists)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// Create request with expired token
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)

	// Execute request
	c.Request = req
	r.ServeHTTP(w, req)

	// Assert response - should still be OK
	assert.Equal(t, http.StatusOK, w.Code)
}

// Benchmark tests
func BenchmarkAuthMiddleware(b *testing.B) {
	// Setup
	os.Setenv("JWT_SECRET", "benchmark-secret")
	jwtSecret = []byte("benchmark-secret")

	userID := uuid.New()
	token, err := GenerateToken(userID, "benchmark@example.com", "Benchmark User")
	if err != nil {
		b.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
	}
}

func BenchmarkOptionalAuthMiddleware(b *testing.B) {
	// Setup
	os.Setenv("JWT_SECRET", "benchmark-secret")
	jwtSecret = []byte("benchmark-secret")

	userID := uuid.New()
	token, err := GenerateToken(userID, "benchmark@example.com", "Benchmark User")
	if err != nil {
		b.Fatal(err)
	}

	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		r.ServeHTTP(w, req)
	}
}

// Test middleware integration
func TestMiddlewareIntegration(t *testing.T) {
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

	os.Setenv("JWT_SECRET", "integration-secret")
	jwtSecret = []byte("integration-secret")

	userID := uuid.New()
	email := "integration@example.com"
	name := "Integration User"

	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	r := gin.New()

	// Protected route
	r.GET("/protected", AuthMiddleware(), func(c *gin.Context) {
		userIDFromContext, _ := c.Get("user_id")
		c.JSON(http.StatusOK, gin.H{
			"message": "protected",
			"user_id": userIDFromContext,
		})
	})

	// Optional protected route
	r.GET("/optional", OptionalAuthMiddleware(), func(c *gin.Context) {
		userIDFromContext, exists := c.Get("user_id")
		if exists {
			c.JSON(http.StatusOK, gin.H{
				"message": "authenticated",
				"user_id": userIDFromContext,
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": "anonymous",
			})
		}
	})

	tests := []struct {
		name            string
		path            string
		authHeader      string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "protected route with valid token",
			path:            "/protected",
			authHeader:      "Bearer " + validToken,
			expectedStatus:  http.StatusOK,
			expectedMessage: "protected",
		},
		{
			name:            "protected route without token",
			path:            "/protected",
			authHeader:      "",
			expectedStatus:  http.StatusUnauthorized,
			expectedMessage: "",
		},
		{
			name:            "optional route with valid token",
			path:            "/optional",
			authHeader:      "Bearer " + validToken,
			expectedStatus:  http.StatusOK,
			expectedMessage: "authenticated",
		},
		{
			name:            "optional route without token",
			path:            "/optional",
			authHeader:      "",
			expectedStatus:  http.StatusOK,
			expectedMessage: "anonymous",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedMessage != "" {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedMessage, response["message"])
			}
		})
	}
}

// TestAuthMiddleware_ErrorMessages - Testa mensagens de erro específicas
func TestAuthMiddleware_ErrorMessages(t *testing.T) {
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

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedError  string
		expectedStatus int
	}{
		{
			name:           "missing authorization header",
			authHeader:     "",
			expectedError:  "Authorization header required",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid bearer format",
			authHeader:     "Invalid token",
			expectedError:  "Bearer token required",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "no bearer prefix",
			authHeader:     "some-token",
			expectedError:  "Bearer token required",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid-token",
			expectedError:  "Invalid token",
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			c.Request = req
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]string
			err := json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedError, response["error"])
		})
	}
}

// TestAuthMiddleware_ContextIntegrity - Testa se o contexto é mantido íntegro
func TestAuthMiddleware_ContextIntegrity(t *testing.T) {
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

	userID := uuid.New()
	email := "context@example.com"
	name := "Context User"

	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, r := gin.CreateTestContext(w)

	r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
		// Verificar se o contexto não foi corrompido
		userIDFromContext, exists := c.Get("user_id")
		assert.True(t, exists)
		assert.Equal(t, userID, userIDFromContext)

		userEmailFromContext, exists := c.Get("user_email")
		assert.True(t, exists)
		assert.Equal(t, email, userEmailFromContext)

		userNameFromContext, exists := c.Get("user_name")
		assert.True(t, exists)
		assert.Equal(t, name, userNameFromContext)

		// Verificar se não há outros valores inesperados no contexto
		_, exists = c.Get("unexpected_key")
		assert.False(t, exists)

		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+validToken)

	c.Request = req
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestOptionalAuthMiddleware_ContextHandling - Testa contexto no middleware opcional
func TestOptionalAuthMiddleware_ContextHandling(t *testing.T) {
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

	userID := uuid.New()
	email := "optional-context@example.com"
	name := "Optional Context User"

	validToken, err := GenerateToken(userID, email, name)
	require.NoError(t, err)

	gin.SetMode(gin.TestMode)

	// Teste com token válido
	t.Run("with valid token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
			// Verificar se o contexto foi definido corretamente
			userIDFromContext, exists := c.Get("user_id")
			assert.True(t, exists)
			assert.Equal(t, userID, userIDFromContext)

			userEmailFromContext, exists := c.Get("user_email")
			assert.True(t, exists)
			assert.Equal(t, email, userEmailFromContext)

			userNameFromContext, exists := c.Get("user_name")
			assert.True(t, exists)
			assert.Equal(t, name, userNameFromContext)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	// Teste sem token
	t.Run("without token", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, r := gin.CreateTestContext(w)

		r.GET("/test", OptionalAuthMiddleware(), func(c *gin.Context) {
			// Verificar que nenhum contexto de usuário foi definido
			_, exists := c.Get("user_id")
			assert.False(t, exists)

			_, exists = c.Get("user_email")
			assert.False(t, exists)

			_, exists = c.Get("user_name")
			assert.False(t, exists)

			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)

		c.Request = req
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

// TestAuthMiddleware_HeaderVariations - Testa variações de headers de autorização
func TestAuthMiddleware_HeaderVariations(t *testing.T) {
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

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "bearer with extra spaces",
			authHeader:     "Bearer  token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "bearer with tabs",
			authHeader:     "Bearer\ttoken",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
		},
		{
			name:           "bearer with newlines",
			authHeader:     "Bearer\ntoken",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
		},
		{
			name:           "multiple authorization headers",
			authHeader:     "Bearer token1, Bearer token2",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "case insensitive bearer",
			authHeader:     "bearer token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
		},
		{
			name:           "mixed case bearer",
			authHeader:     "BeArEr token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Bearer token required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			c.Request = req
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}

// TestAuthMiddleware_EdgeCases - Testa casos extremos
func TestAuthMiddleware_EdgeCases(t *testing.T) {
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

	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "empty bearer token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "bearer with only spaces",
			authHeader:     "Bearer   ",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "bearer with special characters",
			authHeader:     "Bearer !@#$%^&*()",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "bearer with unicode characters",
			authHeader:     "Bearer 🚀🔥💯",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
		{
			name:           "very long bearer token",
			authHeader:     "Bearer " + string(make([]byte, 10000)),
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Invalid token",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, r := gin.CreateTestContext(w)

			r.GET("/test", AuthMiddleware(), func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{"message": "success"})
			})

			req, _ := http.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)

			c.Request = req
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedError != "" {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}
		})
	}
}
