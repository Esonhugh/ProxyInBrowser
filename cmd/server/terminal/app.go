package terminal

import (
	"bytes"
	"os"
	"time"

	"github.com/esonhugh/proxyinbrowser/cmd/server/define"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type ApplicationSpec struct {
	ConsoleLogBuffer *bytes.Buffer
	Rch              define.RelayChan
	CloseCh          chan os.Signal
}

type Application struct {
	Spec    ApplicationSpec
	UI      *tview.Application
	LogArea *tview.TextView
}

func CreateApplication(Spec ApplicationSpec) Application {
	return Application{
		Spec:    Spec,
		UI:      tview.NewApplication(),
		LogArea: tview.NewTextView(),
	}
}

func (app *Application) Run() {
	LogUI := app.CreateLogUI()
	TermUI := app.CreateTermUI()
	app.LogArea = LogUI
	Flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(LogUI, 0, 1, false).
		AddItem(TermUI, 3, 1, true)
	app.UI.SetRoot(Flex, true)
	app.UI.EnableMouse(true)
	if err := app.UI.Run(); err != nil {
		panic(err)
	}
}

func (app *Application) Stop() {
	app.UI.Stop()
}

func (app Application) CreateLogUI() *tview.TextView {
	LogUI := tview.NewTextView().
		SetDynamicColors(true).
		SetRegions(true).
		SetWordWrap(true).
		SetChangedFunc(func() {
			app.UI.Draw()
		})
	LogUI.SetBorder(true).
		SetTitle("Console Log").
		SetTitleAlign(tview.AlignCenter)

	go func() {
		oldByte := []byte{}
		for {
			b := app.Spec.ConsoleLogBuffer.Bytes()
			if bytes.Equal(b, oldByte) {
				time.Sleep(100 * time.Millisecond)
				continue // don't change
			}
			// b = bytes.ReplaceAll(b, oldByte, []byte(""))
			// reconstruct the data
			LogUI.Clear()
			if _, err := LogUI.Write(b); err != nil {
				LogUI.Write([]byte(err.Error()))
			}
			LogUI.ScrollToEnd()
			oldByte = bytes.Clone(b)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return LogUI
}

func (app Application) CreateTermUI() *tview.InputField {
	TermUI := tview.NewInputField().SetLabel("Console> ")
	TermUI.SetBorder(true).SetTitle("Command").SetTitleAlign(tview.AlignCenter)
	cmdhistory := []string{}
	uptime := 0
	TermUI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			cmdinput := TermUI.GetText()
			if cmdinput == "" {
				return event
			}
			cmdhistory = append(cmdhistory, cmdinput)
			app.ExecuteCommand(cmdinput)
			TermUI.SetText("") // clean
			uptime = 0
		} else if event.Key() == tcell.KeyUp {
			if len(cmdhistory) < 1+uptime || uptime < 0 {
				return event
			}
			lastx := cmdhistory[len(cmdhistory)-1-uptime]
			TermUI.SetText(lastx)
			uptime++
			return event
		} else if event.Key() == tcell.KeyDown {
			uptime--
			if len(cmdhistory) < 1+uptime || uptime < 0 {
				return event
			}
			lastx := cmdhistory[len(cmdhistory)-1-uptime]
			TermUI.SetText(lastx)
			return event
		} else {
			uptime = 0
		}
		return event
	})
	return TermUI
}
