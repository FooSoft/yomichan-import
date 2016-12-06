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
	"io/ioutil"
	"regexp"
)

type epwingEntry struct {
	Heading string `json:"heading"`
	Text    string `json:"text"`
}

type epwingSubbook struct {
	Title     string        `json:"title"`
	Copyright string        `json:"copyright"`
	Entries   []epwingEntry `json:"entries"`
}

type epwingBook struct {
	CharCode string          `json:"charCode"`
	DiscCode string          `json:"discCode"`
	Subbooks []epwingSubbook `json:"subbooks"`
}

// 3934
// (?P<kana>[^（【<]+)(?:【(?P<kanji>.*)】)?(?:<\?>(?P<native>.*)<\?>)?(?:（(?P<tag>.*)）)?
// (?P<kana>[^（【〖]+)(?:【(?P<expression>.*)】)?(?:〖(?P<native>.*)〗)?(?:（(?P<tag>.*)）)?
// "heading": "きれ‐あが・る【切れ上がる】",
// "text": "きれ‐あが・る【切れ上がる】\n［動ラ五（四）］上の方まで切れる。また、目尻や額の生え際などが上の方へ上がっている。「―・った目元」\n"
// },
// {
// "heading": "きれ‐あじ【切れ味】‐あぢ",
// "text": "きれ‐あじ【切れ味】‐あぢ\n<?>刃物の切れぐあい。「―のいいナイフ」<?>才能・技などの鋭さ。「鋭い―の批評」「―のいいショット」\n"
// },
// {
// "heading": "き‐れい【×綺麗・奇麗】",
// "text": "き‐れい【×綺麗・奇麗】\n［形動］<?>［ナリ］<?>色・形などが華やかな美しさをもっているさま。「―な花」「―に着飾る」<?>姿・顔かたちが整っていて美しいさま。「―な脚」「―な女性」<?>声などが快く聞こえるさま。「―な発音」<?>よごれがなく清潔なさま。「手を―に洗う」「―な空気」「―な選挙」<?>男女間に肉体的な交渉がないさま。清純。「―な関係」<?>乱れたところがないさま。整然としているさま。「机の上を―に片づける」<?>（「きれいに」の形で）残りなく物事が行われるさま。すっかり。「―に忘れる」「―にたいらげる」→美しい［用法］\n［派生］きれいさ［名］\n［類語］（<?>）美しい・美美（びび）しい・煌（きら）やか・鮮やか・美麗・華麗・華美・鮮麗・流麗・優美・美的／（<?>）麗（うるわ）しい・見目よい・端整・端麗・秀麗・佳麗（かれい）・艶美（えんび）・艶麗（えんれい）・あでやか／（<?>）清い・清らか・清潔・清浄（せいじよう・しようじよう）・清澄・清冽（せいれつ）・無垢（むく）・純潔・潔白（けつぱく）\n"
// },

func extractDaijirinTerms(entry epwingEntry) []dbTerm {
	exp := regexp.MustCompile(`(?P<kana>[^（【〖]+)(?:【(?P<expression>.*)】)?(?:〖(?P<native>.*)〗)?(?:（(?P<tag>.*)）)?`)
	matches := exp.FindStringSubmatch(entry.Heading)

	results := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i > 0 {
			results[name] = matches[i]
		}
	}

	return nil
}

func extractDaijirinKanji(entry epwingEntry) []dbKanji {
	return nil
}

func exportEpwingDb(outputDir, title string, reader io.Reader, flags int) error {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}

	var book epwingBook
	if err := json.Unmarshal(data, &book); err != nil {
		return err
	}

	termExtractors := map[string]func(epwingEntry) []dbTerm{
		"三省堂　スーパー大辞林": extractDaijirinTerms,
	}

	var terms dbTermList
	for _, subbook := range book.Subbooks {
		if extractor, ok := termExtractors[subbook.Title]; ok {
			for _, entry := range subbook.Entries {
				terms = append(terms, extractor(entry)...)
			}
		}
	}

	kanjiExtractors := map[string]func(epwingEntry) []dbKanji{}

	var kanji dbKanjiList
	for _, subbook := range book.Subbooks {
		if extractor, ok := kanjiExtractors[subbook.Title]; ok {
			for _, entry := range subbook.Entries {
				kanji = append(kanji, extractor(entry)...)
			}
		}
	}

	return writeDb(
		outputDir,
		title,
		terms.crush(),
		kanji.crush(),
		nil,
		flags&flagPretty == flagPretty,
	)
}
