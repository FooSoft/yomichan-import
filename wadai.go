/*
 * Copyright (c) 2017 Alex Yatskov <alex@foosoft.net>, ajyliew
 * Author: Alex Yatskov <alex@foosoft.net>, ajyliew
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

type wadaiExtractor struct {
	partsExp        *regexp.Regexp
	literalPartsExp *regexp.Regexp
	readPartsExp    *regexp.Regexp
	quotedExp       *regexp.Regexp
	alphaExp        *regexp.Regexp
}

func makeWadaiExtractor() epwingExtractor {
	return &wadaiExtractor{
		partsExp:        regexp.MustCompile(`([^＜]+)(?:＜([^＞【]+)(?:【([^】]+)】)?＞)?`),
		literalPartsExp: regexp.MustCompile(`(¶)?(.*)`),
		readPartsExp:    regexp.MustCompile(`([^１２３４５６７８９０]+)(.*)`),
		quotedExp:       regexp.MustCompile(`「?([^」]+)`),
		alphaExp:        regexp.MustCompile(`[a-z]+`),
	}
}

func (e *wadaiExtractor) extractTerms(entry epwingEntry, sequence int) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	preset := false
	literal := matches[1]
	if literalMatches := e.literalPartsExp.FindStringSubmatch(literal); literalMatches != nil {
		preset = len(literalMatches[1]) > 0
		literal = literalMatches[2]
	}

	reading := matches[2]
	if readMatches := e.readPartsExp.FindStringSubmatch(reading); readMatches != nil {
		reading = readMatches[1]
	}

	expressions := strings.Split(matches[3], "・")
	if len(expressions) == 0 {
		expressions = append(expressions, "")
	}

	var terms []dbTerm
	for _, expression := range expressions {
		if preset {
			expression = literal
			reading = ""
		} else if len(expression) == 0 {
			expression = literal
		}

		if quotedMatches := e.quotedExp.FindStringSubmatch(reading); quotedMatches != nil {
			reading = quotedMatches[1]
		}

		if alphaMatches := e.alphaExp.FindStringSubmatch(expression); alphaMatches != nil && len(reading) > 0 {
			expression = reading
			reading = ""
		}

		term := dbTerm{
			Expression: expression,
			Reading:    reading,
			Glossary:   []string{entry.Text},
			Sequence:   sequence,
		}

		terms = append(terms, term)
	}

	return terms
}

func (e *wadaiExtractor) extractKanji(entry epwingEntry) []dbKanji {
	return nil
}

func (*wadaiExtractor) getRevision() string {
	return "wadai1"
}

func (*wadaiExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		41267: "﹢",
		41269: "*",
		41270: "ᐦ",
		41284: "Á",
		41285: "É",
		41287: "Ó",
		41288: "Ú",
		41290: "á",
		41291: "é",
		41292: "í",
		41293: "ó",
		41294: "ú",
		41295: "ý",
		41313: "À",
		41314: "È",
		41319: "à",
		41320: "è",
		41321: "ì",
		41322: "ò",
		41323: "ù",
		41505: "Ö",
		41506: "Ü",
		41508: "ä",
		41509: "ë",
		41510: "ï",
		41511: "ö",
		41512: "ü",
		41513: "ÿ",
		41515: "Â",
		41516: "Ê",
		41517: "Î",
		41520: "â",
		41521: "ê",
		41522: "î",
		41523: "ô",
		41524: "û",
		41525: "ā",
		41526: "ē",
		41527: "ī",
		41528: "ō",
		41529: "ū",
		41530: "ȳ",
		41532: "Ç",
		41533: "ç",
		41534: "ɘ́",
		41538: "ɔ́",
		41561: "˜",
		41566: "ã",
		41567: "ñ",
		41581: "ʌ",
		41582: "ø",
		41583: "ə",
		41585: "ε",
		41587: "ɔ",
		41588: "℧",
		41590: "ð",
		41593: "ŋ",
		41594: "ː",
		41596: "Ø",
		41762: "\\",
		41768: "˘",
		41773: "Ŭ",
		41775: "ă",
		41776: "ĕ",
		41777: "ğ",
		41778: "ĭ",
		41779: "ŏ",
		41780: "ŭ",
		41784: "Č",
		41788: "Š",
		41791: "č",
		41792: "ě",
		41794: "ň",
		41795: "ř",
		41796: "š",
		41797: "ž",
		41804: "ą",
		41805: "ę",
		41811: "ș",
		41812: "ț",
		41822: "Ś",
		41823: "ć",
		41824: "ń",
		41825: "ś",
		41826: "ź",
		42061: "‘",
		42063: "Ł",
		42068: "ł",
		42071: "õ",
		42075: "Å",
		42076: "å",
		42077: "ů",
		42081: "Ḥ",
		42089: "ḍ",
		42090: "ḥ",
		42092: "ṃ",
		42093: "ṇ",
		42095: "ṣ",
		42102: "İ",
		42104: "Ż",
		42109: "ṅ",
		42287: "‴",
		42316: "Ō",
		42322: "b̄",
		42324: "d̅",
		42325: "h̄",
		42327: "s̅",
		42330: "z̅",
		42344: "〚",
		42345: "〛",
		42356: "ǔ",
		42357: "ż",
		42358: "Ž",
		42359: "ž",
	}
}

func (*wadaiExtractor) getFontWide() map[int]string {
	return map[int]string{
		45380: "☞",
		45397: "æ",
		45402: "œ",
		45406: "Æ",
		45429: "©",
		45613: "<",
		45614: ">",
		45629: "┏",
		45653: "⛤",
		45662: "嗉",
		45665: "圳",
		45666: "拼",
		45667: "攩",
		45671: "烤",
		45673: "玢",
		45674: "癤",
		45675: "皶",
		45676: "磠",
		45677: "稃",
		45681: "蔲",
		45684: "顬",
		45685: "骶",
		45689: "榍",
		45857: "倻",
		45870: "噯",
		45876: "垜",
		45898: "愷",
		45900: "擤",
		45906: "晷",
		45909: "枘",
		45910: "不",
		45913: "楣",
		45916: "梲",
		45919: "桛",
		45921: "楤",
		45922: "橅",
		45923: "檉",
		45933: "淄",
		46125: "煆",
		46135: "珅",
		46137: "琛",
		46141: "痤",
		46142: "癭",
		46143: "瘭",
		46152: "窠",
		46154: "笯",
		46155: "筠",
		46156: "簎",
		46157: "糝",
		46161: "翟",
		46163: "翮",
		46166: "腊",
		46168: "舢",
		46169: "芷",
		46177: "蒴",
		46181: "蕙",
		46190: "蚉",
		46191: "蝲",
		46197: "豇",
		46198: "跑",
		46200: "跗",
		46201: "跆",
		46202: "蒁",
		46372: "鄱",
		46374: "鄧",
		46388: "卍",
		46390: "𨫤",
		46391: "鈹",
		46398: "顥",
		46404: "駃",
		46405: "騠",
		46406: "髁",
		46409: "魳",
		46410: "鱏",
		46411: "鱓",
		46414: "鱮",
		46415: "鰶",
		46416: "魬",
		46417: "𩸽",
		46418: "鯥",
		46419: "鰙",
		46422: "鮄",
		46423: "鱵",
		46424: "鷴",
		46425: "鶍",
		46426: "鵟",
		46428: "鼯",
		46449: "▶",
		46459: "㧍",
		46460: "嘈",
		46461: "愈",
		46462: "淝",
		46634: "灤",
		46635: "焮",
		46636: "獮",
		46637: "瓚",
		46638: "絓",
		46639: "芎",
		46650: "薏",
		46651: "辶",
		46652: "醞",
		46653: "挵",
		46654: "飥",
		46655: "鬐",
		46656: "俏",
		46657: "啐",
		46658: "塼",
		46659: "濰",
		46660: "磲",
		46661: "篊",
		46662: "菀",
		46663: "芩",
		46664: "𧿹",
		46665: "鈸",
		46666: "驎",
		46667: "硨",
		46668: "蘞",
		46669: "梣",
		46670: "槵",
		46671: "橉",
		46672: "莧",
		46682: "彔",
		46683: "噦",
		46684: "袘",
		46685: "餺",
		46686: "►",
		46688: "棈",
		46689: "▷",
		46695: "[ローマ字]",
		46699: "◧",
		46700: "◨",
	}
}
