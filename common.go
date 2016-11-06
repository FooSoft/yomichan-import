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
	"os"
	"path"
	"strings"
)

type dbTerm struct {
	Expression string
	Reading    string
	Tags       []string
	Glossary   []string
}

type dbTermList []dbTerm

func (term *dbTerm) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, term.Tags) {
			term.Tags = append(term.Tags, tag)
		}
	}
}

func (term *dbTerm) addTagsPri(tags ...string) {
	for _, tag := range tags {
		switch tag {
		case "news1", "ichi1", "spec1", "gai1":
			term.addTags("P")
			fallthrough
		case "news2", "ichi2", "spec2", "gai2":
			term.addTags(tag[:len(tag)-1])
			break
		}
	}
}

func (terms dbTermList) crush() [][]string {
	var results [][]string
	for _, t := range terms {
		result := []string{
			t.Expression,
			t.Reading,
			strings.Join(t.Tags, " "),
		}

		result = append(result, t.Glossary...)
		results = append(results, result)
	}

	return results
}

type dbKanji struct {
	Character string
	Onyomi    []string
	Kunyomi   []string
	Tags      []string
	Meanings  []string
}

type dbKanjiList []dbKanji

func (kanji *dbKanji) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, kanji.Tags) {
			kanji.Tags = append(kanji.Tags, tag)
		}
	}
}

func (kanji dbKanjiList) crush() [][]string {
	var results [][]string
	for _, k := range kanji {
		result := []string{
			k.Character,
			strings.Join(k.Onyomi, " "),
			strings.Join(k.Kunyomi, " "),
			strings.Join(k.Tags, " "),
		}

		result = append(result, k.Meanings...)
		results = append(results, result)
	}

	return results
}

func writeDb(outputDir, title string, termRecords [][]string, kanjiRecords [][]string, entities map[string]string, pretty bool) error {
	const DB_VERSION = 1
	const BANK_STRIDE = 50000

	marshalJson := func(obj interface{}, pretty bool) ([]byte, error) {
		if pretty {
			return json.MarshalIndent(obj, "", "    ")
		}

		return json.Marshal(obj)
	}

	writeDbRecords := func(prefix string, records [][]string) (int, error) {
		recordCount := len(records)
		bankCount := 0

		for i := 0; i < recordCount; i += BANK_STRIDE {
			indexSrc := i
			indexDst := i + BANK_STRIDE
			if indexDst > recordCount {
				indexDst = recordCount
			}

			bytes, err := marshalJson(records[indexSrc:indexDst], pretty)
			if err != nil {
				return 0, err
			}

			fp, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s_bank_%d.json", prefix, i/BANK_STRIDE+1)))
			if err != nil {
				return 0, err
			}
			defer fp.Close()

			if _, err = fp.Write(bytes); err != nil {
				return 0, err
			}

			bankCount += 1
		}

		return bankCount, nil
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	var err error
	var db struct {
		Title      string            `json:"title"`
		Version    int               `json:"version"`
		Entities   map[string]string `json:"entities"`
		TermBanks  int               `json:"termBanks"`
		KanjiBanks int               `json:"kanjiBanks"`
	}

	db.Title = title
	db.Version = DB_VERSION
	db.Entities = entities

	if db.TermBanks, err = writeDbRecords("term", termRecords); err != nil {
		return err
	}

	if db.KanjiBanks, err = writeDbRecords("kanji", kanjiRecords); err != nil {
		return err
	}

	bytes, err := marshalJson(db, pretty)
	if err != nil {
		return err
	}

	fp, err := os.Create(path.Join(outputDir, "index.json"))
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err := fp.Write(bytes); err != nil {
		return err
	}

	return nil
}

func hasString(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}
