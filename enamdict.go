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
	"os"

	"github.com/FooSoft/jmdict"
)

const jmnedictRevision = "jmnedict1"

func jmnedictBuildTagMeta(entities map[string]string) dbTagList {
	var tags dbTagList

	for name, value := range entities {
		tag := dbTag{Name: name, Notes: value}

		switch name {
		case "company", "fem", "given", "masc", "organization", "person", "place", "product", "station", "surname", "unclass", "work":
			tag.Category = "name"
			tag.Order = 4
		}

		tags = append(tags, tag)
	}

	return tags
}

func jmnedictExtractTerms(enamdictEntry jmdict.JmnedictEntry) []dbTerm {
	var terms []dbTerm

	convert := func(reading jmdict.JmnedictReading, kanji *jmdict.JmnedictKanji) {
		if kanji != nil && hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var term dbTerm
		term.Sequence = enamdictEntry.Sequence
		term.addTermTags(reading.Information...)

		if kanji == nil {
			term.Expression = reading.Reading
		} else {
			term.Expression = kanji.Expression
			term.Reading = reading.Reading
			term.addTermTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					term.addTermTags(priority)
				}
			}
		}

		for _, trans := range enamdictEntry.Translations {
			term.Glossary = append(term.Glossary, trans.Translations...)
			term.addDefinitionTags(trans.NameTypes...)
		}

		terms = append(terms, term)
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

	return terms
}

func jmnedictExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dict, entities, err := jmdict.LoadJmnedictNoTransform(reader)
	if err != nil {
		return err
	}

	var terms dbTermList
	for _, entry := range dict.Entries {
		terms = append(terms, jmnedictExtractTerms(entry)...)
	}

	if title == "" {
		title = "JMnedict"
	}

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  jmnedictBuildTagMeta(entities).crush(),
	}

	return writeDb(
		outputPath,
		title,
		jmnedictRevision,
		true,
		recordData,
		stride,
		pretty,
	)
}
