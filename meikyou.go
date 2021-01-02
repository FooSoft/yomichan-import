/*
 * Copyright (c) 2017-2021 Alex Yatskov <alex@foosoft.net>, ajyliew
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

package yomichan

import (
	"regexp"
	"strings"

	zig "github.com/FooSoft/zero-epwing-go"
)

type meikyouExtractor struct {
	partsExp          *regexp.Regexp
	expForeignMetaExp *regexp.Regexp
	expShapesExp      *regexp.Regexp
	expBracketedExp   *regexp.Regexp
	expTermsExp       *regexp.Regexp
	readGroupExp      *regexp.Regexp
	metaExp           *regexp.Regexp
}

func makeMeikyouExtractor() epwingExtractor {
	var foreignMeta = []string{
		"和製",
		"中国",
		"朝鮮",
		"イタリア",
		"スペイン",
		"ドイツ",
		"フランス",
		"オランダ",
		"ポルトガル",
		"ギリシア",
		"アラビア",
		"チベット",
		"タガログ",
		"ヘブライ",
		"ヒンディー",
		"マレーシア",
		"ラテン",
		"アフリカーンス",
		"ロシア",
		"ハワイ",
		"マレー",
		"スウェーデン",
		"ノルウェー・デンマーク",
		"フィンランド",
		"サンスクリット",
		"ポーランド",
	}
	return &meikyouExtractor{
		partsExp:          regexp.MustCompile(`([^（【〖[]+)(?:【(.*)】)?(?:\[(.*)\])?(?:（(.*)）)?`),
		expForeignMetaExp: regexp.MustCompile(strings.Join(foreignMeta, "|")),
		expShapesExp:      regexp.MustCompile(`[▼▽]+`),
		expBracketedExp:   regexp.MustCompile(`(?:[〈《])([^〉》]*)(?:[〉》])`),
		expTermsExp:       regexp.MustCompile(`([^（]*)?(?:（(.*)）)?`),
		readGroupExp:      regexp.MustCompile(`[‐・]+`),
		metaExp:           regexp.MustCompile(`〘([^〙]*)〙`),
	}
}

func (e *meikyouExtractor) extractTerms(entry zig.BookEntry, sequence int) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	var expressions, readings []string
	if expression := matches[2]; len(expression) > 0 {
		expression = e.expShapesExp.ReplaceAllLiteralString(expression, "")
		expression = e.expBracketedExp.ReplaceAllString(expression, "$1")
		if termsMatches := e.expTermsExp.FindStringSubmatch(expression); termsMatches != nil {
			termsMatches = termsMatches[1:]
			for _, terms := range termsMatches {
				if len(terms) > 0 {
					for _, split := range strings.Split(terms, "・") {
						expressions = append(expressions, split)
					}
				}
			}
		}
	}

	if expression := matches[3]; len(expression) > 0 {
		expression = e.expForeignMetaExp.ReplaceAllLiteralString(expression, "")
		expression = strings.Replace(expression, "＋", " ", -1)
		for _, split := range strings.Split(expression, "・") {
			expressions = append(expressions, split)
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

func (e *meikyouExtractor) extractKanji(entry zig.BookEntry) []dbKanji {
	return nil
}

func (e *meikyouExtractor) exportRules(term *dbTerm, tags []string) {
	for _, tag := range tags {
		if tag == "名" {
			term.addRules("n")
		} else if tag == "代" {
			term.addRules("pn")
		} else if tag == "連体" {
			term.addRules("adj-pn")
		} else if tag == "副" {
			term.addRules("adv")
		} else if tag == "副ト" || tag == "副トニ" || tag == "トニ" {
			term.addRules("adv-to")
		} else if tag == "副助" {
			term.addRules("adv")
			term.addRules("prt")
		} else if tag == "格助" || tag == "終助" {
			term.addRules("prt")
		} else if tag == "接" {
			term.addRules("conj")
		} else if tag == "接助" {
			term.addRules("conj")
			term.addRules("prt")
		} else if tag == "接尾" {
			term.addRules("suf")
		} else if tag == "接頭" {
			term.addRules("pref")
		} else if tag == "補形" {
			term.addRules("aux-adj")
		} else if strings.HasPrefix(tag, "助動") || strings.HasPrefix(tag, "補動") {
			term.addRules("aux-v")
		} else if tag == "形動トタル" {
			term.addRules("adj-t")
			term.addRules("adv-to")
		} else if tag == "形" {
			term.addRules("adj-i")
		} else if tag == "形動" {
			term.addRules("adj-na")
		}

		if strings.Contains(tag, "他") {
			term.addRules("vt")
		}
		if strings.Contains(tag, "自") {
			term.addRules("vi")
		}
		if strings.Contains(tag, "一") {
			term.addRules("v1")
		} else if strings.ContainsAny(tag, "二四五") && tag != "二" && tag != " トニ" {
			term.addRules("v5")
		} else if strings.Contains(tag, "サ変") && (strings.HasSuffix(term.Expression, "する") || strings.HasSuffix(term.Expression, "為る")) {
			term.addRules("vs")
		} else if term.Expression == "来る" {
			term.addRules("vk")
		}
	}
}

func (*meikyouExtractor) getRevision() string {
	return "meikyou1"
}

func (*meikyouExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		41249: " ",
		41250: "¡",
		41251: "¢",
		41252: "£",
		41253: "¤",
		41254: "¥",
		41255: "¦",
		41256: "§",
		41257: "¨",
		41258: "©",
		41259: "ª",
		41260: "«",
		41261: "¬",
		41262: "­",
		41263: "®",
		41264: "¯",
		41265: "°",
		41266: "±",
		41267: "²",
		41268: "³",
		41269: "´",
		41270: "µ",
		41271: "¶",
		41272: "·",
		41273: "¸",
		41274: "¹",
		41275: "º",
		41276: "»",
		41277: "¼",
		41278: "½",
		41279: "¾",
		41280: "¿",
		41281: "À",
		41282: "Á",
		41283: "Â",
		41284: "Ã",
		41285: "Ä",
		41286: "Å",
		41287: "Æ",
		41288: "Ç",
		41289: "È",
		41290: "É",
		41291: "Ê",
		41292: "Ë",
		41293: "Ì",
		41294: "Í",
		41295: "Î",
		41296: "Ï",
		41297: "Ð",
		41298: "Ñ",
		41299: "Ò",
		41300: "Ó",
		41301: "Ô",
		41302: "Õ",
		41303: "Ö",
		41304: "×",
		41305: "Ø",
		41306: "Ù",
		41307: "Ú",
		41308: "Û",
		41309: "Ü",
		41310: "Ý",
		41311: "Þ",
		41312: "ß",
		41313: "à",
		41314: "á",
		41315: "â",
		41316: "ã",
		41317: "ä",
		41318: "å",
		41319: "æ",
		41320: "ç",
		41321: "è",
		41322: "é",
		41323: "ê",
		41324: "ë",
		41325: "ì",
		41326: "í",
		41327: "î",
		41328: "ï",
		41329: "ð",
		41330: "ñ",
		41331: "ò",
		41332: "ó",
		41333: "ô",
		41334: "õ",
		41335: "ö",
		41336: "÷",
		41337: "ø",
		41338: "ù",
		41339: "ú",
		41340: "û",
		41341: "ü",
		41342: "ý",
		41505: "þ",
		41506: "ÿ",
		41507: "Ā",
		41508: "ā",
		41509: "Ă",
		41510: "ă",
		41511: "Ą",
		41512: "ą",
		41513: "Ć",
		41514: "ć",
		41515: "Ĉ",
		41516: "ĉ",
		41517: "Ċ",
		41518: "ċ",
		41519: "Č",
		41520: "č",
		41521: "Ď",
		41522: "ď",
		41523: "Đ",
		41524: "đ",
		41525: "Ē",
		41526: "ē",
		41527: "Ĕ",
		41528: "ĕ",
		41529: "Ė",
		41530: "ė",
		41531: "Ę",
		41532: "ę",
		41533: "Ě",
		41534: "ě",
		41535: "Ĝ",
		41536: "ĝ",
		41537: "Ğ",
		41538: "ğ",
		41539: "Ġ",
		41540: "ġ",
		41541: "Ģ",
		41542: "ģ",
		41543: "Ĥ",
		41544: "ĥ",
		41545: "Ħ",
		41546: "ħ",
		41547: "Ĩ",
		41548: "ĩ",
		41549: "Ī",
		41550: "ī",
		41551: "Ĭ",
		41552: "ĭ",
		41553: "Į",
		41554: "į",
		41555: "İ",
		41556: "ı",
		41557: "Ĳ",
		41558: "ĳ",
		41559: "Ĵ",
		41560: "ĵ",
		41561: "Ķ",
		41562: "ķ",
		41563: "ĸ",
		41564: "Ĺ",
		41565: "ĺ",
		41566: "Ļ",
		41567: "ļ",
		41568: "Ľ",
		41569: "ľ",
		41570: "Ŀ",
		41571: "ŀ",
		41572: "Ł",
		41573: "ł",
		41574: "Ń",
		41575: "ń",
		41576: "Ņ",
		41577: "ņ",
		41578: "Ň",
		41579: "ň",
		41580: "ŉ",
		41581: "Ŋ",
		41582: "ŋ",
		41583: "Ō",
		41584: "ō",
		41585: "Ŏ",
		41586: "ŏ",
		41587: "Ő",
		41588: "ő",
		41589: "Œ",
		41590: "œ",
		41591: "Ŕ",
		41592: "ŕ",
		41593: "Ŗ",
		41594: "ŗ",
		41595: "Ř",
		41596: "ř",
		41597: "Ś",
		41598: "ś",
		41761: "Ŝ",
		41762: "ŝ",
		41763: "Ş",
		41764: "ş",
		41765: "Š",
		41766: "š",
		41767: "Ţ",
		41768: "ţ",
		41769: "Ť",
		41770: "ť",
		41771: "Ŧ",
		41772: "ŧ",
		41773: "Ũ",
		41774: "ũ",
		41775: "Ū",
		41776: "ū",
		41777: "Ŭ",
		41778: "ŭ",
		41779: "Ů",
		41780: "ů",
		41781: "Ű",
		41782: "ű",
		41783: "Ų",
		41784: "ų",
		41785: "Ŵ",
		41786: "ŵ",
		41787: "Ŷ",
		41788: "ŷ",
		41789: "Ÿ",
		41790: "Ź",
		41791: "ź",
		41792: "Ż",
		41793: "ż",
		41794: "Ž",
		41795: "ž",
		41796: "ſ",
		41797: "Ǎ",
		41798: "ǎ",
		41799: "Ǐ",
		41800: "ǐ",
		41801: "Ǒ",
		41802: "ǒ",
		41803: "Ǔ",
		41804: "ǔ",
		41805: "ƒ",
		41806: "ˆ",
		41807: "˜",
		41808: "ɔ",
		41809: "ɔ̀",
		41810: "ɔ́",
		41811: "ǝ",
		41812: "ǝ̀",
		41813: "ǝ́",
		41814: "ʌ",
		41815: "ʌ̀",
		41816: "ʌ́",
		41817: "",
		41818: "ɑ",
		41819: "ɑ̀",
		41820: "ɑ́",
		41821: "ʃ",
		41822: "ʊ",
		41823: "θ",
		41824: "ʒ",
		41825: "ɒ",
		41826: "ǽ",
		41827: "ɚ",
		41828: "ɡ",
		41829: "ʤ",
		41830: "ʧ",
		41831: "-",
		41832: ".",
		41833: "¯",
		41834: "℉",
		41835: "Ⅰ",
		41836: "Ⅱ",
		41837: "Ⅲ",
		41838: "Ⅳ",
		41839: "Ⅴ",
		41840: "Ⅹ",
		41841: "↕",
		41842: "■",
		41843: "°",
		41844: "∛",
		41845: "∜",
		41846: "∥",
		41847: "〻",
		41848: "≣",
		41849: "≺",
		41850: "≻",
		41851: "∧",
		41852: "",
		41853: "♠",
		41854: "♣",
		42017: "♥",
		42018: "♦",
		42019: "♩",
		42020: "♮",
		42021: "√",
	}
}

func (*meikyouExtractor) getFontWide() map[int]string {
	return map[int]string{
		45089: "鄧",
		45090: "疒",
		45091: "©",
		45092: "æ",
		45093: "æ̀",
		45094: "ǽ",
		45095: "①",
		45096: "②",
		45097: "③",
		45098: "④",
		45099: "⑤",
		45100: "⑥",
		45101: "⑦",
		45102: "⑧",
		45103: "⑨",
		45104: "⑩",
		45105: "⑪",
		45106: "⑫",
		45107: "⑬",
		45108: "⑭",
		45109: "⑮",
		45110: "⑯",
		45111: "⑰",
		45112: "⑱",
		45113: "⑲",
		45114: "⑳",
		45115: "⑴",
		45116: "⑵",
		45117: "⑶",
		45118: "〘",
		45119: "〙",
		45120: "＼",
		45121: "／",
		45122: "㋐",
		45123: "㋑",
		45124: "㋒",
		45125: "㋓",
		45126: "㋔",
		45127: "㋕",
		45128: "㋖",
		45129: "㋗",
		45130: "㋘",
		45131: "㋙",
		45132: "㋚",
		45133: "㋛",
		45134: "㋜",
		45135: "㋝",
		45136: "㋞",
		45137: "㋟",
		45138: "㋠",
		45139: "㋡",
		45140: "㋢",
		45141: "㋣",
		45142: "丰",
		45143: "仐",
		45144: "你",
		45145: "俏",
		45146: "俠",
		45147: "偓",
		45148: "儈",
		45149: "",
		45150: "厴",
		45151: "呍",
		45152: "啞",
		45153: "嘻",
		45154: "噦",
		45155: "噯",
		45156: "嚙",
		45157: "嚢",
		45158: "埵",
		45159: "塡",
		45160: "增",
		45161: "壔",
		45162: "妤",
		45163: "婟",
		45164: "孒",
		45165: "尩",
		45166: "屢",
		45167: "弴",
		45168: "彽",
		45169: "德",
		45170: "憍",
		45171: "扌",
		45172: "挍",
		45173: "挘",
		45174: "挵",
		45175: "捥",
		45176: "搔",
		45177: "摑",
		45178: "撿",
		45179: "擊",
		45180: "擤",
		45181: "攙",
		45182: "攩",
		45345: "昻",
		45346: "晳",
		45347: "枘",
		45348: "栱",
		45349: "桛",
		45350: "梂",
		45351: "梘",
		45352: "梣",
		45353: "梲",
		45354: "梻",
		45355: "棰",
		45356: "楉",
		45357: "楤",
		45358: "榨",
		45359: "樏",
		45360: "樝",
		45361: "橅",
		45362: "橐",
		45363: "橫",
		45364: "檝",
		45365: "檞",
		45366: "櫧",
		45367: "氵",
		45368: "洄",
		45369: "湑",
		45370: "潑",
		45371: "濹",
		45372: "瀆",
		45373: "瀨",
		45374: "灬",
		45375: "炷",
		45376: "炻",
		45377: "焰",
		45378: "煆",
		45379: "煠",
		45380: "熅",
		45381: "牓",
		45382: "玕",
		45383: "瑇",
		45384: "疒",
		45385: "痀",
		45386: "痎",
		45387: "痹",
		45388: "瘙",
		45389: "瘦",
		45390: "瘭",
		45391: "癤",
		45392: "皂",
		45393: "盬",
		45394: "眴",
		45395: "眶",
		45396: "睺",
		45397: "矠",
		45398: "矻",
		45399: "硨",
		45400: "磲",
		45401: "祆",
		45402: "禱",
		45403: "稭",
		45404: "穇",
		45405: "窠",
		45406: "笧",
		45407: "筕",
		45408: "篊",
		45409: "篖",
		45410: "簎",
		45411: "簶",
		45412: "籡",
		45413: "籹",
		45414: "粑",
		45415: "糈",
		45416: "糗",
		45417: "糝",
		45418: "絇",
		45419: "綠",
		45420: "緖",
		45421: "縕",
		45422: "繇",
		45423: "繡",
		45424: "繫",
		45425: "胳",
		45426: "腭",
		45427: "舢",
		45428: "苆",
		45429: "萁",
		45430: "萊",
		45431: "蒴",
		45432: "蔞",
		45433: "蔣",
		45434: "蔲",
		45435: "蕺",
		45436: "薰",
		45437: "蘞",
		45438: "蘩",
		45601: "虯",
		45602: "蛽",
		45603: "蜱",
		45604: "蜾",
		45605: "蝲",
		45606: "螈",
		45607: "蟎",
		45608: "蟖",
		45609: "蠃",
		45610: "蠆",
		45611: "蠊",
		45612: "蠟",
		45613: "袘",
		45614: "袪",
		45615: "裑",
		45616: "襬",
		45617: "豇",
		45618: "賴",
		45619: "跆",
		45620: "跑",
		45621: "踠",
		45622: "軀",
		45623: "辨",
		45624: "邌",
		45625: "醞",
		45626: "醱",
		45627: "鈸",
		45628: "鎺",
		45629: "雞",
		45630: "韛",
		45631: "頰",
		45632: "顖",
		45633: "顚",
		45634: "顬",
		45635: "飥",
		45636: "餺",
		45637: "駃",
		45638: "騠",
		45639: "驎",
		45640: "骶",
		45641: "魬",
		45642: "魳",
		45643: "鮄",
		45644: "鮧",
		45645: "鮬",
		45646: "鮸",
		45647: "鯁",
		45648: "鯎",
		45649: "鯥",
		45650: "鯧",
		45651: "鰙",
		45652: "鰶",
		45653: "鱁",
		45654: "鱏",
		45655: "鱓",
		45656: "鱝",
		45657: "鱩",
		45658: "鱪",
		45659: "鱮",
		45660: "鱰",
		45661: "鱲",
		45662: "鱵",
		45663: "鵇",
		45664: "鵼",
		45665: "鶀",
		45666: "鷉",
		45667: "鷗",
		45668: "鸊",
		45669: "鹼",
		45670: "麨",
		45671: "麬",
		45672: "麴",
		45673: "黑",
		45674: "鼯",
		45675: "鼹",
		45676: "爛",
		45677: "朗",
		45678: "塚",
		45679: "神",
		45680: "祥",
		45681: "福",
		45682: "﨟",
		45683: "諸",
		45684: "都",
		45685: "-",
		45686: "~",
		45687: "¢",
		45688: "£",
		45689: "〓",
		45690: "〰",
		45691: "㊀",
		45692: "㊁",
		45693: "㊂",
		45694: "㊃",
		45857: "㊙",
		45858: "㋤",
		45859: "懀",
		45860: "杮",
		45861: "〓",
		45862: "",
		45863: "○",
		45864: "",
	}
}
