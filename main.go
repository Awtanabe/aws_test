package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name  string `json:"name"`
	Price int    `json:"price"`
}


func connectDB() *gorm.DB {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port, name)

	maxRetry := 10
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetry; i++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("Database connection successful")
			return db
		}

		log.Printf("Database connection failed: %v. Retrying... (%d/%d)\n", err, i+1, maxRetry)
		time.Sleep(1 * time.Second)
	}

	log.Fatal("Failed to connect to database after multiple attempts")
	return nil
}



func main() {
	e := echo.New()

	db := connectDB()

	db.AutoMigrate(&Product{})

	e.GET("/", func(c echo.Context) error {
		products := []Product{}
		db.Find(&products)
		return c.JSON(http.StatusOK, products)
	})
	e.Logger.Fatal(e.Start(":8080"))
}