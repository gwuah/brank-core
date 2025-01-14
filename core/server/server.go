package core

import (
	"brank/core"
	"brank/core/auth"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
	config *core.Config
}

func NewHTTPServer(cfg *core.Config, auth *auth.Auth) *Server {
	engine := gin.Default()
	engine.Use(auth.CORS())
	engine.Use(auth.ExtractTokenFromAuthHeader(cfg))

	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	})
	return &Server{
		config: cfg,
		Engine: engine,
	}
}

func (s *Server) Start(cleanup func()) {

	h := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.config.PORT),
		Handler: s.Engine,
	}

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		if err := h.Close(); err != nil {
			log.Println("failed To ShutDown Server", err)
		}
		log.Println("Shut Down Server")
	}()

	if err := h.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server Closed After Interruption")
		} else {
			log.Println("Unexpected Server Shutdown")
		}
	}

	cleanup()
}
