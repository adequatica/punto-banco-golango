package main

import (
	"fmt"
	"os"
	"time"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
	"github.com/adequatica/punto-banco-golango/internal/rendering"
	"github.com/adequatica/punto-banco-golango/internal/statistics"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding

	Stats key.Binding
	Reset key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Stats, k.Reset, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter},    // first column
		{k.Stats, k.Reset, k.Quit}, // second column
	}
}

var defaultKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k", "K", "л", "Л"),
		key.WithHelp("↑/K", "— up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j", "J", "о", "О"),
		key.WithHelp("↓/J", "— down"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("ENTER/SPACE", "— select"),
	),
	Stats: key.NewBinding(
		key.WithKeys("s", "S", "ы", "Ы"),
		key.WithHelp("S", "— show/hide statistics"),
	),
	Reset: key.NewBinding(
		key.WithKeys("r", "R", "к", "К"),
		key.WithHelp("R", "— reset the game"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "Q", "й", "Й", "ctrl+c", "esc"),
		key.WithHelp("Q/CTRL+C", "— quit"),
	),
}

var defaultAfterRoundOptions = []string{"Next round", "Reset the game", "Quit"}

type UIstate int

const (
	stateIsBetting UIstate = iota
	stateIsProgress
	stateIsAfterRound
)

type model struct {
	stateUI           UIstate
	stateGame         puntobanco.GameResultState
	statistics        statistics.SessionStatistics
	showStatistics    bool
	cursor            int
	bettingOptions    []string
	afterRoundOptions []string
	selectedOption    string
	keys              keyMap
	help              help.Model
	spinner           spinner.Model
	spinnerStartTime  time.Time
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	return model{
		stateUI:           stateIsBetting,
		stateGame:         puntobanco.GetNewGameResultState(),
		statistics:        statistics.NewSessionStatistics(),
		showStatistics:    false,
		cursor:            0,
		bettingOptions:    puntobanco.GetBettingOptions(),
		afterRoundOptions: defaultAfterRoundOptions,
		selectedOption:    "",
		keys:              defaultKeys,
		help:              help.New(),
		spinner:           s,
	}
}

// Timer tick for spinner animation
type tickMsg time.Time

func tick() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(10 * time.Millisecond) // 10ms per tick
		return tickMsg(time.Now())
	}
}

var spinnerTimeout = 500 * time.Millisecond // 500ms total duration

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch {

		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit

		case key.Matches(msg, m.keys.Up):
			switch m.stateUI {
			case stateIsBetting:
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.bettingOptions) - 1
				}
			case stateIsAfterRound:
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.afterRoundOptions) - 1
				}
			}

		case key.Matches(msg, m.keys.Down):
			switch m.stateUI {
			case stateIsBetting:
				if m.cursor < len(m.bettingOptions)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			case stateIsAfterRound:
				if m.cursor < len(m.afterRoundOptions)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			}

		case key.Matches(msg, m.keys.Enter):
			switch m.stateUI {
			case stateIsBetting:
				// Store the selected betting choice
				m.selectedOption = m.bettingOptions[m.cursor]

				// Switch to in progress state and start animation
				m.stateUI = stateIsProgress
				m.spinnerStartTime = time.Now()
				return m, tea.Batch(
					m.spinner.Tick,
					tick(),
				)

			case stateIsAfterRound:
				switch m.afterRoundOptions[m.cursor] {
				case "Next round":
					// Switch to betting state
					m.stateUI = stateIsBetting
					m.cursor = 0
					m.selectedOption = ""
				case "Reset the game":
					// Switch to betting state with a new game session
					m.stateUI = stateIsBetting
					m.stateGame = puntobanco.GetNewGameResultState()
					m.statistics.ResetStatistics()
					m.cursor = 0
					m.selectedOption = ""
					m.spinnerStartTime = time.Time{} // Reset spinner timeout
				case "Quit":
					return m, tea.Quit
				}
			}

		case key.Matches(msg, m.keys.Reset):
			// Switch to betting state with a new game session
			m.stateUI = stateIsBetting
			m.stateGame = puntobanco.GetNewGameResultState()
			m.statistics.ResetStatistics()
			m.cursor = 0
			m.selectedOption = ""

		case key.Matches(msg, m.keys.Stats):
			m.showStatistics = !m.showStatistics
		}

	// Spinner tick
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// The programm of the game works extremely fast!
	// The progress animation emulates "thinking" of the program,
	// so that it won't be boring to get results instantly.
	case tickMsg:
		if m.stateUI == stateIsProgress {
			// Check if timeout have passed
			if time.Since(m.spinnerStartTime) >= spinnerTimeout {
				// Animation complete, play the game and switch to after round state
				gameResult, err := puntobanco.PlayPuntoBanco(m.stateGame.GetShoe())
				if err != nil {
					fmt.Printf("Alas, game error has happened: %v\n", err)
					// Reset game's session
					m.stateGame = puntobanco.GetNewGameResultState()
					m.statistics.ResetStatistics()
				} else {
					m.stateGame = gameResult

					if gameResult.GetResult() != nil {
						m.statistics.UpdateStatistics(*gameResult.GetResult(), puntobanco.BetType(m.selectedOption))
					}
				}

				m.stateUI = stateIsAfterRound
				m.cursor = 0
				return m, nil
			}

			return m, tick()
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	switch m.stateUI {
	case stateIsBetting:
		// Header
		s += "Make your bet:\n\n"

		for i, choice := range m.bettingOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// Show statistics if enabled
		if m.showStatistics {
			s += fmt.Sprintf("\n%s", rendering.RenderStatisticsTable(&m.statistics))
		}

	case stateIsProgress:
		// Show spinner
		s += fmt.Sprintf("%s Drawing cards...\n", m.spinner.View())

	case stateIsAfterRound:
		// Header
		s += fmt.Sprintf("You bet on %s", m.selectedOption)

		// Show game result state
		s += rendering.RenderGameResultState(&m.stateGame, m.selectedOption)

		for i, choice := range m.afterRoundOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}

		// Show statistics if enabled
		if m.showStatistics {
			s += fmt.Sprintf("\n%s", rendering.RenderStatisticsTable(&m.statistics))
		}
	}

	// Footer with help
	s += fmt.Sprintf("\n\n%s", m.help.FullHelpView(m.keys.FullHelp()))

	return s
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, UI error has happened: %v\n", err)
		os.Exit(1)
	}
}
