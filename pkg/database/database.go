package database

import (
	"database/sql"
	"goblog/pkg/logger"
	"time"

	"github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// Initialize 初始化数据库
func Initialize() {
    initDB()
    createTables()
}

func initDB() {
    var err error
    config := mysql.Config{
        User : "root",
        Passwd: "",
        Addr: "127.0.0.1:3306",
        Net: "tcp",
        DBName: "goblog",
        AllowNativePasswords: true,
    }
    //连接数据库
    DB, err = sql.Open("mysql", config.FormatDSN())
    logger.LogError(err)

    DB.SetMaxOpenConns(25) //最大连接数
    DB.SetMaxIdleConns(25) //最大空闲数
    DB.SetConnMaxLifetime(5 * time.Minute)//每个连接的过期时间

    err = DB.Ping()
    logger.LogError(err)
}

func createTables() {
    createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
    );`
    _, err := DB.Exec(createArticlesSQL)
    logger.LogError(err)
}