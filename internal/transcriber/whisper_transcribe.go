package transcriber

import (
	"context"
	"github.com/asparkoffire/whisper-go/pkg/whisper"
	"github.com/asparkoffire/whisper-speech-transcribe/internal/utils"
)

type Transcriber interface {
	TranscribeFromFile(ctx context.Context, filePath string) ([]whisper.Segment, error)
	TranscribeFromBytes(ctx context.Context, fileData []byte) ([]whisper.Segment, error)
	Unload(ctx context.Context) error
}
type WhisperTranscriber struct {
	model whisper.Model
}

func NewWhisperTranscriber(modelPath string) (Transcriber, error) {
	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, err
	}
	return &WhisperTranscriber{model: model}, nil
}

func (t *WhisperTranscriber) TranscribeFromFile(ctx context.Context, filePath string) ([]whisper.Segment, error) {
	context, err := t.model.NewContext()
	if err != nil {
		return nil, err
	}

	// Read and preprocess the WAV file: convert to mono, 16kHz
	intSamples, _, err := utils.ReadWAVFromFile(filePath, 16000, true)
	if err != nil {
		return nil, err
	}

	// Convert int samples to float32
	samples := intToFloat32(intSamples)

	if err := context.Process(samples, nil, nil, nil); err != nil {
		return nil, err
	}

	var segments []whisper.Segment
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (t *WhisperTranscriber) TranscribeFromBytes(ctx context.Context, fileData []byte) ([]whisper.Segment, error) {
	context, err := t.model.NewContext()
	if err != nil {
		return nil, err
	}

	// Read and preprocess the WAV file: convert to mono, 16kHz
	intSamples, _, err := utils.ReadWAVFromBytes(fileData, 16000, true)
	if err != nil {
		return nil, err
	}

	// Convert int samples to float32
	samples := intToFloat32(intSamples)

	if err := context.Process(samples, nil, nil, nil); err != nil {
		return nil, err
	}

	var segments []whisper.Segment
	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		segments = append(segments, segment)
	}

	return segments, nil
}

func (t *WhisperTranscriber) Unload(ctx context.Context) error {
	return t.model.Close()
}

// intToFloat32 converts int16 PCM samples to float32 for Whisper processing.
func intToFloat32(intSamples []int) []float32 {
	samples := make([]float32, len(intSamples))
	for i, sample := range intSamples {
		// Convert int16 to float32 by dividing by the max int16 value
		samples[i] = float32(sample) / float32(1<<15)
	}
	return samples
}
