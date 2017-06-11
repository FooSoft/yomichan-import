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
	"path/filepath"
	"strings"
)

type dbTagMeta struct {
	Category string `json:"category,omitempty"`
	Notes    string `json:"notes,omitempty"`
	Order    int    `json:"order,omitempty"`
}

type dbTerm struct {
	Expression string
	Reading    string
	Tags       []string
	Rules      []string
	Score      int
	Glossary   []string
}

type dbTermList []dbTerm

func (term *dbTerm) addTags(tags ...string) {
	term.Tags = appendStringUnique(term.Tags, tags...)
}

func (term *dbTerm) addRules(rules ...string) {
	term.Rules = appendStringUnique(term.Rules, rules...)
}

func (terms dbTermList) crush() [][]interface{} {
	var results [][]interface{}
	for _, t := range terms {
		result := []interface{}{
			t.Expression,
			t.Reading,
			strings.Join(t.Tags, " "),
			strings.Join(t.Rules, " "),
			t.Score,
		}

		for _, gloss := range t.Glossary {
			result = append(result, gloss)
		}

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

func (kanji dbKanjiList) crush() [][]interface{} {
	var results [][]interface{}
	for _, k := range kanji {
		result := []interface{}{
			k.Character,
			strings.Join(k.Onyomi, " "),
			strings.Join(k.Kunyomi, " "),
			strings.Join(k.Tags, " "),
		}

		for _, meaning := range k.Meanings {
			result = append(result, meaning)
		}

		results = append(results, result)
	}

	return results
}

func writeDb(outputDir, title, revision string, termRecords [][]interface{}, kanjiRecords [][]interface{}, tagMeta map[string]dbTagMeta, stride int, pretty bool) error {
	const DB_VERSION = 1

	marshalJson := func(obj interface{}, pretty bool) ([]byte, error) {
		if pretty {
			return json.MarshalIndent(obj, "", "    ")
		}

		return json.Marshal(obj)
	}

	writeDbRecords := func(prefix string, records [][]interface{}) (int, error) {
		recordCount := len(records)
		bankCount := 0

		for i := 0; i < recordCount; i += stride {
			indexSrc := i
			indexDst := i + stride
			if indexDst > recordCount {
				indexDst = recordCount
			}

			bytes, err := marshalJson(records[indexSrc:indexDst], pretty)
			if err != nil {
				return 0, err
			}

			fp, err := os.Create(path.Join(outputDir, fmt.Sprintf("%s_bank_%d.json", prefix, i/stride+1)))
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
		Title      string               `json:"title"`
		Version    int                  `json:"version"`
		Revision   string               `json:"revision"`
		TagMeta    map[string]dbTagMeta `json:"tagMeta"`
		TermBanks  int                  `json:"termBanks"`
		KanjiBanks int                  `json:"kanjiBanks"`
	}

	db.Title = title
	db.Version = DB_VERSION
	db.Revision = revision
	db.TagMeta = tagMeta

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

func appendStringUnique(target []string, source ...string) []string {
	for _, str := range source {
		if !hasString(str, target) {
			target = append(target, str)
		}
	}

	return target
}

func hasString(needle string, haystack []string) bool {
	for _, value := range haystack {
		if needle == value {
			return true
		}
	}

	return false
}

func detectFormat(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return ""
	}

	if info.IsDir() {
		_, err := os.Stat(filepath.Join(path, "CATALOGS"))
		if err == nil {
			return "epwing"
		}
	} else {
		base := filepath.Base(path)
		switch base {
		case "JMdict":
		case "JMdict.xml":
		case "JMdict_e":
		case "JMdict_e.xml":
			return "edict"
		case "JMnedict":
		case "JMnedict.xml":
			return "enamdict"
		case "kanjidic2":
		case "kanjidic2.xml":
			return "kanjidic"
		}
	}

	return ""
}
