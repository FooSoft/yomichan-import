package yomichan

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

const (
	DefaultFormat   = ""
	DefaultLanguage = ""
	DefaultPretty   = false
	DefaultStride   = 10000
	DefaultTitle    = ""
)

type dbRecord []any
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
	Data       any
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
	Glossary       []any
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

type dbIndex struct {
	Title       string `json:"title"`
	Format      int    `json:"format"`
	Revision    string `json:"revision"`
	Sequenced   bool   `json:"sequenced"`
	Author      string `json:"author"`
	Url         string `json:"url"`
	Description string `json:"description"`
	Attribution string `json:"attribution"`
}

func (index *dbIndex) setDefaults() {
	if index.Format == 0 {
		index.Format = 3
	}
	if index.Author == "" {
		index.Author = "yomichan-import"
	}
	if index.Url == "" {
		index.Url = "https://github.com/FooSoft/yomichan-import"
	}
}

func writeDb(outputPath string, index dbIndex, recordData map[string]dbRecordList, stride int, pretty bool) error {
	var zbuff bytes.Buffer
	zip := zip.NewWriter(&zbuff)

	marshalJSON := func(obj any, pretty bool) ([]byte, error) {
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

	for recordType, recordEntries := range recordData {
		if _, err := writeDbRecords(recordType, recordEntries); err != nil {
			return err
		}
	}

	index.setDefaults()
	bytes, err := marshalJSON(index, pretty)
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

func intersection(s1, s2 []string) []string {
	s := []string{}
	m := make(map[string]bool)
	for _, e := range s1 {
		m[e] = true
	}
	for _, e := range s2 {
		if m[e] {
			s = append(s, e)
			m[e] = false
		}
	}
	return s
}

func union(s1, s2 []string) []string {
	s := []string{}
	m := make(map[string]bool)
	for _, e := range s1 {
		if !m[e] {
			s = append(s, e)
			m[e] = true
		}
	}
	for _, e := range s2 {
		if !m[e] {
			s = append(s, e)
			m[e] = true
		}
	}
	return s
}

func detectFormat(path string) (string, error) {
	switch filepath.Ext(path) {
	case ".sqlite":
		return "rikai", nil
	case ".kanjifreq":
		return "kanjifreq", nil
	case ".termfreq":
		return "termfreq", nil
	}

	switch filepath.Base(path) {
	case "JMdict", "JMdict.xml", "JMdict_e", "JMdict_e.xml", "JMdict_e_examp":
		return "edict", nil
	case "JMnedict", "JMnedict.xml":
		return "enamdict", nil
	case "kanjidic2", "kanjidic2.xml":
		return "kanjidic", nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		_, err := os.Stat(filepath.Join(path, "CATALOGS"))
		if err == nil {
			return "epwing", nil
		}

		_, err = os.Stat(filepath.Join(path, "catalogs"))
		if err == nil {
			return "epwing", nil
		}
	}

	return "", errors.New("unrecognized dictionary format")
}

func ExportDb(inputPath, outputPath, format, language, title string, stride int, pretty bool) error {
	handlers := map[string]func(string, string, string, string, int, bool) error{
		"edict":     jmdExportDb,
		"forms":     formsExportDb,
		"enamdict":  jmnedictExportDb,
		"epwing":    epwingExportDb,
		"kanjidic":  kanjidicExportDb,
		"rikai":     rikaiExportDb,
		"kanjifreq": frequencyKanjiExportDb,
		"termfreq":  frequencyTermsExportDb,
	}

	var err error
	if format == DefaultFormat {
		if format, err = detectFormat(inputPath); err != nil {
			return err
		}
	}

	handler, ok := handlers[strings.ToLower(format)]
	if !ok {
		return errors.New("unrecognized dictionary format")
	}

	return handler(inputPath, outputPath, strings.ToLower(language), title, stride, pretty)
}
