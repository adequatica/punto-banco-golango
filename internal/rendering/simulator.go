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
		{Title: "Statistics category", Width: 36},
		{Title: "Value", Width: 10},
	}

	rows := []table.Row{
		// Games played statistics
		{"Mean rounds per game", FormatFloat(stats.AvgRoundsPerGame)},
		{"Minimum played rounds per game", fmt.Sprintf("%d", stats.MinRoundsPlayed)},
		{"Maximum played rounds per game", fmt.Sprintf("%d", stats.MaxRoundsPlayed)},
		{"", ""},
		// Wins statistics
		{"Mean wins per game", FormatFloat(stats.AvgWinsPerGames)},
		{"Minimum wins per game", fmt.Sprintf("%d", stats.MinWins)},
		{"Maximum wins per game", fmt.Sprintf("%d", stats.MaxWins)},
		// Win rate statistics
		{"Win rate", FormatPercentage(stats.WinRate)},
		{"Rate of zero-wins games", FormatPercentage(stats.ZeroWinsRate)},
		{"", ""},
		// Streaks statistics
		{"Mean winning streak", FormatFloat(stats.AvgMaxWinsStreak)},
		{"Maximum winning streak", fmt.Sprintf("%d", stats.MaxWinsStreak)},
		{"Mean losing streak", FormatFloat(stats.AvgMaxLossStreak)},
		{"Maximum losing streak", fmt.Sprintf("%d", stats.MaxLossStreak)},
		{"", ""},
		// Bankroll statistics
		{"Mean peak bankroll per game", FormatCurrency(stats.AvgMaxBankrollReached)},
		{"Maximum recorded bankroll", FormatCurrency(stats.MaxBankrollReacorded)},
		{"Profitable games", FormatPercentage(stats.ProfitableBankrollRate)},
		{"Profitably ended games", FormatPercentage(stats.ProfitableEndGamesRate)},
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
