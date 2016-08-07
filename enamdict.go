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

func convertJmnedictEntry(enamdictEntry jmdict.JmnedictEntry) []vocabSource {
	var entries []vocabSource

	convert := func(reading jmdict.JmnedictReading, kanji *jmdict.JmnedictKanji) {
		if kanji != nil && hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var entry vocabSource
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

func processJmnedict(writer io.Writer, reader io.Reader, flags int) error {
	dict, entities, err := jmdict.LoadJmnedictNoTransform(reader)
	if err != nil {
		return err
	}

	var entries []vocabSource
	for _, e := range dict.Entries {
		entries = append(entries, convertJmnedictEntry(e)...)
	}

	return outputVocabJson(writer, entries, entities, flags&flagPrettyJson == flagPrettyJson)
}
