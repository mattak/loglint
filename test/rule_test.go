package test

import (
	"github.com/mattak/loglint/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

type RuleTestContext struct{}

func (RuleTestContext) setup() {
	_ = os.RemoveAll("/tmp/test")
	_ = os.MkdirAll("/tmp/test", 0777)
}

func (RuleTestContext) tearDown() {
	os.RemoveAll("/tmp/test")
}

func (RuleTestContext) writeFile(filepath string, json string) {
	ioutil.WriteFile(filepath, []byte(json), 0644)
}

func TestLoad(t *testing.T) {
	context := RuleTestContext{}
	context.setup()
	defer context.tearDown()

	{
		context.writeFile("/tmp/test/1.json", "[]")

		rules, err := internal.LoadRules("/tmp/test/1.json")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(rules))
	}

	{
		context.writeFile("/tmp/test/2.json", `[
	{
		"detections": ["sample"],
		"help": "see the docs"
	}
]`)
		rules, err := internal.LoadRules("/tmp/test/2.json")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(rules))
		assert.Equal(t, 1, len(rules[0].Detections))
		assert.Equal(t, "sample", rules[0].Detections[0])
		assert.Equal(t, "see the docs", rules[0].Help)
	}

	{
		context.writeFile("/tmp/test/3.json", `[
	{
		"detections": ["sample1"],
		"help": "help1"
	},
	{
		"detections": ["sample2"],
		"help": "help2"
	}
]`)
		rules, err := internal.LoadRules("/tmp/test/3.json")
		assert.Nil(t, err)
		assert.Equal(t, 2, len(rules))
	}

	// invalid json
	{
		context.writeFile("/tmp/test/4.json", `[ {} {} ]`)
		_, err := internal.LoadRules("/tmp/test/4.json")
		assert.NotNil(t, err)
	}

	// un-exist file
	{
		_, err := internal.LoadRules("/tmp/test/x.json")
		assert.NotNil(t, err)
	}
}

func TestMatches(t *testing.T) {
	context := RuleTestContext{}
	context.setup()
	defer context.tearDown()

	{
		context.writeFile("/tmp/test/1.json", `[
	{
		"detections": ["^Error: "],
		"help": "help1"
	}
]`)

		rules, _ := internal.LoadRules("/tmp/test/1.json")
		rules[0].Build()

		{
			matched, isMatched := rules[0].Matches(strings.Split(`
hello
world
`, "\n"))
			assert.False(t, isMatched)
			assert.Nil(t, matched)
		}

		{
			matched, isMatched := rules[0].Matches(strings.Split(`
hello
 Error: 
world
`, "\n"))
			assert.False(t, isMatched)
			assert.Nil(t, matched)
		}

		{
			matched, isMatched := rules[0].Matches(strings.Split(`
hello
Error: 
world
`, "\n"))
			assert.True(t, isMatched)
			assert.NotNil(t, matched)
			assert.Equal(t, 1, len(matched.Matches))
			assert.Equal(t, 2, matched.Matches[0].StartIndex)
			assert.Equal(t, 3, matched.Matches[0].EndIndex)
			assert.Equal(t, "Error: ", matched.Matches[0].Message)
			assert.Equal(t, "help1", matched.Help)
		}
	}

	// multi line matches
	{
		context.writeFile("/tmp/test/2.json", `[
	{
		"detections": ["^start", "^middle", "^end"],
		"help": "help2"
	}
]`)

		rules, _ := internal.LoadRules("/tmp/test/2.json")
		rule := rules[0]
		rule.Build()

		{
			matched, isMatched := rule.Matches(strings.Split(`
 start
middle
end
`, "\n"))
			assert.False(t, isMatched)
			assert.Nil(t, matched)
		}

		{
			matched, isMatched := rule.Matches(strings.Split(`
start
 middle
end
`, "\n"))
			assert.False(t, isMatched)
			assert.Nil(t, matched)
		}

		{
			matched, isMatched := rule.Matches(strings.Split(`
start
middle
 end
`, "\n"))
			assert.False(t, isMatched)
			assert.Nil(t, matched)
		}

		{
			matched, isMatched := rule.Matches(strings.Split(`
start
aaa
middle
bbb
end
`, "\n"))
			assert.True(t, isMatched)
			assert.NotNil(t, matched)
			assert.Equal(t, 1, len(matched.Matches))
			assert.Equal(t, 1, matched.Matches[0].StartIndex)
			assert.Equal(t, 6, matched.Matches[0].EndIndex)
			assert.Equal(t, "start\naaa\nmiddle\nbbb\nend", matched.Matches[0].Message)
			assert.Equal(t, "help2", matched.Help)
		}

		{
			matched, isMatched := rule.Matches(strings.Split(`
start
aaa
middle
bbb
end
start
middle
end
`, "\n"))
			assert.True(t, isMatched)
			assert.NotNil(t, matched)
			assert.Equal(t, 2, len(matched.Matches))
			assert.Equal(t, 1, matched.Matches[0].StartIndex)
			assert.Equal(t, 6, matched.Matches[0].EndIndex)
			assert.Equal(t, 6, matched.Matches[1].StartIndex)
			assert.Equal(t, 9, matched.Matches[1].EndIndex)
			assert.Equal(t, "start\naaa\nmiddle\nbbb\nend", matched.Matches[0].Message)
		assert.Equal(t, "start\nmiddle\nend", matched.Matches[1].Message)
	assert.Equal(t, "help2", matched.Help)
		}
	}
}
