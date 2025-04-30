package utils

import (
	"bytes"
	"errors"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"io"
	"os"
)

var ErrInvalidWAV = errors.New("invalid WAV file")

// ReadWAVAndConvert reads and converts WAV data from an io.ReadSeeker.
// It resamples to targetSampleRate and converts to mono if mono=true.
func ReadWAVAndConvert(r io.ReadSeeker, targetSampleRate int, mono bool) ([]int, *audio.Format, error) {
	decoder := wav.NewDecoder(r)
	if !decoder.IsValidFile() {
		return nil, nil, ErrInvalidWAV
	}

	buffer, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, nil, err
	}

	if mono && buffer.Format.NumChannels > 1 {
		buffer = toMono(buffer)
	}

	if buffer.Format.SampleRate != targetSampleRate {
		buffer = naiveDownsample(buffer, targetSampleRate)
	}

	return buffer.Data, buffer.Format, nil
}

// ReadWAVFromFile reads and converts WAV from a file path.
func ReadWAVFromFile(path string, targetSampleRate int, mono bool) ([]int, *audio.Format, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	return ReadWAVAndConvert(file, targetSampleRate, mono)
}

// ReadWAVFromBytes reads and converts WAV from an in-memory []byte buffer.
func ReadWAVFromBytes(data []byte, targetSampleRate int, mono bool) ([]int, *audio.Format, error) {
	reader := bytes.NewReader(data)
	return ReadWAVAndConvert(reader, targetSampleRate, mono)
}

// SaveWAVFile writes int samples to a .wav file in 16-bit PCM format.
func SaveWAVFile(path string, data []int, format *audio.Format) error {
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := wav.NewEncoder(outFile, format.SampleRate, 16, format.NumChannels, 1)
	buffer := &audio.IntBuffer{
		Data:           data,
		Format:         format,
		SourceBitDepth: 16,
	}
	return encoder.Write(buffer)
}

// toMono converts multichannel audio to mono by averaging all channels.
func toMono(buf *audio.IntBuffer) *audio.IntBuffer {
	numChannels := buf.Format.NumChannels
	monoData := make([]int, 0, len(buf.Data)/numChannels)

	for i := 0; i < len(buf.Data); i += numChannels {
		sum := 0
		for j := 0; j < numChannels; j++ {
			sum += buf.Data[i+j]
		}
		monoData = append(monoData, sum/numChannels)
	}

	return &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: 1,
			SampleRate:  buf.Format.SampleRate,
		},
		SourceBitDepth: buf.SourceBitDepth,
		Data:           monoData,
	}
}

// naiveDownsample reduces sample rate by simple decimation (no interpolation).
func naiveDownsample(buf *audio.IntBuffer, targetRate int) *audio.IntBuffer {
	ratio := float64(buf.Format.SampleRate) / float64(targetRate)
	if ratio <= 1.0 {
		return buf // no need to downsample
	}

	step := int(ratio)
	newData := make([]int, 0, len(buf.Data)/step)
	for i := 0; i < len(buf.Data); i += step {
		newData = append(newData, buf.Data[i])
	}

	return &audio.IntBuffer{
		Format: &audio.Format{
			NumChannels: buf.Format.NumChannels,
			SampleRate:  targetRate,
		},
		SourceBitDepth: buf.SourceBitDepth,
		Data:           newData,
	}
}
