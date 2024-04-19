
package handlers

import (
    "database/sql"
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/arindam923/api-generator/db"
    "github.com/arindam923/api-generator/models"
)

// GetUser retrieves a User by ID
func GetUser(c *gin.Context) {
    id := c.Param("id")
    user := &models.User{}
    err := db.Get(db.DB, user, "SELECT * FROM  WHERE id = $1", id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, user)
}

// ListUsers retrieves all Users
func ListUsers(c *gin.Context) {
    var users []*models.User
    err := db.Select(db.DB, &users, "SELECT * FROM ")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, users)
}

// CreateUser creates a new User
func CreateUser(c *gin.Context) {
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := db.Insert(db.DB, &user, "INSERT INTO  (ID, Name, Email, CreatedAt, UpdatedAt) VALUES ($1, $2, $3, $4, $5) RETURNING id", [ID Name Email CreatedAt UpdatedAt])
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, user)
}

// UpdateUser updates an existing User
func UpdateUser(c *gin.Context) {
    id := c.Param("id")
    var user models.User
    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    user.ID, _ = strconv.ParseUint(id, 10, 64)
    _, err := db.Update(db.DB, &user, "UPDATE  SET ID = $1, Name = $2, Email = $3, CreatedAt = $4, UpdatedAt = $5 WHERE id = $1", [ID Name Email CreatedAt UpdatedAt 0])
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, user)
}

// DeleteUser deletes a User by ID
func DeleteUser(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec(db.DB, "DELETE FROM  WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
