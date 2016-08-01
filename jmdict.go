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
	"io"
	"log"

	"github.com/FooSoft/jmdict"
)

type dictEntry struct {
	Expression string
	Reading    string
	Glossary   []string
	Tags       map[string]bool
}

func (d *dictEntry) addTags(tags []string) {
	for _, tag := range tags {
		d.Tags[tag] = true
	}
}

func findString(needle string, haystack []string) int {
	for index, value := range haystack {
		if needle == value {
			return index
		}
	}

	return -1
}

func convertEdictEntry(edictEntry jmdict.EdictEntry) []dictEntry {
	var entries []dictEntry

	convert := func(reading jmdict.EdictReading, kanji *jmdict.EdictKanji) {
		if kanji != nil && findString(kanji.Expression, reading.Restrictions) != -1 {
			return
		}

		entry := dictEntry{Tags: make(map[string]bool)}

		if kanji == nil {
			entry.Expression = reading.Reading
		} else {
			entry.Expression = kanji.Expression
			entry.Reading = reading.Reading

			entry.addTags(kanji.Information)
			entry.addTags(kanji.Priority)
		}

		entry.addTags(reading.Information)
		entry.addTags(reading.Priority)

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

			entry.addTags(sense.PartOfSpeech)
			entry.addTags(sense.Fields)
			entry.addTags(sense.Misc)
			entry.addTags(sense.Dialect)
		}

		entries = append(entries, entry)
	}

	if len(edictEntry.Kanji) > 0 {
		for _, kanji := range edictEntry.Kanji {
			for _, reading := range edictEntry.Reading {
				convert(reading, &kanji)
			}
		}
	} else {
		for _, reading := range edictEntry.Reading {
			convert(reading, nil)
		}
	}

	return entries
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

	log.Print(entries[42])

	return nil
}
