package post

import (
	"blog/data"
	"blog/logMnt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePost(c *gin.Context) {
	//获取文章信息
	var post data.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid create post parameter"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//插入文章信息
	if err := db.Create(&post).Error; err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to create post"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}

	zap.L().Info("create post",
		zap.Uint("post_id", post.ID),
		zap.String("title", post.Title),
		zap.String("content", post.Content),
		zap.Time("created_at", post.CreatedAt),
	)
	//返回创建文章成功的信息给客户端
	c.JSON(http.StatusCreated, gin.H{
		"message":    "Post created successfully",
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"created_at": post.CreatedAt.Format(time.RFC3339),
	})
}

func GetPosts(c *gin.Context) {
	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid get posts parameter"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//从数据表获取所有文章信息
	var storedPosts []data.Post
	if err := db.Find(&storedPosts).Error; err != nil {
		zap.L().Error(logMnt.ErrDatabaseConnection.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	zap.L().Info("get post list",
		zap.Int("count", len(storedPosts)),
	)
	//将所有文章信息发送给客户端
	c.JSON(http.StatusOK, gin.H{
		"message": "All posts found successfully",
		"count":   len(storedPosts),
		"data":    storedPosts,
	})
}

func GetPost(c *gin.Context) {
	//获取文章信息
	var post data.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid get post parameter"))
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

	//从数据表获取文章信息
	var storedPost data.Post
	if err := db.Where("id = ?", post.ID).First(&storedPost).Error; err != nil {
		zap.L().Error(logMnt.ErrNotFound.Message, zap.String("error", "Post not found"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	zap.L().Info("get post details",
		zap.Uint("post_id", storedPost.ID),
		zap.String("title", storedPost.Title),
		zap.String("content", storedPost.Content),
		zap.Time("created_at", storedPost.CreatedAt),
	)
	//将文章信息发送给客户端
	c.JSON(http.StatusOK, gin.H{
		"message": "Post found successfully",
		"data":    storedPost,
	})
}

func UpdatePost(c *gin.Context) {
	//获取文章信息
	var updatePost data.Post
	if err := c.ShouldBindJSON(&updatePost); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid update post parameter"))
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

	//查询文章并检查是否存在
	var post data.Post
	if err := db.First(&post, updatePost.ID).Error; err != nil {
		zap.L().Error(logMnt.ErrNotFound.Message, zap.String("error", "Post not found"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	//权限检查：只有文章作者才能更新
	if post.UserID != updatePost.UserID {
		zap.L().Error(logMnt.ErrForbidden.Message, zap.String("error", "User does not match post user"))
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not match post user"})
		return
	}

	//将文章更新至数据库表
	if err := db.Model(&post).Updates(&updatePost).Error; err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to update post"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}

	zap.L().Info("update post",
		zap.Uint("post_id", post.ID),
		zap.String("title", post.Title),
		zap.String("content", post.Content),
		zap.Time("created_at", post.CreatedAt),
		zap.Time("updated_at", post.UpdatedAt),
	)
	//将更新后的文章信息返回给客户端
	c.JSON(http.StatusOK, gin.H{
		"message":    "Post updated successfully",
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"created_at": post.CreatedAt.Format(time.RFC3339),
		"updated_at": post.UpdatedAt.Format(time.RFC3339),
	})
}

func DeletePost(c *gin.Context) {
	//获取文章信息
	var deletePost data.Post
	if err := c.ShouldBindJSON(&deletePost); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid delete post parameter"))
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

	//查询文章并检查是否存在
	var post data.Post
	if err := db.Unscoped().Where("id = ?", deletePost.ID).First(&post).Error; err != nil {
		zap.L().Error(logMnt.ErrNotFound.Message, zap.String("error", "Post not found"))
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	//权限检查：只有文章作者才能删除
	if post.UserID != deletePost.UserID {
		zap.L().Error(logMnt.ErrForbidden.Message, zap.String("error", "User does not match post user"))
		c.JSON(http.StatusForbidden, gin.H{"error": "User does not match post user"})
		return
	}

	//将文章从数据库表中删除
	if err := db.Unscoped().Model(&post).Delete(&deletePost).Error; err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to delete post"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}

	zap.L().Info("delete post",
		zap.Uint("post_id", post.ID),
		zap.String("title", post.Title),
		zap.String("content", post.Content),
		zap.Time("created_at", post.CreatedAt),
		zap.Time("deleted_at", post.DeletedAt.Time),
	)
	//将删除了的文章信息返回给客户端
	c.JSON(http.StatusOK, gin.H{
		"message":    "Post deleted successfully",
		"id":         post.ID,
		"title":      post.Title,
		"content":    post.Content,
		"created_at": post.CreatedAt.Format(time.RFC3339),
		"deleted_at": post.DeletedAt.Time.Format(time.RFC3339),
	})
}
