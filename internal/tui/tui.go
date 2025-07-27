package tui

import "github.com/charmbracelet/lipgloss"

// TUIService defines the interface for terminal UI operations
type TUIService interface {
	Start() error
	UpdateProgress(message string)
	UpdateProgressWithBar(message string, current, total int)
	CompleteCurrentTask(message string)
	Finish(message string)
	Stop()
}

const (
	// Dracula Theme. Such an amazing colour palette.
	draculaBackground  = lipgloss.Color("#282A36")
	draculaCurrentLine = lipgloss.Color("#44475A")
	draculaForeground  = lipgloss.Color("#F8F8F2")
	draculaComment     = lipgloss.Color("#6272A4")
	draculaCyan        = lipgloss.Color("#8BE9FD")
	draculaGreen       = lipgloss.Color("#50FA7B")
	draculaOrange      = lipgloss.Color("#FFB86C")
	draculaPink        = lipgloss.Color("#FF79C6")
	draculaRed         = lipgloss.Color("#FF5555")
	draculaYellow      = lipgloss.Color("#F1FA8C")
)

// EmptyTUI is a TUIService implementation that does nothing
// Useful for when you want to disable the TUI output
type EmptyTUI struct{}

// NewEmptyTUI creates a new no-operation TUI service
// This is for when the user wants verbose logging,
// and in that case, we don't show the progress bar.
func NewEmpty() TUIService { return &EmptyTUI{} }

func (n *EmptyTUI) Start() error { return nil }

func (n *EmptyTUI) UpdateProgress(message string) {}

func (n *EmptyTUI) UpdateProgressWithBar(message string, current, total int) {}

func (n *EmptyTUI) CompleteCurrentTask(message string) {}

func (n *EmptyTUI) Finish(message string) {}

func (n *EmptyTUI) Stop() {}
