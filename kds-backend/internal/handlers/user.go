package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/kodland-community/internal/models"
)

func Register(r *gin.RouterGroup, db *pgx.Conn) {
	r.GET("/users", func(c *gin.Context) {
		GetUsers(c, db)
	})
	r.GET("/users/:id", func(c *gin.Context) {
		GetUserByID(c, db)
	})
	r.POST("/users", func(c *gin.Context) {
		CreateUser(c, db)
	})
	r.PUT("/users/:id", func(c *gin.Context) {
		UpdateUser(c, db)
	})
	r.DELETE("/users/:id", func(c *gin.Context) {
		DeleteUser(c, db)
	})
	r.POST("/users/:id/like", func(c *gin.Context) {
		LikeUser(c, db)
	})
}

func GetUsers(c *gin.Context, db *pgx.Conn) {
	// наши параметры запроса
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "12"))
	search := c.Query("search")
	role := c.Query("role")
	age := c.Query("age")
	country := c.Query("country")
	achievement := c.Query("achievement")
	sort := c.DefaultQuery("sort", "recent")

	offset := (page - 1) * limit

	query := `
		SELECT id, nickname, avatar_url, links, achievements, roles, medals, likes, joined_at, created_at
		FROM users
		WHERE 
			($1 = '' OR nickname ILIKE '%' || $1 || '%')
			AND ($2 = '' OR roles @> jsonb_build_array($2))
			AND ($3 = '' OR roles @> jsonb_build_array($3))
			AND ($4 = '' OR roles @> jsonb_build_array($4))
			AND ($5 = '' OR achievements @> jsonb_build_array($5))
	`
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE 
			($1 = '' OR nickname ILIKE '%' || $1 || '%')
			AND ($2 = '' OR roles @> jsonb_build_array($2))
			AND ($3 = '' OR roles @> jsonb_build_array($3))
			AND ($4 = '' OR roles @> jsonb_build_array($4))
			AND ($5 = '' OR achievements @> jsonb_build_array($5))
	`

	// Сортировка
	switch sort {
	case "recent":
		query += " ORDER BY created_at DESC"
	case "nickname":
		query += " ORDER BY nickname ASC"
	case "time":
		query += " ORDER BY joined_at DESC"
	case "achievements":
		query += " ORDER BY jsonb_array_length(achievements) DESC"
	case "likes":
		query += " ORDER BY likes DESC"
	case "medals":
		query += " ORDER BY (medals->>'gold')::int * 3 + (medals->>'silver')::int * 2 + (medals->>'bronze')::int DESC"
	default:
		query += " ORDER BY created_at DESC"
	}

	query += " LIMIT $6 OFFSET $7"

	// Запрос на получение пользователей
	rows, err := db.Query(context.Background(), query, search, role, age, country, achievement, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Nickname,
			&user.AvatarURL,
			&user.Links,
			&user.Achievements,
			&user.Roles,
			&user.Medals,
			&user.Likes,
			&user.JoinedAt,
			&user.CreatedAt,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		users = append(users, user)
	}

	// Получаем общее количество записей
	var total int
	err = db.QueryRow(context.Background(), countQuery, search, role, age, country, achievement).Scan(&total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Отправляем ответ
	c.JSON(http.StatusOK, models.UsersResponse{
		Users: users,
		Total: total,
	})
}

func GetUserByID(c *gin.Context, db *pgx.Conn) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	query := `
		SELECT id, nickname, avatar_url, links, achievements, roles, medals, likes, joined_at, created_at
		FROM users
		WHERE id = $1
	`

	var user models.User
	err = db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Links,
		&user.Achievements,
		&user.Roles,
		&user.Medals,
		&user.Likes,
		&user.JoinedAt,
		&user.CreatedAt,
	)

	if err == pgx.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context, db *pgx.Conn) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// значения по умолчанию если нету
	if req.Medals == nil {
		req.Medals = map[string]int{"gold": 0, "silver": 0, "bronze": 0}
	}
	if req.Links == nil {
		req.Links = []string{}
	}
	if req.Achievements == nil {
		req.Achievements = []string{}
	}
	if req.Roles == nil {
		req.Roles = []string{}
	}
	if req.JoinedAt.IsZero() {
		req.JoinedAt = time.Now()
	}

	query := `
		INSERT INTO users (nickname, avatar_url, links, achievements, roles, medals, likes, joined_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, nickname, avatar_url, links, achievements, roles, medals, likes, joined_at, created_at
	`

	var user models.User
	err := db.QueryRow(context.Background(), query,
		req.Nickname,
		req.AvatarURL,
		req.Links,
		req.Achievements,
		req.Roles,
		req.Medals,
		req.Likes,
		req.JoinedAt,
	).Scan(
		&user.ID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Links,
		&user.Achievements,
		&user.Roles,
		&user.Medals,
		&user.Likes,
		&user.JoinedAt,
		&user.CreatedAt,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_nickname_key" (SQLSTATE 23505)` {
			c.JSON(http.StatusConflict, gin.H{"error": "nickname already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context, db *pgx.Conn) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Проверяем, существует ли пользователь
	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	// Формируем запрос с динамическими обновлениями
	query := `
		UPDATE users
		SET
			nickname = COALESCE($1, nickname),
			avatar_url = COALESCE($2, avatar_url),
			links = COALESCE($3, links),
			achievements = COALESCE($4, achievements),
			roles = COALESCE($5, roles),
			medals = COALESCE($6, medals),
			likes = COALESCE($7, likes),
			joined_at = COALESCE($8, joined_at)
		WHERE id = $9
		RETURNING id, nickname, avatar_url, links, achievements, roles, medals, likes, joined_at, created_at
	`

	var user models.User
	var joinedAtArg interface{}
	if req.JoinedAt != nil {
		joinedAtArg = *req.JoinedAt
	}
	var likesArg interface{}
	if req.Likes != nil {
		likesArg = *req.Likes
	}

	err = db.QueryRow(context.Background(), query,
		req.Nickname, req.AvatarURL, req.Links, req.Achievements, req.Roles, req.Medals, likesArg, joinedAtArg, id,
	).Scan(
		&user.ID,
		&user.Nickname,
		&user.AvatarURL,
		&user.Links,
		&user.Achievements,
		&user.Roles,
		&user.Medals,
		&user.Likes,
		&user.JoinedAt,
		&user.CreatedAt,
	)

	if err != nil {
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_nickname_key" (SQLSTATE 23505)` {
			c.JSON(http.StatusConflict, gin.H{"error": "nickname already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func DeleteUser(c *gin.Context, db *pgx.Conn) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// существует ли пользователь
	var exists bool
	err = db.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	_, err = db.Exec(context.Background(), "DELETE FROM users WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

// LikeUser обрабатывает лайк пользователя
func LikeUser(c *gin.Context, db *pgx.Conn) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID пользователя"})
		return
	}

	var likes int
	err = db.QueryRow(context.Background(), `
		UPDATE users 
		SET likes = likes + 1 
		WHERE id = $1 
		RETURNING likes
	`, userID).Scan(&likes)
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении лайков"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"likes": likes})
}