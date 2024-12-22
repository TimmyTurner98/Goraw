package goraw

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        handler,
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	// Логируем запуск сервера
	log.Println("Starting server on https://localhost:8443")

	err := http.ListenAndServeTLS(":8443", ".tls/server.crt", ".tls/server.key", nil)
	if err != nil {
		log.Fatalf("Failed to start HTTPS server: %v", err)
	}
	return err
}
