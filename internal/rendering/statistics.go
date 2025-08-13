package rendering

import (
	"fmt"

	"github.com/adequatica/punto-banco-golango/internal/statistics"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

// Removes the fractional part when it is .0
func FormatFloat(value float64) string {
	if value == float64(int(value)) {
		return fmt.Sprintf("%.0f", value)
	}

	return fmt.Sprintf("%.1f", value)
}

var (
	resetStyle = lipgloss.NewStyle().UnsetForeground()
)

func FormatUserWinsPercentage(value float64) string {
	userWinsPercentage := fmt.Sprintf("%s%%", FormatFloat(value))

	// Color green, if user wins over 50%
	if value > 50 {
		return greenStyle.Render(userWinsPercentage) + resetStyle.Render("")
	}
	// Color red, if user loses over 50%
	if value < 50 {
		return redStyle.Render(userWinsPercentage) + resetStyle.Render("")
	}
	return userWinsPercentage
}

var baseStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	Padding(1, 2)

const noGamesPlayedYet = "No games played yet"

func RenderStatisticsTable(s *statistics.SessionStatistics) string {
	if s == nil || s.TotalRounds == 0 {
		return noGamesPlayedYet
	}

	columns := []table.Column{
		{Title: "Game session statistics", Width: 26},
		{Title: "Value", Width: 8},
		{Title: "Percentage", Width: 12},
	}

	rows := []table.Row{
		{
			"Total rounds",
			fmt.Sprintf("%d", s.TotalRounds),
		},
		{
			"Punto wins",
			fmt.Sprintf("%d", s.PuntoWins),
			fmt.Sprintf("%s%%", FormatFloat(s.GetPuntoWinsPercentage())),
		},
		{
			"Banco wins",
			fmt.Sprintf("%d", s.BancoWins),
			fmt.Sprintf("%s%%", FormatFloat(s.GetBancoWinsPercentage())),
		},
		{
			"Ties",
			fmt.Sprintf("%d", s.Ties),
			fmt.Sprintf("%s%%", FormatFloat(s.GetTiesPercentage())),
		},
		{
			"Your wins",
			fmt.Sprintf("%d", s.UserWins),
			FormatUserWinsPercentage(s.GetUserWinsPercentage()),
		},
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

	return baseStyle.Render(t.View())
}
