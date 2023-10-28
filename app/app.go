package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
)

type App struct {
	listenAddr string
}

func NewApp(listenAddr string) App {
	return App{listenAddr: listenAddr}
}

func (s App) Run() {
	// Setup
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "OK")
	})

	// Start server
	go func() {
		if err := e.Start(s.listenAddr); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
