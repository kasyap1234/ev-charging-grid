package main

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/kasyap1234/ev-charging-grid/internal/config"
	"github.com/kasyap1234/ev-charging-grid/internal/handler"
	"github.com/kasyap1234/ev-charging-grid/internal/storage/postgres"
	echo "github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	pool, err := postgres.NewPool(ctx, cfg.GetDatabaseURL())
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	e := echo.New()
	e.HideBanner = true
	handler.NewHandler(postgres.NewChargingRepository(pool)).Register(e)

	go func() {
		if err := e.Start(":8080"); err != nil {
			log.Fatal(err)
		}
	}()
	<-ctx.Done()
	shutDownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := e.Shutdown(shutDownCtx); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("database: %s:%d/%s\n", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	fmt.Printf("redis: %s:%d\n", cfg.Redis.Host, cfg.Redis.Port)
	fmt.Printf("rabbitmq: %s:%d\n", cfg.RabbitMQ.Host, cfg.RabbitMQ.Port)
}
