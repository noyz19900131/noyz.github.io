package comment

import (
	"blog/data"
	"blog/logMnt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreateComment(c *gin.Context) {
	//获取评论信息
	var comment data.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid create comment parameter"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrDatabaseConnection.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//插入评论信息
	if err := db.Create(&comment).Error; err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to create comment"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	zap.L().Info("create comment",
		zap.Uint("comment_id", comment.ID),
		zap.String("content", comment.Content),
		zap.Time("created_at", comment.CreatedAt),
	)
	//返回发布评论成功的信息给客户端
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Comment created successfully",
		"id":         comment.ID,
		"content":    comment.Content,
		"created_at": comment.CreatedAt.Format(time.RFC3339),
	})
}

func GetComments(c *gin.Context) {
	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrDatabaseConnection.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//从数据表获取所有评论信息
	var storedComments []data.Comment
	if err := db.Find(&storedComments).Error; err != nil {
		zap.L().Error(logMnt.ErrDatabaseConnection.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get comments"})
		return
	}

	zap.L().Info("get comments",
		zap.Int("count", len(storedComments)),
	)
	//将所有评论信息发送给客户端
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment found successfully",
		"count":   len(storedComments),
		"data":    storedComments,
	})
}
