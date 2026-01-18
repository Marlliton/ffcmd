package ffmpeg

import (
	"fmt"
	"strings"
)

type Filter interface {
	String() string
	NeedsComplex() bool
}

type FilterStage interface{}

type filterCtx struct{ b *ffmpegBuilder }

type AtomicFilter struct {
	Name   string
	Params []string
}

func (f AtomicFilter) String() string {
	if len(f.Params) == 0 {
		return f.Name
	}
	return fmt.Sprintf("%s=%s", f.Name, strings.Join(f.Params, ":"))
}

func (f AtomicFilter) NeedsComplex() bool {
	return false
}

type Chaing struct {
	Inputs []string
	Filter AtomicFilter
	Output string
}

func (c Chaing) String() string {
	var sb strings.Builder

	for _, in := range c.Inputs {
		sb.WriteString("[")
		sb.WriteString(in)
		sb.WriteString("]")
	}

	sb.WriteString(c.Filter.String())

	if c.Output != "" {
		sb.WriteString("[")
		sb.WriteString(c.Output)
		sb.WriteString("]")
	}

	return sb.String()
}

func (c Chaing) NeedsComplex() bool {
	return true
}

type Pipeline struct {
	Nodes []Filter
}

func (p Pipeline) String() string {
	var parts []string
	for _, n := range p.Nodes {
		parts = append(parts, n.String())
	}
	return strings.Join(parts, ";")
}

func (p Pipeline) NeedsComplex() bool {
	for _, n := range p.Nodes {
		if n.NeedsComplex() {
			return true
		}
	}
	return false
}
