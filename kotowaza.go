/*
 * Copyright (c) 2017 Alex Yatskov <alex@foosoft.net>, ajyliew
 * Author: ajyliew
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so,
 * subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"regexp"
	"strings"
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

func (e *kotowazaExtractor) extractTerms(entry epwingEntry) []dbTerm {
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
			}

			terms = append(terms, term)
		}

	}

	return terms
}

func (e *kotowazaExtractor) extractKanji(entry epwingEntry) []dbKanji {
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
