[🇧🇷 Português](./README.pt-br.md)

# fflow - Fluent FFmpeg Command Builder for Go

[![Go Report Card](https://goreportcard.com/badge/github.com/Marlliton/ffcmd)](https://goreportcard.com/report/github.com/Marlliton/ffcmd)

`fflow` is a Go library that provides a fluent and semantic interface for building `ffmpeg` commands programmatically. Say goodbye to string concatenation and flag ordering errors.

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

### Example 5: Complex Filter Graph

A more advanced example that demonstrates a complex filter graph with multiple inputs, branches, and overlays.

```go
package main

import (
"fmt"

"github.com/Marlliton/fflow"
)

func main() {
cmd := fflow.New().
Override().
Input("input.mkv").
Input("train.jpg").
Filter().
Complex().
// Split ORIGINAL video into two branches
Chaing([]string{"0:v"}, fflow.AtomicFilter{Name: "split", Params: []string{"2"}}, []string{"v_main", "v_blur"}).
// Blur background branch
Chaing([]string{"v_blur"}, fflow.AtomicFilter{Name: "boxblur", Params: []string{"20:1"}}, []string{"v_bg"}).
// Scale ONLY the foreground
Chaing([]string{"v_main"}, fflow.AtomicFilter{Name: "scale", Params: []string{"960", "-1"}}, []string{"v_fg"}).
// Overlay foreground centered on blurred background
Chaing([]string{"v_bg", "v_fg"}, fflow.AtomicFilter{Name: "overlay", Params: []string{"(W-w)/2", "(H-h)/2"}}, []string{"v_base"}).
// Scale logo
Chaing([]string{"1:v"}, fflow.AtomicFilter{Name: "scale", Params: []string{"200", "-1"}}, []string{"logo"}).
// Overlay logo
Chaing([]string{"v_base", "logo"}, fflow.AtomicFilter{Name: "overlay", Params: []string{"W-w-20", "H-h-20"}}, []string{"v_logo"}).
// Draw text
Chaing([]string{"v_logo"}, fflow.AtomicFilter{Name: "drawtext", Params: []string{
"text=Builder Test",
"x=(w-text_w)/2",
"y=h-80",
"fontsize=42",
"fontcolor=white",
}}, []string{"v"}).
// Audio mix (fan-in)
Chaing([]string{"0:a", "0:a"}, fflow.AtomicFilter{Name: "amix", Params: []string{"inputs=2"}}, []string{"a"}).
Done(). // filter end
Map("v").
Map("a").
VideoCodec("libx264").
Preset("medium").
CRF(23).
AudioCodec("aac").
Raw("-b:a 192k").
Raw("-movflags +faststart").
Output("final_video.mp4").
Build()

fmt.Println(cmd)
}
```

## 📖 API Overview

The builder is divided into stages to ensure a logical and semantic command construction.

1. **`GlobalStage`**: Entry point (`New()`). Allows setting global options like `-y` (overwrite).
2. **`ReadStage`**: Defines the inputs (`Input()`) and their options, such as `-ss` (seek) or `-t` (duration).
3. **`FilterStage`**: Allows the creation of simple (`Simple()`) or complex (`Complex()`) filters.
4. **`WriteStage`**: Defines the output (`Output()`) and all its options, such as codecs (`-c:v`), presets (`-preset`), CRF, etc. It is the final stage before building the command with `Build()`).

## File Breakdown

*   **`ffmpeg.go`**: The entry point of the `fflow` package. It initializes the FFmpeg command builder and defines core types like `StreamType`.
*   **`global.go`**: Manages global FFmpeg command options such as `-y` (override output files without asking) and custom raw arguments before input specification.
*   **`read.go`**: Handles all input-related FFmpeg arguments, including adding input files (`-i`) and managing input seek and duration parameters (`-ss`, `-to`, `-t`).
*   **`filter.go`**: Contains the logic for building FFmpeg filter graphs, supporting both simple filters (like `-vf` and `-af`) and complex filter chains using `-filter_complex`.
*   **`write.go`**: Deals with output settings, including output file specification, video/audio/subtitle codecs, quality parameters (CRF), encoding presets, and stream mapping. This file also includes the `Build()` method, which constructs the final FFmpeg command string.
*   **`utils.go`**: Provides helper functions, such as `fmtDuration` for formatting `time.Duration` objects into FFmpeg-compatible time strings.

## Testing Files

*   **`ffmpeg_test.go`**: Contains tests for the basic initialization and fluent nature of the FFmpeg builder.
*   **`global_test.go`**: Tests the functionality of global FFmpeg options.
*   **`read_test.go`**: Tests the correct application of input-related options and handling of multiple inputs.
*   **`filter_test.go`**: Ensures the proper construction of filter strings for atomic filters, complex chains, and pipelines.
*   **`write_test.go`**: Verifies the correct generation of FFmpeg command arguments for output settings, codecs, quality, and complex filter integration.