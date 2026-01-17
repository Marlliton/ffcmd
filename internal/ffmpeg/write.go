package ffmpeg

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type WriteStage interface {
	Ss(time.Duration) WriteStage
	To(time.Duration) WriteStage
	T(time.Duration) WriteStage
	VideoCodec(codec string) WriteStage
	AudioCodec(codec string) WriteStage
	SubtitleCodec(codec string) WriteStage
	CopyVideo() WriteStage
	CopyAudio() WriteStage
	CRF(value int) WriteStage

	Output(path string) WriteStage

	Build() string
}

type writeCtx struct{ b *ffmpegBuilder }

// T adiciona a flag -t após os inputs (-i), limitando a duração do output.
//
// T adds the -t flag after the inputs (-i), limiting the output duration.
func (c *writeCtx) T(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-t", fmtDuration(d))
	return c
}

// Ss adiciona a flag -ss após os inputs (-i), realizando um seek preciso no output.
//
// Ss adds the -ss flag after the inputs (-i), performing a precise seek on the output.
func (c *writeCtx) Ss(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-ss", fmtDuration(d))
	return c
}

// To adiciona a flag -to após os inputs (-i), definindo o tempo final absoluto do output.
//
// To adds the -to flag after the inputs (-i), defining the absolute end time of the output.
func (c *writeCtx) To(d time.Duration) WriteStage {
	c.b.write = append(c.b.write, "-to", fmtDuration(d))
	return c
}

// VideoCodec define o codec de vídeo do output (-c:v).
//
// VideoCodec sets the output video codec (-c:v).
func (c *writeCtx) VideoCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:v", codec)
	return c
}

// AudioCodec define o codec de áudio do output (-c:a).
//
// AudioCodec sets the output audio codec (-c:a).
func (c *writeCtx) AudioCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:a", codec)
	return c
}

// SubtitleCodec define o codec de legenda do output (-c:s).
//
// SubtitleCodec sets the output subtitle codec (-c:s).
func (c *writeCtx) SubtitleCodec(codec string) WriteStage {
	c.b.write = append(c.b.write, "-c:s", codec)
	return c
}

// CodecFor define o codec de um stream específico do output (-c:<stream>:<index>).
//
// CodecFor sets the codec for a specific output stream (-c:<stream>:<index>).
func (c *writeCtx) CodecFor(stream StreamType, index int, codec string) WriteStage {
	c.b.write = append(c.b.write, fmt.Sprintf("-c:%s:%d", stream, index), codec)
	return c
}

// CopyVideo copia o stream de vídeo sem recodificar (-c:v copy).
//
// CopyVideo copies the video stream without re-encoding (-c:v copy).
func (c *writeCtx) CopyVideo() WriteStage {
	c.b.write = append(c.b.write, "-c:v", "copy")
	return c
}

// CopyAudio copia o stream de áudio sem recodificar (-c:a copy).
//
// CopyAudio copies the audio stream without re-encoding (-c:a copy).
func (c *writeCtx) CopyAudio() WriteStage {
	c.b.write = append(c.b.write, "-c:a", "copy")
	return c
}

// CRF define o fator de qualidade constante para encoders de vídeo.
//
// CRF sets the constant quality factor for video encoders.
func (c *writeCtx) CRF(value int) WriteStage {
	c.b.write = append(c.b.write, "-crf", strconv.Itoa(value))
	return c
}

// Output define o arquivo de saída.
//
// Output sets the output file path.
func (c *writeCtx) Output(path string) WriteStage {
	c.b.output = path
	return c
}

// Build monta o comando FFmpeg final respeitando a ordem semântica.
//
// Build assembles the final FFmpeg command respecting semantic order.
func (c *writeCtx) Build() string {
	var args []string

	args = append(args, c.b.global...)
	args = append(args, c.b.read...)
	args = append(args, c.b.write...)
	args = append(args, c.b.output)
	return strings.Join(args, " ")
}
