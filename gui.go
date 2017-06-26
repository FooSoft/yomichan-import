/*
 * Copyright (c) 2017 Alex Yatskov <alex@foosoft.net>
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
	"log"
	"path/filepath"
	"strings"

	"github.com/andlabs/ui"
)

type logger struct {
	entry *ui.Entry
}

func (l *logger) Write(p []byte) (n int, err error) {
	ui.QueueMain(func() {
		l.entry.SetText(strings.Trim(string(p), "\n"))
	})

	return len(p), nil
}

func gui() error {
	return ui.Main(func() {
		pathSourceEntry := ui.NewEntry()
		pathSourceButton := ui.NewButton("Browse...")
		pathSourceBox := ui.NewHorizontalBox()
		pathSourceBox.Append(pathSourceEntry, true)
		pathSourceBox.Append(pathSourceButton, false)

		pathTargetEntry := ui.NewEntry()
		pathTargetButton := ui.NewButton("Browse...")
		pathTargetBox := ui.NewHorizontalBox()
		pathTargetBox.Append(pathTargetEntry, true)
		pathTargetBox.Append(pathTargetButton, false)

		titleEntry := ui.NewEntry()
		languageEntry := ui.NewEntry()
		outputEntry := ui.NewEntry()
		importButton := ui.NewButton("Import dictionary...")

		mainBox := ui.NewVerticalBox()
		mainBox.Append(ui.NewLabel("Path to dictionary source (CATALOGS file for EPWING)"), false)
		mainBox.Append(pathSourceBox, false)
		mainBox.Append(ui.NewLabel("Path to dictionary target ZIP file"), false)
		mainBox.Append(pathTargetBox, false)
		mainBox.Append(ui.NewLabel("Dictionary display title (blank for default)"), false)
		mainBox.Append(titleEntry, false)
		mainBox.Append(ui.NewLabel("Dictionary glossary language (blank for English)"), false)
		mainBox.Append(languageEntry, false)
		mainBox.Append(ui.NewLabel("Application output"), false)
		mainBox.Append(outputEntry, false)
		mainBox.Append(ui.NewVerticalBox(), true)
		mainBox.Append(importButton, false)

		window := ui.NewWindow("Yomichan Import", 640, 320, false)
		window.SetMargined(true)
		window.SetChild(mainBox)

		pathSourceButton.OnClicked(func(*ui.Button) {
			if path := ui.OpenFile(window); len(path) > 0 {
				pathSourceEntry.SetText(path)
			}
		})

		pathTargetButton.OnClicked(func(*ui.Button) {
			if path := ui.SaveFile(window); len(path) > 0 {
				if len(filepath.Ext(path)) == 0 {
					path += ".zip"
				}

				pathTargetEntry.SetText(path)
			}
		})

		log.SetOutput(&logger{outputEntry})

		importButton.OnClicked(func(*ui.Button) {
			importButton.Disable()
			outputEntry.SetText("")

			inputPath := pathSourceEntry.Text()
			if len(inputPath) == 0 {
				ui.MsgBoxError(window, "Error", "You must specify a dictionary source path")
				importButton.Enable()
				return
			}

			outputPath := pathTargetEntry.Text()
			if len(outputPath) == 0 {
				ui.MsgBoxError(window, "Error", "You must specify a dictionary target path")
				importButton.Enable()
				return
			}

			format, err := detectFormat(inputPath)
			if err != nil {
				ui.MsgBoxError(window, "Error", "Unable to detect dictionary format")
				importButton.Enable()
				return
			}

			if format == "epwing" {
				inputPath = filepath.Dir(inputPath)
			}

			title := titleEntry.Text()
			language := languageEntry.Text()

			go func() {
				defer ui.QueueMain(func() {
					importButton.Enable()
				})

				if err := exportDb(inputPath, outputPath, format, language, title, DEFAULT_STRIDE, false); err != nil {
					log.Print(err)
				}
			}()
		})

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()
	})
}
