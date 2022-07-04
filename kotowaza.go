package yomichan

import (
	"regexp"
	"strings"

	zig "foosoft.net/projects/zero-epwing-go"
)

type kotowazaExtractor struct {
	readGroupExp       *regexp.Regexp
	readGroupAltsExp   *regexp.Regexp
	readGroupNoAltsExp *regexp.Regexp
	wordGroupExp       *regexp.Regexp
}

func makeKotowazaExtractor() epwingExtractor {
	return &kotowazaExtractor{
		readGroupExp:       regexp.MustCompile(`([^ぁ-ゖァ-ヺ]*)(\([^)]*\))`),
		readGroupAltsExp:   regexp.MustCompile(`\(([^)]*)\)`),
		readGroupNoAltsExp: regexp.MustCompile(`\(([^・)]*)\)`),
		wordGroupExp:       regexp.MustCompile(`＝([^〔＝]*)〔＝([^〕]*)〕`),
	}
}

func (e *kotowazaExtractor) extractTerms(entry zig.BookEntry, sequence int) []dbTerm {
	heading := entry.Heading

	queue := []string{heading}
	reducedExpressions := []string{}

	for len(queue) > 0 {
		expression := queue[0]
		queue = queue[1:]

		matches := e.wordGroupExp.FindStringSubmatch(expression)
		if matches == nil {
			reducedExpressions = append(reducedExpressions, expression)
		} else {
			replacements := []string{matches[1]}
			replacements = append(replacements, strings.Split(matches[2], "・")...)
			for _, replacement := range replacements {
				queue = append(queue, strings.Replace(expression, matches[0], replacement, -1))
			}
		}
	}

	var terms []dbTerm
	for _, reducedExpression := range reducedExpressions {
		expression := e.readGroupExp.ReplaceAllString(reducedExpression, "$1")
		readAltsExpression := e.readGroupExp.ReplaceAllString(reducedExpression, "$2")
		readAltsExpression = e.readGroupNoAltsExp.ReplaceAllString(readAltsExpression, "$1")

		var readings []string
		queue = []string{readAltsExpression}
		for len(queue) > 0 {
			readExpression := queue[0]
			queue = queue[1:]

			matches := e.readGroupAltsExp.FindStringSubmatch(readExpression)
			if matches == nil {
				readings = append(readings, readExpression)
			} else {
				replacements := strings.Split(matches[1], "・")
				for _, replacement := range replacements {
					queue = append(queue, strings.Replace(readExpression, matches[0], replacement, -1))
				}
			}
		}

		for _, reading := range readings {
			term := dbTerm{
				Expression: expression,
				Reading:    reading,
				Glossary:   []string{entry.Text},
				Sequence:   sequence,
			}

			terms = append(terms, term)
		}

	}

	return terms
}

func (e *kotowazaExtractor) extractKanji(entry zig.BookEntry) []dbKanji {
	return nil
}

func (e *kotowazaExtractor) exportRules(term *dbTerm, tags []string) {
}

func (*kotowazaExtractor) getRevision() string {
	return "kotowaza1"
}

func (*kotowazaExtractor) getFontNarrow() map[int]string {
	return map[int]string{}
}

func (*kotowazaExtractor) getFontWide() map[int]string {
	return map[int]string{}
}
