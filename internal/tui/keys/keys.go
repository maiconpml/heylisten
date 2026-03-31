package keys

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	PlayPause  key.Binding
	NextTrack  key.Binding
	PrevTrack  key.Binding
	NextTab    key.Binding
	TabHome    key.Binding
	TabLibrary key.Binding
	TabSearch  key.Binding
	PrevTab    key.Binding
	VolUp      key.Binding
	VolDown    key.Binding
	Help       key.Binding
	Quit       key.Binding
	Back       key.Binding
}

// ShortHelp returns keybindings to be shown in the mini help view.
// Ordered by importance, ending with the toggle help command.
func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.PlayPause, k.NextTrack, k.Quit, k.Help}
}

// FullHelp returns keybindings for the expanded help view.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextTab, k.PrevTab, k.Back},
		{k.PlayPause, k.NextTrack, k.PrevTrack, k.VolUp, k.VolDown},
		{k.Help, k.Quit},
	}
}

var Keys = KeyMap{
	TabHome: key.NewBinding(
		key.WithKeys("1"),
		key.WithHelp("1", "go to home")),
	TabLibrary: key.NewBinding(
		key.WithKeys("2"),
		key.WithHelp("2", "go to library")),
	TabSearch: key.NewBinding(
		key.WithKeys("3"),
		key.WithHelp("3", "go to search")),
	NextTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "next tab")),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "prev tab")),
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
	Back: key.NewBinding(
		key.WithKeys("esc", "backspace"),
		key.WithHelp("esc", "back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
