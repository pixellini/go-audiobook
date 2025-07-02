# go-audiobook

A Go-based tool for converting text files into high-quality audiobooks using text-to-speech technology.

## 📚 Overview

This repository hosts a work-in-progress tool that transforms plain text and EPUB files into natural-sounding audiobooks using advanced TTS synthesis.

## 🚧 Current Status

⚠️ **Work in Progress** – This project is still under active development. Functionality is incomplete, and things may break or change frequently.

### 🛠️ Feature List

- 🗣️ Default speaker option  
- 📄 Support for reading plain .txt files  
- ⚙️ Concurrent TTS file generation for faster processing  
- 🚀 Enhanced Coqui integration, including GPU acceleration if available  
- 🧠 Coqui VITS support for quicker audiobook creation  
- 🔁 Automatic retries for failed paragraph TTS synthesis  
- 📦 Containerise inside Docker  
- 🎧 MP3 output and conversion support  
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

### 🎤 Providing a Narrator Voice (.wav)

You can customise your audiobook by providing a short audio sample of your chosen narrator's voice. This sample will be used as the voice for the entire audiobook, allowing you to generate content in any of the available languages. Coqui TTS synthesises speech using the unique vocal characteristics from the sample, not the spoken content.

- **Recommended length:** 1–3 minutes of clear, uninterrupted speech is usually enough for high-quality results. Samples longer than 3 minutes may increase processing time. If you have a powerful machine, you can experiment with longer inputs, but 3 minutes is sufficient for most cases. 
- **File format:** A `.wav` file is preferred for best audio quality. `.mp3` is also supported, and while it can speed up processing time, it may reduce quality.
- **Setup:** Place your audio file in the `speakers/` directory and update the `speaker_wav` path in your `config.json` accordingly.  

> **Note:** The spoken language in the sample doesn't affect the output. Coqui will synthesise your audiobook in the target language using the voice's characteristics from the sample.

## ⚖️ Terms of Use & Disclaimer

By using this tool, you agree to the following:

- **Consent Required:** You must only use voice samples for which you have explicit permission from the person whose voice is featured. Do not use this tool to generate speech using the voice of any individual without their informed consent.
- **User Responsibility:** You are solely responsible for ensuring that your use of this tool and any voice samples you provide comply with all applicable laws and regulations, including those relating to copyright, privacy, and the right of publicity.
- **No Liability:** The developers and contributors of this project are not responsible for any misuse of this tool or any legal consequences arising from the unauthorised use of voice samples.

