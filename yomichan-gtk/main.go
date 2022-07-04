package main

import (
	"fmt"
	"path/filepath"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"

	yomichan "foosoft.net/projects/yomichan-import"
)

func main() {
	ui.Main(func() {
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

		importButton := ui.NewButton("Import dictionary...")

		titleEntry := ui.NewEntry()
		titleEntry.SetText(yomichan.DefaultTitle)

		languageEntry := ui.NewEntry()
		languageEntry.SetText(yomichan.DefaultLanguage)

		mainBox := ui.NewVerticalBox()
		mainBox.Append(ui.NewLabel("Path to dictionary source (CATALOGS file for EPWING)"), false)
		mainBox.Append(pathSourceBox, false)
		mainBox.Append(ui.NewLabel("Path to dictionary target ZIP file"), false)
		mainBox.Append(pathTargetBox, false)
		mainBox.Append(ui.NewLabel("Dictionary display title (blank for default)"), false)
		mainBox.Append(titleEntry, false)
		mainBox.Append(ui.NewLabel("Dictionary glossary language (blank for English)"), false)
		mainBox.Append(languageEntry, false)
		mainBox.Append(ui.NewVerticalBox(), true)
		mainBox.Append(importButton, false)

		window := ui.NewWindow("Yomichan Import", 640, 280, false)
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

		setBusyState := func(busy bool) {
			if busy {
				importButton.Disable()
				importButton.SetText("Importing, please wait...")

			} else {
				importButton.SetText("Start dictionary import")
				importButton.Enable()
			}
		}

		importButton.OnClicked(func(*ui.Button) {
			setBusyState(true)

			inputPath := pathSourceEntry.Text()
			if filepath.Base(inputPath) == "CATALOGS" {
				inputPath = filepath.Dir(inputPath)
			}

			if len(inputPath) == 0 {
				ui.MsgBoxError(window, "Error", "You must specify a dictionary source path")
				setBusyState(false)
				return
			}

			outputPath := pathTargetEntry.Text()
			if len(outputPath) == 0 {
				ui.MsgBoxError(window, "Error", "You must specify a dictionary target path")
				setBusyState(false)
				return
			}

			go func() {
				err := yomichan.ExportDb(
					inputPath,
					outputPath,
					yomichan.DefaultFormat,
					languageEntry.Text(),
					titleEntry.Text(),
					yomichan.DefaultStride,
					yomichan.DefaultPretty,
				)

				ui.QueueMain(func() {
					setBusyState(false)
					if err == nil {
						ui.MsgBox(window, "Success", "Conversion process complete!")
					} else {
						ui.MsgBox(window, "Error", fmt.Sprintf("Conversion process failed: %s", err.Error()))
					}
				})
			}()
		})

		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})

		window.Show()
	})
}
