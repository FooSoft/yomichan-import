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

	"github.com/FooSoft/jmdict"
)

func convertEdictEntry(edictEntry jmdict.JmdictEntry) []termSource {
	var entries []termSource

	convert := func(reading jmdict.JmdictReading, kanji *jmdict.JmdictKanji) {
		if kanji != nil && hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var entryBase termSource
		entryBase.addTags(reading.Information...)
		entryBase.addTagsPri(reading.Priorities...)

		if kanji == nil {
			entryBase.Expression = reading.Reading
		} else {
			entryBase.Expression = kanji.Expression
			entryBase.Reading = reading.Reading
			entryBase.addTags(kanji.Information...)
			entryBase.addTagsPri(kanji.Priorities...)
		}

		for _, sense := range edictEntry.Sense {
			if hasString(reading.Reading, sense.RestrictedReadings) {
				continue
			}

			if kanji != nil && hasString(kanji.Expression, sense.RestrictedKanji) {
				continue
			}

			entry := entryBase
			entry.Tags = append(entry.Tags, entryBase.Tags...)

			for _, glossary := range sense.Glossary {
				entry.Glossary = append(entry.Glossary, glossary.Content)
			}

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

func outputEdictJson(outputDir string, reader io.Reader, flags int) error {
	dict, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	var entries []termSource
	for _, e := range dict.Entries {
		entries = append(entries, convertEdictEntry(e)...)
	}

	return outputTermIndex(outputDir, entries, entities, flags&flagPrettyJson == flagPrettyJson)
}
