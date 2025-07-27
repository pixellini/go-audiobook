package tui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Message types for the Bubble Tea model
type progressMsg struct{ message string }

type progressWithBarMsg struct {
	message string
	current int
	total   int
}

type completeMsg struct{ message string }

type finishMsg struct{ message string }

type stopMsg struct{}

// Model represents the Bubble Tea model state
type bubbleTeaModel struct {
	spinner      spinner.Model
	progress     progress.Model
	message      string
	completed    []string
	finished     bool
	stopped      bool
	showProgress bool
	current      int
	total        int
	width        int
	height       int
}

type BubbleTerminalUI struct {
	p  *tea.Program
	bm *bubbleTeaModel
}

const (
	progressBarWidth  uint8 = 50
	progressBarGutter uint8 = 4

	exitHint = "\nPress q or Ctrl+C to quit\n"
)

func NewBubbleTeaUI() TUIService {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = s.Style.Foreground(draculaPink)

	progressBar := progress.WithGradient(string(draculaOrange), string(draculaPink))
	p := progress.New(progressBar)
	// Default width, will be updated based on terminal size
	p.Width = int(progressBarWidth)

	bm := &bubbleTeaModel{
		spinner:      s,
		progress:     p,
		message:      "Initializing...",
		completed:    make([]string, 0),
		finished:     false,
		stopped:      false,
		showProgress: false,
		current:      0,
		total:        0,
		width:        80,
		height:       24,
	}

	return &BubbleTerminalUI{
		p:  tea.NewProgram(bm),
		bm: bm,
	}
}

func (b *BubbleTerminalUI) Start() error {
	go func() {
		_, _ = b.p.Run()
	}()
	return nil
}

func (b *BubbleTerminalUI) UpdateProgress(message string) {
	if b.p != nil {
		b.p.Send(progressMsg{message: message})
	}
}

func (b *BubbleTerminalUI) UpdateProgressWithBar(message string, current, total int) {
	if b.p != nil {
		b.p.Send(progressWithBarMsg{
			message: message,
			current: current,
			total:   total,
		})
	}
}

func (b *BubbleTerminalUI) CompleteCurrentTask(message string) {
	if b.p != nil {
		b.p.Send(completeMsg{message: message})
	}
}

func (b *BubbleTerminalUI) Finish(message string) {
	if b.p != nil {
		b.p.Send(finishMsg{message: message})
		time.Sleep(100 * time.Millisecond) // Give time to display final message
		b.p.Quit()
	}
}

func (b *BubbleTerminalUI) Stop() {
	if b.p != nil {
		b.p.Send(stopMsg{})
		b.p.Quit()
	}
}

func (bm bubbleTeaModel) Init() tea.Cmd {
	return tea.Batch(bm.spinner.Tick, tea.EnterAltScreen)
}

func (bm bubbleTeaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		return bm.updateProgressBarWidth(msg)

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			bm.stopped = true
			return bm, tea.Quit
		}

	case spinner.TickMsg:
		if !bm.finished && !bm.stopped {
			var cmd tea.Cmd
			bm.spinner, cmd = bm.spinner.Update(msg)
			return bm, cmd
		}

	case progressMsg:
		bm.message = msg.message
		bm.showProgress = false
		return bm, nil

	case progressWithBarMsg:
		bm.message = msg.message
		bm.current = msg.current
		bm.total = msg.total
		bm.showProgress = true
		return bm, nil

	case completeMsg:
		if bm.message != "" {
			bm.completed = append(bm.completed, "✓ "+bm.message+" - Complete")
		}
		if msg.message != "" {
			bm.message = msg.message
		}
		return bm, nil

	case finishMsg:
		bm.finished = true
		if bm.message != "" {
			bm.completed = append(bm.completed, "✓ "+bm.message+" - Complete")
		}
		if msg.message != "" {
			bm.completed = append(bm.completed, msg.message)
		}
		return bm, tea.Quit

	case stopMsg:
		bm.stopped = true
		return bm, tea.Quit
	}

	return bm, nil
}

func (bm bubbleTeaModel) View() string {
	if bm.stopped {
		return ""
	}

	if bm.finished {
		var output string
		for _, completed := range bm.completed {
			output += completed + "\n"
		}
		return output
	}

	var out string

	// Show completed tasks
	for _, completed := range bm.completed {
		out += completed + "\n"
	}

	// Show current progress
	if bm.message != "" {
		if bm.showProgress && bm.total > 0 {
			// Show progress bar
			percent := float64(bm.current) / float64(bm.total)
			progressBar := bm.progress.ViewAs(percent)
			out += fmt.Sprintf("%s %s\n", bm.spinner.View(), bm.message)
			out += fmt.Sprintf("\n%s %d/%d\n", progressBar, bm.current, bm.total)
		} else {
			// Show regular spinner
			out += fmt.Sprintf("%s %s\n", bm.spinner.View(), bm.message)
		}
	}

	out += exitHint

	return out
}

func (bm bubbleTeaModel) updateProgressBarWidth(msg tea.WindowSizeMsg) (tea.Model, tea.Cmd) {
	bm.width = msg.Width
	bm.height = msg.Height

	dw := int(progressBarWidth)
	// Set progress bar width to be ~75% of terminal width, with minimum of 50
	progressWidth := int(float64(msg.Width) * 0.75)
	if progressWidth < dw {
		progressWidth = dw
	}

	maxWidth := msg.Width - int(progressBarGutter)
	if progressWidth > maxWidth {
		progressWidth = maxWidth
	}

	bm.progress.Width = progressWidth
	return bm, nil
}
