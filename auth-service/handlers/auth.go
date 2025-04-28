package handlers

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"auth-service/config"
	"auth-service/email"
)

type AuthHandler struct {
	db     *sql.DB
	cfg    *config.Config
	email  *email.EmailService
}

func NewAuthHandler(db *sql.DB, cfg *config.Config, email *email.EmailService) *AuthHandler {
	return &AuthHandler{
		db:    db,
		cfg:   cfg,
		email: email,
	}
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	TwoFactorCode string `json:"two_factor_code,omitempty"`
}

// Проверка сложности пароля
func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("пароль должен содержать минимум 8 символов")
	}
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну заглавную букву")
	}
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну строчную букву")
	}
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return errors.New("пароль должен содержать хотя бы одну цифру")
	}
	return nil
}

// Проверка количества попыток входа
func (h *AuthHandler) checkLoginAttempts(email string) error {
	var attempts int
	err := h.db.QueryRow(
		"SELECT COUNT(*) FROM login_attempts WHERE email = $1 AND created_at > NOW() - INTERVAL '15 minutes'",
		email,
	).Scan(&attempts)

	if err != nil {
		return err
	}

	if attempts >= h.cfg.Security.MaxLoginAttempts {
		return errors.New("слишком много попыток входа. Попробуйте позже")
	}

	return nil
}

// Запись попытки входа
func (h *AuthHandler) recordLoginAttempt(email string, success bool) error {
	_, err := h.db.Exec(
		"INSERT INTO login_attempts (email, success) VALUES ($1, $2)",
		email, success,
	)
	return err
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка сложности пароля
	if err := validatePassword(req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Генерация соли
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации соли"})
		return
	}

	// Хеширование пароля с солью
	hashedPassword, err := bcrypt.GenerateFromPassword(
		append([]byte(req.Password), salt...),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	var userID int
	err = h.db.QueryRow(
		"INSERT INTO users (email, password_hash, password_salt) VALUES ($1, $2, $3) RETURNING id",
		req.Email, string(hashedPassword), base64.StdEncoding.EncodeToString(salt),
	).Scan(&userID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Пользователь с таким email уже существует"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании пользователя"})
		return
	}

	// Генерация токена верификации
	token := generateSecureToken()
	_, err = h.db.Exec(
		"INSERT INTO verification_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		userID, token, time.Now().Add(24*time.Hour),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	// Отправка email с токеном верификации
	if err := h.email.SendVerificationEmail(req.Email, token); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при отправке email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Пользователь успешно зарегистрирован"})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверка количества попыток входа
	if err := h.checkLoginAttempts(req.Email); err != nil {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": err.Error()})
		return
	}

	var (
		userID       int
		passwordHash string
		passwordSalt string
		verified     bool
		isLocked     bool
		lockUntil    time.Time
		twoFactorEnabled bool
	)

	err := h.db.QueryRow(
		`SELECT id, password_hash, password_salt, verified, is_locked, lock_until, two_factor_enabled 
		FROM users WHERE email = $1`,
		req.Email,
	).Scan(&userID, &passwordHash, &passwordSalt, &verified, &isLocked, &lockUntil, &twoFactorEnabled)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			h.recordLoginAttempt(req.Email, false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске пользователя"})
		return
	}

	// Проверка блокировки аккаунта
	if isLocked && time.Now().Before(lockUntil) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Аккаунт заблокирован до " + lockUntil.Format(time.RFC3339)})
		return
	}

	if !verified {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email не подтвержден"})
		return
	}

	// Декодирование соли
	salt, err := base64.StdEncoding.DecodeString(passwordSalt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обработке пароля"})
		return
	}

	// Проверка пароля с солью
	if err := bcrypt.CompareHashAndPassword(
		[]byte(passwordHash),
		append([]byte(req.Password), salt...),
	); err != nil {
		h.recordLoginAttempt(req.Email, false)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный email или пароль"})
		return
	}

	// Проверка двухфакторной аутентификации
	if twoFactorEnabled {
		if req.TwoFactorCode == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Требуется код двухфакторной аутентификации"})
			return
		}
		// Здесь должна быть проверка кода 2FA
	}

	// Создание JWT токена
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(h.cfg.JWTExpiry).Unix(),
		"iat":     time.Now().Unix(),
		"iss":     "auth-service",
	})

	tokenString, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании токена"})
		return
	}

	// Запись успешной попытки входа
	h.recordLoginAttempt(req.Email, true)

	// Сброс счетчика неудачных попыток
	_, err = h.db.Exec(
		"UPDATE users SET failed_login_attempts = 0 WHERE id = $1",
		userID,
	)

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"expires_in": h.cfg.JWTExpiry.Seconds(),
	})
}

func (h *AuthHandler) VerifyEmail(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Токен не указан"})
		return
	}

	var userID int
	err := h.db.QueryRow(
		"SELECT user_id FROM verification_tokens WHERE token = $1 AND expires_at > NOW()",
		token,
	).Scan(&userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Недействительный или просроченный токен"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при проверке токена"})
		return
	}

	_, err = h.db.Exec("UPDATE users SET verified = TRUE WHERE id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении статуса пользователя"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email успешно подтвержден"})
}

// Генерация безопасного токена
func generateSecureToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(b)
} 