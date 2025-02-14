#include <cstring>
#include <iostream>
#include <portaudio.h>
#include <vosk_api.h>

#define SAMPLE_RATE 16000
#define FRAMES_PER_BUFFER 512

// Structure for audio data
typedef struct {
  char buffer[FRAMES_PER_BUFFER * 2]; // Audio buffer (16-bit PCM)
  int bufferSize;
} AudioData;

// PortAudio callback function to capture microphone input
static int audioCallback(const void *inputBuffer, void *outputBuffer,
                         unsigned long framesPerBuffer,
                         const PaStreamCallbackTimeInfo *timeInfo,
                         PaStreamCallbackFlags statusFlags, void *userData) {
  AudioData *data = (AudioData *)userData;
  if (inputBuffer == nullptr)
    return paContinue;

  memcpy(data->buffer, inputBuffer, framesPerBuffer * 2);
  data->bufferSize = framesPerBuffer * 2;
  return paContinue;
}

int main() {
  Pa_Initialize();

  // Initialize Vosk model and recognizer
  VoskModel *model = vosk_model_new("/path/to/vosk-model");
  VoskRecognizer *recognizer = vosk_recognizer_new(model, SAMPLE_RATE);

  AudioData audioData = {{0}, 0};

  // Open PortAudio stream for microphone input
  PaStream *stream;
  Pa_OpenDefaultStream(&stream,
                       1,       // Input channels
                       0,       // Output channels
                       paInt16, // 16-bit PCM format
                       SAMPLE_RATE, FRAMES_PER_BUFFER, audioCallback,
                       &audioData);

  Pa_StartStream(stream);

  std::cout << "Listening for keyword 'iris'...\n";

  while (true) {
    if (audioData.bufferSize > 0) {
      if (vosk_recognizer_accept_waveform(recognizer, audioData.buffer,
                                          audioData.bufferSize)) {
        std::string result = vosk_recognizer_result(recognizer);
        if (result.find("iris") != std::string::npos) {
          std::cout << "Keyword 'iris' detected!\n";
          // Trigger your AI assistant here
        }
      }
      audioData.bufferSize = 0; // Reset buffer size for the next batch
    }
    Pa_Sleep(100); // Sleep to avoid busy-waiting
  }

  Pa_StopStream(stream);
  Pa_CloseStream(stream);
  Pa_Terminate();

  vosk_recognizer_free(recognizer);
  vosk_model_free(model);

  return 0;
}
