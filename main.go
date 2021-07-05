package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type myChan chan string

var mainWindow fyne.Window

func main() {
	// create new Fyne app
	myApp := app.NewWithID("blog.letscode.fynegoroutines")

	// create main window
	mainWindow = myApp.NewWindow("Fyne goroutines")

	// create button which will create new channel, start background worker and listen for results
	mainWindow.SetContent(widget.NewButton("Click Me", func() {
		// create channel to communicate with background goroutine
		chn := make(myChan)

		// start background worker in separate goroutine
		go backgroundWorker(chn)

		// listen on data received from background worker
		go listen(chn)
	}))

	// set window size and start app
	mainWindow.Resize(fyne.Size{
		Width:  400,
		Height: 225,
	})
	mainWindow.ShowAndRun()
}

// listen function to receive data from background worker
func listen(chn myChan) {
	for {
		select {
		case msg := <-chn:
			if msg == "get" {
				showDialog(chn)
			}
			if msg == "ok" {
				return
			}
		}
	}
}

// helper function to spawn new dialog and send data back to background worker
func showDialog(chn myChan) {
	items := make([]*widget.FormItem, 1)
	pinEntry := widget.NewEntry()
	items[0] = widget.NewFormItem("PIN: ", pinEntry)
	dialog.ShowForm("Enter PIN", "OK", "Cancel", items,
		func(bool) {
			// send data back to background worker
			chn <- pinEntry.Text
		},
		mainWindow,
	)
}

// background worker running in separate goroutine
func backgroundWorker(chn myChan) {
	for {
		chn <- "get"
		pin := <-chn
		if pin == "1234" {
			chn <- "ok"
			return
		}
	}
}
