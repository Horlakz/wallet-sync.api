package database

import (
	"fmt"
	"time"

	"github.com/horlakz/wallet-sync.api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DatabaseFacade *gorm.DB

type MySQLClientInterface interface {
	Connection() *gorm.DB
}

type mysqlClient struct {
	database *gorm.DB
}

func NewMySQLClient(env config.Env) MySQLClientInterface {
	dsn := env.DB_USER + ":" + env.DB_PASSWORD + "@tcp(" + env.DB_HOST + ":" + env.DB_PORT + ")/" + env.DB_NAME + "?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
	})

	if err != nil {
		fmt.Printf("Database connection error: %v\n", err)
		panic(err)
	}

	sqlDb, err := database.DB()
	if err != nil {
		fmt.Printf("Failed to get SQL database: %v\n", err)
		panic(err)
	}

	// Configure connection pool
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetMaxOpenConns(10)
	sqlDb.SetConnMaxLifetime(time.Hour)

	fmt.Println("MySQL database connection is successful")

	DatabaseFacade = database

	return &mysqlClient{
		database: database,
	}
}

func (conn mysqlClient) Connection() *gorm.DB {
	return conn.database
}
