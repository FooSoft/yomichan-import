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

const jmdictRevision = "jmdict4"

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
	for _, tag := range term.Tags {
		switch tag {
		case "arch":
			term.Score -= 100
		}
	}
	for _, tag := range term.TermTags {
		switch tag {
		case "news", "ichi", "spec", "gai1":
			term.Score += 100
		case "P":
			term.Score += 500
		case "iK":
			term.Score -= 100
		}
	}
}

func jmdictAddPriorities(term *dbTerm, priorities ...string) {
	for _, priority := range priorities {
		switch priority {
		case "news1", "ichi1", "spec1", "gai1":
			term.addTermTags("P")
			fallthrough
		case "news2", "ichi2", "spec2", "gai2":
			term.addTermTags(priority[:len(priority)-1])
		}
	}
}

func jmdictBuildTagMeta(entities map[string]string) dbTagList {
	tags := dbTagList{
		dbTag{Name: "news", Notes: "appears frequently in Mainichi Shimbun", Category: "frequent", Order: -2},
		dbTag{Name: "ichi", Notes: "listed as common in Ichimango Goi Bunruishuu", Category: "frequent", Order: -2},
		dbTag{Name: "spec", Notes: "common words not included in frequency lists", Category: "frequent", Order: -2},
		dbTag{Name: "gai", Notes: "common loanword", Category: "frequent", Order: -2},
		dbTag{Name: "P", Notes: "popular term", Category: "popular", Order: -10},
	}

	for name, value := range entities {
		tag := dbTag{Name: name, Notes: value}

		switch name {
		case "exp", "id":
			tag.Category = "expression"
			tag.Order = -5
		case "arch", "iK":
			tag.Category = "archaism"
			tag.Order = -4
		case "adj-f", "adj-i", "adj-ix", "adj-ku", "adj-na", "adj-nari", "adj-no", "adj-pn", "adj-shiku", "adj-t", "adv", "adv-to", "aux-adj",
			"aux", "aux-v", "conj", "cop-da", "ctr", "int", "n-adv", "n", "n-pref", "n-pr", "n-suf", "n-t", "num", "pn", "pref", "prt", "suf",
			"unc", "v1", "v1-s", "v2a-s", "v2b-k", "v2d-s", "v2g-k", "v2g-s", "v2h-k", "v2h-s", "v2k-k", "v2k-s", "v2m-s", "v2n-s", "v2r-k",
			"v2r-s", "v2s-s", "v2t-k", "v2t-s", "v2w-s", "v2y-k", "v2y-s", "v2z-s", "v4b", "v4h", "v4k", "v4m", "v4r", "v4s", "v4t", "v5aru",
			"v5b", "v5g", "v5k", "v5k-s", "v5m", "v5n", "v5r-i", "v5r", "v5s", "v5t", "v5u", "v5u-s", "vi", "vk", "vn", "vr", "vs-c", "vs-i",
			"vs", "vs-s", "vt", "vz":
			tag.Category = "pos"
			tag.Order = -3
		}

		tags = append(tags, tag)
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
		termBase.addTermTags(reading.Information...)

		if kanji == nil {
			termBase.Expression = reading.Reading
			jmdictAddPriorities(&termBase, reading.Priorities...)
		} else {
			termBase.Expression = kanji.Expression
			termBase.Reading = reading.Reading
			termBase.addTermTags(kanji.Information...)

			for _, priority := range kanji.Priorities {
				if hasString(priority, reading.Priorities) {
					jmdictAddPriorities(&termBase, priority)
				}
			}
		}

		var partsOfSpeech []string
		for index, sense := range edictEntry.Sense {

			if len(sense.PartsOfSpeech) != 0 {
				partsOfSpeech = sense.PartsOfSpeech
			}

			if sense.RestrictedReadings != nil && !hasString(reading.Reading, sense.RestrictedReadings) {
				continue
			}

			if kanji != nil && sense.RestrictedKanji != nil && !hasString(kanji.Expression, sense.RestrictedKanji) {
				continue
			}

			term := dbTerm{
				Reading:    termBase.Reading,
				Expression: termBase.Expression,
				Score:      len(edictEntry.Sense) - index,
				Sequence:   edictEntry.Sequence,
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
			term.addTermTags(termBase.TermTags...)
			term.addTags(partsOfSpeech...)
			term.addTags(sense.Fields...)
			term.addTags(sense.Misc...)
			term.addTags(sense.Dialects...)

			jmdictBuildRules(&term)
			jmdictBuildScore(&term)

			terms = append(terms, term)
		}
	}

	if len(edictEntry.Kanji) > 0 {
		for _, kanji := range edictEntry.Kanji {
			for _, reading := range edictEntry.Readings {
				if reading.NoKanji == nil {
					convert(reading, &kanji)
				}
			}
		}
		for _, reading := range edictEntry.Readings {
			if reading.NoKanji != nil  {
				convert(reading, nil)
			}
		}
	} else {
		for _, reading := range edictEntry.Readings {
			convert(reading, nil)
		}
	}

	return terms
}

func jmdictExportDb(inputPath, outputPath, language, title string, stride int, pretty bool) error {
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

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  jmdictBuildTagMeta(entities).crush(),
	}

	return writeDb(
		outputPath,
		title,
		jmdictRevision,
		recordData,
		stride,
		pretty,
	)
}
