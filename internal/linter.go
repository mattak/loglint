package internal

import (
	"encoding/json"
	"fmt"
	"log"
	regexp2 "regexp"
)

type Linter struct {
	Errors   []Rule `json:"errors"`
	Warnings []Rule `json:"warnings"`
}

type LintResult struct {
	Passed   bool     `json:"passed"`
	Errors   []Result `json:"errors"`
	Warnings []Result `json:"warnings"`
}

func (linter *Linter) Prepare(filepath string) {
	rules, err := LoadRules(filepath)

	if err != nil {
		log.Fatal("Loading rules error. ", err)
	}

	var errorRules []Rule
	var warningRules []Rule

	for i := 0; i < len(rules); i++ {
		rules[i].Build()

		if rules[i].IsError() {
			errorRules = append(errorRules, rules[i])
		} else {
			warningRules = append(warningRules, rules[i])
		}
	}

	linter.Errors = errorRules
	linter.Warnings = warningRules
}

func (linter *Linter) Run(content string) LintResult {
	lines := SplitLines(content)
	result := LintResult{}
	result.Errors = Analyze(lines, linter.Errors)
	result.Warnings = Analyze(lines, linter.Warnings)
	result.Passed = len(result.Errors) < 1
	return result
}

func SplitLines(content string) []string {
	regexp := regexp2.MustCompile("\\r?\\n")
	return regexp.Split(content, -1)
}

func Analyze(lines []string, rules []Rule) []Result {
	var results []Result

	for _, rule := range rules {
		result, matched := rule.Matches(lines)

		if matched {
			results = append(results, *result)
		}
	}

	return results
}

func (result *LintResult) Print() {
	data, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Marshal result json error. ", err)
	}

	fmt.Print(string(data))
}
