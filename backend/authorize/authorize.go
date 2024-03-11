package authorize

import (
	"OnlineChatApplication/MYSQLDB"
	_ "OnlineChatApplication/MYSQLDB"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

var (
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	SecretKey = []byte("dived24")
)

func RegisterHandler(c *gin.Context) {
	MYSQLDB.Register(c)
}

func LoginHandler(c *gin.Context) {
	// 这里应该有一些用户验证的逻辑
	var user MYSQLDB.User

	// 使用 c.ShouldBindJSON() 方法获取 JSON 请求体数据
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if MYSQLDB.UserLogin(user) {
		// 如果用户验证成功，生成一个JWT
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = "username"
		claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

		t, err := token.SignedString(SecretKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
			return
		}

		// 将JWT存储在Redis中
		ctx := context.Background()
		err = Rdb.Set(ctx, "token", t, 0).Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not store token"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": t})

		// 进行登录逻辑处理，比如验证凭证等
		log.Println(user.Email + user.Password)

		// 假设登录成功，返回成功消息
		c.JSON(http.StatusOK, gin.H{"message": "登录成功"})

		// 在这里执行登录逻辑,例如验证用户凭证
		c.JSON(200, gin.H{"message": "Login successful"})

	} else {
		log.Println("LoginFail")
	}
}

func ProtectedHandler(c *gin.Context) {
	// 从请求中获取JWT
	tokenString := c.Request.Header.Get("Authorization")

	// 从Redis中获取JWT
	ctx := context.Background()
	redisToken, err := Rdb.Get(ctx, "token").Result()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// 验证JWT
	if tokenString != redisToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return SecretKey, nil
	})

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		c.JSON(http.StatusOK, gin.H{"name": claims["name"]})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}
