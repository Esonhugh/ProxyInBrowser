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

func RunApplication(Spec ApplicationSpec) {
	app := CreateApplication(Spec)
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
		for {
			LogUI.Clear()
			if _, err := LogUI.Write(app.Spec.ConsoleLogBuffer.Bytes()); err != nil {
				LogUI.Write([]byte(err.Error()))
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()

	return LogUI
}

func (app Application) CreateTermUI() *tview.InputField {
	TermUI := tview.NewInputField().SetLabel("Console> ")
	TermUI.SetBorder(true).SetTitle("Command").SetTitleAlign(tview.AlignCenter)
	TermUI.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEnter {
			app.ExecuteCommand(TermUI.GetText())
			TermUI.SetText("")
		}
		return event
	})
	return TermUI
}
