package handlers

import (
    "context"
    "encoding/json"
    "log"
    "strconv"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/jackc/pgx/v5"

    "github.com/kodland-community/internal/models"
)

func CreateTable(c *gin.Context, db *pgx.Conn) {
    var req models.TableRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    if req.Data == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data cannot be empty"})
        return
    }

    if len(req.Data) > 0 && req.Data[0] == '{' {
        var temp interface{}
        if err := json.Unmarshal([]byte(req.Data), &temp); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Data must be a valid JSON string"})
            return
        }
    }

    var tableID int
    err := db.QueryRow(context.Background(),
        `INSERT INTO tables (name, description, data, created_by, is_public) 
         VALUES ($1, $2, $3, $4, $5) 
         RETURNING id`,
        req.Name, req.Description, req.Data, userID, req.IsPublic,
    ).Scan(&tableID)
    if err != nil {
        log.Println("Failed to create table:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
        return
    }

    table := models.Table{
        ID:          tableID,
        Name:        req.Name,
        Description: req.Description,
        Data:        req.Data,
        CreatedBy:   userID.(int),
        IsPublic:    req.IsPublic,
    }

    c.JSON(http.StatusCreated, table)
}

func UpdateTable(c *gin.Context, db *pgx.Conn) {
    tableID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var req models.TableRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    if req.Data == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Data cannot be empty"})
        return
    }

    if len(req.Data) > 0 && req.Data[0] == '{' {
        var temp interface{}
        if err := json.Unmarshal([]byte(req.Data), &temp); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Data must be a valid JSON string"})
            return
        }
    }

    result, err := db.Exec(context.Background(),
        `UPDATE tables 
         SET name = $1, description = $2, data = $3, is_public = $4 
         WHERE id = $5 AND created_by = $6`,
        req.Name, req.Description, req.Data, req.IsPublic, tableID, userID,
    )
    if err != nil {
        log.Println("Failed to update table:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table"})
        return
    }

    if result.RowsAffected() == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Table not found or not authorized"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Table updated"})
}

func GetTables(c *gin.Context, db *pgx.Conn) {
    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    rows, err := db.Query(context.Background(),
        `SELECT id, name, description, data, created_by, is_public, created_at 
         FROM tables 
         WHERE created_by = $1 
         ORDER BY created_at DESC`,
        userID,
    )
    if err != nil {
        log.Println("Failed to fetch tables:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
        return
    }
    defer rows.Close()

    var tables []models.Table
    for rows.Next() {
        var table models.Table
        if err := rows.Scan(&table.ID, &table.Name, &table.Description, &table.Data, &table.CreatedBy, &table.IsPublic, &table.CreatedAt); err != nil {
            log.Println("Failed to scan table:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
            return
        }
        tables = append(tables, table)
    }

    c.JSON(http.StatusOK, tables)
}

func GetTable(c *gin.Context, db *pgx.Conn) {
    tableID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    var table models.Table
    err = db.QueryRow(context.Background(),
        `SELECT id, name, description, data, created_by, is_public, created_at 
         FROM tables 
         WHERE id = $1 AND created_by = $2`,
        tableID, userID,
    ).Scan(&table.ID, &table.Name, &table.Description, &table.Data, &table.CreatedBy, &table.IsPublic, &table.CreatedAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
        } else {
            log.Println("Failed to fetch table:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch table"})
        }
        return
    }

    c.JSON(http.StatusOK, table)
}

func DeleteTable(c *gin.Context, db *pgx.Conn) {
    tableID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
        return
    }

    userID, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
        return
    }

    result, err := db.Exec(context.Background(),
        `DELETE FROM tables WHERE id = $1 AND created_by = $2`,
        tableID, userID,
    )
    if err != nil {
        log.Println("Failed to delete table:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table"})
        return
    }

    if result.RowsAffected() == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Table deleted"})
}

func GetPublicTables(c *gin.Context, db *pgx.Conn) {
    rows, err := db.Query(context.Background(),
        `SELECT id, name, description, data, created_by, is_public, created_at 
         FROM tables 
         WHERE is_public = true 
         ORDER BY created_at DESC`,
    )
    if err != nil {
        log.Println("Failed to fetch public tables:", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
        return
    }
    defer rows.Close()

    var tables []models.Table
    for rows.Next() {
        var table models.Table
        if err := rows.Scan(&table.ID, &table.Name, &table.Description, &table.Data, &table.CreatedBy, &table.IsPublic, &table.CreatedAt); err != nil {
            log.Println("Failed to scan table:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
            return
        }
        tables = append(tables, table)
    }

    c.JSON(http.StatusOK, tables)
}

func GetTableData(c *gin.Context, db *pgx.Conn) {
    tableID, err := strconv.Atoi(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid table ID"})
        return
    }

    var table models.Table
    err = db.QueryRow(context.Background(),
        `SELECT id, name, description, data, created_by, is_public, created_at 
         FROM tables 
         WHERE id = $1 AND is_public = true`,
        tableID,
    ).Scan(&table.ID, &table.Name, &table.Description, &table.Data, &table.CreatedBy, &table.IsPublic, &table.CreatedAt)
    if err != nil {
        if err == pgx.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Table not found or not public"})
        } else {
            log.Println("Failed to fetch table:", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch table"})
        }
        return
    }

    c.JSON(http.StatusOK, table)
}