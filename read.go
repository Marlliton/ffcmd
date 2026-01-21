package fflow

import "time"

type readStagee interface {
	// Ss adiciona a flag -ss após os inputs (-i), realizando um seek preciso no output.
	//
	// Ss adds the -ss flag after the inputs (-i), performing a precise seek on the output.
	Ss(time.Duration) readStagee

	// To adiciona a flag -to após os inputs (-i), definindo o tempo final absoluto do output.
	//
	// To adds the -to flag after the inputs (-i), defining the absolute end time of the output.
	To(time.Duration) readStagee

	// T adiciona a flag -t após os inputs (-i), limitando a duração do output.
	//
	// T adds the -t flag after the inputs (-i), limiting the output duration.
	T(time.Duration) readStagee

	// Input adiciona um arquivo de entrada (-i).
	//
	// Input adds an input file (-i).
	Input(path string) readStagee

	// Filter transiciona para a etapa de filtros da entrada atual.
	//
	// Filter transitions to the filter stage for the current input.
	Filter() filterStage

	// Output define o arquivo de saída e transiciona para o WriteStage.
	//
	// Output sets the output file and transitions to WriteStage.
	Output(path string) writeStage
}

type readCtx struct{ b *ffmpegBuilder }

func (c *readCtx) T(d time.Duration) readStagee {
	c.b.read = append(c.b.read, "-t", fmtDuration(d))
	return c
}

func (c *readCtx) Ss(d time.Duration) readStagee {
	c.b.read = append(c.b.read, "-ss", fmtDuration(d))
	return c
}

func (c *readCtx) To(d time.Duration) readStagee {
	c.b.read = append(c.b.read, "-to", fmtDuration(d))
	return c
}

func (c *readCtx) Input(path string) readStagee {
	c.b.read = append(c.b.read, "-i", path)
	return c
}

func (c *readCtx) Filter() filterStage {
	return &filterCtx{c.b}
}

func (c *readCtx) Output(path string) writeStage {
	write := &writeCtx{c.b}
	write.Output(path)
	return write
}
