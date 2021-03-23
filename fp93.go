package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	appTitle      string = `[green]FloppyPunk`
	appHeaderText string = `[yellow::b]Traversing GlitchSpace in Relative Safety and Style since '93`
	asciiFloppy   string = `
	,'";-------------------;"'.
	;[]; ................. ;[];
	;  ; ................. ;  ;
	;  ; ................. ;  ;
	;  ; ................. ;  ;
	;  ; ................. ;  ;
	;  ; ................. ;  ;
	;  ; ................. ;  ;
	;  '.                 ,'  ;
	;    """""""""""""""""    ;
	;    ,-------------.---.  ;
	;    ;  ;"";       ;   ;  ;
	;    ;  ;  ;       ;   ;  ;
	;    ;  ;  ;       ;   ;  ;
	;//||;  ;  ;       ;   ;||;
	;\\||;  ;__;       ;   ;\/;
	'. _;          _  ;  _;  ;
	" """"""""""" """"" """

	Welcome to FloppyPunk

	[yellow]Press Enter to continue
`
	landingBodyText string = `
	In the near retrofuture, the GlitchSpace has revolutionized every field of human endeavor, enabling FTL travel, magic, esoteric mecha, and bringing transdimensional beings into our everyday lives.

	Among the vast stars and the less-than-empty night, millions of people make their living and try to stay one step ahead of the breakdown of reality.
`
)

// Shorthand function for error handling
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Instantiate the application and pages so they're available across functions
var app = tview.NewApplication()
var pages = tview.NewPages()

// Set the header for the app
func newFPHeader(title string, text string) tview.Primitive {
	header := tview.NewTextView().SetText(text).
		SetTextAlign(1).
		SetDynamicColors(true)
	header.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple).
		SetTitle(title)
	return header
}

func newIntroText() tview.Primitive {
	landingBody := tview.NewTextView().
		SetWordWrap(true).
		SetText(landingBodyText)
	landingBody.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple).
		SetTitle("[green]Introduction")
	return landingBody
}

var introText = newIntroText()

func newRulesText(path string) tview.Primitive {
	rulesText, err := ioutil.ReadFile(path)
	check(err)
	rulesBody := tview.NewTextView().
		SetWordWrap(true).
		SetDynamicColors(true).
		SetRegions(true).
		SetText(string(rulesText))
	rulesBody.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple).
		SetTitle("[green]Rules")
	return rulesBody
}

var rulesText = newRulesText("./rules.txt")

// The main menu controls
var mainMenu = tview.NewList().
	AddItem("Home", "Return to start", 'h', func() {
		pages.SwitchToPage("main")
		app.SetFocus(introText)
	}).
	AddItem("Rules", "Read the rules", 'r', func() {
		pages.SwitchToPage("rules")
		app.SetFocus(rulesText)
	}).
	AddItem("Create Character", "Create & save a PC", 'c', nil).
	AddItem("Load Character", "Load a saved PC", 'l', nil).
	AddItem("Quit", "Press to exit", 'q', func() { app.Stop() })

func newFPMainMenu(menu tview.Primitive) tview.Primitive {
	mainMenu := tview.NewFlex().AddItem(menu, 0, 1, false)
	mainMenu.SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple).
		SetTitle("[green]Menu")
	return mainMenu
}

func newLoadingPage(menu tview.Primitive) (textview tview.Primitive, flex tview.Primitive) {
	frontTextView := tview.NewTextView().
		SetDynamicColors(true).
		SetChangedFunc(func() {
			app.Draw()
		}).
		SetTextAlign(tview.AlignCenter).
		SetDoneFunc(func(key tcell.Key) {
			if key == tcell.KeyEnter {
				pages.SwitchToPage("main")
				app.SetFocus(menu)
			}
		})

	go func() {
		for _, word := range strings.Split(asciiFloppy, "\n") {
			fmt.Fprintf(frontTextView, "%s\n", word)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	frontTextView.
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple)

	frontFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newFPHeader(appTitle, appHeaderText), 0, 1, false).
		AddItem(frontTextView, 0, 6, true)

	return frontTextView, frontFlex
}

func newContextMenu(title string) tview.Primitive {
	controlPanel := tview.NewBox().
		SetBorder(true).
		SetBorderAttributes(tcell.AttrBold).
		SetBorderColor(tcell.ColorPurple).
		SetTitle(title)
	return controlPanel
}

func newContentPage(body tview.Primitive) tview.Primitive {
	middle := tview.NewFlex().
		AddItem(newFPMainMenu(mainMenu), 0, 1, false).
		AddItem(body, 0, 3, false)
	flex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(newFPHeader(appTitle, appHeaderText), 0, 1, false).
		AddItem(middle, 0, 8, false).
		AddItem(newContextMenu("[green]Controls"), 0, 1, false)
	return flex
}

func main() {
	frontText, frontFlex := newLoadingPage(mainMenu)
	pages.AddPage("front", frontFlex, true, true)
	pages.AddPage("main", newContentPage(introText), true, false).Focus(func(p tview.Primitive) {
		app.SetFocus(mainMenu)
	})
	pages.AddPage("rules", newContentPage(rulesText), true, false).Focus(func(p tview.Primitive) {
		app.SetFocus(mainMenu)
	})

	if err := app.SetRoot(pages, true).SetFocus(frontText).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
