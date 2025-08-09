package database

import (
	"gorm.io/gorm"

	"github.com/horlakz/wallet-sync.api/internal/config"
)

type DatabaseInterface interface {
	Connection() *gorm.DB
	Cache() RedisClientInterface
}

type connection struct {
	mysql MySQLClientInterface
	cache RedisClientInterface
}

func StartDatabaseClient(env config.Env) DatabaseInterface {
	return &connection{
		mysql: NewMySQLClient(env),
		cache: NewRedisClient(env),
	}
}

func (conn connection) Connection() *gorm.DB {
	return conn.mysql.Connection()
}

func (conn connection) Cache() RedisClientInterface {
	return conn.cache
}
