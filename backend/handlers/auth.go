package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/AbdelrahmanKhaled21/musicauniversalis/config"
	"github.com/AbdelrahmanKhaled21/musicauniversalis/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"gorm.io/gorm"
)

type AuthHandler struct {
	config *config.Config
	db     *gorm.DB
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthHandler(config *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{
		config: config,
		db:     db,
	}
}

// InitiateOAuth starts the GitHub OAuth flow
func (h *AuthHandler) InitiateOAuth(c *gin.Context) {
	oauth2Config := &oauth2.Config{
		ClientID:     h.config.GitHubClientID,
		ClientSecret: h.config.GitHubClientSecret,
		RedirectURL:  "http://localhost:8080/api/v1/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Generate state parameter for security
	state := uuid.New().String()

	// Store state in session or cache (simplified for now)
	// In production, use Redis or similar for state management

	authURL := oauth2Config.AuthCodeURL(state)
	c.JSON(http.StatusOK, gin.H{
		"auth_url": authURL,
		"state":    state,
	})
}

// OAuthCallback handles the OAuth callback from GitHub
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not provided"})
		return
	}

	// Verify state parameter (simplified for now)
	if state == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "State parameter missing"})
		return
	}

	oauth2Config := &oauth2.Config{
		ClientID:     h.config.GitHubClientID,
		ClientSecret: h.config.GitHubClientSecret,
		RedirectURL:  "http://localhost:8080/api/v1/auth/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}

	// Exchange code for token
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Printf("Failed to exchange code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate"})
		return
	}

	// Get user info from GitHub
	client := oauth2Config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Printf("Failed to get user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var githubUser GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		log.Printf("Failed to decode user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode user info"})
		return
	}

	// Get user email if not provided
	if githubUser.Email == "" {
		emailsResp, err := client.Get("https://api.github.com/user/emails")
		if err == nil {
			defer emailsResp.Body.Close()
			var emails []struct {
				Email   string `json:"email"`
				Primary bool   `json:"primary"`
			}
			if json.NewDecoder(emailsResp.Body).Decode(&emails) == nil {
				for _, email := range emails {
					if email.Primary {
						githubUser.Email = email.Email
						break
					}
				}
			}
		}
	}

	// Find or create user in database
	var user models.User
	result := h.db.Where("oauth_provider = ? AND oauth_id = ?", "github", fmt.Sprintf("%d", githubUser.ID)).First(&user)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			// Create new user
			user = models.User{
				Email:         githubUser.Email,
				Name:          githubUser.Name,
				OAuthProvider: "github",
				OAuthID:       fmt.Sprintf("%d", githubUser.ID),
			}

			if err := h.db.Create(&user).Error; err != nil {
				log.Printf("Failed to create user: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
				return
			}
		} else {
			log.Printf("Database error: %v", result.Error)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
	}

	// Generate JWT token
	tokenString, err := h.generateJWT(user.ID, user.Email)
	if err != nil {
		log.Printf("Failed to generate JWT: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// Redirect to frontend with token as query parameter
	frontendURL := "http://localhost:3000"
	redirectURL := fmt.Sprintf("%s?token=%s&user_id=%d&email=%s&name=%s",
		frontendURL, tokenString, user.ID, user.Email, user.Name)

	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// generateJWT creates a JWT token for the user
func (h *AuthHandler) generateJWT(userID uint, email string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.config.JWTSecret))
}

// AuthMiddleware validates JWT tokens
func (h *AuthHandler) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(h.config.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Add user info to context
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
