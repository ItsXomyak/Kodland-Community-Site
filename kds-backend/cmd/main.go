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

    // Базовые middleware
    r.Use(middleware.LoggingMiddleware())
    r.Use(middleware.CORSMiddleware())
    r.Use(middleware.XSSMiddleware())
    r.Use(middleware.SecurityHeadersMiddleware())
    r.Use(middleware.RateLimitMiddleware(100))
    r.Use(middleware.RequestSizeMiddleware(10 * 1024 * 1024)) // 10MB

    // Публичные маршруты
    public := r.Group("/api")
    {
        public.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
        public.GET("/users/:id", func(c *gin.Context) { handlers.GetUserByID(c, db) })
        public.POST("/users/:id/like", func(c *gin.Context) { handlers.LikeUser(c, db) })

        public.GET("/tables/public", func(c *gin.Context) { handlers.GetPublicTables(c, db) })
        public.GET("/tables/:id/data", func(c *gin.Context) { handlers.GetTableData(c, db) })

        public.POST("/login", auth.Login)
        public.POST("/refresh", auth.RefreshToken)
    }

    // Защищенные маршруты
    admin := r.Group("/api/admin")
    admin.Use(auth.JWTAuthMiddleware())
    admin.Use(middleware.CSRFMiddleware())
    admin.Use(middleware.IPWhitelistMiddleware([]string{
        "127.0.0.1",
        "::1",
        // Добавьте здесь IP адреса администраторов
    }))
    {
        admin.GET("/users", func(c *gin.Context) { handlers.GetUsers(c, db) })
        admin.POST("/users", func(c *gin.Context) { handlers.CreateUser(c, db) })
        admin.PUT("/users/:id", func(c *gin.Context) { handlers.UpdateUser(c, db) })
        admin.DELETE("/users/:id", func(c *gin.Context) { handlers.DeleteUser(c, db) })

        admin.POST("/tables", func(c *gin.Context) { handlers.CreateTable(c, db) })
        admin.GET("/tables", func(c *gin.Context) { handlers.GetTables(c, db) })
        admin.GET("/tables/:id", func(c *gin.Context) { handlers.GetTable(c, db) })
        admin.PUT("/tables/:id", func(c *gin.Context) { handlers.UpdateTable(c, db) })
        admin.DELETE("/tables/:id", func(c *gin.Context) { handlers.DeleteTable(c, db) })
    }

    port := os.Getenv("APP_PORT")
    if port == "" {
        port = "8080"
    }
    r.Run(":" + port)
}