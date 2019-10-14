package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

type Rule struct {
	DetectionsRegex []regexp.Regexp
	Detections      []string `json:"detections"`
	Help            string   `json:"help"`
	Type            string   `json:"type"`
}

type Result struct {
	Matches []Match `json:"matches"`
	Help    string  `json:"help"`
}

type Match struct {
	StartIndex int    `json:"start"`
	EndIndex   int    `json:"end"`
	Message    string `json:"message"`
}

func LoadRules(filename string) ([]Rule, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	if err := json.Unmarshal(bytes, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

func (rule *Rule) IsError() bool {
	return rule.Type == "error";
}

func (rule *Rule) Build() {
	for _, detection := range rule.Detections {
		regex, err := regexp.Compile(detection)

		if err != nil {
			log.Fatal(err)
		}

		rule.DetectionsRegex = append(rule.DetectionsRegex, *regex)
	}
}

func (rule *Rule) Matches(lines []string) (*Result, bool) {
	matches := []Match{}
	pass := true
	index := 0

	for pass {
		match, _pass := rule.Match(lines, index)
		pass = _pass

		if _pass {
			index = match.EndIndex
			matches = append(matches, *match)
		}
	}

	if len(matches) < 1 {
		return nil, false
	}

	return &Result{
		Matches: matches,
		Help:    rule.Help,
	}, len(matches) > 0
}

func (rule *Rule) Match(lines []string, fromIndex int) (*Match, bool) {
	matchStartIndex := -1
	matchLastIndex := fromIndex - 1
	lastRuleIndex := -1

	if len(rule.DetectionsRegex) < 1 {
		return nil, false
	}

	if fromIndex >= len(lines) {
		return nil, false
	}

	for ruleIndex, regex := range rule.DetectionsRegex {
		for i := matchLastIndex + 1; i < len(lines); i++ {
			line := lines[i]

			if regex.MatchString(line) {
				if matchStartIndex < 0 {
					matchStartIndex = i
				}
				matchLastIndex = i
				lastRuleIndex = ruleIndex
				break
			}
		}

		if lastRuleIndex != ruleIndex {
			return nil, false
		}
	}

	matched := strings.Join(lines[matchStartIndex:(matchLastIndex+1)], "\n")

	return &Match{
		StartIndex: matchStartIndex,
		EndIndex:   matchLastIndex + 1,
		Message:    matched,
	}, true
}
