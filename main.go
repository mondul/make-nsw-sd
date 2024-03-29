package main

import (
	"errors"
	"fmt"
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type dos_type struct {
	atmosphere bool
	hekate     bool
	payload    bool
	bootdat    bool
	lockpick   bool
	sps        bool
	dbi        bool
}

var log_add func(string)

func myTitle(icon fyne.Resource, title string, fg_color color.Color) *fyne.Container {
	text := canvas.NewText(title, fg_color)
	text.TextSize = 11
	text.TextStyle.Bold = true
	return container.NewHBox(
		widget.NewIcon(icon),
		text,
	)
}

func newOutdir() string {
	return fmt.Sprintf("SD_%X", time.Now().Unix())
}

func main() {
	a := app.New()
	a.Settings().SetTheme(&myTheme{})

	w := a.NewWindow("Make NSW SD")

	w.Resize(fyne.NewSize(432, 336))
	w.SetFixedSize(true)

	// Create output dir name
	folder_entry_data := binding.NewString()
	folder_entry_data.Set(newOutdir())
	folder_entry := widget.NewEntryWithData(folder_entry_data)
	folder_entry.Disable()

	atmosphere_check := widget.NewCheck("Atmosphère", nil)
	atmosphere_check.Checked = true

	var hekate_row_1 *fyne.Container
	var hekate_row_2 *fyne.Container

	hekate_check := widget.NewCheck("Hekate", func(b bool) {
		if b {
			hekate_row_1.Show()
			hekate_row_2.Show()
		} else {
			hekate_row_1.Hide()
			hekate_row_2.Hide()
		}
	})
	hekate_check.Checked = true

	payload_check_data := binding.NewBool()
	bootdat_check_data := binding.NewBool()

	payload_check_data.AddListener(binding.NewDataListener(func() {
		if checked, _ := payload_check_data.Get(); checked {
			if checked, _ = bootdat_check_data.Get(); checked {
				bootdat_check_data.Set(false)
			}
		}
	}))

	bootdat_check_data.AddListener(binding.NewDataListener(func() {
		if checked, _ := bootdat_check_data.Get(); checked {
			if checked, _ = payload_check_data.Get(); checked {
				payload_check_data.Set(false)
			}
		}
	}))

	payload_check := widget.NewCheckWithData("payload.bin from Hekate", payload_check_data)
	bootdat_check := widget.NewCheckWithData("boot.dat from SX Gear", bootdat_check_data)
	lockpick_check := widget.NewCheck("Lockpick_RCM", nil)

	emsps := canvas.NewText("  ", color.Transparent)

	hekate_row_1 = container.NewHBox(emsps, payload_check, bootdat_check)
	hekate_row_2 = container.NewHBox(emsps, lockpick_check)

	sps_check := widget.NewCheck("SPs", nil)
	sps_check.Checked = true

	dbi_check := widget.NewCheck("DBI", nil)

	/* App containers */

	var home_container *fyne.Container

	log_txt_close := widget.NewButton("Close", func() {
		w.SetContent(home_container)
	})

	log_txt := widget.NewTextGrid()
	log_txt_scroll := container.NewScroll(log_txt)

	log_add = func(txt string) {
		log_txt.SetText(log_txt.Text() + txt)
		log_txt_scroll.ScrollToBottom()
	}

	log_container := container.NewBorder(
		nil,
		container.NewCenter(log_txt_close),
		nil,
		nil,
		log_txt_scroll,
	)

	/* Action buttons */

	ntd_err := errors.New(" Nothing to do! ")

	start_btn := widget.NewButton("Start", func() {
		if !atmosphere_check.Checked &&
			!hekate_check.Checked &&
			!sps_check.Checked &&
			!dbi_check.Checked {
			dialog.ShowError(ntd_err, w)
			return
		}

		// Avoid closing log when processing
		log_txt_close.Disable()

		// Show log
		w.SetContent(log_container)

		// Start process
		do_payload, _ := payload_check_data.Get()
		do_bootdat, _ := bootdat_check_data.Get()

		go start(dos_type{
			atmosphere: atmosphere_check.Checked,
			hekate:     hekate_check.Checked,
			payload:    do_payload,
			bootdat:    do_bootdat,
			lockpick:   lockpick_check.Checked,
			sps:        sps_check.Checked,
			dbi:        dbi_check.Checked,
		}, folder_entry_data, log_txt_close)
	})

	browse_btn := widget.NewButton(" … ", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if list == nil {
				folder_entry_data.Set(newOutdir())
				return
			}
			folder_entry_data.Set(list.Path())
		}, w)
	})

	/* Put everything together */

	fg_color := theme.DefaultTheme().Color(theme.ColorNameForeground, a.Settings().ThemeVariant())

	home_container = container.NewBorder(
		// Top
		nil,
		// Bottom
		container.NewVBox(
			widget.NewSeparator(),
			container.NewGridWithColumns(
				5,
				emsps,
				start_btn,
				emsps,
				widget.NewButton("Quit", w.Close),
				emsps,
			),
		),
		// Left
		nil,
		// Right
		nil,
		// Content
		container.NewVBox(
			myTitle(theme.FolderOpenIcon(), "Output folder", fg_color),
			container.NewBorder(nil, nil, nil, browse_btn, folder_entry),
			widget.NewSeparator(),
			myTitle(theme.DownloadIcon(), "Download & extract latest…", fg_color),
			// Checkboxes container without inner vertical padding
			container.New(
				newMyLayout(),
				atmosphere_check,
				hekate_check,
				hekate_row_1,
				hekate_row_2,
				sps_check,
				dbi_check,
			),
		),
	)

	w.SetContent(home_container)
	w.CenterOnScreen()
	w.ShowAndRun()
}
