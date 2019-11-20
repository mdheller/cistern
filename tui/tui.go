package tui

import (
	"context"
	"errors"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/nbedos/citop/cache"
	"github.com/nbedos/citop/text"
)

type ExecCmd struct {
	name string
	args []string
}

var ErrNoProvider = errors.New("list of providers must not be empty")

func RunApplication(ctx context.Context, newScreen func() (tcell.Screen, error), repositoryURL string, providers []cache.Provider, loc *time.Location) (err error) {
	if len(providers) == 0 {
		return ErrNoProvider
	}
	// FIXME Discard log until the status bar is implemented in order to hide the "Unsolicited response received on
	//  idle HTTP channel" from GitLab's HTTP client
	log.SetOutput(ioutil.Discard)
	encoding.Register()

	tmpDir, err := ioutil.TempDir("", "citop")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	defaultStyle := tcell.StyleDefault
	styleSheet := text.StyleSheet{
		text.TableHeader: func(s tcell.Style) tcell.Style {
			return s.Bold(true).Reverse(true)
		},
		text.ActiveRow: func(s tcell.Style) tcell.Style {
			return s.Background(tcell.ColorSilver).Foreground(tcell.ColorBlack).Bold(false).Underline(false).Blink(false)
		},
		text.Provider: func(s tcell.Style) tcell.Style {
			return s.Bold(true)
		},
		text.StatusFailed: func(s tcell.Style) tcell.Style {
			return s.Foreground(tcell.ColorMaroon)
		},
		text.StatusPassed: func(s tcell.Style) tcell.Style {
			return s.Foreground(tcell.ColorGreen)
		},
		text.StatusRunning: func(s tcell.Style) tcell.Style {
			return s.Foreground(tcell.ColorOlive)
		},
	}
	defaultStatus := "j:Down  k:Up  oO:Open  cC:Close  /:Search  b:Browser  ?:Help  q:Quit"

	ctx, cancel := context.WithCancel(ctx)
	cacheDB := cache.NewCache(providers)
	source := cacheDB.BuildsByCommit()

	ui, err := NewTUI(newScreen, defaultStyle, styleSheet)
	if err != nil {
		return err
	}
	defer func() {
		ui.Finish()
	}()

	controller, err := NewTableController(&ui, &source, loc, tmpDir, defaultStatus)
	if err != nil {
		return err
	}

	errCache := make(chan error)
	updates := make(chan time.Time)
	go func() {
		errCache <- cacheDB.UpdateFromProviders(ctx, repositoryURL, 30, updates)
	}()

	errController := make(chan error)
	go func() {
		errController <- controller.Run(ctx, updates)
	}()

	var e error
	errSet := false
	for i := 0; i < 2; i++ {
		select {
		case e = <-errCache:
			if e != nil && !errSet {
				cancel()
				err = e
				errSet = true
			}
		case e = <-errController:
			if !errSet {
				cancel()
				err = e
				errSet = true
			}
		}
	}

	return err
}

type TUI struct {
	newScreen    func() (tcell.Screen, error)
	screen       tcell.Screen
	defaultStyle tcell.Style
	styleSheet   text.StyleSheet
	eventc       chan tcell.Event
}

func NewTUI(newScreen func() (tcell.Screen, error), defaultStyle tcell.Style, styleSheet text.StyleSheet) (TUI, error) {
	ui := TUI{
		newScreen:    newScreen,
		defaultStyle: defaultStyle,
		styleSheet:   styleSheet,
		eventc:       make(chan tcell.Event),
	}
	err := ui.init()

	return ui, err
}

func (t *TUI) init() error {
	var err error
	t.screen, err = t.newScreen()
	if err != nil {
		return err
	}

	if err = t.screen.Init(); err != nil {
		return err
	}
	t.screen.SetStyle(t.defaultStyle)
	//screen.EnableMouse()
	t.screen.Clear()

	go t.poll()

	return nil
}

func (t TUI) Finish() {
	t.screen.Fini()
}

func (t TUI) Events() <-chan tcell.Event {
	return t.eventc
}

func (t TUI) poll() {
	// Exits when t.Finish() is called
	for {
		event := t.screen.PollEvent()
		if event == nil {
			break
		}
		t.eventc <- event
	}
}

func (t TUI) Draw(texts ...text.LocalizedStyledString) {
	t.screen.Clear()
	text.Draw(texts, t.screen, t.defaultStyle, t.styleSheet)
	t.screen.Show()
}

func (t *TUI) Exec(ctx context.Context, e ExecCmd) error {
	var err error
	t.Finish()
	defer func() {
		if e := t.init(); err == nil {
			err = e
		}
	}()

	cmd := exec.CommandContext(ctx, e.name, e.args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	return err
}
