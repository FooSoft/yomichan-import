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
	"strconv"
	"strings"

	"github.com/FooSoft/jmdict"
)

//
//	Edict and Enamdict processing
//

type vocabDictJson struct {
	Indices  map[string]string `json:"i"`
	Entities map[string]string `json:"e"`
	Defs     [][]string        `json:"d"`
}

type vocabDictSource struct {
	Expression string
	Reading    string
	Tags       []string
	Glossary   []string
}

func (d *vocabDictSource) addTags(tags []string) {
	for _, tag := range tags {
		if findString(tag, d.Tags) == -1 {
			d.Tags = append(d.Tags, tag)
		}
	}
}

func appendIndex(indices map[string]string, key string, value int) {
	def, _ := indices[key]
	if len(def) > 0 {
		def += " "
	}
	def += strconv.Itoa(value)
	indices[key] = def
}

func buildVocabDictJson(entries []vocabDictSource, entities map[string]string) vocabDictJson {
	dict := vocabDictJson{
		Indices:  make(map[string]string),
		Entities: entities,
	}

	for i, e := range entries {
		entry := []string{e.Expression, e.Reading, strings.Join(e.Tags, " ")}
		entry = append(entry, e.Glossary...)
		dict.Defs = append(dict.Defs, entry)

		appendIndex(dict.Indices, e.Expression, i)
		if len(e.Reading) > 0 {
			appendIndex(dict.Indices, e.Reading, i)
		}
	}

	return dict
}

func outputVocabDictJson(writer io.Writer, entries []vocabDictSource, entities map[string]string, pretty bool) error {
	dict := buildVocabDictJson(entries, entities)

	var (
		bytes []byte
		err   error
	)

	if pretty {
		bytes, err = json.MarshalIndent(dict, "", "    ")
	} else {
		bytes, err = json.Marshal(dict)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}

func findString(needle string, haystack []string) int {
	for index, value := range haystack {
		if needle == value {
			return index
		}
	}

	return -1
}

func convertEnamdictEntry(enamdictEntry jmdict.EnamdictEntry) []vocabDictSource {
	var entries []vocabDictSource

	convert := func(reading jmdict.EnamdictReading, kanji *jmdict.EnamdictKanji) {
		if kanji != nil && findString(kanji.Expression, reading.Restrictions) != -1 {
			return
		}

		var entry vocabDictSource
		if kanji == nil {
			entry.Expression = reading.Reading
		} else {
			entry.Expression = kanji.Expression
			entry.Reading = reading.Reading

			entry.addTags(kanji.Information)
			entry.addTags(kanji.Priorities)
		}

		entry.addTags(reading.Information)
		entry.addTags(reading.Priorities)

		for _, trans := range enamdictEntry.Translations {
			entry.Glossary = append(entry.Glossary, trans.Translations...)
			entry.addTags(trans.NameTypes)
		}

		entries = append(entries, entry)
	}

	if len(enamdictEntry.Kanji) > 0 {
		for _, kanji := range enamdictEntry.Kanji {
			for _, reading := range enamdictEntry.Readings {
				convert(reading, &kanji)
			}
		}
	} else {
		for _, reading := range enamdictEntry.Readings {
			convert(reading, nil)
		}
	}

	return entries
}

func convertEdictEntry(edictEntry jmdict.EdictEntry) []vocabDictSource {
	var entries []vocabDictSource

	convert := func(reading jmdict.EdictReading, kanji *jmdict.EdictKanji) {
		if kanji != nil && findString(kanji.Expression, reading.Restrictions) != -1 {
			return
		}

		var entry vocabDictSource
		if kanji == nil {
			entry.Expression = reading.Reading
		} else {
			entry.Expression = kanji.Expression
			entry.Reading = reading.Reading

			entry.addTags(kanji.Information)
			entry.addTags(kanji.Priorities)
		}

		entry.addTags(reading.Information)
		entry.addTags(reading.Priorities)

		for _, sense := range edictEntry.Sense {
			if findString(reading.Reading, sense.RestrictedReadings) != -1 {
				continue
			}

			if kanji != nil && findString(kanji.Expression, sense.RestrictedKanji) != -1 {
				continue
			}

			for _, glossary := range sense.Glossary {
				entry.Glossary = append(entry.Glossary, glossary.Content)
			}

			entry.addTags(sense.PartsOfSpeech)
			entry.addTags(sense.Fields)
			entry.addTags(sense.Misc)
			entry.addTags(sense.Dialects)
		}

		entries = append(entries, entry)
	}

	if len(edictEntry.Kanji) > 0 {
		for _, kanji := range edictEntry.Kanji {
			for _, reading := range edictEntry.Readings {
				convert(reading, &kanji)
			}
		}
	} else {
		for _, reading := range edictEntry.Readings {
			convert(reading, nil)
		}
	}

	return entries
}

func processEnamdict(writer io.Writer, reader io.Reader, flags int) error {
	enamdictEntries, entities, err := jmdict.LoadEnamdict(reader, false)
	if err != nil {
		return err
	}

	var entries []vocabDictSource
	for _, enamdictEntry := range enamdictEntries {
		entries = append(entries, convertEnamdictEntry(enamdictEntry)...)
	}

	return outputVocabDictJson(writer, entries, entities, flags&flagPrettyJson == flagPrettyJson)
}

func processEdict(writer io.Writer, reader io.Reader, flags int) error {
	edictEntries, entities, err := jmdict.LoadEdict(reader, false)
	if err != nil {
		return err
	}

	var entries []vocabDictSource
	for _, edictEntry := range edictEntries {
		entries = append(entries, convertEdictEntry(edictEntry)...)
	}

	return outputVocabDictJson(writer, entries, entities, flags&flagPrettyJson == flagPrettyJson)
}

//
//	Kanjidic processing
//

type characterDictJson struct {
	Characters map[string][]string `json:"c"`
}

type characterDictSource struct {
	Character string
	Kunyomi   []string
	Onyomi    []string
	Meanings  []string
}

func buildCharacterDictJson(characters []characterDictSource) characterDictJson {
	dict := characterDictJson{make(map[string][]string)}

	for _, c := range characters {
		var params []string
		params = append(params, strings.Join(c.Onyomi, " "))
		params = append(params, strings.Join(c.Kunyomi, " "))
		params = append(params, c.Meanings...)
		dict.Characters[c.Character] = params
	}

	return dict
}

func outputCharacterDictJson(writer io.Writer, characters []characterDictSource, pretty bool) error {
	dict := buildCharacterDictJson(characters)

	var (
		bytes []byte
		err   error
	)

	if pretty {
		bytes, err = json.MarshalIndent(dict, "", "    ")
	} else {
		bytes, err = json.Marshal(dict)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}

func convertKanjidicCharacter(kanjidicCharacter jmdict.KanjidicCharacter) characterDictSource {
	var character characterDictSource

	character.Character = kanjidicCharacter.Literal

	return character
}

func processKanjidic(writer io.Writer, reader io.Reader, flags int) error {
	kanjidicCharacters, err := jmdict.LoadKanjidic(reader)
	if err != nil {
		return err
	}

	var characters []characterDictSource
	for _, kanjidicCharacter := range kanjidicCharacters {
		characters = append(characters, convertKanjidicCharacter(kanjidicCharacter))
	}

	return nil
}
