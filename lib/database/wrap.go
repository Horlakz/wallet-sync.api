package database

import "gorm.io/gorm"

type txWrapper struct {
	tx *gorm.DB
}

// Cache implements DatabaseInterface.
func (w *txWrapper) Cache() RedisClientInterface {
	panic("unimplemented")
}

func (w *txWrapper) Connection() *gorm.DB {
	return w.tx
}

func Wrap(tx *gorm.DB) DatabaseInterface {
	return &txWrapper{tx: tx}
}
