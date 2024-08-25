package core

import (
    "os"
	"fmt"
    "time"
	"context"
	"encoding/binary"
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

        // TODO: Pick voice
        Voice: &texttospeechpb.VoiceSelectionParams{
            LanguageCode: "en-US",
            Name:         "en-US-Polyglot-1",
            SsmlGender:   texttospeechpb.SsmlVoiceGender_MALE,
        },
        // Voice: &texttospeechpb.VoiceSelectionParams{
        //     LanguageCode: "en-AU",
        //     Name:         "en-AU-Wavenet-C",
        //     SsmlGender:   texttospeechpb.SsmlVoiceGender_FEMALE,
        // },        

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

// Helper Method to play audio using PortAudio
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

// TODO: Move into utility methods
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