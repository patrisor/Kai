package core

import (
	"fmt"
	"github.com/gordonklaus/portaudio"
)

// Method captures audio input from the microphone using the portaudio library.
func (kai *Kai) Listen(stop <-chan struct{}) ([]byte, error) {
    // Initialize the PortAudio library
    err := portaudio.Initialize()
    if err != nil {
        return nil, fmt.Errorf("failed to initialize PortAudio: %v", err)
    }
    defer portaudio.Terminate()
    // Create an input buffer to store audio samples
    in := make([]int16, 64)
    // Open a default stream for audio input
    stream, err := portaudio.OpenDefaultStream(1, 0, 44100, len(in), in)
    if err != nil {
        return nil, fmt.Errorf("failed to open default stream: %v", err)
    }
    defer stream.Close()
    // Start the audio stream
    err = stream.Start()
    if err != nil {
        return nil, fmt.Errorf("failed to start the stream: %v", err)
    }
    defer stream.Stop()
    // Record audio data
    var audioData []int16
    for {
        select {
        case <-stop:
            // Stop recording when a signal is received on the stop channel
            return convertToBytes(audioData), nil
        default:
            // Read audio samples into the input buffer
            err := stream.Read()
            if err != nil {
                return nil, fmt.Errorf("failed to read from stream: %v", err)
            }
            // Append the samples to the audio data slice
            audioData = append(audioData, in...)
        }
    }
}

// Method converts audio data from int16 to a byte slice.
func convertToBytes(audioData []int16) []byte {
    audioBytes := make([]byte, len(audioData)*2)
    for i, sample := range audioData {
        audioBytes[i*2] = byte(sample)
        audioBytes[i*2+1] = byte(sample >> 8)
    }
    return audioBytes
}