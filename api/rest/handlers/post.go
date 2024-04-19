
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

// GetPost retrieves a Post by ID
func GetPost(c *gin.Context) {
    id := c.Param("id")
    post := &models.Post{}
    err := db.Get(db.DB, post, "SELECT * FROM  WHERE id = $1", id)
    if err != nil {
        if err == sql.ErrNoRows {
            c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }
    c.JSON(http.StatusOK, post)
}

// ListPosts retrieves all Posts
func ListPosts(c *gin.Context) {
    var posts []*models.Post
    err := db.Select(db.DB, &posts, "SELECT * FROM ")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, posts)
}

// CreatePost creates a new Post
func CreatePost(c *gin.Context) {
    var post models.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := db.Insert(db.DB, &post, "INSERT INTO  (ID, Content, AuthorId, TimeStamp) VALUES ($1, $2, $3, $4) RETURNING id", [ID Content AuthorId TimeStamp])
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, post)
}

// UpdatePost updates an existing Post
func UpdatePost(c *gin.Context) {
    id := c.Param("id")
    var post models.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    post.ID, _ = strconv.ParseUint(id, 10, 64)
    _, err := db.Update(db.DB, &post, "UPDATE  SET ID = $1, Content = $2, AuthorId = $3, TimeStamp = $4 WHERE id = $1", [ID Content AuthorId TimeStamp 0])
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, post)
}

// DeletePost deletes a Post by ID
func DeletePost(c *gin.Context) {
    id := c.Param("id")
    _, err := db.Exec(db.DB, "DELETE FROM  WHERE id = $1", id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
