package core

import (
    "context"
    "fmt"
    speech "cloud.google.com/go/speech/apiv1"
    "cloud.google.com/go/speech/apiv1/speechpb"
)

// Method sends recorded audio to the Google Cloud Speech-to-Text API for 
// transcription.
func (kai *Kai) Recognize(audioData []byte) (string, error) {
    ctx := context.Background()
    client, err := speech.NewClient(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to create speech client: %v", err)
    }
    defer client.Close()
    // Configure the request with the correct audio encoding and sample rate
    req := &speechpb.RecognizeRequest{
        Config: &speechpb.RecognitionConfig{
            Encoding:        speechpb.RecognitionConfig_LINEAR16,
            SampleRateHertz: 44100,
            LanguageCode:    "en-US",
        },
        Audio: &speechpb.RecognitionAudio{
            AudioSource: &speechpb.RecognitionAudio_Content{
                Content: audioData,
            },
        },
    }
    // Send the request and get the response
    resp, err := client.Recognize(ctx, req)
    if err != nil {
        return "", fmt.Errorf("failed to recognize speech: %v", err)
    }
    // Process the response and extract the transcribed text
    if len(resp.Results) > 0 && len(resp.Results[0].Alternatives) > 0 {
        return resp.Results[0].Alternatives[0].Transcript, nil
    }
    return "", fmt.Errorf("no transcription results")
}