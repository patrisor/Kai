package core

import (
	"fmt"
    "time"
	"context"
	// Text-to-Speech
	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
    // Audio
	"github.com/gordonklaus/portaudio"
)

// Synthesizes speech from the input text and plays it.
func (kai *Kai) Speak(text string) error {
    ctx := context.Background()
    client, err := texttospeech.NewClient(ctx)
    if err != nil {
        return fmt.Errorf("failed to create text-to-speech client: %v", err)
    }
    defer client.Close()
    // Perform the text-to-speech request
    req := &texttospeechpb.SynthesizeSpeechRequest{
        Input: &texttospeechpb.SynthesisInput{
            InputSource: &texttospeechpb.SynthesisInput_Text{
                Text: text,
            },
        },
        Voice: &texttospeechpb.VoiceSelectionParams{
            LanguageCode: "en-US",
            Name:         "en-US-Polyglot-1",
            SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
        },
        AudioConfig: &texttospeechpb.AudioConfig{
            AudioEncoding:    texttospeechpb.AudioEncoding_LINEAR16,
            SpeakingRate:     1.0,
            Pitch:            0.0,
            SampleRateHertz:  int32(kai.SampleRate),
        },
    }
    resp, err := client.SynthesizeSpeech(ctx, req)
    if err != nil {
        return fmt.Errorf("failed to synthesize speech: %v", err)
    }
    // Play the audio (you can use any audio playing library)
    if err := playAudio(resp.AudioContent, kai.SampleRate); err != nil {
        return fmt.Errorf("failed to play audio: %v", err)
    }
    return nil
}

// Helper Mmthod to play audio using PortAudio
func playAudio(audioData []byte, sampleRate int) error {
    // Initialize PortAudio
    err := portaudio.Initialize()
    if err != nil {
        return fmt.Errorf("failed to initialize PortAudio: %v", err)
    }
    defer portaudio.Terminate()
    // Define a buffer size and create a buffer
    bufferSize := 1024
    dataSize := len(audioData) / 2
    currentIndex := 0
    // Open a stream for audio playback
    stream, err := portaudio.OpenDefaultStream(
        0, 1, float64(sampleRate), bufferSize, 
        func(out []int16) {
            for i := range out {
                if currentIndex < dataSize {
                    out[i] = int16(audioData[2 * currentIndex]) | 
                             int16(audioData[2 * currentIndex + 1]) << 8
                    currentIndex++
                } else {
                    out[i] = 0 // Fill with silence if the data is finished
                }
            }
        },
    )
    if err != nil {
        return fmt.Errorf("failed to open stream: %v", err)
    }
    defer stream.Close()
    // Start the audio stream
    if err := stream.Start(); err != nil {
        return fmt.Errorf("failed to start stream: %v", err)
    }
    // Wait for the audio to finish playing
    for currentIndex < dataSize {
        time.Sleep(10 * time.Millisecond)
    }
    // Stop the audio stream
    if err := stream.Stop(); err != nil {
        return fmt.Errorf("failed to stop stream: %v", err)
    }
    return nil
}