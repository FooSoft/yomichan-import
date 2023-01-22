package yomichan

import (
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"

	"foosoft.net/projects/jmdict"
	"golang.org/x/exp/slices"
)

type headword struct {
	Expression   string
	Reading      string
	TermTags     []string
	Index        int
	IsPriority   bool
	IsIrregular  bool
	IsOutdated   bool
	IsRareKanji  bool
	IsSearchOnly bool
	IsAteji      bool
	IsGikun      bool
}

type hash uint64

func (h *headword) Hash() hash {
	return hashText(h.Expression + "␞" + h.Reading)
}

func (h *headword) ExpHash() hash {
	return hashText(h.Expression + "␞" + h.Expression)
}

func (h *headword) ReadingHash() hash {
	return hashText(h.Reading + "␞" + h.Reading)
}

func hashText(s string) hash {
	h := fnv.New64a()
	h.Write([]byte(s))
	return hash(h.Sum64())
}

func (h *headword) IsKanaOnly() bool {
	if h.Expression != h.Reading {
		return false
	}
	for _, char := range h.Expression {
		if char >= 'ぁ' && char <= 'ヿ' {
			// hiragana and katakana range
			continue
		} else if char >= '･' && char <= 'ﾟ' {
			// halfwidth katakana range
			continue
		} else if char == '〜' {
			continue
		} else {
			return false
		}
	}
	return true
}

func (h *headword) Score() int {
	score := 0
	if h.IsPriority {
		score += 1
	}
	if h.IsIrregular {
		score -= 5
	}
	if h.IsOutdated {
		score -= 5
	}
	if h.IsRareKanji {
		score -= 5
	}
	if h.IsSearchOnly {
		score -= 5
	}
	return score
}

func (h *headword) ToInternalLink(includeReading bool) any {
	if !includeReading || h.Expression == h.Reading {
		return contentInternalLink(
			contentAttr{lang: ISOtoHTML["jpn"]},
			h.Expression,
		)
	} else {
		return contentSpan(
			contentAttr{lang: ISOtoHTML["jpn"]},
			contentInternalLink(contentAttr{}, h.Expression),
			"（",
			contentInternalLink(contentAttr{}, h.Reading),
			"）",
		)
	}
}

func (h *headword) SetFlags(infoTags, freqTags []string) {
	priorityTags := []string{"ichi1", "news1", "gai1", "spec1", "spec2"}
	for _, priorityTag := range priorityTags {
		if slices.Contains(freqTags, priorityTag) {
			h.IsPriority = true
			break
		}
	}
	for _, infoTag := range infoTags {
		switch infoTag {
		case "iK", "ik", "io":
			h.IsIrregular = true
		case "oK", "ok":
			h.IsOutdated = true
		case "sK", "sk":
			h.IsSearchOnly = true
		case "rK":
			h.IsRareKanji = true
		case "ateji":
			h.IsAteji = true
		case "gikun":
			h.IsGikun = true
		}
	}
	if h.IsOutdated && h.IsRareKanji {
		h.IsRareKanji = false
	}
}

func (h *headword) SetTermTags(freqTags []string) {
	h.TermTags = []string{}
	if h.IsPriority {
		h.TermTags = append(h.TermTags, priorityTagName)
	}
	for _, tag := range freqTags {
		isNewsFreqTag, _ := regexp.MatchString(`nf\d\d`, tag)
		if isNewsFreqTag {
			// nf tags are divided into ranks of 500
			// (nf01 to nf48), but it will be easier
			// for the user to read 1k, 2k, etc.
			var i int
			if _, err := fmt.Sscanf(tag, "nf%2d", &i); err == nil {
				i = (i + (i % 2)) / 2
				newsTag := "news" + strconv.Itoa(i) + "k"
				h.TermTags = append(h.TermTags, newsTag)
			}
		} else if tag == "news1" || tag == "news2" {
			continue
		} else {
			tagWithoutTheNumber := tag[:len(tag)-1] // "ichi", "gai", or "spec"
			h.TermTags = append(h.TermTags, tagWithoutTheNumber)
		}
	}
	if h.IsIrregular {
		h.TermTags = append(h.TermTags, irregularTagName)
	}
	if h.IsOutdated {
		h.TermTags = append(h.TermTags, outdatedTagName)
	}
	if h.IsRareKanji {
		h.TermTags = append(h.TermTags, rareKanjiTagName)
	}
	if h.IsAteji {
		h.TermTags = append(h.TermTags, atejiTagName)
	}
	if h.IsGikun {
		h.TermTags = append(h.TermTags, gikunTagName)
	}
}

func newHeadword(kanji *jmdict.JmdictKanji, reading *jmdict.JmdictReading) headword {
	h := headword{}
	infoTags := []string{}
	freqTags := []string{}
	if kanji == nil {
		h.Expression = reading.Reading
		h.Reading = reading.Reading
		infoTags = reading.Information
		freqTags = reading.Priorities
	} else if reading == nil {
		// should only apply to search-only kanji terms
		h.Expression = kanji.Expression
		h.Reading = ""
		infoTags = kanji.Information
		freqTags = kanji.Priorities
	} else {
		h.Expression = kanji.Expression
		h.Reading = reading.Reading
		infoTags = union(kanji.Information, reading.Information)
		freqTags = intersection(kanji.Priorities, reading.Priorities)
	}
	h.SetFlags(infoTags, freqTags)
	h.SetTermTags(freqTags)
	return h
}

func areAllKanjiIrregular(allKanji []jmdict.JmdictKanji) bool {
	// If every kanji form is rare or irregular, then we'll make
	// kana-only headwords for each kana form.
	if len(allKanji) == 0 {
		return false
	}
	for _, kanji := range allKanji {
		h := newHeadword(&kanji, nil)
		kanjiIsIrregular := h.IsRareKanji || h.IsIrregular || h.IsOutdated || h.IsSearchOnly
		if !kanjiIsIrregular {
			return false
		}
	}
	return true
}

func extractHeadwords(entry jmdict.JmdictEntry) []headword {
	headwords := []headword{}
	allKanjiAreIrregular := areAllKanjiIrregular(entry.Kanji)

	if allKanjiAreIrregular {
		// Adding the reading-only terms before kanji+reading
		// terms here for the sake of the Index property,
		// which affects the yomichan term ranking.
		for _, reading := range entry.Readings {
			h := newHeadword(nil, &reading)
			h.Index = len(headwords)
			headwords = append(headwords, h)
		}
	}

	for _, kanji := range entry.Kanji {
		if slices.Contains(kanji.Information, "sK") {
			// Search-only kanji forms do not have associated readings.
			h := newHeadword(&kanji, nil)
			h.Index = len(headwords)
			headwords = append(headwords, h)
			continue
		}
		for _, reading := range entry.Readings {
			if reading.NoKanji != nil {
				continue
			} else if slices.Contains(reading.Information, "sk") {
				// Search-only kana forms do not have associated kanji forms.
				continue
			} else if reading.Restrictions != nil && !slices.Contains(reading.Restrictions, kanji.Expression) {
				continue
			} else {
				h := newHeadword(&kanji, &reading)
				h.Index = len(headwords)
				headwords = append(headwords, h)
			}
		}
	}

	if !allKanjiAreIrregular {
		noKanjiInEntry := (len(entry.Kanji) == 0)
		for _, reading := range entry.Readings {
			if reading.NoKanji != nil || noKanjiInEntry || slices.Contains(reading.Information, "sk") {
				h := newHeadword(nil, &reading)
				h.Index = len(headwords)
				headwords = append(headwords, h)
			}
		}
	}

	return headwords
}
