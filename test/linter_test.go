package test

import (
	"github.com/mattak/loglint/internal"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type LinterTestContext struct{}

func (LinterTestContext) setup() {
	_ = os.RemoveAll("/tmp/test")
	_ = os.MkdirAll("/tmp/test", 0777)
}

func (LinterTestContext) tearDown() {
	os.RemoveAll("/tmp/test")
}

func TestSplitLines(t *testing.T) {
	context := LinterTestContext{}
	context.setup()
	defer context.tearDown()

	{
		lines := internal.SplitLines("hello\r\nworld\r\n")
		assert.Equal(t, lines, []string{"hello", "world", ""})
	}

	{
		lines := internal.SplitLines("hello\nworld\n")
		assert.Equal(t, lines, []string{"hello", "world", ""})
	}
}

func TestAnalyze(t *testing.T) {
	context := LinterTestContext{}
	context.setup()
	defer context.tearDown()

	{
		rule := internal.Rule{
			DetectionsRegex: nil,
			Detections:      []string{"1","3"},
			Help:            "",
			Type:            "error",
		}
		rule.Build()

		results := internal.Analyze([]string{"1", "2", "3"}, []internal.Rule{rule})

		assert.Equal(t, 1, len(results))
		assert.Equal(t, 1, len(results[0].Matches))
		assert.Equal(t, 0, results[0].Matches[0].StartIndex)
		assert.Equal(t, 3, results[0].Matches[0].EndIndex)
		assert.Equal(t, "1\n2\n3", results[0].Matches[0].Message)
	}
}
