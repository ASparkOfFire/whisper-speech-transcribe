package handlers

import (
	"errors"
	"github.com/asparkoffire/whisper-speech-transcribe/api"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/config"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/transcriber"
	"github.com/gorilla/websocket"
	"io"

	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// HandleFileUploadTranscribe handles audio file uploads and returns transcription
func HandleFileUploadTranscribe(c *gin.Context) error {
	// Load Whisper model
	whisperTranscriber, err := transcriber.NewWhisperTranscriber(config.AppConfig.ModelPath)
	if err != nil {
		log.Println("Error creating WhisperTranscriber:", err)
		return api.Error{
			Code: http.StatusInternalServerError,
			Err:  "Failed to initialize transcriber",
		}
	}
	defer whisperTranscriber.Unload(c)

	// Parse the uploaded file
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("Error retrieving file:", err)
		return api.Error{
			Code: http.StatusBadRequest,
			Err:  "Invalid audio file",
		}
	}
	defer file.Close()

	// Read file content into memory
	audioData, err := io.ReadAll(file)
	if err != nil {
		log.Println("Error reading audio file:", err)
		return api.Error{
			Code: http.StatusInternalServerError,
			Err:  "Failed to read audio file",
		}
	}

	// Transcribe audio
	segments, err := whisperTranscriber.TranscribeFromBytes(c, audioData)
	if err != nil {
		log.Println("Error during transcription:", err)
		return api.Error{
			Code: http.StatusInternalServerError,
			Err:  "Error during transcription",
		}
	}

	// Respond with transcription
	return api.Response{
		Code: http.StatusOK,
		Msg:  "Successfully transcribed",
		Data: segments,
	}
}

// HandleWebSocketTranscribe WebSocket handler for transcription that processes audio chunks in real-time
func HandleWebSocketTranscribe(c *gin.Context) {
	whisperTranscriber, err := transcriber.NewWhisperTranscriber(config.AppConfig.ModelPath)
	if err != nil {
		log.Println("Error creating WhisperTranscriber:", err)
		return
	}
	defer whisperTranscriber.Unload(c)

	// Upgrade the connection to WebSocket using Gorilla WebSocket
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all connections (be cautious with production)
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	// Channel for handling errors in reading or writing WebSocket messages
	errChan := make(chan error)

	// Goroutine to continuously read WebSocket messages (audio chunks)
	go func() {
		for {
			// Read a chunk of audio data from the WebSocket
			_, p, err := conn.ReadMessage()
			if err != nil {
				if !errors.Is(err, websocket.ErrCloseSent) {
					log.Println("Error reading message from WebSocket:", err)
					errChan <- err // send the error to the main goroutine
				}
				break
			}

			// Process the audio chunk immediately in a separate goroutine
			// Transcribe the audio chunk
			segments, err := whisperTranscriber.TranscribeFromBytes(c, p)
			if err != nil {
				log.Println("Error during transcription:", err)
				errChan <- err // send the error to the main goroutine
				return
			}

			// Send the transcription result back for this chunk
			for _, segment := range segments {
				err := conn.WriteJSON(segment)
				if err != nil {
					log.Println("Error sending WebSocket message:", err)
					errChan <- err // send the error to the main goroutine
					return
				}
			}
		}
	}()

	// Wait until the WebSocket connection is closed or an error occurs
	select {
	case <-c.Request.Context().Done():
		log.Println("WebSocket connection closed")
	case err := <-errChan:
		if err != nil {
			log.Println("Error in WebSocket communication:", err)
			conn.WriteJSON(api.Error{
				Code: http.StatusInternalServerError,
				Err:  "Internal server error during transcription",
			})
		}
	}
}
