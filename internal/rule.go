package internal

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"
)

type Rule struct {
	DetectionsRegex []regexp.Regexp
	Detections      []string `json:"detections"`
	Help            string   `json:"help"`
	Type            string   `json:"type"`
}

type Result struct {
	Matches []string `json:"matches"`
	Help    string   `json:"help"`
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
	matches := []string{}
	startIndex := -1
	lastIndex := -1
	lastRuleIndex := -1

	for ruleIndex, regex := range rule.DetectionsRegex {
		for i := lastIndex + 1; i < len(lines); i++ {
			line := lines[i]

			if regex.MatchString(line) {
				if startIndex < 0 {
					startIndex = i
				}
				lastIndex = i
				lastRuleIndex = ruleIndex
				break
			}
		}

		if lastRuleIndex != ruleIndex {
			return nil, false
		}
	}

	for i := startIndex; i <= lastIndex; i++ {
		matches = append(matches, lines[i])
	}

	return &Result{
		Matches: matches,
		Help:    rule.Help,
	}, true
}
