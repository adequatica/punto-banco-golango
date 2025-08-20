package rendering

import (
	"strings"
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/simulator"
)

func TestFormatCurrency(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "zero value",
			value: 0.0,
			want:  "$0.00",
		},
		{
			name:  "positive integer",
			value: 100.0,
			want:  "$100.00",
		},
		{
			name:  "positive decimal",
			value: 123.45,
			want:  "$123.45",
		},
		{
			name:  "negative value",
			value: -67.89,
			want:  "$-67.89",
		},
		{
			name:  "large number",
			value: 9999999.99,
			want:  "$9999999.99",
		},
		{
			name:  "many decimal places",
			value: 10.1234567890,
			want:  "$10.12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCurrency(tt.value)
			if result != tt.want {
				t.Errorf("FormatCurrency(%f) = %s, want %s", tt.value, result, tt.want)
			}
		})
	}
}

func TestFormatPercentage(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{
			name:  "zero value",
			value: 0.0,
			want:  "0.00%",
		},
		{
			name:  "positive integer",
			value: 50.0,
			want:  "50.00%",
		},
		{
			name:  "positive decimal",
			value: 25.75,
			want:  "25.75%",
		},
		{
			name:  "negative value",
			value: -10.5,
			want:  "-10.50%",
		},
		{
			name:  "large number",
			value: 999.99,
			want:  "999.99%",
		},
		{
			name:  "many decimal places",
			value: 33.333333,
			want:  "33.33%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPercentage(tt.value)
			if result != tt.want {
				t.Errorf("FormatPercentage(%f) = %s, want %s", tt.value, result, tt.want)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name    string
		seconds float64
		want    string
	}{
		{
			name:    "zero seconds",
			seconds: 0.0,
			want:    "0.00 seconds",
		},
		{
			name:    "less than 1 second",
			seconds: 0.5,
			want:    "0.50 seconds",
		},
		{
			name:    "exactly 1 second",
			seconds: 1.0,
			want:    "1.00 seconds",
		},
		{
			name:    "less than 60 seconds",
			seconds: 43.21,
			want:    "43.21 seconds",
		},
		{
			name:    "exactly 60 seconds",
			seconds: 60.0,
			want:    "1 minutes",
		},
		{
			name:    "1 minute with remaining seconds",
			seconds: 90.5,
			want:    "1 minutes 30.50 seconds",
		},
		{
			name:    "large number of minutes",
			seconds: 3661.0,
			want:    "61 minutes 1.00 seconds",
		},
		{
			name:    "fractional minutes",
			seconds: 61.123,
			want:    "1 minutes 1.12 seconds",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatDuration(tt.seconds)
			if result != tt.want {
				t.Errorf("FormatDuration(%f) = %s, want %s", tt.seconds, result, tt.want)
			}
		})
	}
}

func TestRenderSimulatorStatistics(t *testing.T) {
	tests := []struct {
		name           string
		stats          *simulator.MultipleSimulationsStats
		strategy       simulator.StrategyType
		numSimulations int
		duration       float64
		wantContains   []string
		wantNotContain []string
	}{
		{
			name:           "with stats data",
			stats:          &simulator.MultipleSimulationsStats{TotalSimulations: 100},
			strategy:       simulator.BetOnPunto,
			numSimulations: 100,
			duration:       5.5,
			wantContains: []string{
				"Results for Bet on Punto (player) strategy (100 simulations)",
				"Simulation completed in: 5.50 seconds",
			},
		},
		{
			name:           "without stats data - nil stats",
			stats:          nil,
			strategy:       simulator.BetOnBanco,
			numSimulations: 100,
			duration:       5.5,
			wantContains: []string{
				noSimulationsYet,
			},
		},
		{
			name:           "without stats data - zero simulations",
			stats:          &simulator.MultipleSimulationsStats{TotalSimulations: 0},
			strategy:       simulator.BetOnEgalite,
			numSimulations: 10,
			duration:       0,
			wantContains: []string{
				noSimulationsYet,
			},
		},
		{
			name:           "without number of simulations",
			stats:          &simulator.MultipleSimulationsStats{TotalSimulations: 1},
			strategy:       simulator.BetOnLastHand,
			numSimulations: 0,
			duration:       0,
			wantContains: []string{
				noSimulationsYet,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSimulatorStatistics(tt.stats, tt.strategy, tt.numSimulations, tt.duration)

			for _, expected := range tt.wantContains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderSimulatorStatistics() result should contain: %s", expected)
				}
			}
		})
	}
}

func TestRenderSimulatorTable(t *testing.T) {
	tests := []struct {
		name           string
		stats          *simulator.MultipleSimulationsStats
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "with stats data",
			stats: &simulator.MultipleSimulationsStats{
				TotalSimulations:            100,
				AvgRoundsPerGame:            50.5,
				MinRoundsPlayed:             30,
				MaxRoundsPlayed:             70,
				AvgWinsPerGames:             25.4,
				MinWins:                     15,
				MaxWins:                     35,
				WinRate:                     50.3,
				GamesWithZeroWins:           1,
				ZeroWinsRate:                1.0,
				AvgMaxWinsStreak:            3.2,
				MaxWinsStreak:               4,
				AvgMaxLossStreak:            2.1,
				MaxLossStreak:               1,
				AvgMaxBankrollReached:       1010.0,
				MaxBankrollReacorded:        1100.0,
				GamesWithProfitableBankroll: 1,
				ProfitableBankrollRate:      1.0,
				GamesWithProfitableEnd:      1,
				ProfitableEndGamesRate:      1.0,
			},
			wantContains: []string{
				"Statistics category",
			},
		},
		{
			name:  "without stats data - nil stats",
			stats: nil,
			wantContains: []string{
				noSimulationsYet,
			},
		},
		{
			name:  "without stats data - zero simulations",
			stats: &simulator.MultipleSimulationsStats{TotalSimulations: 0},
			wantContains: []string{
				noSimulationsYet,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RenderSimulatorTable(tt.stats)

			for _, expected := range tt.wantContains {
				if !strings.Contains(result, expected) {
					t.Errorf("RenderSimulatorTable() result should contain: %s", expected)
				}
			}
		})
	}
}
