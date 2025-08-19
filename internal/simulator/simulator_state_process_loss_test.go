package simulator

import (
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestSimulatorStateProcessLoss_FlatBettingStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "Bet on Punto",
			strategy: BetOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "Bet on Banco",
			strategy: BetOnBanco,
			betType:  puntobanco.BancoBanker,
		},
		{
			name:     "Bet on Égalité",
			strategy: BetOnEgalite,
			betType:  puntobanco.EgaliteTie,
		},
		{
			name:     "Bet on Last Hand",
			strategy: BetOnLastHand,
			betType:  puntobanco.PuntoPlayer,
		},
		// Don't test "Bet on random" cause it is a case of "Bet on Punto/Banco"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			initialConsecutiveLosses := state.ConsecutiveLosses
			initialBetAmount := state.BetAmount

			state.BettingOn = tt.betType
			state.ProcessLoss(tt.strategy)

			if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
				t.Errorf("ConsecutiveLosses should increment: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
			}
			if state.ConsecutiveWins != 0 {
				t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
			}
			if state.MaxConsecutiveLosses != 1 {
				t.Errorf("MaxConsecutiveLosses should updat: got %d, want 1", state.MaxConsecutiveLosses)
			}
			if state.BetAmount != initialBetAmount {
				t.Errorf("BetAmount should not change: got %.2f, want %.2f", state.BetAmount, initialBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessLoss_MartingaleStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "Martingale on Punto",
			strategy: MartingaleOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "Martingale on Banco",
			strategy: MartingaleOnBanco,
			betType:  puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			initialConsecutiveLosses := state.ConsecutiveLosses
			initialBetAmount := state.BetAmount

			state.BettingOn = tt.betType
			state.ProcessLoss(tt.strategy)

			if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
				t.Errorf("ConsecutiveLosses should increment: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
			}
			if state.ConsecutiveWins != 0 {
				t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
			}
			if state.MaxConsecutiveLosses != 1 {
				t.Errorf("MaxConsecutiveLosses should update: got %d, want 1", state.MaxConsecutiveLosses)
			}
			if state.BetAmount != initialBetAmount*2 {
				t.Errorf("BetAmount should not double: got %.2f, want %.2f", state.BetAmount, initialBetAmount*2)
			}
		})
	}
}

func TestSimulatorStateProcessLoss_FibonacciStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "Fibonacci on Punto",
			strategy: FibonacciOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "Fibonacci on Banco",
			strategy: FibonacciOnBanco,
			betType:  puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			initialConsecutiveLosses := state.ConsecutiveLosses
			initialFibonacciIndex := state.FibonacciSequenceIndex
			initialFibonacciProfit := state.FibonacciProfit

			state.BettingOn = tt.betType
			state.ProcessLoss(tt.strategy)

			if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
				t.Errorf("ConsecutiveLosses should incremente: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
			}
			if state.ConsecutiveWins != 0 {
				t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
			}
			if state.FibonacciSequenceIndex != initialFibonacciIndex+1 {
				t.Errorf("Fibonacci sequence index should incremente: got %d, want %d", state.FibonacciSequenceIndex, initialFibonacciIndex+1)
			}

			expectedLossUnits := state.BetAmount / MinimumBet
			expectedProfit := initialFibonacciProfit - expectedLossUnits
			if state.FibonacciProfit != expectedProfit {
				t.Errorf("Fibonacci profit should decrease: got %.2f, want %.2f", state.FibonacciProfit, expectedProfit)
			}

			expectedBetAmount := float64(GetFibonacciValue(state.FibonacciSequenceIndex)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should update: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessLoss_ParoliStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "Paroli on Punto",
			strategy: ParoliOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "Paroli on Banco",
			strategy: ParoliOnBanco,
			betType:  puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			initialConsecutiveLosses := state.ConsecutiveLosses

			// Precondition of Paroli progression
			state.IsInParoliProgression = true
			state.ParoliProgressionLevel = 2
			state.BetAmount = state.BaseBetAmount * 2

			state.BettingOn = tt.betType
			state.ProcessLoss(tt.strategy)

			if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
				t.Errorf("ConsecutiveLosses should increment: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
			}
			if state.ConsecutiveWins != 0 {
				t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
			}
			if state.IsInParoliProgression {
				t.Error("Paroli progression should reset")
			}
			if state.ParoliProgressionLevel != 0 {
				t.Errorf("Paroli progression level should reset: got %d, want 0", state.ParoliProgressionLevel)
			}
			if state.BetAmount != state.BaseBetAmount {
				t.Errorf("BetAmount should reset to base bet: got %.2f, want %.2f", state.BetAmount, state.BaseBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessLoss_DAlembertStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "D'Alembert on Punto",
			strategy: DAlembertOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "D'Alembert on Banco",
			strategy: DAlembertOnBanco,
			betType:  puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			initialConsecutiveLosses := state.ConsecutiveLosses
			initialDAlembertLevel := state.DAlembertLevel

			state.BettingOn = tt.betType
			state.ProcessLoss(tt.strategy)

			if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
				t.Errorf("ConsecutiveLosses should increment: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
			}
			if state.ConsecutiveWins != 0 {
				t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
			}
			if state.DAlembertLevel != initialDAlembertLevel+1 {
				t.Errorf("D'Alembert level should increment: got %d, want %d", state.DAlembertLevel, initialDAlembertLevel+1)
			}

			expectedBetAmount := state.DAlembertUnitSize + (float64(state.DAlembertLevel) * state.DAlembertUnitSize)
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should update: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessLoss_OneThreeTwoSixStrategies(t *testing.T) {
	tests := []struct {
		name     string
		strategy StrategyType
		betType  puntobanco.BetType
	}{
		{
			name:     "1-3-2-6 on Punto",
			strategy: OneThreeTwoSixOnPunto,
			betType:  puntobanco.PuntoPlayer,
		},
		{
			name:     "1-3-2-6 on Banco",
			strategy: OneThreeTwoSixOnBanco,
			betType:  puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("Loss from initial position (1 unit)", func(t *testing.T) {
				state := NewSimulatorState()
				initialConsecutiveLosses := state.ConsecutiveLosses

				state.BettingOn = tt.betType
				state.ProcessLoss(tt.strategy)

				if state.ConsecutiveLosses != initialConsecutiveLosses+1 {
					t.Errorf("ConsecutiveLosses should increment: got %d, want %d", state.ConsecutiveLosses, initialConsecutiveLosses+1)
				}
				if state.ConsecutiveWins != 0 {
					t.Errorf("ConsecutiveWins should reset: got %d, want 0", state.ConsecutiveWins)
				}
				if state.OneThreeTwoSixSequenceIndex != 0 {
					t.Errorf("Sequence index should reset: got %d, want 0", state.OneThreeTwoSixSequenceIndex)
				}
				expectedBetAmount := float64(GetOneThreeTwoSixValue(0)) * MinimumBet
				if state.BetAmount != expectedBetAmount {
					t.Errorf("BetAmount should reset: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
				}
			})

			t.Run("Loss from non first position", func(t *testing.T) {
				state := NewSimulatorState()

				// Precondition of position 2 (2 units)
				state.OneThreeTwoSixSequenceIndex = 2
				state.BetAmount = float64(GetOneThreeTwoSixValue(2)) * MinimumBet

				state.ProcessLoss(tt.strategy)

				if state.OneThreeTwoSixSequenceIndex != 0 {
					t.Errorf("Sequence index should reset: got %d, want 0", state.OneThreeTwoSixSequenceIndex)
				}
				expectedBetAmount := float64(GetOneThreeTwoSixValue(0)) * MinimumBet
				if state.BetAmount != expectedBetAmount {
					t.Errorf("BetAmount should reset: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
				}
			})
		})
	}
}
