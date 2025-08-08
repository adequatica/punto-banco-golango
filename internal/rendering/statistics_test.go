package rendering

import (
	"strings"
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
	"github.com/adequatica/punto-banco-golango/internal/statistics"
)

func TestFormatFloat(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "whole number with .0 (zero)",
			input: 0.0,
			want:  "0",
		},
		{
			name:  "whole number with .0 (non-zero)",
			input: 100.0,
			want:  "100",
		},
		{
			name:  "decimal with trailing zeros",
			input: 10.200,
			want:  "10.2",
		},
		{
			name:  "decimal with non-zero fractional part",
			input: 12.34,
			want:  "12.3",
		},
		{
			name:  "decimal with non-zero fractional part with rounding up",
			input: 56.78,
			want:  "56.8",
		},
		{
			name:  "decimal with non-zero fractional part starts with .0",
			input: 10.02,
			want:  "10.0",
		},
		{
			name:  "decimal with non-zero fractional part starts with .0 and rounding up",
			input: 10.06,
			want:  "10.1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatFloat(tt.input)
			if result != tt.want {
				t.Errorf("FormatFloat(%f) = %s, got %s", tt.input, result, tt.want)
			}
		})
	}
}

func TestFormatUserWinsPercentage(t *testing.T) {
	tests := []struct {
		name  string
		input float64
		want  string
	}{
		{
			name:  "zero percentage (red)",
			input: 0.0,
			want:  "0%",
		},
		{
			name:  "lose percentage (red)",
			input: 25.5,
			want:  "25.5%",
		},
		{
			name:  "exactly 50 percent",
			input: 50.0,
			want:  "50%",
		},
		{
			name:  "win percentage percent (green)",
			input: 75.556,
			want:  "75.6%",
		},
		{
			name:  "100 percentage (green)",
			input: 100.0,
			want:  "100%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatUserWinsPercentage(tt.input)

			if !strings.Contains(result, tt.want) {
				t.Errorf("FormatUserWinsPercentage(%f) should contains %s", tt.input, tt.want)
			}
		})
	}
}

func TestRenderStatisticsTable(t *testing.T) {
	t.Run("new game session statistics", func(t *testing.T) {
		newSessionStats := statistics.NewSessionStatistics()

		tests := []struct {
			name  string
			stats *statistics.SessionStatistics
			want  string
		}{
			{
				name:  "nil statistics",
				stats: nil,
				want:  noGamesPlayedYet,
			},
			{
				name:  "empty statistics",
				stats: &statistics.SessionStatistics{},
				want:  noGamesPlayedYet,
			},
			{
				name:  "new session statistics",
				stats: &newSessionStats,
				want:  noGamesPlayedYet,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := RenderStatisticsTable(tt.stats)
				if result != tt.want {
					t.Errorf("RenderStatisticsTable() = %q should be %q", result, tt.want)
				}
			})
		}
	})

	t.Run("continuing game session statistics", func(t *testing.T) {
		stats := &statistics.SessionStatistics{
			TotalRounds: 10,
			PuntoWins:   5,
			BancoWins:   3,
			Ties:        2,
			UserWins:    1,
			UserBets:    make(map[puntobanco.BetType]int),
		}

		result := RenderStatisticsTable(stats)

		if result == "" {
			t.Errorf("RenderStatisticsTable() should not return empty string")
		}
		if result == noGamesPlayedYet {
			t.Errorf("RenderStatisticsTable() should not return no games played yet")
		}
	})
}
