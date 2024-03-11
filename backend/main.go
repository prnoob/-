package main

import (
	"OnlineChatApplication/MYSQLDB"
	_ "OnlineChatApplication/MYSQLDB"
	"OnlineChatApplication/authorize"
	"OnlineChatApplication/websocket"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	log.Println("Gin start")
	r := gin.Default()

	// 使用中间件传递消息
	// 使用默认CORS配置
	r.Use(cors.Default())

	MYSQLDB.InitDB()

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/login")
	})

	// 使用中间件的路由组
	authorized := r.Group("/")
	authorized.Use(RedirectMiddleware())
	{
		authorized.GET("/protected", authorize.ProtectedHandler)
		authorized.GET("/online-chat-page")
	}

	// 不使用中间件的路由组
	noAuth := r.Group("/")
	{
		noAuth.GET("/login", func(c *gin.Context) {
			log.Println("跳转到登陆界面")
		})
		noAuth.POST("/login", authorize.LoginHandler)
		noAuth.POST("/register", authorize.RegisterHandler)
		// 其他不需要认证的路由...
	}
	setupRoutes(r)

	r.Run(":8080")
}

// 首先尝试升级HTTP连接到WebSocket连接，如果失败则返回错误。
//创建一个新的WebSocket客户端，并将其注册到连接池中，最后开始读取来自客户端的消息。
func serveWs(pool *websocket.Pool, c *gin.Context) {
	log.Println("WebSocket Endpoint Hit")
	conn, err := websocket.Upgrade(c.Writer, c.Request)
	if err != nil {
		fmt.Fprintf(c.Writer, "%+v\n", err)
		return
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

//设置路由，创建一个websocket 连接池，并启动
// 在/ws路径设置了一个get路由，调用serveWs 处理ws连接
func setupRoutes(router *gin.Engine) {
	pool := websocket.NewPool()
	go pool.Start()

	router.GET("/ws", func(c *gin.Context) {
		serveWs(pool, c)
	})
}

func RedirectMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取JWT
		tokenString := c.GetHeader("Authorization")

		// 验证JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("invalid signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return authorize.SecretKey, nil
		})

		if err != nil || !token.Valid {
			// 如果JWT无效，重定向到登录页面
			c.Redirect(http.StatusTemporaryRedirect, "/login")
			c.Abort()
		} else {
			// 如果JWT有效，重定向到聊天页面
			c.Redirect(http.StatusTemporaryRedirect, "/chat")
			c.Abort()
		}
	}
}
