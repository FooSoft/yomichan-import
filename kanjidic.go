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
	"strconv"
	"strings"

	"github.com/FooSoft/jmdict"
)

type kanjiIndex struct {
	Version int `json:"version"`
	Banks   int `json:"banks"`
	defs    [][]string
}

type kanjiSource struct {
	Character string
	Onyomi    []string
	Kunyomi   []string
	Tags      []string
	Meanings  []string
}

func (s *kanjiSource) addTags(tags ...string) {
	for _, tag := range tags {
		if !hasString(tag, s.Tags) {
			s.Tags = append(s.Tags, tag)
		}
	}
}

func buildKanjiIndex(entries []kanjiSource) kanjiIndex {
	dict := kanjiIndex{
		Version: DB_VERSION,
		Banks:   bankCount(len(entries)),
	}

	for _, e := range entries {
		def := []string{e.Character, strings.Join(e.Onyomi, " "), strings.Join(e.Kunyomi, " "), strings.Join(e.Tags, " ")}
		def = append(def, e.Meanings...)
		dict.defs = append(dict.defs, def)
	}

	return dict
}

func outputKanjiIndex(outputDir string, entries []kanjiSource, pretty bool) error {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return err
	}

	outputIndex, err := os.Create(path.Join(outputDir, "index.json"))
	if err != nil {
		return err
	}
	defer outputIndex.Close()

	dict := buildKanjiIndex(entries)
	indexBytes, err := marshalJson(dict, pretty)
	if err != nil {
		return err
	}

	if _, err = outputIndex.Write(indexBytes); err != nil {
		return err
	}

	defCnt := len(dict.defs)
	for i := 0; i < defCnt; i += BANK_STRIDE {
		outputRef, err := os.Create(path.Join(outputDir, fmt.Sprintf("bank_%d.json", i/BANK_STRIDE+1)))
		if err != nil {
			return err
		}
		defer outputRef.Close()

		indexSrc := i
		indexDst := i + BANK_STRIDE
		if indexDst > defCnt {
			indexDst = defCnt
		}

		refBytes, err := marshalJson(dict.defs[indexSrc:indexDst], pretty)
		if err != nil {
			return err
		}

		if _, err = outputRef.Write(refBytes); err != nil {
			return err
		}
	}

	return nil
}

func convertKanjidicCharacter(kanjidicCharacter jmdict.KanjidicCharacter) kanjiSource {
	character := kanjiSource{Character: kanjidicCharacter.Literal}

	if level := kanjidicCharacter.Misc.JlptLevel; level != nil {
		character.addTags(fmt.Sprintf("jlpt:%s", *level))
	}

	if grade := kanjidicCharacter.Misc.Grade; grade != nil {
		character.addTags(fmt.Sprintf("grade:%s", *grade))
		if gradeInt, err := strconv.Atoi(*grade); err == nil {
			if gradeInt >= 1 && gradeInt <= 8 {
				character.addTags("jouyou")
			} else if gradeInt >= 9 && gradeInt <= 10 {
				character.addTags("jinmeiyou")
			}
		}
	}

	for _, number := range kanjidicCharacter.DictionaryNumbers {
		if number.Type == "heisig" {
			character.addTags(fmt.Sprintf("heisig:%s", number.Value))
		}
	}

	if counts := kanjidicCharacter.Misc.StrokeCounts; len(counts) > 0 {
		character.addTags(fmt.Sprintf("strokes:%s", counts[0]))
	}

	if kanjidicCharacter.ReadingMeaning != nil {
		for _, m := range kanjidicCharacter.ReadingMeaning.Meanings {
			if m.Language == nil || *m.Language == "en" {
				character.Meanings = append(character.Meanings, m.Meaning)
			}
		}

		for _, r := range kanjidicCharacter.ReadingMeaning.Readings {
			switch r.Type {
			case "ja_on":
				character.Onyomi = append(character.Onyomi, r.Value)
				break
			case "ja_kun":
				character.Kunyomi = append(character.Kunyomi, r.Value)
				break
			}
		}
	}

	return character
}

func outputKanjidicJson(outputDir string, reader io.Reader, flags int) error {
	dict, err := jmdict.LoadKanjidic(reader)
	if err != nil {
		return err
	}

	var kanji []kanjiSource
	for _, kanjidicCharacter := range dict.Characters {
		kanji = append(kanji, convertKanjidicCharacter(kanjidicCharacter))
	}

	return outputKanjiIndex(outputDir, kanji, flags&flagPrettyJson == flagPrettyJson)
}
