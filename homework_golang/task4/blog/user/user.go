package user

import (
	"blog/cfg"
	"blog/data"
	"blog/logMnt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 认证中间件：验证用户是否已登录
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头获取Authorization字段
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "No authorization header"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthenticated, please log in first"})
			//终止请求处理
			c.Abort()
			return
		}

		//验证Authorization格式(必须是 "Bearer <token>")
		if len(authHeader) < 7 || !strings.HasPrefix(authHeader, "Bearer ") {
			zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Invalid authorization header"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format; \"Bearer <token>\""})
		}
		tokenString := authHeader[7:]

		//解析并验证token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			//验证签名方法是否为HS256
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(cfg.CFG.Jwt.Secret), nil
		})

		//检查token有效性
		if err != nil || !token.Valid {
			zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Invalid token or token has expired"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token or token has expired"})
			c.Abort()
			return
		}

		//提取token中的用户ID(需与生成token时的字段一致)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Failed to parse token"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Failed to parse token"})
			c.Abort()
			return
		}

		//将用户ID存入上下文，供后续接口使用
		userID, ok := claims["id"].(float64)
		if !ok {
			zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Invalid user information in token"))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user information in token"})
			c.Abort()
			return
		}
		//存储用户ID(uint类型)
		c.Set("userID", uint(userID))

		//继续处理请求
		c.Next()
	}
}

func Register(c *gin.Context) {
	//获取用户注册信息
	var user data.User
	if err := c.ShouldBindJSON(&user); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid register parameter"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Failed to hash password"))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = string(hashedPassword)

	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//插入用户信息
	if err := db.Create(&user).Error; err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to create user"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	zap.L().Info("user register",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("password", user.Password),
	)
	//返回用户注册成功的信息给客户端
	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	//获取用户登录信息
	var user data.User
	if err := c.ShouldBindJSON(&user); err != nil {
		zap.L().Error(logMnt.ErrBadRequest.Message, zap.String("error", "Invalid login parameter"))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	//连接数据库
	db := data.ConnectDatabase()
	if db == nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to connect to database"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to database"})
		return
	}

	//从数据表获取用户信息
	var storedUser data.User
	if err := db.Where("username = ?", user.Username).First(&storedUser).Error; err != nil {
		zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Invalid username or password"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	//验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(storedUser.Password), []byte(user.Password)); err != nil {
		zap.L().Error(logMnt.ErrUnauthorized.Message, zap.String("error", "Invalid username or password"))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	//创建JWT声明
	claims := jwt.MapClaims{
		"id":       storedUser.ID,
		"username": storedUser.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //24小时后过期
	}

	//生成JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//签名令牌
	tokenString, err := token.SignedString([]byte(cfg.CFG.Jwt.Secret))
	if err != nil {
		zap.L().Error(logMnt.ErrInternalServerError.Message, zap.String("error", "Failed to generate token"))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
	}

	zap.L().Info("user login",
		zap.Uint("user_id", storedUser.ID),
		zap.String("username", storedUser.Username),
		zap.String("token", tokenString),
	)
	//返回令牌给客户端
	c.JSON(http.StatusOK, gin.H{
		"message": "User login successfully",
		"user": gin.H{
			"id":       storedUser.ID,
			"username": storedUser.Username,
		},
		"token":   tokenString,
		"expires": claims["exp"],
	})
}
