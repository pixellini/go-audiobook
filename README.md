# go-audiobook

A Go-based tool for converting text files into high-quality audiobooks using text-to-speech technology.

## 📚 Overview

This repository hosts a work-in-progress tool that transforms plain text and EPUB files into natural-sounding audiobooks using advanced TTS synthesis.

## 🚧 Current Status

⚠️ **Work in Progress** – This project is still under active development. Functionality is incomplete, and things may break or change frequently.

### 🛠️ Feature List

- 🗣️ Additional speaker options  
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

## 📄 License

This project is licensed under the MIT License – see the LICENSE file for details.
