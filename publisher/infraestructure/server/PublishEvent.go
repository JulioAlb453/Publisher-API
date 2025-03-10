package server

import (
	"encoding/json"
	"log"
	"net/http"
	"publisher/publisher/application"
	"publisher/publisher/domain"
)

type Server struct {
    service *application.EventService
}

func NewServer(repo domain.EventPublisher) *Server {
    return &Server{service: application.NewEventService(repo)}
}

func (s *Server) Start() error {
    http.HandleFunc("/publish", s.handlerPublishEvent)
    log.Println("Servidor iniciando en :8081")
    return http.ListenAndServe(":8081", nil)
}

func (s *Server) handlerPublishEvent(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
            http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
            return
    }

    var event domain.Event
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
            http.Error(w, "Error decodificando evento", http.StatusBadRequest)
            return
    }
    if err := s.service.PublishEvent(event); err != nil {
            log.Printf("Error al publicar el evento: %v", err) 
            http.Error(w, "Error publicando evento", http.StatusInternalServerError)
            return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Evento publicado con exito"))
}