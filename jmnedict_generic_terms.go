package yomichan

import (
	"golang.org/x/exp/slices"
)

type genericTermMap map[string]map[string][]string

type genericTermInfo struct {
	expressionToTagToGlosses genericTermMap
	usedSequences            map[sequence]bool
	currentSequence          sequence
}

func newGenericTermInfo() genericTermInfo {
	return genericTermInfo{
		expressionToTagToGlosses: genericTermMap{},
		usedSequences:            map[sequence]bool{},
	}
}

func (i *genericTermInfo) NewSequence() sequence {
	seq := i.currentSequence + 1
	for i.usedSequences[seq] {
		seq += 1
	}
	i.AddUsedSequence(seq)
	i.currentSequence = seq
	return seq
}

func (i *genericTermInfo) AddUsedSequence(s sequence) {
	i.usedSequences[s] = true
}

func (i *genericTermInfo) AddGlosses(exp string, tags []string, gloss string) {
	if i.expressionToTagToGlosses[exp] == nil {
		i.expressionToTagToGlosses[exp] = map[string][]string{}
	}
	for _, tag := range tags {
		glosses := i.expressionToTagToGlosses[exp][tag]
		if !slices.Contains(glosses, gloss) {
			glosses = append(glosses, gloss)
			i.expressionToTagToGlosses[exp][tag] = glosses
		}
	}
}

func (i *genericTermInfo) IsGenericName(headword headword, definitions []string) bool {
	if headword.IsKanaOnly() {
		// No reason to process these terms.
		return false
	}
	isGenericName := true
	for _, definition := range definitions {
		if !isTransliteration(definition, headword.Reading) {
			isGenericName = false
			break
		}
	}
	return isGenericName
}

func (i *genericTermInfo) Terms() (terms []dbTerm) {
	for expression, tagToGlosses := range i.expressionToTagToGlosses {
		seq := i.NewSequence()
		for tag, glosses := range tagToGlosses {
			term := dbTerm{
				Expression: expression,
				Sequence:   seq,
			}
			for _, gloss := range glosses {
				term.Glossary = append(term.Glossary, gloss)
			}
			term.addDefinitionTags(tag)
			terms = append(terms, term)
		}
	}
	return terms
}
