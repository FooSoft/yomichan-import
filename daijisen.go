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
)

type daijisenExtractor struct {
	partsExp     *regexp.Regexp
	readGroupExp *regexp.Regexp
	expVarExp    *regexp.Regexp
	metaExp      *regexp.Regexp
	v5Exp        *regexp.Regexp
	v1Exp        *regexp.Regexp
}

func makeDaijisenExtractor() epwingExtractor {
	return &daijisenExtractor{
		partsExp:     regexp.MustCompile(`([^（【〖]+)(?:【(.*)】)?(?:〖(.*)〗)?(?:（(.*)）)?`),
		readGroupExp: regexp.MustCompile(`[-・]+`),
		expVarExp:    regexp.MustCompile(`\(([^\)]*)\)`),
		metaExp:      regexp.MustCompile(`（([^）]*)）`),
		v5Exp:        regexp.MustCompile(`(動.[四五](［[^］]+］)?)|(動..二)`),
		v1Exp:        regexp.MustCompile(`(動..一)`),
	}
}

func (e *daijisenExtractor) extractTerms(entry epwingEntry) []dbTerm {
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
				}

				e.exportRules(&term, tags)
				terms = append(terms, term)
			}
		}
	}

	return terms
}

func (*daijisenExtractor) extractKanji(entry epwingEntry) []dbKanji {
	return nil
}

func (e *daijisenExtractor) exportRules(term *dbTerm, tags []string) {
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

func (*daijisenExtractor) getRevision() string {
	return "daijisen1"
}

func (*daijisenExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		0xa121: " ",
		0xa122: "¡",
		0xa123: "¢",
		0xa124: "£",
		0xa125: "¤",
		0xa126: "¥",
		0xa127: "¦",
		0xa128: "§",
		0xa129: "¨",
		0xa12a: "©",
		0xa12b: "ª",
		0xa12c: "«",
		0xa12d: "¬",
		0xa12e: "­",
		0xa12f: "®",
		0xa130: "¯",
		0xa131: "°",
		0xa132: "±",
		0xa133: "²",
		0xa134: "³",
		0xa135: "´",
		0xa136: "µ",
		0xa137: "¶",
		0xa138: "·",
		0xa139: "¸",
		0xa13a: "¹",
		0xa13b: "º",
		0xa13c: "»",
		0xa13d: "¼",
		0xa13e: "½",
		0xa13f: "¾",
		0xa140: "¿",
		0xa141: "À",
		0xa142: "Á",
		0xa143: "Â",
		0xa144: "Ã",
		0xa145: "Ä",
		0xa146: "Å",
		0xa147: "Æ",
		0xa148: "Ç",
		0xa149: "È",
		0xa14a: "É",
		0xa14b: "Ê",
		0xa14c: "Ë",
		0xa14d: "Ì",
		0xa14e: "Í",
		0xa14f: "Î",
		0xa150: "Ï",
		0xa151: "Ð",
		0xa152: "Ñ",
		0xa153: "Ò",
		0xa154: "Ó",
		0xa155: "Ô",
		0xa156: "Õ",
		0xa157: "Ö",
		0xa158: "×",
		0xa159: "Ø",
		0xa15a: "Ù",
		0xa15b: "Ú",
		0xa15c: "Û",
		0xa15d: "Ü",
		0xa15e: "Ý",
		0xa15f: "Þ",
		0xa160: "ß",
		0xa161: "à",
		0xa162: "á",
		0xa163: "â",
		0xa164: "ã",
		0xa165: "ä",
		0xa166: "å",
		0xa167: "æ",
		0xa168: "ç",
		0xa169: "è",
		0xa16a: "é",
		0xa16b: "ê",
		0xa16c: "ë",
		0xa16d: "ì",
		0xa16e: "í",
		0xa16f: "î",
		0xa170: "ï",
		0xa171: "ð",
		0xa172: "ñ",
		0xa173: "ò",
		0xa174: "ó",
		0xa175: "ô",
		0xa176: "õ",
		0xa177: "ö",
		0xa178: "÷",
		0xa179: "ø",
		0xa17a: "ù",
		0xa17b: "ú",
		0xa17c: "û",
		0xa17d: "ü",
		0xa17e: "ý",
		0xa221: "þ",
		0xa222: "ÿ",
		0xa223: "Ā",
		0xa224: "ā",
		0xa225: "Ă",
		0xa226: "ă",
		0xa227: "Ą",
		0xa228: "ą",
		0xa229: "Ć",
		0xa22a: "ć",
		0xa22b: "Ĉ",
		0xa22c: "ĉ",
		0xa22d: "Ċ",
		0xa22e: "ċ",
		0xa22f: "Č",
		0xa230: "č",
		0xa231: "Ď",
		0xa232: "ď",
		0xa233: "Đ",
		0xa234: "đ",
		0xa235: "Ē",
		0xa236: "ē",
		0xa237: "Ĕ",
		0xa238: "ĕ",
		0xa239: "Ė",
		0xa23a: "ė",
		0xa23b: "Ę",
		0xa23c: "ę",
		0xa23d: "Ě",
		0xa23e: "ě",
		0xa23f: "Ĝ",
		0xa240: "ĝ",
		0xa241: "Ğ",
		0xa242: "ğ",
		0xa243: "Ġ",
		0xa244: "ġ",
		0xa245: "Ģ",
		0xa246: "ģ",
		0xa247: "Ĥ",
		0xa248: "ĥ",
		0xa249: "Ħ",
		0xa24a: "ħ",
		0xa24b: "Ĩ",
		0xa24c: "ĩ",
		0xa24d: "Ī",
		0xa24e: "ī",
		0xa24f: "Ĭ",
		0xa250: "ĭ",
		0xa251: "Į",
		0xa252: "į",
		0xa253: "İ",
		0xa254: "ı",
		0xa255: "Ĳ",
		0xa256: "ĳ",
		0xa257: "Ĵ",
		0xa258: "ĵ",
		0xa259: "Ķ",
		0xa25a: "ķ",
		0xa25b: "ĸ",
		0xa25c: "Ĺ",
		0xa25d: "ĺ",
		0xa25e: "Ļ",
		0xa25f: "ļ",
		0xa260: "Ľ",
		0xa261: "ľ",
		0xa262: "Ŀ",
		0xa263: "ŀ",
		0xa264: "Ł",
		0xa265: "ł",
		0xa266: "Ń",
		0xa267: "ń",
		0xa268: "Ņ",
		0xa269: "ņ",
		0xa26a: "Ň",
		0xa26b: "ň",
		0xa26c: "ŉ",
		0xa26d: "Ŋ",
		0xa26e: "ŋ",
		0xa26f: "Ō",
		0xa270: "ō",
		0xa271: "Ŏ",
		0xa272: "ŏ",
		0xa273: "Ő",
		0xa274: "ő",
		0xa275: "Œ",
		0xa276: "œ",
		0xa277: "Ŕ",
		0xa278: "ŕ",
		0xa279: "Ŗ",
		0xa27a: "ŗ",
		0xa27b: "Ř",
		0xa27c: "ř",
		0xa27d: "Ś",
		0xa27e: "ś",
		0xa321: "Ŝ",
		0xa322: "ŝ",
		0xa323: "Ş",
		0xa324: "ş",
		0xa325: "Š",
		0xa326: "š",
		0xa327: "Ţ",
		0xa328: "ţ",
		0xa329: "Ť",
		0xa32a: "ť",
		0xa32b: "Ŧ",
		0xa32c: "ŧ",
		0xa32d: "Ũ",
		0xa32e: "ũ",
		0xa32f: "Ū",
		0xa330: "ū",
		0xa331: "Ŭ",
		0xa332: "ŭ",
		0xa333: "Ů",
		0xa334: "ů",
		0xa335: "Ű",
		0xa336: "ű",
		0xa337: "Ų",
		0xa338: "ų",
		0xa339: "Ŵ",
		0xa33a: "ŵ",
		0xa33b: "Ŷ",
		0xa33c: "ŷ",
		0xa33d: "Ÿ",
		0xa33e: "Ź",
		0xa33f: "ź",
		0xa340: "Ż",
		0xa341: "ż",
		0xa342: "Ž",
		0xa343: "ž",
		0xa344: "ſ",
		0xa34d: "ƒ",
		0xa34e: "ˆ",
		0xa34f: "˜",
	}
}

func (*daijisenExtractor) getFontWide() map[int]string {
	return map[int]string{
		0xb322: "㋘",
		0xb323: "㋙",
		0xb324: "㋚",
		0xb325: "㋛",
		0xb326: "㋜",
		0xb327: "㋝",
		0xb424: "↔",
		0xb646: "㋐",
		0xb647: "㋑",
		0xb648: "㋒",
		0xb649: "㋓",
		0xb64a: "㋔",
		0xb64b: "㋕",
		0xb64c: "㋖",
		0xb64d: "㋗",
		0xb852: "⇒",
		0xbc2c: "･",
		0xc36e: "❶",
		0xc36f: "❷",
		0xc370: "❸",
		0xc371: "❹",
		0xc372: "❺",
		0xc373: "①",
		0xc374: "②",
		0xc375: "③",
		0xc376: "④",
		0xc377: "⑤",
		0xc378: "⑥",
		0xc379: "⑦",
		0xc37a: "⑧",
		0xc37b: "⑨",
		0xc37c: "⑩",
		0xc37d: "⑪",
		0xc37e: "⑫",
		0xc421: "⑬",
		0xc422: "⑭",
		0xc423: "⑮",
		0xc424: "⑯",
		0xc425: "⑰",
		0xc426: "⑱",
		0xc427: "⑲",
		0xc428: "⑳",
		0xc429: "㉑",
		0xc42a: "㉒",
		0xc42b: "㉓",
		0xc42c: "㉔",
		0xc42d: "㉕",
		0xc431: "Ⅰ",
		0xc432: "Ⅱ",
		0xc437: "㊀",
		0xc438: "㊁",
		0xc439: "㊂",
		0xc43a: "㊃",
		0xc43b: "㊄",
		0xc43c: "㊅",
		0xc43d: "㊆",
		0xc43e: "㊇",
		0xc43f: "㊈",
		0xc440: "㉖",
		0xc441: "㉗",
		0xc442: "㉘",
		0xc443: "㉙",
		0xc444: "㉚",
		0xc445: "㉛",
		0xc446: "㉜",
		0xc447: "㉜",
		0xc448: "㉝",
		0xc449: "㉞",
		0xc44a: "㉟",
		0xc455: "[",
		0xc463: "[",
		0xc464: "[",
		0xc465: "♪",
	}
}
