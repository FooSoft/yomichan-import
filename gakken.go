package yomichan

import (
	"regexp"
	"strings"

	zig "foosoft.net/projects/zero-epwing-go"
)

type gakkenExtractor struct {
	partsExp     *regexp.Regexp
	readGroupExp *regexp.Regexp
	expVarExp    *regexp.Regexp
	metaExp      *regexp.Regexp
	v5Exp        *regexp.Regexp
	v1Exp        *regexp.Regexp
}

func makeGakkenExtractor() epwingExtractor {
	return &gakkenExtractor{
		partsExp:     regexp.MustCompile(`([\p{Hiragana}\p{Katakana}ー‐・]*)?(?:【(.*)】)?`),
		readGroupExp: regexp.MustCompile(`[‐・]+`),
		expVarExp:    regexp.MustCompile(`\(([^\)]*)\)`),
		metaExp:      regexp.MustCompile(`（([^）]*)）`),
		v5Exp:        regexp.MustCompile(`(動.[四五](［[^］]+］)?)|(動..二)`),
		v1Exp:        regexp.MustCompile(`(動..一)`),
	}
}

var cosmetics = strings.NewReplacer("(1)", "①", "(2)", "②", "(3)", "③", "(4)", "④", "(5)", "⑤", "(6)", "⑥", "(7)", "⑦", "(8)", "⑧", "(9)", "⑨", "(10)", "⑩", "(11)", "⑪", "(12)", "⑫", "(13)", "⑬", "(14)", "⑭", "(15)", "⑮", "(16)", "⑯", "(17)", "⑰", "(18)", "⑱", "(19)", "⑲", "(20)", "⑳",
	"カ゛", "ガ",
	"キ゛", "ギ",
	"ク゛", "グ",
	"ケ゛", "ゲ",
	"コ゛", "ゴ",
	"タ゛", "ダ",
	"チ゛", "ヂ",
	"ツ゛", "ヅ",
	"テ゛", "デ",
	"ト゛", "ド",
	"ハ゛", "バ",
	"ヒ゛", "ビ",
	"フ゛", "ブ",
	"ヘ゛", "ベ",
	"ホ゛", "ボ",
	"サ゛", "ザ",
	"シ゛", "ジ",
	"ス゛", "ズ",
	"セ゛", "ゼ",
	"ソ゛", "ゾ")

func (e *gakkenExtractor) extractTerms(entry zig.BookEntry, sequence int) []dbTerm {
	matches := e.partsExp.FindStringSubmatch(entry.Heading)
	if matches == nil {
		return nil
	}

	var expressions, readings []string
	if expression := matches[2]; len(expression) > 0 {
		expression = e.metaExp.ReplaceAllLiteralString(expression, "")
		for _, split := range regexp.MustCompile("(・|】【)").Split(expression, -1) {
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

	entryText := cosmetics.Replace(entry.Text)

	for _, split := range strings.Split(entryText, "\n") {
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
				Glossary:   []string{entryText},
				Sequence:   sequence,
			}

			e.exportRules(&term, tags)
			terms = append(terms, term)
		}

	} else {
		if len(readings) == 0 {
			readings = append(readings, "")
		}
		for _, expression := range expressions {
			for _, reading := range readings {
				term := dbTerm{
					Expression: expression,
					Reading:    reading,
					Glossary:   []string{entryText},
					Sequence:   sequence,
				}

				e.exportRules(&term, tags)
				terms = append(terms, term)
			}
		}
	}

	return terms
}

func (*gakkenExtractor) extractKanji(entry zig.BookEntry) []dbKanji {
	return nil
}

func (e *gakkenExtractor) exportRules(term *dbTerm, tags []string) {
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

func (*gakkenExtractor) getRevision() string {
	return "gakken"
}

func (*gakkenExtractor) getFontNarrow() map[int]string {
	return map[int]string{
		41550: "ī",
	}
}

func (*gakkenExtractor) getFontWide() map[int]string {
	return map[int]string{
		42017: "国",
		42018: "古",
		42019: "故",
		42021: "(拡)",
		42020: "漢",
		42033: "",
		42034: "",
		42070: "㋐",
		42071: "㋑",
		42072: "㋒",
		42073: "㋓",
		42074: "㋔",
		42075: "㋕",
		42076: "㋖",
		42077: "㋗",
		42078: "㋘",
		42079: "㋙",
		42080: "㋚",
		42081: "㋛",
		42082: "㋜",
		42083: "㋝",
		42084: "🈩",
		42085: "🈔",
		42086: "🈪",
		42087: "[四]",
		42088: "[五]",
		42089: "❶",
		42090: "❷",
		42091: "❸",
		42092: "❹",
		42093: "❺",
		42094: "❻",
		42095: "❼",
		42096: "❽",
		42097: "❾",
		42098: "❿",
		42099: "⓫",
		42100: "⓬",
		42101: "⓭",
		42102: "⓮",
		42103: "⓯",
		42104: "⓰",
		42105: "⓱",
		42106: "⓲",
		42107: "㊀",
		42108: "㊁",
		42109: "㊂",
		42110: "㊃",
		43599: "咍",
		46176: "(扌)",
		48753: "灾",
		48936: "烖",
		58176: "(呉)",
		58177: "(漢)",
	}
}
