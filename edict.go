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

import "io"

type edictKanjiElement struct {
	Expression  string   `xml:"keb"`
	Information []string `xml:"ke_inf"`
	Priority    []string `xml:"ke_pri"`
}

type edictReadingElement struct {
	Reading      string   `xml:"reb"`
	NoKanji      string   `xml:"re_nokanji"`
	Restrictions []string `xml:"re_restr"`
	Information  []string `xml:"re_inf"`
	Priority     []string `xml:"re_pri"`
}

type edictSense struct {
	RestrictKanji   []string `xml:"stagk"`
	RestrictReading []string `xml:"stagr"`
	References      []string `xml:"xref"`
	Antonyms        []string `xml:"ant"`
	PartOfSpeech    []string `xml:"pos"`
	Field           []string `xml:"field"`
	Misc            []string `xml:"misc"`
	SourceLanguage  []string `xml:"lsource"`
	Dialect         []string `xml:"dial"`
	Information     []string `xml:"s_inf"`
	Glossary        []string `xml:"gloss"`
}

type edictEntry struct {
	Sequence int                   `xml:"ent_seq"`
	Kanji    []edictKanjiElement   `xml:"k_ele"`
	Reading  []edictReadingElement `xml:"r_ele"`
	Sense    []edictSense          `xml:"sense"`
}

func processEdict(reader io.Reader, writer io.Writer) error {
	return nil
}
