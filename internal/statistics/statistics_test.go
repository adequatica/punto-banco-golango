package statistics

import (
	"reflect"
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestNewSessionStatistics(t *testing.T) {
	stats := NewSessionStatistics()

	if stats.TotalRounds != 0 {
		t.Errorf("TotalRounds should be 0, got %d", stats.TotalRounds)
	}
	if stats.PuntoWins != 0 {
		t.Errorf("PuntoWins should be 0, got %d", stats.PuntoWins)
	}
	if stats.BancoWins != 0 {
		t.Errorf("BancoWins should be 0, got %d", stats.BancoWins)
	}
	if stats.Ties != 0 {
		t.Errorf("Ties should be 0, got %d", stats.Ties)
	}
	if stats.UserWins != 0 {
		t.Errorf("UserWins should be 0, got %d", stats.UserWins)
	}
	if stats.UserBets == nil {
		t.Error("UserBets should be initialized, got nil")
	}
	if len(stats.UserBets) != 0 {
		t.Errorf("UserBets should be empty, got %d items", len(stats.UserBets))
	}
}

func TestUpdateStatistics(t *testing.T) {
	tests := []struct {
		name       string
		gameResult puntobanco.BetType
		userBet    puntobanco.BetType
		want       SessionStatistics
	}{
		{
			name:       "Punto wins, user bets Punto (win)",
			gameResult: puntobanco.PuntoPlayer,
			userBet:    puntobanco.PuntoPlayer,
			want: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   1,
				BancoWins:   0,
				Ties:        0,
				UserWins:    1,
				UserBets:    map[puntobanco.BetType]int{puntobanco.PuntoPlayer: 1},
			},
		},
		{
			name:       "Banco wins, user bets Banco (win)",
			gameResult: puntobanco.BancoBanker,
			userBet:    puntobanco.BancoBanker,
			want: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   0,
				BancoWins:   1,
				Ties:        0,
				UserWins:    1,
				UserBets:    map[puntobanco.BetType]int{puntobanco.BancoBanker: 1},
			},
		},
		{
			name:       "Punto wins, user bets Banco (lose)",
			gameResult: puntobanco.PuntoPlayer,
			userBet:    puntobanco.BancoBanker,
			want: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   1,
				BancoWins:   0,
				Ties:        0,
				UserWins:    0,
				UserBets:    map[puntobanco.BetType]int{puntobanco.BancoBanker: 1},
			},
		},
		{
			name:       "Banco wins, user bets Punto (lose)",
			gameResult: puntobanco.BancoBanker,
			userBet:    puntobanco.PuntoPlayer,
			want: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   0,
				BancoWins:   1,
				Ties:        0,
				UserWins:    0,
				UserBets:    map[puntobanco.BetType]int{puntobanco.PuntoPlayer: 1},
			},
		},
		{
			name:       "Tie, user bets Égalité (win)",
			gameResult: puntobanco.EgaliteTie,
			userBet:    puntobanco.EgaliteTie,
			want: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   0,
				BancoWins:   0,
				Ties:        1,
				UserWins:    1,
				UserBets:    map[puntobanco.BetType]int{puntobanco.EgaliteTie: 1},
			},
		},
		{
			name:       "invalid game result (statistics doesn't count)",
			gameResult: "invalid",
			userBet:    puntobanco.PuntoPlayer,
			want: SessionStatistics{
				TotalRounds: 0,
				PuntoWins:   0,
				BancoWins:   0,
				Ties:        0,
				UserWins:    0,
				UserBets:    map[puntobanco.BetType]int{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stats := NewSessionStatistics()
			stats.UpdateStatistics(tt.gameResult, tt.userBet)

			if stats.TotalRounds != tt.want.TotalRounds {
				t.Errorf("TotalRounds = %d, got %d", stats.TotalRounds, tt.want.TotalRounds)
			}
			if stats.PuntoWins != tt.want.PuntoWins {
				t.Errorf("PuntoWins = %d, got %d", stats.PuntoWins, tt.want.PuntoWins)
			}
			if stats.BancoWins != tt.want.BancoWins {
				t.Errorf("BancoWins = %d, got %d", stats.BancoWins, tt.want.BancoWins)
			}
			if stats.Ties != tt.want.Ties {
				t.Errorf("Ties = %d, got %d", stats.Ties, tt.want.Ties)
			}
			if stats.UserWins != tt.want.UserWins {
				t.Errorf("UserWins = %d, got %d", stats.UserWins, tt.want.UserWins)
			}
			if len(stats.UserBets) != len(tt.want.UserBets) {
				t.Errorf("UserBets length = %d, got %d", len(stats.UserBets), len(tt.want.UserBets))
			}
			for betType, count := range tt.want.UserBets {
				if stats.UserBets[betType] != count {
					t.Errorf("UserBets[%s] = %d, got %d", betType, stats.UserBets[betType], count)
				}
			}
		})
	}
}

func TestUpdateStatisticsMultipleRounds(t *testing.T) {
	stats := NewSessionStatistics()

	rounds := []struct {
		gameResult puntobanco.BetType
		userBet    puntobanco.BetType
	}{
		{puntobanco.PuntoPlayer, puntobanco.PuntoPlayer}, // Win
		{puntobanco.PuntoPlayer, puntobanco.BancoBanker}, // Lose
		{puntobanco.PuntoPlayer, puntobanco.EgaliteTie},  // Lose

		{puntobanco.BancoBanker, puntobanco.PuntoPlayer}, // Lose
		{puntobanco.BancoBanker, puntobanco.BancoBanker}, // Win
		{puntobanco.BancoBanker, puntobanco.EgaliteTie},  // Lose

		{puntobanco.EgaliteTie, puntobanco.PuntoPlayer}, // Lose
		{puntobanco.EgaliteTie, puntobanco.BancoBanker}, // Lose
		{puntobanco.EgaliteTie, puntobanco.EgaliteTie},  // Tie (not win)
	}

	for _, round := range rounds {
		stats.UpdateStatistics(round.gameResult, round.userBet)
	}

	want := SessionStatistics{
		TotalRounds: 9,
		PuntoWins:   3,
		BancoWins:   3,
		Ties:        3,
		UserWins:    3,
		UserBets: map[puntobanco.BetType]int{
			puntobanco.PuntoPlayer: 3,
			puntobanco.BancoBanker: 3,
			puntobanco.EgaliteTie:  3,
		},
	}

	if !reflect.DeepEqual(stats, want) {
		t.Errorf("should have correct statistics after multiple rounds, got %+v, want %+v", stats, want)
	}
}

func TestGetPuntoWinsPercentage(t *testing.T) {
	tests := []struct {
		name  string
		stats SessionStatistics
		want  float64
	}{
		{
			name:  "no rounds played",
			stats: NewSessionStatistics(),
			want:  0.0,
		},
		{
			name: "100% Punto wins",
			stats: SessionStatistics{
				TotalRounds: 10,
				PuntoWins:   10,
				BancoWins:   0,
				Ties:        0,
			},
			want: 100.0,
		},
		{
			name: "50% Punto wins",
			stats: SessionStatistics{
				TotalRounds: 10,
				PuntoWins:   5,
				BancoWins:   3,
				Ties:        2,
			},
			want: 50.0,
		},

		{
			name: "0% Punto wins",
			stats: SessionStatistics{
				TotalRounds: 5,
				PuntoWins:   0,
				BancoWins:   3,
				Ties:        2,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.GetPuntoWinsPercentage()
			if result != tt.want {
				t.Errorf("GetPuntoWinsPercentage() = %.1f, got %.1f", result, tt.want)
			}
		})
	}
}

func TestGetBancoWinsPercentage(t *testing.T) {
	tests := []struct {
		name  string
		stats SessionStatistics
		want  float64
	}{
		{
			name:  "no rounds played",
			stats: NewSessionStatistics(),
			want:  0.0,
		},
		{
			name: "100% Banco wins",
			stats: SessionStatistics{
				TotalRounds: 10,
				PuntoWins:   0,
				BancoWins:   10,
				Ties:        0,
			},
			want: 100.0,
		},
		{
			name: "60% Banco wins",
			stats: SessionStatistics{
				TotalRounds: 10,
				PuntoWins:   2,
				BancoWins:   6,
				Ties:        2,
			},
			want: 60.0,
		},
		{
			name: "0% Banco wins",
			stats: SessionStatistics{
				TotalRounds: 6,
				PuntoWins:   4,
				BancoWins:   0,
				Ties:        2,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.GetBancoWinsPercentage()
			if result != tt.want {
				t.Errorf("GetBancoWinsPercentage() = %.1f, got %.1f", result, tt.want)
			}
		})
	}
}

func TestGetTiesPercentage(t *testing.T) {
	tests := []struct {
		name  string
		stats SessionStatistics
		want  float64
	}{
		{
			name:  "no rounds played",
			stats: NewSessionStatistics(),
			want:  0.0,
		},
		{
			name: "100% Ties",
			stats: SessionStatistics{
				TotalRounds: 1,
				PuntoWins:   0,
				BancoWins:   0,
				Ties:        1,
			},
			want: 100.0,
		},
		{
			name: "25% Ties",
			stats: SessionStatistics{
				TotalRounds: 8,
				PuntoWins:   4,
				BancoWins:   4,
				Ties:        2,
			},
			want: 25.0,
		},
		{
			name: "0% Ties",
			stats: SessionStatistics{
				TotalRounds: 2,
				PuntoWins:   1,
				BancoWins:   1,
				Ties:        0,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.GetTiesPercentage()
			if result != tt.want {
				t.Errorf("GetTiesPercentage() = %.1f, got %.1f", result, tt.want)
			}
		})
	}
}

func TestGetUserWinsPercentage(t *testing.T) {
	tests := []struct {
		name  string
		stats SessionStatistics
		want  float64
	}{
		{
			name:  "no rounds played",
			stats: NewSessionStatistics(),
			want:  0.0,
		},
		{
			name: "100% User wins",
			stats: SessionStatistics{
				TotalRounds: 9,
				UserWins:    9,
			},
			want: 100.0,
		},
		{
			name: "33.3% User wins",
			stats: SessionStatistics{
				TotalRounds: 9,
				UserWins:    3,
			},
			want: 33.33333333333333,
		},
		{
			name: "0% User wins",
			stats: SessionStatistics{
				TotalRounds: 9,
				UserWins:    0,
			},
			want: 0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.GetUserWinsPercentage()
			if result != tt.want {
				t.Errorf("GetUserWinsPercentage() = %.1f, got %.1f", result, tt.want)
			}
		})
	}
}

func TestGetUserBetsDistribution(t *testing.T) {
	tests := []struct {
		name  string
		stats SessionStatistics
		want  map[puntobanco.BetType]float64
	}{
		{
			name:  "no rounds played",
			stats: NewSessionStatistics(),
			want:  map[puntobanco.BetType]float64{},
		},
		{
			name: "mixed betting distribution",
			stats: SessionStatistics{
				TotalRounds: 10,
				UserBets: map[puntobanco.BetType]int{
					puntobanco.PuntoPlayer: 5,
					puntobanco.BancoBanker: 3,
					puntobanco.EgaliteTie:  2,
				},
			},
			want: map[puntobanco.BetType]float64{
				puntobanco.PuntoPlayer: 50.0,
				puntobanco.BancoBanker: 30.0,
				puntobanco.EgaliteTie:  20.0,
			},
		},
		{
			name: "equal distribution",
			stats: SessionStatistics{
				TotalRounds: 6,
				UserBets: map[puntobanco.BetType]int{
					puntobanco.PuntoPlayer: 2,
					puntobanco.BancoBanker: 2,
					puntobanco.EgaliteTie:  2,
				},
			},
			want: map[puntobanco.BetType]float64{
				puntobanco.PuntoPlayer: 33.33333333333333,
				puntobanco.BancoBanker: 33.33333333333333,
				puntobanco.EgaliteTie:  33.33333333333333,
			},
		},
		{
			name: "only one item distribution",
			stats: SessionStatistics{
				TotalRounds: 5,
				UserBets: map[puntobanco.BetType]int{
					puntobanco.PuntoPlayer: 5,
				},
			},
			want: map[puntobanco.BetType]float64{
				puntobanco.PuntoPlayer: 100.0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.stats.GetUserBetsDistribution()

			if len(result) != len(tt.want) {
				t.Errorf("distribution length = %d, got %d", len(result), len(tt.want))
			}
		})
	}
}

func TestResetStatistics(t *testing.T) {
	stats := SessionStatistics{
		TotalRounds: 5,
		PuntoWins:   4,
		BancoWins:   3,
		Ties:        2,
		UserWins:    1,
		UserBets: map[puntobanco.BetType]int{
			puntobanco.PuntoPlayer: 3,
			puntobanco.BancoBanker: 2,
			puntobanco.EgaliteTie:  1,
		},
	}

	stats.ResetStatistics()

	want := NewSessionStatistics()

	if want.TotalRounds != stats.TotalRounds {
		t.Errorf("TotalRounds should be %d after reset, got %d", stats.TotalRounds, want.TotalRounds)
	}
	if want.PuntoWins != stats.PuntoWins {
		t.Errorf("PuntoWins should be %d after reset, got %d", stats.PuntoWins, want.PuntoWins)
	}
	if want.BancoWins != stats.BancoWins {
		t.Errorf("BancoWins should be %d after reset, got %d", stats.BancoWins, want.BancoWins)
	}
	if want.Ties != stats.Ties {
		t.Errorf("Ties should be %d after reset, got %d", stats.Ties, want.Ties)
	}
	if want.UserWins != stats.UserWins {
		t.Errorf("UserWins should be %d after reset, got %d", stats.UserWins, want.UserWins)
	}
	if want.UserBets == nil {
		t.Error("UserBets should be initialized after reset, got nil")
	}
	if len(want.UserBets) != 0 {
		t.Errorf("UserBets should be empty after reset, got %d items", len(want.UserBets))
	}
}
