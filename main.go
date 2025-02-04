package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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

		// クエリパラメータ "workers" で並列数を指定（デフォルト: 10）
		workersStr := c.QueryParam("workers")
		workers, err := strconv.Atoi(workersStr)
		if err != nil || workers <= 0 {
			workers = 10 // デフォルトの並列数
		}

		var wg sync.WaitGroup
		var mu sync.Mutex
		resultMsg := ""

		start := time.Now()

		for i := 0; i < workers; i++ {
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				log.Info().Msgf("Worker %d started", id)

				// 100MB のスライスを作成してメモリを確保
				data := make([]byte, 100*1024*1024)

				// ダミーデータを入れる
				for i := range data {
					data[i] = byte(i % 256)
				}

				log.Info().Msgf("Worker %d finished", id)

				// 結果を保存
				mu.Lock()
				resultMsg += fmt.Sprintf("Worker %d completed\n", id)
				mu.Unlock()
			}(i)
		}

		wg.Wait()
		duration := time.Since(start)

		// 全ワーカーの完了メッセージを返す
		return c.String(http.StatusOK, fmt.Sprintf("Allocated %d x 100MB memory\nTime taken: %s\n%s", workers, duration, resultMsg))
	})

	e.GET("/test2", func(c echo.Context) error {
		log.Info().Msg("Accessed /test2")
		return c.String(http.StatusOK, "test2")
	})

	e.Logger.Fatal(e.Start(":8080"))
}
