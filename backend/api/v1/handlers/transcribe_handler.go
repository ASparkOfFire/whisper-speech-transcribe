package handlers

import (
	"errors"
	"github.com/asparkoffire/whisper-go/pkg/whisper"
	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/asparkoffire/whisper-speech-transcribe/internal/config"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/service"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/transcriber"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	jobIDCounter int32 = 0
)

type TranscribeHandler struct {
	Pool *service.Pool
	Jobs chan service.Job
}

func NewTranscribeHandler() (*TranscribeHandler, error) {
	tr, err := transcriber.NewWhisperTranscriber(config.AppConfig.ModelPath)
	if err != nil {
		log.Printf("Error creating transcriber handler: %v", err)
		return nil, err
	}

	jobs := make(chan service.Job, 100)
	pool := service.NewPool(tr)
	pool.Start(4, jobs) // Start 4 workers

	return &TranscribeHandler{
		Jobs: jobs,
		Pool: pool,
	}, nil
}

func (h *TranscribeHandler) HandleWebSocketTranscribe(c *gin.Context) {
	// Init transcriber
	whisperTranscriber, err := transcriber.NewWhisperTranscriber(config.AppConfig.ModelPath)
	if err != nil {
		log.Println("Error creating WhisperTranscriber:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize transcriber"})
		return
	}
	defer whisperTranscriber.Unload(c)

	// WebSocket upgrade
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	// Listen for audio chunks
	for {
		_, audioData, err := conn.ReadMessage()
		if err != nil {
			if !errors.Is(err, websocket.ErrCloseSent) {
				log.Println("WebSocket read error:", err)
			}
			break
		}

		resultChan := make(chan []whisper.Segment)
		errorChan := make(chan error)

		job := service.Job{
			ID:         int(atomic.AddInt32(&jobIDCounter, 1)),
			Payload:    audioData,
			Ctx:        c, // Pass request context
			ResultChan: resultChan,
			ErrorChan:  errorChan,
		}

		h.Jobs <- job

		// Handle job result or error
		select {
		case result := <-resultChan:
			for _, segment := range result {
				if err := conn.WriteJSON(segment); err != nil {
					log.Println("Error sending transcription segment:", err)
					return
				}
			}
		case err := <-errorChan:
			log.Println("Transcription error:", err)
			conn.WriteJSON(api.Error{
				Code: http.StatusInternalServerError,
				Err:  "Transcription failed",
				Data: err.Error(),
			})
			continue
		}
	}
}
