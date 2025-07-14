# 🐸 go-coqui

text-to-speech for Go.

<!-- [![GoDoc](https://pkg.go.dev/badge/github.com/pixellini/go-audiobook/pkg/coqui)](https://pkg.go.dev/github.com/pixellini/go-audiobook/pkg/coqui) -->

## Installation

```bash
go get -u github.com/pixellini/go-audiobook/pkg/coqui
```

### Prerequisites

go-coqui requires [Coqui TTS](https://github.com/coqui-ai/TTS) (Python) to be installed and available in your PATH.

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

In contexts where performance is nice, but not critical, use XTTS with a custom speaker sample. It supports multiple languages and provides excellent voice cloning capabilities.

```go
package main

import (
    "context"
    "log"
    "github.com/pixellini/go-audiobook/pkg/coqui"
)

func main() {
    tts, err := coqui.NewWithXtts(
        "./speakers/speaker.wav", // Path to speaker sample
        coqui.WithLanguage(coqui.English),
        coqui.WithDevice(coqui.Auto),
        coqui.WithMaxRetries(2),
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

When performance and speed are critical, use VITS. It's significantly faster than XTTS but only supports English and predefined speakers.

```go
tts, err := coqui.NewWithVits(
    "p287", // Speaker index (e.g., "p287")
    coqui.WithDevice(coqui.CPU),
)
```

---

Released under the [MIT License](./LICENSE).