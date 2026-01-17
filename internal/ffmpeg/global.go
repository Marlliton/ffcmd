package ffmpeg

type GlobalStage interface {
	Override() GlobalStage
	// LogLevel(level string) GlobalStage

	Input(path string) ReadStage
}

type globalCtx struct{ b *ffmpegBuilder }

// Input adiciona um arquivo de entrada (-i) e transiciona para o ReadStage.
//
// Input adds an input file (-i) and transitions to ReadStage.
func (c *globalCtx) Input(path string) ReadStage {
	read := &readCtx{c.b}
	read.Input(path)
	return read
}

// Override adiciona a flag global -y.
//
// Override adds the global -y flag.
func (c *globalCtx) Override() GlobalStage {
	c.b.global = append(c.b.global, "-y")
	return c
}
