package main

import (
	"fmt"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/go-vgo/robotgo"
)

func StrToTimeUnit(unitSTR string) (error, time.Duration) {
	var Unit time.Duration
	switch unitSTR {
	case "Second":
		Unit = time.Second
	case "Millisecond":
		Unit = time.Millisecond
	case "Nanosecond":
		Unit = time.Nanosecond
	default:
		return fmt.Errorf("error parsing time unit"), time.Second
	}
	return nil, Unit
}

type NeutronXAutoClicker struct {
	Started  bool
	unit     time.Duration
	interval int
}

func (clickr *NeutronXAutoClicker) Start() {
	clickr.Started = true
}

func (clickr *NeutronXAutoClicker) AutoClicker() {
	for {
		if clickr.Started {
			robotgo.MouseClick()
			time.Sleep(clickr.unit * time.Duration(clickr.interval))
		} else {
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func (clickr *NeutronXAutoClicker) Stop() {
	clickr.Started = false
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("NeutronX Auto Clicker")
	myWindow.Resize(fyne.NewSize(620, 320))
	SelectedEntry := widget.NewSelectEntry([]string{"Second", "Millisecond", "Nanosecond"})
	ClickingInterval := widget.NewEntry()
	KeyBind := widget.NewEntry()

	SelectedEntry.PlaceHolder = "Choose an Interval Unit"
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Time Unit", Widget: SelectedEntry}},
		OnSubmit: func() {
			interval, err := strconv.Atoi(ClickingInterval.Text)
			if err != nil {
				dialog.ShowError(err, myWindow)
			} else {
				err, timeUnit := StrToTimeUnit(SelectedEntry.Text)
				if err != nil {
					dialog.ShowError(err, myWindow)
				} else {
					if len(KeyBind.Text) == 1 {
						myWindow.Resize(fyne.NewSize(20, 20))
						dialog.ShowInformation("Started", fmt.Sprintf("To toggle press \"%v\".", KeyBind.Text), myWindow)
						var AutoClicker NeutronXAutoClicker = NeutronXAutoClicker{unit: timeUnit, interval: interval}

						AutoClicker.Stop()
						go AutoClicker.AutoClicker()
						for {
							onkey := robotgo.AddEvent(KeyBind.Text)
							if onkey {
								if AutoClicker.Started {
									AutoClicker.Started = false
								} else {
									AutoClicker.Started = true
								}
							}
						}

					} else {
						dialog.ShowError(fmt.Errorf("invalid keybind (must be 1 characters long)"), myWindow)
					}
				}
			}
		},
	}
	form.Append("Time Interval", ClickingInterval)
	form.Append("Keybind", KeyBind)

	myWindow.SetContent(form)
	myWindow.ShowAndRun()
}
