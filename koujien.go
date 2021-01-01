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
	"regexp"
	"strings"

	zig "github.com/FooSoft/zero-epwing-go"
)

type koujienExtractor struct {
	partsExp     *regexp.Regexp
	readGroupExp *regexp.Regexp
	expVarExp    *regexp.Regexp
	metaExp      *regexp.Regexp
	v5Exp        *regexp.Regexp
	v1Exp        *regexp.Regexp
}

func makeKoujienExtractor() epwingExtractor {
	return &koujienExtractor{
		partsExp:     regexp.MustCompile(`([^（【〖]+)(?:【(.*)】)?(?:〖(.*)〗)?(?:（(.*)）)?`),
		readGroupExp: regexp.MustCompile(`[‐・]+`),
		expVarExp:    regexp.MustCompile(`\(([^\)]*)\)`),
		metaExp:      regexp.MustCompile(`（([^）]*)）`),
		v5Exp:        regexp.MustCompile(`(動.[四五](［[^］]+］)?)|(動..二)`),
		v1Exp:        regexp.MustCompile(`(動..一)`),
	}
}
func makeFuzokuExtractor() epwingExtractor {
	return &koujienExtractor{
		partsExp:     regexp.MustCompile(`([^（【〖]+)(?:【(.*)】)?(?:〖(.*)〗)?(?:（(.*)）)?`),
		readGroupExp: regexp.MustCompile(`[-・]+`),
		expVarExp:    regexp.MustCompile(`\(([^\)]*)\)`),
		metaExp:      regexp.MustCompile(`（([^）]*)）`),
		v5Exp:        regexp.MustCompile(`(動.[四五](［[^］]+］)?)|(動..二)`),
		v1Exp:        regexp.MustCompile(`(動..一)`),
	}
}

func (e *koujienExtractor) extractTerms(entry zig.BookEntry, sequence int) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	var expressions, readings []string
	if expression := matches[2]; len(expression) > 0 {
		expression = e.metaExp.ReplaceAllLiteralString(expression, "")
		for _, split := range strings.Split(expression, "・") {
			splitInc := e.expVarExp.ReplaceAllString(split, "$1")
			expressions = append(expressions, splitInc)
			if split != splitInc {
				splitExc := e.expVarExp.ReplaceAllLiteralString(split, "")
				expressions = append(expressions, splitExc)
			}
		}
	}

	if reading := matches[1]; len(reading) > 0 {
		reading = e.readGroupExp.ReplaceAllLiteralString(reading, "")
		readings = append(readings, reading)
	}

	var tags []string
	for _, split := range strings.Split(entry.Text, "\n") {
		if matches := e.metaExp.FindStringSubmatch(split); matches != nil {
			for _, tag := range strings.Split(matches[1], "・") {
				tags = append(tags, tag)
			}
		}
	}

	var terms []dbTerm
	if len(expressions) == 0 {
		for _, reading := range readings {
			term := dbTerm{
				Expression: reading,
				Glossary:   []string{entry.Text},
				Sequence:   sequence,
			}

			e.exportRules(&term, tags)
			terms = append(terms, term)
		}

	} else {
		for _, expression := range expressions {
			for _, reading := range readings {
				term := dbTerm{
					Expression: expression,
					Reading:    reading,
					Glossary:   []string{entry.Text},
					Sequence:   sequence,
				}

				e.exportRules(&term, tags)
				terms = append(terms, term)
			}
		}
	}

	return terms
}

func (*koujienExtractor) extractKanji(entry zig.BookEntry) []dbKanji {
	return nil
}

func (e *koujienExtractor) exportRules(term *dbTerm, tags []string) {
	for _, tag := range tags {
		if tag == "形" {
			term.addRules("adj-i")
		} else if tag == "動サ変" && (strings.HasSuffix(term.Expression, "する") || strings.HasSuffix(term.Expression, "為る")) {
			term.addRules("vs")
		} else if term.Expression == "来る" {
			term.addRules("vk")
		} else if e.v5Exp.MatchString(tag) {
			term.addRules("v5")
		} else if e.v1Exp.MatchString(tag) {
			term.addRules("v1")
		}
	}
}

func (*koujienExtractor) getRevision() string {
	return "koujien"
}

func (*koujienExtractor) getFontNarrow() map[int]string {
	return map[int]string{}
}

func (*koujienExtractor) getFontWide() map[int]string {
	return map[int]string{
		41531: "⟨",
		41532: "⟩",
		42017: "⇿",
		42018: "🈑",
		42023: "🈩",
		42024: "🈔",
		42025: "㊇",
		42026: "3",
		42027: "❷",
		42028: "❶",
		42031: "❸",
		42037: "❹",
		42043: "❺",
		42045: "❻",
		42057: "❼",
		42083: "❽",
		42284: "❾",
		42544: "❿",
		42561: "鉏",
		43611: "⓫",
		43612: "⓬",
		44142: "𑖀",
		44856: "㉑",
		44857: "㉒",
		46374: "〔",
		46375: "〕",
		46390: "①",
		46391: "②",
		46392: "③",
		46393: "④",
		46394: "⑤",
		46395: "⑥",
		46396: "⑦",
		46397: "⑧",
		46398: "⑨",
		46399: "⑩",
		46400: "⑪",
		46401: "⑫",
		46402: "⑬",
		46403: "⑭",
		46404: "⑮",
		46405: "⑯",
		46406: "⑰",
		46407: "⑱",
		46408: "⑲",
		46409: "⑳",
		46677: "⇀",
		46420: "⇨",
		47175: "(季)",
		56383: "㋐",
		56384: "㋑",
		56385: "㋒",
		56386: "㋓",
		56387: "㋔",
		56388: "㋕",
		56389: "㋖",
		56390: "㋗",
		56391: "㋘",
		56392: "㋙",
		56393: "㋚",
		56394: "㋛",
		56395: "㋜",
		56396: "㋝",
		56397: "㋞",
		56398: "▷",
	}
}
