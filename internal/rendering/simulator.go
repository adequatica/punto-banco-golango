package rendering

import (
	"fmt"

	"github.com/adequatica/punto-banco-golango/internal/simulator"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

var (
	noBorderStyle = lipgloss.NewStyle().BorderStyle(lipgloss.HiddenBorder())
)

func FormatCurrency(value float64) string {
	return fmt.Sprintf("$%.2f", value)
}

func FormatPercentage(value float64) string {
	return fmt.Sprintf("%.2f%%", value)
}

func FormatDuration(seconds float64) string {
	if seconds < 60 {
		return fmt.Sprintf("%.2f seconds", seconds)
	}

	minutes := int(seconds / 60)
	remainingSeconds := seconds - float64(minutes*60)

	if remainingSeconds == 0 {
		return fmt.Sprintf("%d minutes", minutes)
	}

	return fmt.Sprintf("%d minutes %.2f seconds", minutes, remainingSeconds)
}

const noSimulationsYet = "No simulations run yet"

func RenderSimulatorStatistics(stats *simulator.MultipleSimulationsStats, strategy simulator.StrategyType, numSimulations int, duration float64) string {
	if stats == nil || stats.TotalSimulations == 0 || numSimulations == 0 {
		return noSimulationsYet
	}

	header := fmt.Sprintf("Results for %s strategy (%d simulations)\n", strategy, numSimulations)
	header += fmt.Sprintf("Simulation completed in: %s\n", FormatDuration(duration))

	table := RenderSimulatorTable(stats)

	return header + table
}

func RenderSimulatorTable(stats *simulator.MultipleSimulationsStats) string {
	if stats == nil || stats.TotalSimulations == 0 {
		return noSimulationsYet
	}

	columns := []table.Column{
		{Title: "Statistics category", Width: 34},
		{Title: "Value", Width: 12},
	}

	rows := []table.Row{
		// Games played statistics
		{"Rounds played, average per game", FormatFloat(stats.AvgRoundsPlayed)},
		{"Rounds played, min per game", fmt.Sprintf("%d", stats.MinRoundsPlayed)},
		{"Rounds played, max per game", fmt.Sprintf("%d", stats.MaxRoundsPlayed)},
		{"", ""},
		// Wins statistics
		{"Wins, average per game", FormatFloat(stats.AvgWins)},
		{"Wins, min per game", fmt.Sprintf("%d", stats.MinWins)},
		{"Wins, max per game", fmt.Sprintf("%d", stats.MaxWins)},
		// Win rate statistics
		{"Win rate", FormatPercentage(stats.WinRate)},
		{"Zero wins games rate", FormatPercentage(stats.ZeroWinsRate)},
		{"", ""},
		// Consecutive wins statistics
		{"Consecutive wins, average", FormatFloat(stats.AvgMaxConsecutiveWins)},
		{"Maximum consecutive wins", fmt.Sprintf("%d", stats.MaxConsecutiveWins)},
		// Consecutive losses statistics
		{"Consecutive losses, average", FormatFloat(stats.AvgMaxConsecutiveLosses)},
		{"Maximum consecutive losses", fmt.Sprintf("%d", stats.MaxConsecutiveLosses)},
		{"", ""},
		// Budget statistics
		{"Maximum budget reached, average", FormatCurrency(stats.AvgMaxBudgetReached)},
		{"Maximum recorded budget", FormatCurrency(stats.MaxBudgetReacorded)},
		{"Profitable games", FormatPercentage(stats.ProfitableBudgetRate)},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)+1),
	)

	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		Bold(true)
	styles.Selected = styles.Selected.
		// Reset default selected cell styles
		UnsetForeground().
		Bold(false)
	t.SetStyles(styles)

	return noBorderStyle.Render(t.View())
}
