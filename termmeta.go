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
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/FooSoft/jmdict"
)

type termMetaEntry struct {
	Expression string   `json:"exp"`
	Reading    string   `json:"read"`
	Tags       []string `json:"tags"`
}

func (meta *termMetaEntry) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, meta.Tags) {
			meta.Tags = append(meta.Tags, tag)
		}
	}
}

func (meta *termMetaEntry) addTagsPri(tags ...string) {
	for _, tag := range tags {
		switch tag {
		case "news1", "ichi1", "spec1", "gai1":
			meta.addTags("P")
			fallthrough
		case "news2", "ichi2", "spec2", "gai2":
			meta.addTags(tag[:len(tag)-1])
			break
		}
	}
}

type termMetaDb struct {
	Version  int               `json:"version"`
	Banks    int               `json:"banks"`
	Entities map[string]string `json:"entities"`
	entries  []termMetaEntry
}

func newTermMetaIndex(entries []termMetaEntry, entities map[string]string) termMetaDb {
	return termMetaDb{
		Version:  DB_VERSION,
		Banks:    bankCount(len(entries)),
		Entities: entities,
		entries:  entries,
	}
}

func (index *termMetaDb) output(dir string, pretty bool) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	bytes, err := marshalJson(index, pretty)
	if err != nil {
		return err
	}

	fp, err := os.Create(path.Join(dir, "index.json"))
	if err != nil {
		return err
	}
	defer fp.Close()

	if _, err = fp.Write(bytes); err != nil {
		return err
	}

	var entries [][]string
	var entryCount = len(entries)
	for _, e := range index.entries {
		entries = append(
			entries,
			[]string{
				e.Expression,
				e.Reading,
				strings.Join(e.Tags, " "),
			},
		)
	}

	for i := 0; i < entryCount; i += BANK_STRIDE {
		indexSrc := i
		indexDst := i + BANK_STRIDE
		if indexDst > entryCount {
			indexDst = entryCount
		}

		bytes, err := marshalJson(entries[indexSrc:indexDst], pretty)
		if err != nil {
			return err
		}

		fp, err := os.Create(path.Join(dir, fmt.Sprintf("bank_%d.json", i/BANK_STRIDE+1)))
		if err != nil {
			return err
		}
		defer fp.Close()

		if _, err = fp.Write(bytes); err != nil {
			return err
		}
	}

	return nil
}

func extractEdictTermMeta(edictEntry jmdict.JmdictEntry) []termMetaEntry {
	var entries []termMetaEntry

	convert := func(reading jmdict.JmdictReading, kanji *jmdict.JmdictKanji) {
		if kanji != nil && reading.Restrictions != nil && !hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var entryBase termMetaEntry
		entryBase.addTags(reading.Information...)
		entryBase.addTagsPri(reading.Priorities...)

		if kanji == nil {
			entryBase.Expression = reading.Reading
			entryBase.addTagsPri(reading.Priorities...)
		} else {
			entryBase.Expression = kanji.Expression
			entryBase.Reading = reading.Reading
			entryBase.addTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					entryBase.addTagsPri(priority)
				}
			}
		}

		for _, sense := range edictEntry.Sense {
			if sense.RestrictedReadings != nil && !hasString(reading.Reading, sense.RestrictedReadings) {
				continue
			}

			if kanji != nil && sense.RestrictedKanji != nil && !hasString(kanji.Expression, sense.RestrictedKanji) {
				continue
			}

			entry := termMetaEntry{
				Reading:    entryBase.Reading,
				Expression: entryBase.Expression,
			}

			entry.addTags(entryBase.Tags...)
			entry.addTags(sense.PartsOfSpeech...)
			entry.addTags(sense.Fields...)
			entry.addTags(sense.Misc...)
			entry.addTags(sense.Dialects...)

			entries = append(entries, entry)
		}
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

func outputTermMetaJson(dir string, reader io.Reader, flags int) error {
	dict, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	var entries []termMetaEntry
	for _, jmdictEntry := range dict.Entries {
		entries = append(entries, extractEdictTermMeta(jmdictEntry)...)
	}

	index := newTermMetaIndex(entries, entities)
	return index.output(dir, flags&flagPrettyJson == flagPrettyJson)
}
