/*
 * Copyright (c) 2017 Alex Yatskov <alex@foosoft.net>
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
	"database/sql"
	"log"
	"strings"

	_ "github.com/mattn/go-sqlite3"
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

func rikaiExportDb(inputPath, outputDir, title string, stride int, pretty bool) error {
	db, err := sql.Open("sqlite3", inputPath)
	if err != nil {
		return err
	}
	defer db.Close()

	dictRows, err := db.Query("SELECT kanji, kana, entry FROM dict")
	if err != nil {
		return err
	}

	for dictRows.Next() {
		var kanji, kana, entry *string
		if err := dictRows.Scan(&kanji, &kana, &entry); err != nil {
			return err
		}

		log.Print(kanji, kana, entry)
	}

	if title == "" {
		title = "Rikai"
	}

	return nil

	// return writeDb(
	// 	outputDir,
	// 	title,
	// 	RIKAI_REVISION,
	// 	terms.crush(),
	// 	nil,
	// 	rikaiBuildTagMeta(entities),
	// 	stride,
	// 	pretty,
	// )
}
