/*
 * Copyright (c) 2017-2021 Alex Yatskov <alex@foosoft.net>
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
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	application := app.New()
	window := application.NewWindow("Yomichan Import")

	labelDictPath := widget.NewLabel("Path to dictionary source (CATALOGS file for EPWING)")
	entryDictPath := widget.NewEntry()
	buttonDictPath := widget.NewButtonWithIcon("Browse...", theme.SearchIcon(), func() {})
	layoutDictPath := container.NewBorder(
		nil,
		nil,
		nil,
		buttonDictPath,
		entryDictPath,
		buttonDictPath,
	)

	labelZipPath := widget.NewLabel("Path to dictionary target ZIP file")
	entryZipPath := widget.NewEntry()
	buttonZipPath := widget.NewButtonWithIcon("Browse...", theme.SearchIcon(), func() {})
	layoutZipPath := container.NewBorder(
		nil,
		nil,
		nil,
		buttonZipPath,
		entryZipPath,
		buttonZipPath,
	)

	labelDictLang := widget.NewLabel("Dictionary glossary language (for JMDICT only)")
	selectDictLang := widget.NewSelectEntry([]string{"Dutch", "English", "French", "German", "Hungarian", "Russian", "Slovenian", "Spanish"})
	selectDictLang.Text = "English"

	labelDictTitle := widget.NewLabel("Dictionary display title (blank for default)")
	entryDictTitle := widget.NewEntry()

	buttonImport := widget.NewButtonWithIcon("Convert", theme.ConfirmIcon(), func() {})
	buttonExit := widget.NewButtonWithIcon("Exit", theme.CancelIcon(), func() {})
	layoutButtons := container.NewHBox(
		layout.NewSpacer(),
		buttonImport,
		buttonExit,
	)

	layoutMaster := container.NewVBox(
		labelDictPath,
		layoutDictPath,
		labelZipPath,
		layoutZipPath,
		labelDictLang,
		selectDictLang,
		labelDictTitle,
		entryDictTitle,
		layout.NewSpacer(),
		layoutButtons,
	)

	window.SetContent(layoutMaster)
	window.ShowAndRun()
}
