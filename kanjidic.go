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
	"encoding/json"
	"io"
	"strings"

	"github.com/FooSoft/jmdict"
)

type kanjiJson struct {
	Characters map[string][]string `json:"c"`
}

type kanjiSource struct {
	Character string
	Kunyomi   []string
	Onyomi    []string
	Meanings  []string
}

func buildKanjiJson(kanji []kanjiSource) kanjiJson {
	dict := kanjiJson{make(map[string][]string)}

	for _, k := range kanji {
		var params []string
		params = append(params, strings.Join(k.Onyomi, " "))
		params = append(params, strings.Join(k.Kunyomi, " "))
		params = append(params, k.Meanings...)
		dict.Characters[k.Character] = params
	}

	return dict
}

func outputKanjiJson(writer io.Writer, kanji []kanjiSource, pretty bool) error {
	dict := buildKanjiJson(kanji)

	var (
		bytes []byte
		err   error
	)

	if pretty {
		bytes, err = json.MarshalIndent(dict, "", "    ")
	} else {
		bytes, err = json.Marshal(dict)
	}

	if err != nil {
		return err
	}

	_, err = writer.Write(bytes)
	return err
}

func convertKanjidicCharacter(kanjidicCharacter jmdict.KanjidicCharacter) kanjiSource {
	character := kanjiSource{Character: kanjidicCharacter.Literal}

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

func processKanjidic(writer io.Writer, reader io.Reader, flags int) error {
	dict, err := jmdict.LoadKanjidic(reader)
	if err != nil {
		return err
	}

	var kanji []kanjiSource
	for _, kanjidicCharacter := range dict.Characters {
		kanji = append(kanji, convertKanjidicCharacter(kanjidicCharacter))
	}

	return outputKanjiJson(writer, kanji, flags&flagPrettyJson == flagPrettyJson)
}
