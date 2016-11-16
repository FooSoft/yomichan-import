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
	"fmt"
	"io"
	"io/ioutil"
)

type epwingEntry struct {
	Heading string `json:"heading"`
	Text    string `json:"text"`
}

type epwingBook struct {
	Title     string        `json:"title"`
	Copyright string        `json:"copyright"`
	Entries   []epwingEntry `json:"entries"`
}

type epwingDict struct {
	CharacterCode string       `json:"characterCode"`
	DiscCode      string       `json:"discCode"`
	SubBooks      []epwingBook `json:"subBooks"`
}

func extractEpwingTerms(entry epwingEntry) []dbTerm {
	fmt.Print(entry.Heading)
	return nil
}

func extractEpwingKanji(entry epwingEntry) []dbKanji {
	return nil
}

func exportEpwingDb(outputDir, title string, reader io.Reader, flags int) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var dict epwingDict
	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}

	var terms dbTermList
	for _, subBook := range dict.SubBooks {
		for _, entry := range subBook.Entries {
			terms = append(terms, extractEpwingTerms(entry)...)
		}
	}

	var kanji dbKanjiList
	for _, subBook := range dict.SubBooks {
		for _, entry := range subBook.Entries {
			kanji = append(kanji, extractEpwingKanji(entry)...)
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
