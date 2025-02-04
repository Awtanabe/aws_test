package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
			log.Info().Msg("Successfully connected to database")
			return db
		}

		log.Error().Err(err).Msgf("Database connection failed. Retrying... (%d/%d)", i+1, maxRetry)
		time.Sleep(1 * time.Second)
	}

	log.Fatal().Msg("Failed to connect to database after multiple attempts")
	return nil
}

func main() {
	// ログファイルを開く（なければ作成）
	logFile, err := os.OpenFile("/var/log/go_app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open log file")
	}
	defer logFile.Close()

	// zerolog の出力先をファイルに設定
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = zerolog.New(logFile).With().Timestamp().Logger()

	e := echo.New()

	e.GET("/", func(c echo.Context) error {
		log.Info().Msg("Accessed root path")
		return c.String(http.StatusOK, "root path")
	})

	e.GET("/test", func(c echo.Context) error {
		log.Info().Msg("Accessed /test")
		return c.String(http.StatusOK, "test")
	})

	e.GET("/burden_test", func(c echo.Context) error {
		log.Info().Msg("Accessed /burden_test")

		for i:= 0; i < 1000; i++ {
			log.Print("print", i)
		}
		return c.String(http.StatusOK, "burden_test")
	})

	e.GET("/test2", func(c echo.Context) error {
		log.Info().Msg("Accessed /test2")
		return c.String(http.StatusOK, "test2")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
