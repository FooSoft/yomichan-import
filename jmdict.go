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
	"strings"

	"github.com/FooSoft/jmdict"
)

type dictJson struct {
	Indices map[string]string `json:"i"`
	Items   [][]string        `json:"d"`
}

type dictEntry struct {
	Expression string
	Reading    string
	Tags       []string
	Glossary   []string
}

func (d *dictEntry) addTags(tags []string) {
	for _, tag := range tags {
		if findString(tag, d.Tags) == -1 {
			d.Tags = append(d.Tags, tag)
		}
	}
}

func appendIndex(indices map[string]string, key string, value int) {
	items, _ := indices[key]
	if len(indices) == 0 {
		items = fmt.Sprintf(" %d", value)
	} else {
		items = string(value)
	}
	indices[key] = items
}

func buildDictJson(entries []dictEntry) dictJson {
	dict := dictJson{Indices: make(map[string]string)}

	for index, entry := range entries {
		item := []string{entry.Expression, entry.Reading, strings.Join(entry.Tags, " ")}
		item = append(item, entry.Glossary...)
		dict.Items = append(dict.Items, item)

		appendIndex(dict.Indices, entry.Expression, index)
		if len(entry.Reading) > 0 {
			appendIndex(dict.Indices, entry.Reading, index)
		}
	}

	return dict
}

func outputJson(entries []dictEntry, writer io.Writer) error {
	dict := buildDictJson(entries)

	bytes, err := json.MarshalIndent(dict, "", "    ")
	// bytes, err := json.Marshal(dict)
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

func convertEnamdictEntry(enamdictEntry jmdict.EnamdictEntry) []dictEntry {
	var entries []dictEntry

	convert := func(reading jmdict.EnamdictReading, kanji *jmdict.EnamdictKanji) {
		if kanji != nil && findString(kanji.Expression, reading.Restrictions) != -1 {
			return
		}

		var entry dictEntry
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

func convertEdictEntry(edictEntry jmdict.EdictEntry) []dictEntry {
	var entries []dictEntry

	convert := func(reading jmdict.EdictReading, kanji *jmdict.EdictKanji) {
		if kanji != nil && findString(kanji.Expression, reading.Restrictions) != -1 {
			return
		}

		var entry dictEntry
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

func processEnamdict(reader io.Reader, writer io.Writer) error {
	enamdictEntries, _, err := jmdict.LoadEnamdict(reader, false)
	if err != nil {
		return err
	}

	var entries []dictEntry
	for _, enamdictEntry := range enamdictEntries {
		entries = append(entries, convertEnamdictEntry(enamdictEntry)...)
	}

	return outputJson(entries, writer)
}

func processEdict(reader io.Reader, writer io.Writer) error {
	edictEntries, _, err := jmdict.LoadEdict(reader, false)
	if err != nil {
		return err
	}

	var entries []dictEntry
	for _, edictEntry := range edictEntries {
		entries = append(entries, convertEdictEntry(edictEntry)...)
	}

	return outputJson(entries, writer)
}
