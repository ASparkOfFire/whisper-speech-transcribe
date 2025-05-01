package service

import (
	"context"
	"log"

	"github.com/asparkoffire/whisper-go/pkg/whisper"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/transcriber"
)

type Job struct {
	ID         int
	Payload    []byte
	Ctx        context.Context
	ResultChan chan<- []whisper.Segment
	ErrorChan  chan<- error
}

type Pool struct {
	transcriber transcriber.Transcriber
}

func NewPool(transcriber transcriber.Transcriber) *Pool {
	return &Pool{transcriber: transcriber}
}

func (p *Pool) Worker(id int, jobs <-chan Job) {
	for job := range jobs {
		log.Printf("Worker %d: is transcribing job %d\n", id, job.ID)
		result, err := p.transcriber.TranscribeFromBytes(job.Ctx, job.Payload)
		if err != nil {
			log.Printf("Worker %d: error transcribing job %d: %v", id, job.ID, err)
			job.ErrorChan <- err
			continue
		}
		job.ResultChan <- result
	}
}

func (p *Pool) Start(numWorkers int, jobs <-chan Job) {
	for i := 1; i <= numWorkers; i++ {
		go p.Worker(i, jobs)
	}
}
