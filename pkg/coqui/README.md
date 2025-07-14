# 🐸 go-coqui

text-to-speech for Go.

<!-- [![GoDoc](https://pkg.go.dev/badge/github.com/pixellini/go-audiobook/pkg/coqui)](https://pkg.go.dev/github.com/pixellini/go-audiobook/pkg/coqui) -->

## Installation

```bash
go get -u github.com/pixellini/go-audiobook/pkg/coqui
```

### Prerequisites

go-coqui requires [Coqui TTS](https://github.com/coqui-ai/TTS) to be installed and available in your PATH.

#### Install Coqui TTS

```bash
# Using pip
pip install coqui-tts

# Or using conda
conda install -c conda-forge coqui-tts
```

**Verify Installation:**
```bash
# Check Coqui TTS
tts --help
```

## Quick Start

Basic usage with the default configuration:

```go
package main

import (
    "context"
    "log"
    "github.com/pixellini/go-audiobook/pkg/coqui"
)

func main() {
    tts, err := coqui.New(
        coqui.WithModel(coqui.XTTS),
        coqui.WithSpeakerWav("./speakers/speaker.wav"),
        coqui.WithLanguage(coqui.English),
    )
    if err != nil {
        log.Fatalf("failed to initialize TTS: %v", err)
    }
    _, err = tts.SynthesizeContext(context.Background(), "Hello, world!", "output.wav")
    if err != nil {
        log.Fatalf("synthesis failed: %v", err)
    }
}
```

For voice cloning with custom speaker samples, use XTTS:

```go
tts, err := coqui.NewWithXtts(
    "./speakers/speaker.wav", // Path to speaker sample
    coqui.WithLanguage(coqui.English),
)
```

When you need fast synthesis with predefined voices, use VITS.

```go
tts, err := coqui.NewWithVits(
    "p287", // Speaker ID (e.g., "p287")
    coqui.WithDevice(coqui.CPU),
)
```

---

Released under the [MIT License](./LICENSE).