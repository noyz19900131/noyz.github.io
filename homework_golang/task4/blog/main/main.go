package main

import (
	"blog/cfg"
	"blog/comment"
	"blog/logMnt"
	"blog/post"
	"blog/user"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	// 初始化zap日志
	logger, err := logMnt.InitZapLogger()
	if err != nil {
		panic("Failed to init log: " + err.Error())
	}
	defer logger.Sync() // 确保所有日志都被写入输出

	// 替换全局logger
	zap.ReplaceGlobals(logger)

	r := gin.Default()

	// 移除默认的日志中间件，使用自定义日志中间件
	r.Use(gin.Recovery()) // 恢复panic
	r.Use(logMnt.LoggingMiddleware())
	r.Use(logMnt.ErrorHandlingMiddleware())

	apiGroup := r.Group("/api")
	{
		//用户
		apiUserGroup := apiGroup.Group("/users")
		{
			//用户注册
			apiUserGroup.POST("/register", user.Register)
			//用户登录
			apiUserGroup.POST("/login", user.Login)

			//文章
			apiUserPostGroup := apiUserGroup.Group("/posts")
			{
				//用户认证
				apiUserPostGroup.Use(user.JWTAuthMiddleware())
				{
					//创建文章
					apiUserPostGroup.POST("/create", post.CreatePost)
				}

				//读取文章列表
				apiUserPostGroup.GET("/all/get", post.GetPosts)

				//读取单篇文章
				apiUserPostGroup.GET("/get", post.GetPost)

				//更新文章
				apiUserPostGroup.PUT("/update", post.UpdatePost)

				//删除文章
				apiUserPostGroup.DELETE("/delete", post.DeletePost)

				//评论
				apiUserPostCommentGroup := apiUserPostGroup.Group("/comments")
				{
					//用户认证
					apiUserPostCommentGroup.Use(user.JWTAuthMiddleware())
					{
						//对文章发表评论
						apiUserPostCommentGroup.POST("/create", comment.CreateComment)
					}

					//读取某篇文章的所有评论列表
					apiUserPostCommentGroup.GET("/all/get", comment.GetComments)
				}
			}
		}
	}

	zap.L().Info("start server", zap.Uint("port", cfg.CFG.Server.Port))
	err = r.Run(":" + strconv.Itoa(int(cfg.CFG.Server.Port))) // 监听并在 0.0.0.0:8080 上启动服务
	if err != nil {
		panic(err.Error())
		os.Exit(1)
	}
}
