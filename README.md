# go-audiobook

A Go-based tool for converting text files into high-quality audiobooks using text-to-speech technology.

<div align="center">
  <img src="https://raw.githubusercontent.com/pixellini/go-audiobook/main/assets/logo.png" alt="Logo">
</div>



## 📚 Overview

This repository hosts a work-in-progress tool that transforms plain text and EPUB files into natural-sounding audiobooks using advanced TTS synthesis.

## 🚧 Current Status

⚠️ **Work in Progress** – This project is still under active development. Functionality is incomplete, and things may break or change frequently.

### 🛠️ Feature List

- 📄 Support for reading plain .txt files  
- 📦 Containerise inside Docker  
- ⏳ Progress tracking with estimated time remaining  
- 🌐 Basic web interface (potentially using WASM)  

and of course... A better name than "go-audiobook" 😄

## 🚀 Getting Started

Instructions for installing and using this tool will be added once it's ready for public use.

### 📦 Requirements

To run this project, you'll need the following dependencies installed:

#### **[Coqui TTS](https://github.com/coqui-ai/TTS)**
Used for generating natural-sounding speech  

Install via pip:

```bash
pip install TTS
```

#### **[FFmpeg](https://github.com/FFmpeg/FFmpeg)**
Handles audio processing and conversion

On macOS (with Homebrew):

```bash
brew install ffmpeg
```

On Ubuntu/Debian:

```bash
sudo apt install ffmpeg
```


## ⚙️ Configuration

The application is configured using a `config.json` file in the project root. Below are the available options:

| Option                | Type    | Accepted Values                | Description                                                                 |
|-----------------------|---------|--------------------------------|-----------------------------------------------------------------------------|
| `epub_path`           | string  | –                              | Path to the input EPUB file.                                                |
| `image_path`          | string  | –                              | Path to the cover image for the audiobook.                                  |
| `speaker_wav`         | string  | –                              | Path to the narrator's voice sample (`.wav` or `.mp3`).                     |
| `dist_dir`            | string  | –                              | Output directory for generated audiobook files.                             |
| `output_format`       | string  | `m4b`, `mp3`, `m4a`, `wav`     | Output file format: `m4b` (default), `mp3`, `m4a` (AAC), or `wav`. See below for details. |
| `verbose_logs`        | bool    | -                              | If `true`, enables detailed error and debug logs.                           |
| `test_mode`        | bool    | -                              | If `true`, only processes the first 3 chapters for quick testing.            |
| `tts.max_retries`     | int     | –                              | Number of times to retry TTS synthesis on failure.                          |
| `tts.parallel_audio_count`  | int     | -                              | Number of audio files to generate in parallel (see recommendations below).  |
| `tts.use_vits`        | bool    | -                              | If `true`, uses the VITS model for TTS.                                     |
| `tts.vits_voice`      | string  | e.g. `p287`, `p330`            | VITS voice id. Only used if `tts.use_vits` is `true`.                       |
| `tts.device`         | string  | `auto`, `cpu`, `cuda`, `mps`   | Device for Coqui TTS. This helps with GPU acceleration. |

### Example `config.json`

```json
{
  "epub_path": "./book/book.epub",
  "image_path": "./book/cover.png",
  "speaker_wav": "./speakers/speaker.wav",
  "dist_dir": "./.dist",
  "output_format": "m4b",
  "verbose_logs": false,
  "tts": {
    "max_retries": 3,
    "parallel_audio_count": 4,
    "use_vits": false,
    "vits_voice": "p287",
    "device": "auto"
  }
}
```

> **Note:** All paths are relative to the project root unless otherwise specified.

### Command Line Options

| Flag        | Description                                                                                     |
|-------------|-------------------------------------------------------------------------------------------------|
| `--reset`   | Removes all previously processed audiobook files and restarts the creation process. If this flag is not set, processing will resume from where it stopped. |
| `--finish`  | Creates the final audiobook from the chapters processed so far. Useful if you stopped the process partway through and still want a playable, though partial, audiobook. |

## 🧩 Guide & Customisation

### 🎤 Providing a Narrator Voice (.wav)

You can customise your audiobook by providing a short audio sample of your chosen narrator's voice. This sample will be used as the voice for the entire audiobook, allowing you to generate content in any of the available languages. Coqui TTS synthesises speech using the unique vocal characteristics from the sample, not the spoken content.

- **Recommended length:** 1–3 minutes of clear, uninterrupted speech is usually enough for high-quality results. Samples longer than 3 minutes may increase processing time. If you have a powerful machine, you can experiment with longer inputs, but 3 minutes is sufficient for most cases. 
- **File format:** A `.wav` file is preferred for best audio quality. `.mp3` is also supported, and while it can speed up processing time, it may reduce quality.
- **Setup:** Place your audio file in the `speakers/` directory and update the `speaker_wav` path in your `config.json` accordingly.  

> 📝 **Note:** The spoken language in the sample doesn't affect the output. Coqui will synthesise your audiobook in the target language using the voice's characteristics from the sample.

### 🏎️ VITS Model

The VITS model (`tts_models/en/vctk/vits`) offers a much faster and more convenient way to generate audiobooks compared to XTTS with a custom `.wav` speaker file. With VITS, you simply select a voice by id (e.g. `p287`).

> **Performance Note:** VITS can be approximately 80% faster (or more) than XTTS with a `.wav` speaker.

#### 🎤 Choosing a VITS Voice
To see all available voices for the VITS model, run:

```bash
tts --model_name "tts_models/en/vctk/vits" --list_speaker_idxs
```

> Note: The command only lists voice numbers, without names or descriptions. To preview the available VITS voices, refer to the `examples/vits` directory, which includes audio files for each VCTK speaker (e.g. `p225.mp3`, `p226.mp3`, etc.). Below are some recommended voices with their key characteristics:

| Voice # | Gender | Description                                    |
| ------- | ------ | ---------------------------------------------- |
| p340    | Male   | Lively, expressive, clear communicator         |
| p330    | Male   | Authoritative, composed, documentary style     |
| p306    | Female | Composed, articulate, quick paced              |
| p287    | Male   | Deep-toned, resonant, cinematic                |
| p285    | Male   | Soft-spoken, measured, contemplative           |
| p267    | Male   | Intense, dramatic, grand, epic fantasy feel    |
| p262    | Male   | Introspective, philosophical, warm-hearted     |
| p258    | Male   | Formal, precise, analytical, newsreader style  |

Experiment with different voices to find the best fit for your audiobook!

### ⚡️ Optimising Audiobook Creation

#### GPU Acceleration
By default, go-audiobook will detect and use GPU acceleration if available. You can override this by setting `tts.device` to `cpu`, `cuda`, or `mps` in your config.
For best performance, ensure you have the appropriate drivers and dependencies for your hardware.

#### Parallel Processing

Go's concurrency makes it easy to speed up audiobook generation by running multiple text-to-speech processes in parallel. The `tts.parallel_audio_count` setting controls how many TTS operations run at once. Raising this value reduces processing time but increases CPU usage, heat, and fan noise 🔥

##### Choosing the Right Parallel Audio Value
Set `tts.parallel_audio_count` to slightly below your machine's physical core count for the best results (unless you're also mining crypto, in which case... good luck).

**MacBook M1 Pro Example (with 10 cores):**
  - Minimal: `1–2`
  - Balanced: `3–5`
  - Maximum: `6-8`

##### Check Your CPU Core Count:
**macOS:**
```bash
sysctl -n hw.physicalcpu
```
**Linux:**
```bash
nproc --all
```
**Windows (Command Prompt or PowerShell):**
```powershell
WMIC CPU Get NumberOfCores
```

##### Why Not Exceed Core Count?
Setting `tts.parallel_audio_count` higher than your number of physical cores usually doesn't improve performance. It can make your system less responsive, increase heat and fan noise, and may cause CPU throttling.

> 💡 **Tip:** Start low (2–4), then increase if your system handles it well.

## ⚖️ Terms of Use & Disclaimer

By using this tool, you agree to the following:

- **Consent Required:** You must only use voice samples for which you have explicit permission from the person whose voice is featured. Do not use this tool to generate speech using the voice of any individual without their informed consent.
- **User Responsibility:** You are solely responsible for ensuring that your use of this tool and any voice samples you provide comply with all applicable laws and regulations, including those relating to copyright, privacy, and the right of publicity.
- **No Liability:** The developers and contributors of this project are not responsible for any misuse of this tool or any legal consequences arising from the unauthorised use of voice samples.

