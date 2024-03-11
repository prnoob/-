package MYSQLDB

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" // 使用匿名导入导入MySQL驱动
	"log"
	"strings"
)

type User struct {
	UserId   int
	UserName string
	Password string
	Email    string
}

const (
	USERNAME = "root"
	PASSWORD = "123"
	IP       = "localhost"
	PORT     = "3307"
	dbName   = "loginserver"
)

var db *sql.DB

func InitDB() *sql.DB {
	dsn := strings.Join([]string{USERNAME, ":", PASSWORD, "@tcp(", IP, ":", PORT, ")/", dbName, "?charset=utf8"}, "")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	// 尝试与数据库建立连接以验证 DSN 是否正确
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Database connection established")
	return db
}

func InsertUser(user User) bool {
	//开启事务
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	defer tx.Rollback() // 事务回滚

	// 检查用户名是否已存在
	var identifyEmail string
	queryCheck := "SELECT id FROM users WHERE name = ? LIMIT 1"
	err = tx.QueryRow(queryCheck, user.Email).Scan(&identifyEmail)
	if !errors.Is(err, sql.ErrNoRows) {
		fmt.Println("User already exists")
		return false
	}

	// 准备SQL语句
	query := "INSERT INTO users(name, password, email) VALUES(?, ?, ?)"

	// 执行SQL语句
	res, err := db.Exec(query, user.UserName, user.Password, user.Email)
	if err != nil {
		log.Printf("Error inserting new user: %v\n", err)
		return false
	}

	tx.Commit()
	fmt.Println(res.LastInsertId())
	return true
}

func UserLogin(user User) bool {
	tx, err := db.Begin()
	if err != nil {
		fmt.Println("tx fail")
		return false
	}
	defer tx.Rollback() // 如果后续有错误发生，确保回滚事务

	// 修改查询语句，只根据用户名来查找用户
	query := "SELECT id, password FROM users WHERE name = ? LIMIT 1"

	// 由于我们只根据用户名查找，只需传入用户名作为参数
	// 使用Scan来获取用户ID和密码，这里假设你想要检查密码是否匹配
	var id int
	var password string
	err = tx.QueryRow(query, user.UserName).Scan(&id, &password)
	if err != nil {
		if errors.Is(sql.ErrNoRows, err) {
			// 如果没有找到用户，可能想要记录或处理这种情况
			fmt.Println("User not found")
			return false
		}
		log.Printf("Error finding user: %v\n", err)
		return false
	}

	if user.Password == password {
		user.UserId = id
		tx.Commit()
		fmt.Println("Login successful")
		return true

	}
	// 如果需要，这里可以对用户对象做进一步操作，例如比对密码
	// user.Password = password // 如果你想在函数外部比对密码
	// user.ID = id // 保存查找到的用户ID

	return false
}

func Register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if InsertUser(user) {
		// 在这里执行注册逻辑,例如将用户信息存储在数据库中
		c.JSON(200, gin.H{"message": "Registration successful"})
	} else {
		log.Println("register failed")
		c.JSON(401, gin.H{"message": "Registration fail"})
	}
}
