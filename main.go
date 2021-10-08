package main

/*
	unsafe.Pointer(p interface{}): Returns a pointer that can be used as a parameter for syscall.Proc.Call.
*/

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gen2brain/dlgs"
	"github.com/go-vgo/robotgo"
	"github.com/pterm/pterm"
)

var AutoClickerEnabled bool = false

func TimeUnitPromp(p *pterm.ProgressbarPrinter) (time.Duration, string) {
	var Unit time.Duration
	item, _, err := dlgs.List("Time Unit", "Select a time unit from the List:", []string{"Second", "Millisecond", "Nanosecond"})
	if err != nil {
		panic(err)
	}
	switch item {
	case "Second":
		Unit = time.Second
		pterm.Success.Println("Set Time Unit to 'Second'")
	case "Millisecond":
		Unit = time.Millisecond
		pterm.Success.Println("Set Time Unit to 'Millisecond'")
	case "Nanosecond":
		Unit = time.Nanosecond
		pterm.Success.Println("Set Time Unit to 'Nanosecond'")
	default:
		pterm.Warning.Println("No Time Unit selected: Selected 'Millisecond' by Default")
		item = "Millisecond"
		Unit = time.Millisecond
	}
	return Unit, item
}

func TimeInterval(p *pterm.ProgressbarPrinter, unitStr string) int {
	res, _, err := dlgs.Entry("Interval", fmt.Sprintf("How many clicks per %v", unitStr), "")
	if err != nil {
		panic(err)
	}

	if res == "exit" {
		os.Exit(0)
	}

	Parsed, err := strconv.Atoi(res)
	if err != nil {
		pterm.Error.Println("Invalid Number. Prompting again for a valid Number")
		return TimeInterval(p, unitStr)
	}
	pterm.Success.Printf("Set clicks to 1 every %v.", unitStr)
	return Parsed
}

func AutoClicker(unit time.Duration, clicks int) {
	for {
		if AutoClickerEnabled {
			robotgo.MouseClick()
			time.Sleep(unit * time.Duration(clicks))
		} else {
			time.Sleep(time.Millisecond * 1)
		}
	}
}

func main() {
	pterm.DefaultHeader.Println("NeutronX Auto Clicker")
	p, _ := pterm.DefaultProgressbar.WithTotal(2).WithTitle("Set-up").Start()
	timeUnit, timeUnit_str := TimeUnitPromp(p)
	p.Increment()
	timeInterval := TimeInterval(p, timeUnit_str)
	p.Increment()
	for i := 0; i < 10; i++ {
		robotgo.MouseClick()
		time.Sleep(timeUnit * time.Duration(timeInterval))
	}
	fmt.Println(("\033[H\033[2J"))
	pterm.DefaultHeader.Println("NeutronX Auto Clicker")
	pterm.DefaultTable.WithHasHeader().WithData(pterm.TableData{
		{"Key", "Value"},
		{"Time Unit", timeUnit_str},
		{"Time Interval", strconv.Itoa(timeInterval)},
		{"CPS", "1 Click Every " + timeUnit_str},
	}).Render()
	go AutoClicker(timeUnit, timeInterval)
	for {
		onkey := robotgo.AddEvent("`")
		if onkey {
			if AutoClickerEnabled {
				AutoClickerEnabled = false
			} else {
				AutoClickerEnabled = true
			}
		}
	}
}
