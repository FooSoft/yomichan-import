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
)

const RIKAI_REVISION = "rikai1"

func rikaiBuildRules(term *dbTerm) {
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

func rikaiBuildScore(term *dbTerm) {
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

func rikaiBuildTagMeta(entities map[string]string) map[string]dbTagMeta {
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

func rikaiExtractTerms(edictEntry rikai.rikaiEntry) []dbTerm {

	return terms
}

func rikaiExportDb(inputPath, outputDir, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dict, entities, err := rikai.LoadrikaiNoTransform(reader)
	if err != nil {
		return err
	}

	var terms dbTermList
	for _, entry := range dict.Entries {
		terms = append(terms, rikaiExtractTerms(entry)...)
	}

	if title == "" {
		title = "rikai"
	}

	return writeDb(
		outputDir,
		title,
		RIKAI_REVISION,
		terms.crush(),
		nil,
		rikaiBuildTagMeta(entities),
		stride,
		pretty,
	)
}
