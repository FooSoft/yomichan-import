package yomichan

import (
	"fmt"
	"strconv"

	"foosoft.net/projects/jmdict"
)

func glossaryContainsLanguage(glossary []jmdict.JmdictGlossary, language string) bool {
	hasGlosses := false
	for _, gloss := range glossary {
		if glossContainsLanguage(gloss, language) {
			hasGlosses = true
			break
		}
	}
	return hasGlosses
}

func glossContainsLanguage(gloss jmdict.JmdictGlossary, language string) bool {
	if gloss.Language == nil && language != "eng" {
		return false
	} else if gloss.Language != nil && language != *gloss.Language {
		return false
	} else {
		return true
	}
}

func makeGlossListItem(gloss jmdict.JmdictGlossary, language string) any {
	contents := []any{gloss.Content}
	listItem := contentListItem(contentAttr{}, contents...)
	return listItem
}

func makeInfoGlossListItem(gloss jmdict.JmdictGlossary, language string) any {
	// Prepend gloss with "type" (literal, figurative, trademark, etc.)
	glossTypeCode := *gloss.Type
	contents := []any{}
	if name, ok := glossTypeCodeToName[LangCode{language, glossTypeCode}]; ok {
		if name != "" {
			italicStyle := contentAttr{fontStyle: "italic"}
			contents = append(contents, contentSpan(italicStyle, "("+name+")"), " ")
		}
	} else {
		fmt.Println("Unknown glossary type code " + *gloss.Type + " for build language " + language)
		contents = append(contents, "["+glossTypeCode+"] ")
	}
	contents = append(contents, gloss.Content)
	listItem := contentListItem(contentAttr{}, contents...)
	return listItem
}

func makeSourceLangListItem(sourceLanguage jmdict.JmdictSource, language string) any {
	contents := []any{}

	var srcLangCode string
	if sourceLanguage.Language == nil {
		srcLangCode = "eng"
	} else {
		srcLangCode = *sourceLanguage.Language
	}

	// Format: [Language] ([Partial?], [Wasei?]): [Original word?]
	// [Language]
	if langName, ok := langCodeToName[LangCode{language, srcLangCode}]; ok {
		contents = append(contents, langName)
	} else {
		contents = append(contents, srcLangCode)
		fmt.Println("Unable to convert ISO 639 code " + srcLangCode + " to its full name in language " + language)
	}

	// ([Partial?], [Wasei?])
	var sourceLangTypeCode string
	if sourceLanguage.Type == nil {
		sourceLangTypeCode = ""
	} else {
		sourceLangTypeCode = *sourceLanguage.Type
	}
	var sourceLangType string
	if val, ok := sourceLangTypeCodeToType[LangCode{language, sourceLangTypeCode}]; ok {
		sourceLangType = val
	} else {
		sourceLangType = sourceLangTypeCode
		fmt.Println("Unknown source language type code " + sourceLangTypeCode + " for build language " + language)
	}
	if sourceLangType != "" && sourceLanguage.Wasei == "y" {
		contents = append(contents, " ("+sourceLangType+", wasei)")
	} else if sourceLangType != "" {
		contents = append(contents, " ("+sourceLangType+")")
	} else if sourceLanguage.Wasei == "y" {
		contents = append(contents, " (wasei)")
	}

	// : [Original word?]
	if sourceLanguage.Content != "" {
		contents = append(contents, ": ")
		attr := contentAttr{lang: ISOtoHTML[srcLangCode]}
		contents = append(contents, contentSpan(attr, sourceLanguage.Content))
	}

	listItem := contentListItem(contentAttr{}, contents...)
	return listItem
}

func makeReferenceListItem(reference string, refType string, meta jmdictMetadata) any {
	contents := []any{}
	attr := contentAttr{}

	hint := refNoteHint[LangCode{meta.language, refType}]
	contents = append(contents, hint+": ")

	refHeadword, senseNumber, ok := parseReference(reference)
	if !ok {
		contents = append(contents, "【"+reference+"】")
		return contentListItem(attr, contents...)
	}

	sequence, ok := meta.referenceToSeq[reference]
	if !ok {
		contents = append(contents, "【"+reference+"】")
		return contentListItem(attr, contents...)
	}

	targetSense := senseID{
		sequence: sequence,
		number:   senseNumber,
	}

	expHash := refHeadword.ExpHash()
	doDisplayReading := (len(meta.expHashToReadings[expHash]) > 1)
	doDisplaySenseNumber := (meta.seqToSenseCount[targetSense.sequence] > 1)
	refGlossAttr := contentAttr{
		fontSize:      "65%",
		verticalAlign: "middle",
		data:          map[string]string{"content": "refGlosses"},
	}

	contents = append(contents, refHeadword.ToInternalLink(doDisplayReading))
	if doDisplaySenseNumber {
		contents = append(contents, contentSpan(refGlossAttr, " "+strconv.Itoa(targetSense.number)+". "+meta.condensedGlosses[targetSense]))
	} else {
		contents = append(contents, contentSpan(refGlossAttr, " "+meta.condensedGlosses[targetSense]))
	}

	listItem := contentListItem(attr, contents...)
	return listItem
}

func makeExampleListItem(sentence jmdict.JmdictExampleSentence) any {
	if sentence.Lang == "jpn" {
		return contentListItem(contentAttr{}, sentence.Text)
	} else {
		attr := contentAttr{
			lang:          ISOtoHTML[sentence.Lang],
			listStyleType: ISOtoFlag[sentence.Lang],
		}
		return contentListItem(attr, sentence.Text)
	}
}

func listAttr(lang string, listStyleType string, dataContent string) contentAttr {
	return contentAttr{
		lang:          lang,
		listStyleType: listStyleType,
		data:          map[string]string{"content": dataContent},
	}
}

func needsStructuredContent(sense jmdict.JmdictSense, language string) bool {
	for _, gloss := range sense.Glossary {
		if glossContainsLanguage(gloss, language) && gloss.Type != nil {
			return true
		}
	}
	if len(sense.SourceLanguages) > 0 {
		return true
	} else if len(sense.Information) > 0 {
		return true
	} else if len(sense.Antonyms) > 0 {
		return true
	} else if len(sense.References) > 0 {
		return true
	} else if len(sense.Examples) > 0 {
		return true
	} else {
		return false
	}
}

func createGlossaryContent(sense jmdict.JmdictSense, meta jmdictMetadata) any {
	glossaryContents := []any{}

	// Add normal glosses
	glossListItems := []any{}
	for _, gloss := range sense.Glossary {
		if glossContainsLanguage(gloss, meta.language) && gloss.Type == nil {
			listItem := makeGlossListItem(gloss, meta.language)
			glossListItems = append(glossListItems, listItem)
		}
	}
	if len(glossListItems) > 0 {
		attr := listAttr(ISOtoHTML[meta.language], "circle", "glossary")
		list := contentUnorderedList(attr, glossListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add information glosses
	infoGlossListItems := []any{}
	for _, gloss := range sense.Glossary {
		if glossContainsLanguage(gloss, meta.language) && gloss.Type != nil {
			listItem := makeInfoGlossListItem(gloss, meta.language)
			infoGlossListItems = append(infoGlossListItems, listItem)
		}
	}
	if len(infoGlossListItems) > 0 {
		attr := listAttr(ISOtoHTML[meta.language], infoMarker, "infoGlossary")
		list := contentUnorderedList(attr, infoGlossListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add language-of-origin / loanword information
	sourceLangListItems := []any{}
	for _, sourceLanguage := range sense.SourceLanguages {
		listItem := makeSourceLangListItem(sourceLanguage, meta.language)
		sourceLangListItems = append(sourceLangListItems, listItem)
	}
	if len(sourceLangListItems) > 0 {
		attr := listAttr(ISOtoHTML[meta.language], langMarker, "sourceLanguages")
		list := contentUnorderedList(attr, sourceLangListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add sense notes
	noteListItems := []any{}
	for _, information := range sense.Information {
		listItem := contentListItem(contentAttr{}, information)
		noteListItems = append(noteListItems, listItem)
	}
	if len(noteListItems) > 0 {
		attr := listAttr(ISOtoHTML["jpn"], noteMarker, "notes") // notes often contain japanese text
		list := contentUnorderedList(attr, noteListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add antonyms
	antonymListItems := []any{}
	for _, antonym := range sense.Antonyms {
		listItem := makeReferenceListItem(antonym, "ant", meta)
		antonymListItems = append(antonymListItems, listItem)
	}
	if len(antonymListItems) > 0 {
		attr := listAttr(ISOtoHTML[meta.language], antonymMarker, "antonyms")
		list := contentUnorderedList(attr, antonymListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add cross-references
	referenceListItems := []any{}
	for _, reference := range sense.References {
		listItem := makeReferenceListItem(reference, "xref", meta)
		referenceListItems = append(referenceListItems, listItem)
	}
	if len(referenceListItems) > 0 {
		attr := listAttr(ISOtoHTML[meta.language], refMarker, "references")
		list := contentUnorderedList(attr, referenceListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	// Add example sentences
	exampleListItems := []any{}
	for _, example := range sense.Examples {
		for _, sentence := range example.Sentences {
			listItem := makeExampleListItem(sentence)
			exampleListItems = append(exampleListItems, listItem)
		}
	}
	if len(exampleListItems) > 0 {
		attr := listAttr(ISOtoHTML["jpn"], ISOtoFlag["jpn"], "examples")
		list := contentUnorderedList(attr, exampleListItems...)
		glossaryContents = append(glossaryContents, list)
	}

	return contentStructure(glossaryContents...)
}

func createGlossary(sense jmdict.JmdictSense, meta jmdictMetadata) []any {
	glossary := []any{}
	if needsStructuredContent(sense, meta.language) {
		glossary = append(glossary, createGlossaryContent(sense, meta))
	} else {
		for _, gloss := range sense.Glossary {
			if glossContainsLanguage(gloss, meta.language) {
				glossary = append(glossary, gloss.Content)
			}
		}
	}
	return glossary
}
