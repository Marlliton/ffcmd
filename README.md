[🇧🇷 Português](./README.pt-br.md)

# ffcmd - Fluent FFmpeg Command Builder for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/Marlliton/ffcmd)](https://goreportcard.com/report/github.com/Marlliton/ffcmd)

`ffcmd` is a Go library that provides a fluent and semantic interface for building `ffmpeg` commands programmatically. Say goodbye to string concatenation and flag ordering errors.

## ✨ Features

- **Fluent API**: Build complex commands by chaining methods in a readable way.
- **Semantic Order**: The library ensures the correct order of FFmpeg flags (global, input, and output options).
- **Simple and Complex Filters**: Native support for `-vf`, `-af`, and `-filter_complex` in an organized manner.
- **Type-Safe**: Avoid common errors by specifying whether a simple filter is for **video** or **audio**.
- **Clarity**: Clear separation between configuration stages (Global, Read, Filter, Write).

## 📦 Installation

```bash
go get github.com/Marlliton/ffcmd/ffmpeg
```

## 🚀 Usage and Examples

Using the library follows the logic of building an `ffmpeg` command: first global options, then inputs, filters, and finally the output and its options.

### Example 1: Basic Conversion

Convert a video file to a different format.

```go
package main

import (
 "fmt"
 "github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
 cmd := ffmpeg.New().
  Override(). // Adds the global -y flag to overwrite the output file
  Input("input.mp4").
  Output("output.webm").
  VideoCodec("libvpx-vp9").
  AudioCodec("libopus").
  Build()

 fmt.Println(cmd)
 // Output: ffmpeg -y -i input.mp4 -c:v libvpx-vp9 -c:a libopus output.webm
}
```

### Example 2: Trimming a Video

You can use `Ss` (seek) and `T` (duration) for both input (for a fast seek) and output (for a precise cut).

```go
package main

import (
 "fmt"
 "time"
 "github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
    // Fast seek on input and limited duration
 cmd := ffmpeg.New().
  Ss(1 * time.Minute). // Jumps to 1 minute from the start of the input file
  T(30 * time.Second).  // Reads only 30 seconds of the input
  Input("input.mp4").
  Output("output.mp4").
  CopyVideo(). // Copies the video stream without re-encoding
  CopyAudio().  // Copies the audio stream without re-encoding
  Build()

 fmt.Println(cmd)
 // Output: ffmpeg -ss 00:01:00.000 -t 00:00:30.000 -i input.mp4 -c:v copy -c:a copy output.mp4
}
```

### Example 3: Simple Filters (Video and Audio)

Apply filters to a single video (`-vf`) or audio (`-af`) stream.

```go
package main

import (
 "fmt"
 "github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
    // Video filter to resize and flip horizontally
 videoCmd := ffmpeg.New().
  Input("input.mp4").
  Filter().
  Simple(ffmpeg.FilterVideo). // Specifies it's a video filter (-vf)
  Add(ffmpeg.AtomicFilter{Name: "scale", Params: []string{"1280", "-1"}}).
  Add(ffmpeg.AtomicFilter{Name: "hflip"}).
  Done().
  Output("video_filtered.mp4").
  Build()

 fmt.Println(videoCmd)
 // Output: ffmpeg -i input.mp4 -vf scale=1280:-1,hflip video_filtered.mp4

    // Audio filter to adjust the volume
 audioCmd := ffmpeg.New().
  Input("input.mp4").
  Filter().
  Simple(ffmpeg.FilterAudio). // Specifies it's an audio filter (-af)
  Add(ffmpeg.AtomicFilter{Name: "volume", Params: []string{"0.5"}}).
  Done().
  Output("audio_filtered.mp3").
  Build()

 fmt.Println(audioCmd)
    // Output: ffmpeg -i input.mp4 -af volume=0.5 audio_filtered.mp3
}
```

### Example 4: Complex Filter (Real-world Scenario)

A more advanced example: trimming a video, overlaying a watermark, speeding up the audio, and re-encoding with specific presets.

```go
package main

import (
 "fmt"
 "time"
 "github.com/Marlliton/ffcmd/ffmpeg"
)

func main() {
 cmd := ffmpeg.New().
  Override().
  Input("input.mp4"). // Main video input
  Input("logo.png").   // Image input for the watermark
  Filter().
  Complex(). // Starts a -filter_complex
  Chaing(
   []string{"0:v"}, // Takes the video from the first input
   ffmpeg.AtomicFilter{Name: "scale", Params: []string{"1920", "-1"}},
   "scaled", // Names the output for later use
  ).
  Chaing(
   []string{"scaled", "1:v"}, // Takes the resized video and the image from the second input
   ffmpeg.AtomicFilter{Name: "overlay", Params: []string{"W-w-10", "10"}},
   "video_out", // Names the final video output
  ).
  Chaing(
   []string{"0:a"}, // Takes the audio from the first input
   ffmpeg.AtomicFilter{Name: "atempo", Params: []string{"1.5"}},
   "audio_out", // Names the final audio output
  ).
  Done().
  Map("video_out"). // Maps the video output from the complex filter
  Map("audio_out"). // Maps the audio output from the complex filter
  VideoCodec("libx264").
  AudioCodec("aac").
  Preset("fast").
  CRF(23).
  Output("final_video.mp4").
  Build()

 fmt.Println(cmd)
 /*
    Output: ffmpeg -y -i input.mp4 -i logo.png -filter_complex [0:v]scale=1920:-1[scaled];[scaled][1:v]overlay=W-w-10:10[video_out];[0:a]atempo=1.5[audio_out] -map [video_out] -map [audio_out] -c:v libx264 -c:a aac -preset fast -crf 23 final_video.mp4
 */
}
```

## 📖 API Overview

The builder is divided into stages to ensure a logical and semantic command construction.

1. **`GlobalStage`**: Entry point (`New()`). Allows setting global options like `-y` (overwrite).
2. **`ReadStage`**: Defines the inputs (`Input()`) and their options, such as `-ss` (seek) or `-t` (duration).
3. **`FilterStage`**: Allows the creation of simple (`Simple()`) or complex (`Complex()`) filters.
4. **`WriteStage`**: Defines the output (`Output()`) and all its options, such as codecs (`-c:v`), presets (`-preset`), CRF, etc. It is the final stage before building the command with `Build()`.
