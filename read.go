package fflow

import "time"

type readStagee interface {
	// Ss adiciona a flag -ss antes do -i, realizando um seek rápido na entrada.
	//
	// Ss adds the -ss flag before -i, performing a fast seek on the input.
	Ss(d time.Duration) readStagee

	// To adiciona a flag -to antes do -i, definindo o tempo final absoluto da leitura.
	//
	// To adds the -to flag before -i, defining the absolute end time of the input read.
	To(d time.Duration) readStagee

	// T adiciona a flag -t antes do -i, limitando quanto da entrada será lida.
	//
	// T adds the -t flag before -i, limiting how much of the input is read.
	T(d time.Duration) readStagee

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
	c.b.read = append([]string{"-t", fmtDuration(d)}, c.b.read...)
	return c
}

func (c *readCtx) Ss(d time.Duration) readStagee {
	c.b.read = append([]string{"-ss", fmtDuration(d)}, c.b.read...)
	return c
}

func (c *readCtx) To(d time.Duration) readStagee {
	c.b.read = append([]string{"-to", fmtDuration(d)}, c.b.read...)
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
