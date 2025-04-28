package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	"auth-service/internal/logger"
	"auth-service/internal/models"
)

type ProfileHandler struct {
	db *sql.DB
}

func NewProfileHandler(db *sql.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

type ProfileResponse struct {
	ID        uuid.UUID            `json:"id"`
	Email     string               `json:"email"`
	Login     string               `json:"login"`
	Bio       string               `json:"bio"`
	Avatar    string               `json:"avatar"`
	CreatedAt time.Time            `json:"createdAt"`
	UpdatedAt time.Time            `json:"updatedAt"`
	Comments  []models.ProfileComment `json:"comments"`
}

type UpdateProfileRequest struct {
	Login string `json:"login"`
	Bio   string `json:"bio"`
}

type AddCommentRequest struct {
	Content string `json:"content"`
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		logger.Error("Invalid user ID:", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Получаем данные пользователя
	user, err := h.getUserByID(userID)
	if err != nil {
		logger.Error("Error getting user:", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Получаем данные профиля
	profile, err := h.getUserProfile(userID)
	if err != nil {
		logger.Error("Error getting profile:", err)
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	// Получаем комментарии
	comments, err := h.getProfileComments(userID)
	if err != nil {
		logger.Error("Error getting comments:", err)
		http.Error(w, "Failed to get comments", http.StatusInternalServerError)
		return
	}

	response := ProfileResponse{
		ID:        user.ID,
		Email:     user.Email,
		Login:     profile.Login,
		Bio:       profile.Bio,
		Avatar:    profile.AvatarURL,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Comments:  comments,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ProfileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		logger.Error("Invalid user ID:", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Error decoding request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем профиль
	err = h.updateUserProfile(userID, req)
	if err != nil {
		logger.Error("Error updating profile:", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *ProfileHandler) AddComment(w http.ResponseWriter, r *http.Request) {
	userID, err := uuid.Parse(r.Header.Get("X-User-ID"))
	if err != nil {
		logger.Error("Invalid user ID:", err)
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var req AddCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("Error decoding request:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Добавляем комментарий
	err = h.addProfileComment(userID, req.Content)
	if err != nil {
		logger.Error("Error adding comment:", err)
		http.Error(w, "Failed to add comment", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ProfileHandler) getUserByID(id uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, email, is_active, role, created_at, updated_at
		FROM users WHERE id = $1`
	user := &models.User{}
	err := h.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.IsActive,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (h *ProfileHandler) getUserProfile(userID uuid.UUID) (*models.UserProfile, error) {
	query := `
		SELECT user_id, login, avatar_url, bio, created_at, updated_at
		FROM user_profiles WHERE user_id = $1`
	profile := &models.UserProfile{}
	err := h.db.QueryRow(query, userID).Scan(
		&profile.UserID,
		&profile.Login,
		&profile.AvatarURL,
		&profile.Bio,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return profile, nil
}

func (h *ProfileHandler) getProfileComments(userID uuid.UUID) ([]models.ProfileComment, error) {
	query := `
		SELECT id, user_id, profile_id, content, created_at, updated_at
		FROM profile_comments 
		WHERE profile_id = $1
		ORDER BY created_at DESC`
	rows, err := h.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.ProfileComment
	for rows.Next() {
		var comment models.ProfileComment
		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.ProfileID,
			&comment.Content,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}
	return comments, nil
}

func (h *ProfileHandler) addProfileComment(userID uuid.UUID, content string) error {
	query := `
		INSERT INTO profile_comments (id, user_id, profile_id, content)
		VALUES ($1, $2, $3, $4)`
	_, err := h.db.Exec(query, uuid.New(), userID, userID, content)
	return err
}

func (h *ProfileHandler) updateUserProfile(userID uuid.UUID, req UpdateProfileRequest) error {
	query := `
		UPDATE user_profiles
		SET login = $1, bio = $2, updated_at = NOW()
		WHERE user_id = $3`
	_, err := h.db.Exec(query, req.Login, req.Bio, userID)
	return err
}

func (h *ProfileHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/profile", h.GetProfile).Methods("GET")
	router.HandleFunc("/profile", h.UpdateProfile).Methods("PUT")
	router.HandleFunc("/profile/comments", h.AddComment).Methods("POST")
} 