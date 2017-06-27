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
	"strconv"

	"github.com/FooSoft/jmdict"
)

const KANJIDIC_REVISION = "kanjidic2"

func kanjidicExtractKanji(entry jmdict.KanjidicCharacter) dbKanji {
	kanji := dbKanji{
		Character: entry.Literal,
		Stats:     make(map[string]string),
	}

	if level := entry.Misc.JlptLevel; level != nil {
		kanji.Stats["jlpt level"] = *level
	}

	if grade := entry.Misc.Grade; grade != nil {
		kanji.Stats["school grade"] = *grade
		if gradeInt, err := strconv.Atoi(*grade); err == nil {
			if gradeInt >= 1 && gradeInt <= 8 {
				kanji.addTags("jouyou")
			} else if gradeInt >= 9 && gradeInt <= 10 {
				kanji.addTags("jinmeiyou")
			}
		}
	}

	for _, number := range entry.DictionaryNumbers {
		kanji.Stats[number.Type] = number.Value
	}

	if counts := entry.Misc.StrokeCounts; len(counts) > 0 {
		kanji.Stats["stroke count"] = counts[0]
	}

	if entry.ReadingMeaning != nil {
		for _, m := range entry.ReadingMeaning.Meanings {
			if m.Language == nil || *m.Language == "en" {
				kanji.Meanings = append(kanji.Meanings, m.Meaning)
			}
		}

		for _, r := range entry.ReadingMeaning.Readings {
			switch r.Type {
			case "ja_on":
				kanji.Onyomi = append(kanji.Onyomi, r.Value)
			case "ja_kun":
				kanji.Kunyomi = append(kanji.Kunyomi, r.Value)
			}
		}
	}

	return kanji
}

func kanjidicExportDb(inputPath, outputDir, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dict, err := jmdict.LoadKanjidic(reader)
	if err != nil {
		return err
	}

	var kanji dbKanjiList
	for _, entry := range dict.Characters {
		kanji = append(kanji, kanjidicExtractKanji(entry))
	}

	tagMeta := map[string]dbTagMeta{
		"jouyou":           {Notes: "included in list of regular-use characters", Category: "frequent", Order: -5},
		"jinmeiyou":        {Notes: "included in list of characters for use in personal names", Category: "frequent", Order: -5},
		"nelson_c":         {Notes: "Modern Reader's Japanese-English Character Dictionary"},
		"nelson_n":         {Notes: "The New Nelson Japanese-English Character Dictionary"},
		"halpern_njecd":    {Notes: "New Japanese-English Character Dictionary"},
		"halpern_kkd":      {Notes: "Kodansha Kanji Dictionary"},
		"halpern_kkld":     {Notes: "Kanji Learners Dictionary"},
		"halpern_kkld_2ed": {Notes: "Kanji Learners Dictionary"},
		"heisig":           {Notes: "Remembering The  Kanji"},
		"heisig6":          {Notes: "Remembering The  Kanji, Sixth Ed."},
		"gakken":           {Notes: "A New Dictionary of Kanji Usage"},
		"oneill_names":     {Notes: "Japanese Names"},
		"oneill_kk":        {Notes: "Essential Kanji"},
		"moro":             {Notes: "Daikanwajiten"},
		"henshall":         {Notes: "A Guide To Remembering Japanese Characters"},
		"sh_kk":            {Notes: "Kanji and Kana"},
		"sh_kk2":           {Notes: "Kanji and Kana"},
		"sakade":           {Notes: "A Guide To Reading and Writing Japanese"},
		"jf_cards":         {Notes: "Japanese Kanji Flashcards"},
		"henshall3":        {Notes: "A Guide To Reading and Writing Japanese"},
		"tutt_cards":       {Notes: "Tuttle Kanji Cards"},
		"crowley":          {Notes: "The Kanji Way to Japanese Language Power"},
		"kanji_in_context": {Notes: "Kanji in Context"},
		"busy_people":      {Notes: "Japanese For Busy People"},
		"kodansha_compact": {Notes: "Kodansha Compact Kanji Guide"},
		"maniette":         {Notes: "Les Kanjis dans la tete"},
	}

	if title == "" {
		title = "KANJIDIC2"
	}

	return writeDb(
		outputDir,
		title,
		KANJIDIC_REVISION,
		nil,
		kanji.crush(),
		tagMeta,
		stride,
		pretty,
	)
}
