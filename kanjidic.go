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
	"strconv"

	"github.com/FooSoft/jmdict"
)

func extractKanjidicKanji(entry jmdict.KanjidicCharacter) dbKanji {
	kanji := dbKanji{Character: entry.Literal}

	if level := entry.Misc.JlptLevel; level != nil {
		kanji.addTags(fmt.Sprintf("jlpt:%s", *level))
	}

	if grade := entry.Misc.Grade; grade != nil {
		kanji.addTags(fmt.Sprintf("grade:%s", *grade))
		if gradeInt, err := strconv.Atoi(*grade); err == nil {
			if gradeInt >= 1 && gradeInt <= 8 {
				kanji.addTags("jouyou")
			} else if gradeInt >= 9 && gradeInt <= 10 {
				kanji.addTags("jinmeiyou")
			}
		}
	}

	for _, number := range entry.DictionaryNumbers {
		if number.Type == "heisig" {
			kanji.addTags(fmt.Sprintf("heisig:%s", number.Value))
		}
	}

	if counts := entry.Misc.StrokeCounts; len(counts) > 0 {
		kanji.addTags(fmt.Sprintf("strokes:%s", counts[0]))
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
				break
			case "ja_kun":
				kanji.Kunyomi = append(kanji.Kunyomi, r.Value)
				break
			}
		}
	}

	return kanji
}

func exportKanjidicDb(outputDir, title string, reader io.Reader, flags int) error {
	dict, err := jmdict.LoadKanjidic(reader)
	if err != nil {
		return err
	}

	var kanji dbKanjiList
	for _, entry := range dict.Characters {
		kanji = append(kanji, extractKanjidicKanji(entry))
	}

	return writeDb(
		outputDir,
		title,
		nil,
		kanji.crush(),
		nil,
		flags&flagPretty == flagPretty,
	)
}
