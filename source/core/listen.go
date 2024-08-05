package core

import (
	"fmt"
	"os"
	"encoding/binary"
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

// Method saves the recorded audio data to a WAV file.
func saveToWavFile(filename string, audioData []byte) error {
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create file: %v", err)
    }
    defer file.Close()
    // WAV file header
    var header = []byte{
        'R', 'I', 'F', 'F',
        0, 0, 0, 0, // ChunkSize (to be filled later)
        'W', 'A', 'V', 'E',
        'f', 'm', 't', ' ',
        16, 0, 0, 0, // Subchunk1Size (16 for PCM)
        1, 0, // AudioFormat (1 for PCM)
        1, 0, // NumChannels (1 for mono)
        0x44, 0xac, 0, 0, // SampleRate (44100 Hz)
        0x88, 0x58, 1, 0, // ByteRate (SampleRate * NumChannels * BitsPerSample/8)
        2, 0, // BlockAlign (NumChannels * BitsPerSample/8)
        16, 0, // BitsPerSample (16 bits)
        'd', 'a', 't', 'a',
        0, 0, 0, 0, // Subchunk2Size (to be filled later)
    }
    // Fill in the ChunkSize and Subchunk2Size
    chunkSize := 36 + len(audioData)
    subchunk2Size := len(audioData)
    binary.LittleEndian.PutUint32(header[4:], uint32(chunkSize))
    binary.LittleEndian.PutUint32(header[40:], uint32(subchunk2Size))
    // Write the header and audio data to the file
    _, err = file.Write(header)
    if err != nil {
        return fmt.Errorf("failed to write header: %v", err)
    }
    _, err = file.Write(audioData)
    if err != nil {
        return fmt.Errorf("failed to write audio data: %v", err)
    }
    return nil
}