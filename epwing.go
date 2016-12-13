/*
 * Copyright (c) 2016 Alex Yatskov <alex@foosoft.net>
 * Author: Alex Yatskov <alex@foosoft.net>
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
	"encoding/json"
	"io"
	"io/ioutil"
	"regexp"
	"strings"
)

type epwingEntry struct {
	Heading string `json:"heading"`
	Text    string `json:"text"`
}

type epwingSubbook struct {
	Title     string        `json:"title"`
	Copyright string        `json:"copyright"`
	Entries   []epwingEntry `json:"entries"`
}

type epwingBook struct {
	CharCode string          `json:"charCode"`
	DiscCode string          `json:"discCode"`
	Subbooks []epwingSubbook `json:"subbooks"`
}

type epwingExtractor interface {
	extractTerms(entry epwingEntry) []dbTerm
	extractKanji(entry epwingEntry) []dbKanji
}

type daijirinExtractor struct {
	partsExp   *regexp.Regexp
	phonExp    *regexp.Regexp
	variantExp *regexp.Regexp
	annotExp   *regexp.Regexp
}

func makeDaijirinExtractor() epwingExtractor {
	return &daijirinExtractor{
		partsExp:   regexp.MustCompile(`(?P<reading>[^（【〖]+)(?:【(?P<expression>.*)】)?(?:〖(?P<native>.*)〗)?(?:（(?P<tag>.*)）)?`),
		phonExp:    regexp.MustCompile(`[-・]+`),
		variantExp: regexp.MustCompile(`\((.*)\)`),
		annotExp:   regexp.MustCompile(`（(.*)）`),
	}
}

func (e *daijirinExtractor) extractTerms(entry epwingEntry) []dbTerm {
	var expressions, readings, glossary, tags []string

	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	for i, name := range e.partsExp.SubexpNames() {
		value := matches[i]
		if i == 0 || len(value) == 0 {
			continue
		}

		switch name {
		case "expression":
			expression := e.annotExp.ReplaceAllLiteralString(value, "")
			for _, split := range strings.Split(expression, `・`) {
				splitInc := e.variantExp.ReplaceAllString(split, "$1")
				expressions = append(expressions, splitInc)
				if split != splitInc {
					splitExc := e.variantExp.ReplaceAllLiteralString(split, "")
					expressions = append(expressions, splitExc)
				}
			}
		case "reading":
			reading := e.phonExp.ReplaceAllLiteralString(value, "")
			readings = append(readings, reading)
		}
	}

	for i, split := range strings.Split(entry.Text, "\n") {
		if i == 0 {
			matches := e.annotExp.FindStringSubmatch(split)
			if len(matches) >= 1 {
				tags = append(tags, strings.Split(matches[1], `・`)...)
			}
		}

		glossary = append(glossary, split)
	}

	return nil
}

func (e *daijirinExtractor) extractKanji(entry epwingEntry) []dbKanji {
	return nil
}

func exportEpwingDb(outputDir, title string, reader io.Reader, flags int) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var book epwingBook
	if err := json.Unmarshal(data, &book); err != nil {
		return err
	}

	epwingExtractors := map[string]epwingExtractor{
		"三省堂　スーパー大辞林": makeDaijirinExtractor(),
	}

	var terms dbTermList
	var kanji dbKanjiList

	for _, subbook := range book.Subbooks {
		if extractor, ok := epwingExtractors[subbook.Title]; ok {
			for _, entry := range subbook.Entries {
				terms = append(terms, extractor.extractTerms(entry)...)
				kanji = append(kanji, extractor.extractKanji(entry)...)
			}
		}
	}

	return writeDb(
		outputDir,
		title,
		terms.crush(),
		kanji.crush(),
		nil,
		flags&flagPretty == flagPretty,
	)
}
