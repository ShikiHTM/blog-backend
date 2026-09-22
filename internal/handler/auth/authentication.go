package auth

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shikihtm/blog-backend/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type AuthenticationHandler struct {
	users  []model.User
	config Config
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewAuthenticationHandler(cfg Config) *AuthenticationHandler {
	data, err := os.ReadFile(cfg.UserCredentialsFilePath)
	if err != nil {
		panic("Error: " + err.Error())
	}

	var users []model.User
	if err := json.Unmarshal(data, &users); err != nil {
		panic("Error: " + err.Error())
	}

	return &AuthenticationHandler{
		users:  users,
		config: cfg,
	}
}

func (auth *AuthenticationHandler) createJWTToken(user model.User) (string, error) {
	secretKey := auth.config.JWTSecretKey

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &model.JwtClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	signedKey, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("Error when signing token: %w", err)
	}

	return signedKey, nil
}

func (auth *AuthenticationHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "BAD_REQUEST",
			"message": "Invalid request body.",
		})
		return
	}

	var user *model.User
	for i := range auth.users {
		if auth.users[i].Username == req.Username {
			user = &auth.users[i]
			break
		}
	}

	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "INVALID_CREDENTIALS",
			"message": "Incorrect username or password.",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "INVALID_CREDENTIALS",
			"message": "Incorrect username or password.",
		})
		return
	}

	token, err := auth.createJWTToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"Error":   "INTERNAL_SERVER_ERROR",
			"message": err.Error(),
		})
		return
	}

	c.SetCookie("jwt_token", token, 86400, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successfully.",
	})
}

func (auth *AuthenticationHandler) RegisterRoutes(c *gin.RouterGroup) {
	authentication := c.Group("/auth")
	{
		authentication.POST("/login", auth.Login)
	}
}
