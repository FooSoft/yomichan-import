package yomichan

import (
	"os"
	"strings"

	"foosoft.net/projects/jmdict"
	"golang.org/x/exp/slices"
)

func kata2hira(word string) string {
	charMap := func(character rune) rune {
		if (character >= 'ァ' && character <= 'ヶ') || (character >= 'ヽ' && character <= 'ヾ') {
			return character - 0x60
		} else {
			return character
		}
	}
	return strings.Map(charMap, word)
}

func (h *headword) InfoSymbols() string {
	infoSymbols := []string{}
	if h.IsPriority {
		infoSymbols = append(infoSymbols, prioritySymbol)
	}
	if h.IsRareKanji {
		infoSymbols = append(infoSymbols, rareKanjiSymbol)
	}
	if h.IsIrregular {
		infoSymbols = append(infoSymbols, irregularSymbol)
	}
	if h.IsOutdated {
		infoSymbols = append(infoSymbols, outdatedSymbol)
	}
	return strings.Join(infoSymbols[:], " | ")
}

func (h *headword) GlossText() string {
	gloss := h.Expression
	if h.IsAteji {
		gloss = "〈" + gloss + "〉"
	}
	symbolText := h.InfoSymbols()
	if symbolText != "" {
		gloss += "（" + symbolText + "）"
	}
	return gloss
}

func (h *headword) TableColHeaderText() string {
	text := h.KanjiForm()
	if h.IsAteji {
		text = "〈" + text + "〉"
	}
	return text
}

func (h *headword) TableRowHeaderText() string {
	text := h.Reading
	if h.IsGikun {
		text = "〈" + text + "〉"
	}
	return text
}

func (h *headword) TableCellText() string {
	text := h.InfoSymbols()
	if text == "" {
		return defaultSymbol
	} else {
		return text
	}
}

func (h *headword) KanjiForm() string {
	if h.IsKanaOnly() {
		return "∅"
	} else {
		return h.Expression
	}
}

func needsFormTable(headwords []headword) bool {
	// Does the entry contain more than 1 distinct reading?
	// E.g. バカがい and ばかがい are not distinct.
	uniqueReading := ""
	for _, h := range headwords {
		if h.IsGikun {
			return true
		} else if h.IsSearchOnly {
			continue
		} else if h.IsKanaOnly() {
			continue
		} else if uniqueReading == "" {
			uniqueReading = kata2hira(h.Reading)
		} else if uniqueReading != kata2hira(h.Reading) {
			return true
		}
	}
	return false
}

type formTableData struct {
	kanjiForms    []string
	readings      []string
	colHeaderText map[string]string
	rowHeaderText map[string]string
	cellText      map[string]map[string]string
}

func tableData(headwords []headword) formTableData {
	d := formTableData{
		kanjiForms:    []string{},
		readings:      []string{},
		colHeaderText: make(map[string]string),
		rowHeaderText: make(map[string]string),
		cellText:      make(map[string]map[string]string),
	}
	for _, h := range headwords {
		if h.IsSearchOnly {
			continue
		}
		kanjiForm := h.KanjiForm()
		if !slices.Contains(d.kanjiForms, kanjiForm) {
			d.kanjiForms = append(d.kanjiForms, kanjiForm)
			d.colHeaderText[kanjiForm] = h.TableColHeaderText()
		}
		reading := h.Reading
		if !slices.Contains(d.readings, reading) {
			d.readings = append(d.readings, reading)
			d.rowHeaderText[reading] = h.TableRowHeaderText()
			d.cellText[reading] = make(map[string]string)
		}
		d.cellText[reading][kanjiForm] = h.TableCellText()
	}
	return d
}

func formsTableGlossary(headwords []headword) []any {
	d := tableData(headwords)

	attr := contentAttr{}
	centeredAttr := contentAttr{textAlign: "center"}
	leftAttr := contentAttr{textAlign: "left"}

	cornerCell := contentTableHeadCell(attr, "") // empty cell in upper left corner
	headRowCells := []any{cornerCell}
	for _, kanjiForm := range d.kanjiForms {
		content := d.colHeaderText[kanjiForm]
		cell := contentTableHeadCell(centeredAttr, content)
		headRowCells = append(headRowCells, cell)
	}
	headRow := contentTableRow(attr, headRowCells...)
	tableRows := []any{headRow}
	for _, reading := range d.readings {
		rowHeadCellText := d.rowHeaderText[reading]
		rowHeadCell := contentTableHeadCell(leftAttr, rowHeadCellText)
		rowCells := []any{rowHeadCell}
		for _, kanjiForm := range d.kanjiForms {
			text := d.cellText[reading][kanjiForm]
			rowCell := contentTableCell(centeredAttr, text)
			rowCells = append(rowCells, rowCell)
		}
		tableRow := contentTableRow(attr, rowCells...)
		tableRows = append(tableRows, tableRow)
	}
	tableAttr := contentAttr{data: map[string]string{"content": "formsTable"}}
	contentTable := contentTable(tableAttr, tableRows...)
	content := contentStructure(contentTable)
	return []any{content}
}

func formsGlossary(headwords []headword) []any {
	glossary := []any{}
	for _, h := range headwords {
		if h.IsSearchOnly {
			continue
		}
		text := h.GlossText()
		glossary = append(glossary, text)
	}
	return glossary
}

func baseFormsTerm(entry jmdict.JmdictEntry) dbTerm {
	term := dbTerm{Sequence: entry.Sequence}
	headwords := extractHeadwords(entry)
	if needsFormTable(headwords) {
		term.Glossary = formsTableGlossary(headwords)
	} else {
		term.Glossary = formsGlossary(headwords)
	}
	for _, sense := range entry.Sense {
		rules := grammarRules(sense.PartsOfSpeech)
		term.addRules(rules...)
	}
	return term
}

func formsExportDb(inputPath, outputPath, languageName, title string, stride int, pretty bool) error {
	reader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	dictionary, entities, err := jmdict.LoadJmdictNoTransform(reader)
	if err != nil {
		return err
	}

	meta := newJmdictMetadata(dictionary, "english")

	terms := dbTermList{}
	for _, entry := range dictionary.Entries {
		baseTerm := baseFormsTerm(entry)
		headwords := extractHeadwords(entry)
		for _, h := range headwords {
			if h.IsSearchOnly {
				if term, ok := createSearchTerm(h, entry, meta); ok {
					terms = append(terms, term)
				}
				continue
			}
			term := baseTerm
			term.Expression = h.Expression
			term.Reading = h.Reading
			term.addTermTags(h.TermTags...)
			term.Score = calculateTermScore(1, 0, h)
			terms = append(terms, term)
		}
	}

	tags := dbTagList{}
	tags = append(tags, entityTags(entities)...)
	tags = append(tags, newsFrequencyTags()...)
	tags = append(tags, customDbTags()...)

	if title == "" {
		title = "JMdict Forms"
	}

	recordData := map[string]dbRecordList{
		"term": terms.crush(),
		"tag":  tags.crush(),
	}

	jmdictDate := jmdictPublicationDate(dictionary)

	index := dbIndex{
		Title:       title,
		Revision:    "JMdict." + jmdictDate,
		Sequenced:   true,
		Attribution: edrdgAttribution,
	}
	index.setDefaults()

	return writeDb(
		outputPath,
		index,
		recordData,
		stride,
		pretty,
	)
}
