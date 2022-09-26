package main

import (
	"fmt"
	"io"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func ReadFile(w *fyne.Window) []byte {
	dest, err := storage.ParseURI("file://testdata.txt")
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return []byte{}
	}
	ok, err := storage.Exists(dest)
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return []byte{}
	}
	if !ok {
		wr, err := storage.Writer(dest)
		if err != nil {
			d := dialog.NewError(err, *w)
			d.Show()
			return []byte{}
		}
		wr.Write([]byte("The quick brown fox jumps over the lazy dogs back.\ntestapp.53"))
		wr.Close()
	}
	rr, err := storage.Reader(dest)
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return []byte{}
	}
	defer rr.Close()
	data, err := io.ReadAll(rr)
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return []byte{}
	}
	return data
}

func WriteFile(data []byte, w *fyne.Window) {
	dest, err := storage.ParseURI("file://testdata.txt")
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return
	}
	wr, err := storage.Writer(dest)
	if err != nil {
		d := dialog.NewError(err, *w)
		d.Show()
		return
	}
	wr.Write(data)
	wr.Close()
}

func main() {
	APP := app.NewWithID("apollo.testapp_53")
	WIN := APP.NewWindow("Testapp #53")
	WIN.CenterOnScreen()
	WIN.Resize(fyne.NewSize(400, 600))

	tick := time.NewTicker(time.Duration(100) * time.Millisecond)
	stop := make(chan bool, 1)

	clicks := 0
	test1 := widget.NewLabel("Hello and welcome to Testapp #53")

	go func() {
		for {
			select {
			case <-tick.C:
				clicks += 1
				test1.Text = fmt.Sprintf("Hello and welcome to Testapp #53\nClick #%d", clicks)
				test1.Refresh()
			case <-stop:
				break
			}
		}
	}()

	test2 := widget.NewMultiLineEntry()
	test2.Text = string(ReadFile(&WIN))
	test3 := widget.NewButton("Save", func() {
		WriteFile([]byte(test2.Text), &WIN)
	})

	WIN.SetContent(container.NewAppTabs(
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), container.NewVBox(
			test1,
			widget.NewButtonWithIcon("Click Me", theme.ConfirmIcon(), func() {
				clicks += 1
				test1.Text = fmt.Sprintf("Hello and welcome to Testapp #53\nClick #%d", clicks)
				test1.Refresh()
			}),
		)),
		container.NewTabItem("Tab 2", container.NewVBox(
			widget.NewButtonWithIcon("Reset Clicks", theme.MediaReplayIcon(), func() {
				clicks = 0
				test1.Text = fmt.Sprintf("Hello and welcome to Testapp #53\nClick #%d", clicks)
				test1.Refresh()
			}),
			widget.NewButtonWithIcon("Super Click Me", theme.CancelIcon(), func() {
				if clicks != 0 {
					clicks *= 2
				} else {
					clicks += 2
				}
				test1.Text = fmt.Sprintf("Hello and welcome to Testapp #53\nClick #%d", clicks)
				test1.Refresh()
			}),
		)),
		container.NewTabItemWithIcon("", theme.StorageIcon(), container.NewVBox(
			test3,
			test2,
		)),
	))

	WIN.Show()
	WIN.SetOnClosed(func() {
		stop <- true
	})
	APP.Run()
}
