package simulator

import (
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestNewSimulatorState(t *testing.T) {
	newState := NewSimulatorState()

	if newState == nil {
		t.Error("new simulator state should not be nil")
	}
}

func TestGetFibonacciValue(t *testing.T) {
	tests := []struct {
		index int
		want  int
		name  string
	}{
		{0, 1, "Fibonacci(0) should be 1"},
		{1, 1, "Fibonacci(1) should be 1"},
		{2, 2, "Fibonacci(2) should be 2"},
		{3, 3, "Fibonacci(3) should be 3"},
		{4, 5, "Fibonacci(4) should be 5"},
		{5, 8, "Fibonacci(5) should be 8"},
		{6, 13, "Fibonacci(6) should be 13"},
		{7, 21, "Fibonacci(7) should be 21"},
		{8, 34, "Fibonacci(8) should be 34"},
		{9, 55, "Fibonacci(9) should be 55"},
		{10, 89, "Fibonacci(10) should be 89"},
		{11, 144, "Fibonacci(11) should be 144"},
		{12, 233, "Fibonacci(12) should be 233"},
		{13, 377, "Fibonacci(13) should be 377"},
		{14, 610, "Fibonacci(14) should be 610"},
		{15, 987, "Fibonacci(15) should be 987"},
		{-1, 1, "Fibonacci(-1) should be 1 (handles negative)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFibonacciValue(tt.index)
			if result != tt.want {
				t.Errorf("GetFibonacciValue(%d) = %d should be %d", tt.index, result, tt.want)
			}
		})
	}
}

func TestGetOneThreeTwoSixValue(t *testing.T) {
	tests := []struct {
		index int
		want  int
		name  string
	}{
		{0, 1, "1-3-2-6(0) should be 1"},
		{1, 3, "1-3-2-6(1) should be 3"},
		{2, 2, "1-3-2-6(2) should be 2"},
		{3, 6, "1-3-2-6(3) should be 6"},
		{4, 1, "1-3-2-6(4) should be 1 (out of bounds)"},
		{-1, 1, "1-3-2-6(-1) should be 1 (handles negative)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetOneThreeTwoSixValue(tt.index)
			if result != tt.want {
				t.Errorf("GetOneThreeTwoSixValue(%d) = %d should be %d", tt.index, result, tt.want)
			}
		})
	}
}

func TestCalculatePayout(t *testing.T) {
	tests := []struct {
		name      string
		betType   puntobanco.BetType
		betAmount float64
		want      float64
	}{
		{
			name:      "Punto Player bet - even money payout",
			betType:   puntobanco.PuntoPlayer,
			betAmount: 100.0,
			want:      100.0, // 1:1 payout
		},
		{
			name:      "Banco Banker bet - 5% commission",
			betType:   puntobanco.BancoBanker,
			betAmount: 100.0,
			want:      95.0, // 19:20 payout (5% commission)
		},
		{
			name:      "Egalite Tie bet - 8:1 payout",
			betType:   puntobanco.EgaliteTie,
			betAmount: 100.0,
			want:      800.0, // 8:1 payout
		},
		{
			name:      "Zero bet amount",
			betType:   puntobanco.PuntoPlayer,
			betAmount: 0.0,
			want:      0.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePayout(tt.betType, tt.betAmount)
			if result != tt.want {
				t.Errorf("calculatePayout(%v, %.2f) = %.2f, want %.2f",
					tt.betType, tt.betAmount, result, tt.want)
			}
		})
	}
}

func TestSimulatorStateCanPlaceBet(t *testing.T) {
	tests := []struct {
		name          string
		currentBudget float64
		betAmount     float64
		want          bool
	}{
		{
			name:          "can place bet when budget equals bet amount",
			currentBudget: 100.0,
			betAmount:     100.0,
			want:          true,
		},
		{
			name:          "can place bet when budget greater than bet amount",
			currentBudget: 101.0,
			betAmount:     100.0,
			want:          true,
		},
		{
			name:          "cannot place bet when budget less than bet amount",
			currentBudget: 1.0,
			betAmount:     10.0,
			want:          false,
		},
		{
			name:          "cannot place bet when budget is zero",
			currentBudget: 0.0,
			betAmount:     10.0,
			want:          false,
		},
		{
			name:          "can place bet when bet amount is zero",
			currentBudget: 1.0,
			betAmount:     0.0,
			want:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &SimulatorState{
				CurrentBudget: tt.currentBudget,
				BetAmount:     tt.betAmount,
			}
			result := state.CanPlaceBet()
			if result != tt.want {
				t.Errorf("CanPlaceBet() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestSimulatorStatePlaceBet(t *testing.T) {
	tests := []struct {
		name          string
		initialBudget float64
		betAmount     float64
		wantBudget    float64
	}{
		{
			name:          "place bet reduces budget by bet amount",
			initialBudget: 1000.0,
			betAmount:     10.0,
			wantBudget:    990.0,
		},
		{
			name:          "place bet with zero bet amount",
			initialBudget: 1000.0,
			betAmount:     0.0,
			wantBudget:    1000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := &SimulatorState{
				CurrentBudget: tt.initialBudget,
				BetAmount:     tt.betAmount,
			}

			initialBudget := state.CurrentBudget
			state.PlaceBet()

			if state.CurrentBudget != tt.wantBudget {
				t.Errorf("PlaceBet() changed budget from %.2f to %.2f, want %.2f",
					initialBudget, state.CurrentBudget, tt.wantBudget)
			}
		})
	}
}
