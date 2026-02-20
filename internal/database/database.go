package database

import (
	"fmt"
	"sync"

	"github.com/Dawwami/go-order-api/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db   *gorm.DB
	once sync.Once
)

func GetDB() *gorm.DB {
	once.Do(func() {
		config := config.Load()
		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
			config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

		var err error
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed connect to database: " + err.Error())
		}
		fmt.Println("DB initialized!")
	})
	return db
}

func CloseDB() error {
	if db != nil {
		sqlDB, err := db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

// func TestOnce() {
// 	for i := 0; i < 10; i++ {
// 		go func(n int) {
// 			fmt.Printf("goroutine %d calling GetDB()\n", n)
// 			GetDB()
// 		}(i)
// 	}
// 	time.Sleep(2 * time.Second)
// }

// func TestUnsafe() {
// 	db = nil

// 	for i := 0; i < 10; i++ {
// 		go func(n int) {
// 			fmt.Printf("goroutine %d calling GetDBUnsafe()\n", n)
// 			GetDBUnsafe()
// 		}(i)
// 	}
// 	time.Sleep(2 * time.Second)
// }

// func GetDBUnsafe() *gorm.DB {
// 	if db == nil {
// 		config := config.Load()
// 		dsn := fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
// 			config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

// 		var err error
// 		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 		if err != nil {
// 			panic("failed connect to database: " + err.Error())
// 		}
// 		fmt.Println("DB initialized!")
// 	}
// 	return db
// }
