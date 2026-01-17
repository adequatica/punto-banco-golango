package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/adequatica/punto-banco-golango/internal/rendering"
	"github.com/adequatica/punto-banco-golango/internal/simulator"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	defaultNumberOfSimulations = 10000
	// Limit of 1M simulations needs just to prevent too long calculations in case of input mistake
	maxNumberOfSimulations = 999999 // This number of simulations take ~ 25 minutes depends on choosen strategy
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Up, k.Down, k.Enter, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down},    // first column
		{k.Enter, k.Quit}, // second column
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
	Quit: key.NewBinding(
		key.WithKeys("q", "Q", "й", "Й", "ctrl+c", "esc"),
		key.WithHelp("Q/CTRL+C", "— quit"),
	),
}

type UIstate int

const (
	stateSelectStrategy UIstate = iota
	stateEnterSimulations
	stateRunningSimulation
	stateShowResults
)

type model struct {
	stateUI            UIstate
	cursor             int
	strategyOptions    []string
	selectedStrategy   simulator.StrategyType
	textInput          textinput.Model
	numSimulations     int
	saveData           bool
	stats              simulator.MultipleSimulationsStats
	keys               keyMap
	help               help.Model
	spinner            spinner.Model
	simulationStart    time.Time
	simulationDuration time.Duration
}

func InitialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("%d", defaultNumberOfSimulations)
	ti.Focus()
	ti.CharLimit = 6 // Prevent input above maxNumberOfSimulations
	ti.Width = 6

	return model{
		stateUI:         stateSelectStrategy,
		cursor:          0,
		strategyOptions: simulator.GetStrategyOptions(),
		textInput:       ti,
		numSimulations:  0,
		keys:            defaultKeys,
		help:            help.New(),
		spinner:         s,
	}
}

// Simulation completion message
type simulationCompleteMsg struct {
	stats simulator.MultipleSimulationsStats
	err   error
}

func runSimulation(strategy simulator.StrategyType, numSimulations int, saveData bool) tea.Cmd {
	return func() tea.Msg {
		// Run the simulation with error handling
		stats := simulator.RunMultipleSimulations(strategy, numSimulations, saveData)
		// Note: If RunMultipleSimulations could return an error, we would handle it here
		return simulationCompleteMsg{stats: stats, err: nil}
	}
}

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
			case stateSelectStrategy:
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.strategyOptions) - 1
				}
			case stateEnterSimulations:
				// Toggle save data option when Up is pressed
				if num, err := strconv.Atoi(m.textInput.Value()); err == nil && num > 0 && num <= 1000 {
					m.saveData = !m.saveData
				}
			}

		case key.Matches(msg, m.keys.Down):
			switch m.stateUI {
			case stateSelectStrategy:
				if m.cursor < len(m.strategyOptions)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
			case stateEnterSimulations:
				// Toggle save data option when Down is pressed
				if num, err := strconv.Atoi(m.textInput.Value()); err == nil && num > 0 && num <= 1000 {
					m.saveData = !m.saveData
				}
			}

		case key.Matches(msg, m.keys.Enter):
			switch m.stateUI {
			case stateSelectStrategy:
				// Store the selected strategy with bounds checking
				if len(m.strategyOptions) > 0 && m.cursor >= 0 && m.cursor < len(m.strategyOptions) {
					m.selectedStrategy = simulator.StrategyType(m.strategyOptions[m.cursor])
					m.stateUI = stateEnterSimulations
					m.textInput.SetValue(fmt.Sprintf("%d", defaultNumberOfSimulations))
					m.textInput.Focus()
					m.saveData = false // Reset to default NO
				} else {
					// Handle invalid state
					m.cursor = 0
				}

			case stateEnterSimulations:
				// Parse number of simulations with validation
				if num, err := strconv.Atoi(m.textInput.Value()); err == nil && num > 0 && num <= maxNumberOfSimulations {
					m.numSimulations = num
					// Only save data if <= 1000 simulations
					if num > 1000 {
						m.saveData = false
					}
					m.stateUI = stateRunningSimulation
					m.simulationStart = time.Now()
					// Start running simulation
					return m, tea.Batch(
						m.spinner.Tick,
						runSimulation(m.selectedStrategy, m.numSimulations, m.saveData),
					)
				}

			case stateShowResults:
				// Return to strategy selection with complete reset
				m.stateUI = stateSelectStrategy
				m.cursor = 0
				m.selectedStrategy = ""
				m.textInput.SetValue("")
				m.numSimulations = 0
				m.saveData = false
				m.stats = simulator.MultipleSimulationsStats{}
				m.simulationDuration = 0
				m.simulationStart = time.Time{}
			}

		default:
			// Handle text input for number of simulations
			if m.stateUI == stateEnterSimulations {
				var cmd tea.Cmd
				m.textInput, cmd = m.textInput.Update(msg)
				return m, cmd
			}
		}

	// Spinner tick
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// Simulation completion
	case simulationCompleteMsg:
		if m.stateUI == stateRunningSimulation {
			if msg.err != nil {
				// Handle simulation error - could add error state here
				fmt.Printf("Simulation error: %v\n", msg.err)
				m.stateUI = stateSelectStrategy
			} else {
				m.stats = msg.stats
				m.simulationDuration = time.Since(m.simulationStart)
				m.stateUI = stateShowResults
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	var s string

	switch m.stateUI {
	case stateSelectStrategy:
		s += "Select a betting strategy:\n\n"

		for i, strategy := range m.strategyOptions {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, strategy)
		}

	case stateEnterSimulations:
		s += fmt.Sprintf("Selected strategy: %s\n\n", m.selectedStrategy)
		s += "Enter number of simulations to run:\n"
		s += m.textInput.View()

		// Show save data option only when number of simulations is <= 1000
		// Cause storing data for large numbers of simulations may cause memory exhaustion
		if num, err := strconv.Atoi(m.textInput.Value()); err == nil && num > 0 && num <= 1000 {
			saveStatus := "NO"
			if m.saveData {
				saveStatus = "YES"
			}
			s += fmt.Sprintf("\n\nSave data into a file: %s (Press ↑/↓ to toggle)", saveStatus)
		}

		s += "\n\nPress ENTER to start simulation"

	case stateRunningSimulation:
		s += fmt.Sprintf("Running %d simulations for %s\n\n", m.numSimulations, m.selectedStrategy)
		s += fmt.Sprintf("%s Simulation in progress...\n", m.spinner.View())

	case stateShowResults:
		s += rendering.RenderSimulatorStatistics(&m.stats, m.selectedStrategy, m.numSimulations, m.simulationDuration.Seconds())
		s += "\nPress ENTER to run another simulation"
	}

	// Footer with help
	s += fmt.Sprintf("\n\n%s", m.help.FullHelpView(m.keys.FullHelp()))

	return s
}

func main() {
	p := tea.NewProgram(InitialModel())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, UI error has happened: %v\n", err)
		os.Exit(1)
	}
}
