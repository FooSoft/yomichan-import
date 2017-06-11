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
	"strings"

	"github.com/FooSoft/jmdict"
)

const JMDICT_REVISION = "jmdict2"

func jmdictBuildRules(term *dbTerm) {
	for _, tag := range term.Tags {
		switch tag {
		case "adj-i", "v1", "vk":
			term.addRules(tag)
		default:
			if strings.HasPrefix(tag, "v5") {
				term.addRules("v5")
			} else if strings.HasPrefix(tag, "vs") {
				term.addRules("vs")
			}
		}
	}
}

func jmdictBuildScore(term *dbTerm) {
	term.Score = 0
	for _, tag := range term.Tags {
		switch tag {
		case "P":
			term.Score += 5
		case "arch", "iK":
			term.Score -= 1
		}
	}
}

func jmdictAddPriorities(term *dbTerm, priorities ...string) {
	for _, priority := range priorities {
		switch priority {
		case "news1", "ichi1", "spec1", "gai1":
			term.addTags("P")
			fallthrough
		case "news2", "ichi2", "spec2", "gai2":
			term.addTags(priority[:len(priority)-1])
		}
	}
}

func jmdictBuildTagMeta(entities map[string]string) map[string]dbTagMeta {
	tags := map[string]dbTagMeta{
		"news": {Notes: "appears frequently in Mainichi Shimbun", Category: "frequent", Order: -2},
		"ichi": {Notes: "listed as common in Ichimango Goi Bunruishuu", Category: "frequent", Order: -2},
		"spec": {Notes: "common words not included in frequency lists", Category: "frequent", Order: -2},
		"gai":  {Notes: "common loanword", Category: "frequent", Order: -2},
		"P":    {Notes: "popular term", Category: "popular", Order: -10},
	}

	for name, value := range entities {
		tag := dbTagMeta{Notes: value}

		switch name {
		case "exp", "id":
			tag.Category = "expression"
			tag.Order = -5
		case "arch", "iK":
			tag.Category = "archaism"
			tag.Order = -5
		}

		tags[name] = tag
	}

	return tags
}

func jmdictExtractTerms(edictEntry jmdict.JmdictEntry, language string) []dbTerm {
	var terms []dbTerm

	convert := func(reading jmdict.JmdictReading, kanji *jmdict.JmdictKanji) {
		if kanji != nil && reading.Restrictions != nil && !hasString(kanji.Expression, reading.Restrictions) {
			return
		}

		var termBase dbTerm
		termBase.addTags(reading.Information...)

		if kanji == nil {
			termBase.Expression = reading.Reading
			jmdictAddPriorities(&termBase, reading.Priorities...)
		} else {
			termBase.Expression = kanji.Expression
			termBase.Reading = reading.Reading
			termBase.addTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					jmdictAddPriorities(&termBase, priority)
				}
			}
		}

		var partsOfSpeech []string
		for index, sense := range edictEntry.Sense {
			if sense.RestrictedReadings != nil && !hasString(reading.Reading, sense.RestrictedReadings) {
				continue
			}

			if kanji != nil && sense.RestrictedKanji != nil && !hasString(kanji.Expression, sense.RestrictedKanji) {
				continue
			}

			term := dbTerm{
				Reading:    termBase.Reading,
				Expression: termBase.Expression,
			}

			for _, glossary := range sense.Glossary {
				if glossary.Language == nil && language == "" || glossary.Language != nil && language == *glossary.Language {
					term.Glossary = append(term.Glossary, glossary.Content)
				}
			}

			if len(term.Glossary) == 0 {
				continue
			}

			term.addTags(termBase.Tags...)
			term.addTags(sense.PartsOfSpeech...)
			term.addTags(sense.Fields...)
			term.addTags(sense.Misc...)
			term.addTags(sense.Dialects...)

			if index == 0 {
				partsOfSpeech = sense.PartsOfSpeech
			} else {
				term.addTags(partsOfSpeech...)
			}

			jmdictBuildRules(&term)
			jmdictBuildScore(&term)

			terms = append(terms, term)
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

	return terms
}

func jmdictExportDb(inputPath, outputDir, language, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dict, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	var langTag string
	switch language {
	case "dutch":
		langTag = "dut"
	case "french":
		langTag = "fre"
	case "german":
		langTag = "ger"
	case "hungarian":
		langTag = "hun"
	case "italian":
		langTag = "ita"
	case "russian":
		langTag = "rus"
	case "slovenian":
		langTag = "slv"
	case "spanish":
		langTag = "spa"
	case "swedish":
		langTag = "swe"
	}

	var terms dbTermList
	for _, entry := range dict.Entries {
		terms = append(terms, jmdictExtractTerms(entry, langTag)...)
	}

	if title == "" {
		title = "JMdict"
	}

	return writeDb(
		outputDir,
		title,
		JMDICT_REVISION,
		terms.crush(),
		nil,
		jmdictBuildTagMeta(entities),
		stride,
		pretty,
	)
}
