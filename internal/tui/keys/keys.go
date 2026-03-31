package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Select    key.Binding
	Back      key.Binding
	PlayPause key.Binding
	NextTrack key.Binding
	PrevTrack key.Binding
	VolUp     key.Binding
	VolDown   key.Binding
	Help      key.Binding
	Quit      key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
// Ordered by importance, ending with the toggle help command.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Select, k.PlayPause, k.NextTrack, k.Back, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Select, k.Back},
		{k.PlayPause, k.NextTrack, k.PrevTrack, k.VolUp, k.VolDown},
		{k.Help, k.Quit},
	}
}

var Keys = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	PlayPause: key.NewBinding(
		key.WithKeys(" "),
		key.WithHelp("space", "play/pause"),
	),
	NextTrack: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "next track"),
	),
	PrevTrack: key.NewBinding(
		key.WithKeys("N"),
		key.WithHelp("N", "prev track"),
	),
	VolUp: key.NewBinding(
		key.WithKeys("+"),
		key.WithHelp("+", "vol up"),
	),
	VolDown: key.NewBinding(
		key.WithKeys("-"),
		key.WithHelp("-", "vol down"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "more"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
