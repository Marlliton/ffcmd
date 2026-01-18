package ffmpeg

type globalStage interface {
	// Override adiciona a flag global -y.
	//
	// Override adds the global -y flag.
	Override() globalStage

	// LogLevel(level string) GlobalStage

	// Raw adiciona um argumento bruto ao comando FFmpeg, antes do -i
	//
	// Raw adds a raw argument to the FFmpeg command, before -i flag
	Raw(value string) globalStage

	// Input adiciona um arquivo de entrada (-i) e transiciona para o ReadStage.
	//
	// Input adds an input file (-i) and transitions to ReadStage.
	Input(path string) ReadStage
}

type globalCtx struct{ b *ffmpegBuilder }

func (c *globalCtx) Input(path string) ReadStage {
	read := &readCtx{c.b}
	read.Input(path)
	return read
}

func (c *globalCtx) Raw(value string) globalStage {
	c.b.global = append(c.b.global, value)
	return c
}

func (c *globalCtx) Override() globalStage {
	c.b.global = append(c.b.global, "-y")
	return c
}
