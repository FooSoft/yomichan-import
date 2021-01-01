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
	"archive/zip"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const databaseFormat = 3

type dbRecord []interface{}
type dbRecordList []dbRecord

type dbTag struct {
	Name     string
	Category string
	Order    int
	Notes    string
	Score    int
}

type dbTagList []dbTag

func (meta dbTagList) crush() dbRecordList {
	var results dbRecordList
	for _, m := range meta {
		results = append(results, dbRecord{m.Name, m.Category, m.Order, m.Notes, m.Score})
	}

	return results
}

type dbMeta struct {
	Expression string
	Mode       string
	Data       interface{}
}

type dbMetaList []dbMeta

func (freqs dbMetaList) crush() dbRecordList {
	var results dbRecordList
	for _, f := range freqs {
		results = append(results, dbRecord{f.Expression, f.Mode, f.Data})
	}

	return results
}

type dbTerm struct {
	Expression     string
	Reading        string
	DefinitionTags []string
	Rules          []string
	Score          int
	Glossary       []string
	Sequence       int
	TermTags       []string
}

type dbTermList []dbTerm

func (term *dbTerm) addDefinitionTags(tags ...string) {
	term.DefinitionTags = appendStringUnique(term.DefinitionTags, tags...)
}

func (term *dbTerm) addTermTags(tags ...string) {
	term.TermTags = appendStringUnique(term.TermTags, tags...)
}

func (term *dbTerm) addRules(rules ...string) {
	term.Rules = appendStringUnique(term.Rules, rules...)
}

func (terms dbTermList) crush() dbRecordList {
	var results dbRecordList
	for _, t := range terms {
		result := dbRecord{
			t.Expression,
			t.Reading,
			strings.Join(t.DefinitionTags, " "),
			strings.Join(t.Rules, " "),
			t.Score,
			t.Glossary,
			t.Sequence,
			strings.Join(t.TermTags, " "),
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
	Stats     map[string]string
}

type dbKanjiList []dbKanji

func (kanji *dbKanji) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, kanji.Tags) {
			kanji.Tags = append(kanji.Tags, tag)
		}
	}
}

func (kanji dbKanjiList) crush() dbRecordList {
	var results dbRecordList
	for _, k := range kanji {
		result := dbRecord{
			k.Character,
			strings.Join(k.Onyomi, " "),
			strings.Join(k.Kunyomi, " "),
			strings.Join(k.Tags, " "),
			k.Meanings,
			k.Stats,
		}

		results = append(results, result)
	}

	return results
}

func writeDb(outputPath, title, revision string, sequenced bool, recordData map[string]dbRecordList, stride int, pretty bool) error {
	var zbuff bytes.Buffer
	zip := zip.NewWriter(&zbuff)

	marshalJSON := func(obj interface{}, pretty bool) ([]byte, error) {
		if pretty {
			return json.MarshalIndent(obj, "", "    ")
		}

		return json.Marshal(obj)
	}

	writeDbRecords := func(prefix string, records dbRecordList) (int, error) {
		recordCount := len(records)
		bankCount := 0

		for i := 0; i < recordCount; i += stride {
			indexSrc := i
			indexDst := i + stride
			if indexDst > recordCount {
				indexDst = recordCount
			}

			bytes, err := marshalJSON(records[indexSrc:indexDst], pretty)
			if err != nil {
				return 0, err
			}

			zw, err := zip.Create(fmt.Sprintf("%s_bank_%d.json", prefix, i/stride+1))
			if err != nil {
				return 0, err
			}

			if _, err := zw.Write(bytes); err != nil {
				return 0, err
			}

			bankCount++
		}

		return bankCount, nil
	}

	var err error
	var db struct {
		Title     string `json:"title"`
		Format    int    `json:"format"`
		Revision  string `json:"revision"`
		Sequenced bool   `json:"sequenced"`
	}

	db.Title = title
	db.Format = databaseFormat
	db.Revision = revision
	db.Sequenced = sequenced

	for recordType, recordEntries := range recordData {
		if _, err := writeDbRecords(recordType, recordEntries); err != nil {
			return err
		}
	}

	bytes, err := marshalJSON(db, pretty)
	if err != nil {
		return err
	}

	zw, err := zip.Create("index.json")
	if err != nil {
		return err
	}

	if _, err := zw.Write(bytes); err != nil {
		return err
	}

	zip.Close()

	fp, err := os.Create(outputPath)
	if err != nil {
		return err
	}

	if _, err := fp.Write(zbuff.Bytes()); err != nil {
		return err
	}

	return fp.Close()
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

func detectFormat(path *string) (string, error) {
	switch filepath.Ext(*path) {
	case ".sqlite":
		return "rikai", nil
	case ".kanjifreq":
		return "kanjifreq", nil
	case ".termfreq":
		return "termfreq", nil
	}

	switch filepath.Base(*path) {
	case "JMdict", "JMdict.xml", "JMdict_e", "JMdict_e.xml":
		return "edict", nil
	case "JMnedict", "JMnedict.xml":
		return "enamdict", nil
	case "kanjidic2", "kanjidic2.xml":
		return "kanjidic", nil
	case "CATALOGS":
		*path = filepath.Dir(*path)
		return "epwing", nil
	}

	info, err := os.Stat(*path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		_, err := os.Stat(filepath.Join(*path, "CATALOGS"))
		if err == nil {
			return "epwing", nil
		}
	}

	return "", errors.New("unrecognized dictionary format")
}
