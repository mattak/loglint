package internal

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
		log.Fatal(err)
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
	lines := strings.Split(content, "\n")
	result := LintResult{}
	result.Errors = analyze(lines, linter.Errors)
	result.Warnings = analyze(lines, linter.Warnings)
	result.Passed = len(result.Errors) < 1
	return result
}

func analyze(lines []string, rules []Rule) []Result {
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
		log.Fatal(err)
	}

	fmt.Print(string(data))
}
