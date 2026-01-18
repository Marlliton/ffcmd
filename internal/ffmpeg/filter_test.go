package ffmpeg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilterElements(t *testing.T) {
	scale := "scale"

	t.Run("AtomicFilter String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name         string
			filter       AtomicFilter
			expectedStr  string
			needsComplex bool
		}{
			{
				name:         "Scale filter with params",
				filter:       AtomicFilter{Name: scale, Params: []string{"1280", "-1"}},
				expectedStr:  "scale=1280:-1",
				needsComplex: false,
			},
			{
				name:         "Hflip filter without params",
				filter:       AtomicFilter{Name: "hflip", Params: []string{}},
				expectedStr:  "hflip",
				needsComplex: false,
			},
			{
				name:         "Empty filter name",
				filter:       AtomicFilter{Name: "", Params: []string{"param"}},
				expectedStr:  "=param",
				needsComplex: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expectedStr, tc.filter.String())
				assert.Equal(t, tc.needsComplex, tc.filter.NeedsComplex())
			})
		}
	})

	t.Run("Chaing String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name        string
			inputs      []string
			filter      AtomicFilter
			output      string
			expectedStr string
		}{
			{
				name:        "Single input, scale filter, output label",
				inputs:      []string{"0:v"},
				filter:      AtomicFilter{Name: scale, Params: []string{"1280", "-1"}},
				output:      "out",
				expectedStr: "[0:v]scale=1280:-1[out]",
			},
			{
				name:        "Single input, hflip filter, output label",
				inputs:      []string{"0:v"},
				filter:      AtomicFilter{Name: "hflip", Params: []string{}},
				output:      "out",
				expectedStr: "[0:v]hflip[out]",
			},
			{
				name:        "Multiple inputs, overlay filter, output label",
				inputs:      []string{"main", "logo"},
				filter:      AtomicFilter{Name: "overlay", Params: []string{"W-w-10:10"}},
				output:      "final_video",
				expectedStr: "[main][logo]overlay=W-w-10:10[final_video]",
			},
			{
				name:        "No output label",
				inputs:      []string{"0:v"},
				filter:      AtomicFilter{Name: scale, Params: []string{"640", "-1"}},
				output:      "",
				expectedStr: "[0:v]scale=640:-1",
			},
			{
				name:        "No input label, simple filter, output label (less common)",
				inputs:      []string{},
				filter:      AtomicFilter{Name: "null", Params: []string{}},
				output:      "null_out",
				expectedStr: "null[null_out]",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				c := Chaing{Inputs: tc.inputs, Filter: tc.filter, Output: tc.output}
				assert.Equal(t, tc.expectedStr, c.String())
				assert.True(t, c.NeedsComplex())
			})
		}
	})

	t.Run("Pipeline String and NeedsComplex", func(t *testing.T) {
		tests := []struct {
			name         string
			nodes        []Filter
			expectedStr  string
			needsComplex bool
		}{
			{
				name: "Pipeline of simple atomic filters (treated as separate chains)",
				nodes: []Filter{
					AtomicFilter{Name: scale, Params: []string{"1280", "-1"}},
					AtomicFilter{Name: "hflip", Params: []string{}},
				},
				expectedStr:  "scale=1280:-1;hflip", // INFO: Pipeline.String() always joins with ';'
				needsComplex: false,                 // Nodes themselves are not complex
			},
			{
				name: "Pipeline including a complex chain",
				nodes: []Filter{
					AtomicFilter{Name: "format", Params: []string{"yuv420p"}},
					Chaing{Inputs: []string{"0:v"}, Filter: AtomicFilter{Name: "fade", Params: []string{"in", "0", "30"}}, Output: "faded_video"},
					AtomicFilter{Name: "setsar", Params: []string{"1"}},
				},
				expectedStr:  "format=yuv420p;[0:v]fade=in:0:30[faded_video];setsar=1",
				needsComplex: true, // INFO: One node needs complex, so pipeline needs complex
			},
			{
				name: "Pipeline of only complex chains",
				nodes: []Filter{
					Chaing{Inputs: []string{"0:v"}, Filter: AtomicFilter{Name: scale, Params: []string{"640", "-1"}}, Output: "scaled"},
					Chaing{Inputs: []string{"scaled", "1:v"}, Filter: AtomicFilter{Name: "overlay", Params: []string{"W-w-10:10"}}, Output: "final"},
				},
				expectedStr:  "[0:v]scale=640:-1[scaled];[scaled][1:v]overlay=W-w-10:10[final]",
				needsComplex: true, // INFO: All nodes need complex
			},
			{
				name:         "Empty pipeline",
				nodes:        []Filter{},
				expectedStr:  "",
				needsComplex: false,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				p := Pipeline{Nodes: tc.nodes}
				assert.Equal(t, tc.expectedStr, p.String())
				assert.Equal(t, tc.needsComplex, p.NeedsComplex())
			})
		}
	})
}
