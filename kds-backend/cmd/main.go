package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/kodland-community/internal/auth"
	"github.com/kodland-community/internal/handlers"
	"github.com/kodland-community/internal/middleware"
)

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())
	r.Use(middleware.XSSMiddleware())
	r.Use(middleware.SecurityHeadersMiddleware())
	r.Use(middleware.RateLimitMiddleware(100)) 

	// Публичные маршруты
	public := r.Group("/api")
	{
		// Маршруты для index.html
		public.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
		public.GET("/users/:id", func(c *gin.Context) { handlers.GetUserByID(c, db) })
		public.POST("/users/:id/like", func(c *gin.Context) { handlers.LikeUser(c, db) })

		// Маршруты для авторизации
		public.POST("/login", auth.Login)
		public.POST("/refresh", auth.RefreshToken)
	}

	// Защищенные маршруты (только для admin.html)
	admin := r.Group("/api/admin")
	admin.Use(auth.JWTAuthMiddleware())
	{
		admin.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
		admin.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
		admin.PUT("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
		admin.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
} 