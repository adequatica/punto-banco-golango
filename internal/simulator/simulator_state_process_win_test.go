package simulator

import (
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestSimulatorStateProcessWin_FlatBettingStrategies(t *testing.T) {
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
			name:     "Bet on last hand",
			strategy: BetOnLastHand,
			betType:  puntobanco.PuntoPlayer,
		},
		// Don't test "Bet on random" cause it is a case of "Bet on Punto/Banco"
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := NewSimulatorState()
			startingBankroll := state.CurrentBankroll
			initialWins := state.Wins
			initialBetAmount := state.BetAmount

			state.BettingOn = tt.betType
			state.ProcessWin(tt.strategy)

			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}

			expectedPayout := CalculatePayout(tt.betType, initialBetAmount)
			expectedBankroll := startingBankroll + initialBetAmount + expectedPayout
			if state.CurrentBankroll != expectedBankroll {
				t.Errorf("bankroll should update: got %.2f, want %.2f", state.CurrentBankroll, expectedBankroll)
			}
			if state.MaxBankrollReached != state.CurrentBankroll {
				t.Errorf("MaxBankrollReached should update: got %.2f, want %.2f", state.MaxBankrollReached, state.CurrentBankroll)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}
			if state.MaxWinsStreak != 1 {
				t.Errorf("MaxWinsStreak  should increment: got %d, want 1", state.MaxWinsStreak)
			}
			if state.BetAmount != initialBetAmount {
				t.Errorf("BetAmount should not change: got %.2f, want %.2f", state.BetAmount, initialBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessWin_MartingaleStrategies(t *testing.T) {
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
			initialWins := state.Wins

			// Precondition of loss streak
			state.LossStreak = 3
			state.BetAmount = state.BaseBetAmount * 8 // After 3 losses, bet would be 8x base

			state.BettingOn = tt.betType
			state.ProcessWin(tt.strategy)

			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}
			if state.BetAmount != state.BaseBetAmount {
				t.Errorf("BetAmount should reset: got %.2f, want %.2f", state.BetAmount, state.BaseBetAmount)
			}
		})
	}
}

func TestSimulatorStateProcessWin_FibonacciStrategies(t *testing.T) {
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
			initialWins := state.Wins

			// Precondition of Fibonacci sequence position with enough negative profit to avoid reset
			state.FibonacciSequenceIndex = 5                             // This corresponds to fib(5) = 8
			state.BetAmount = float64(GetFibonacciValue(5)) * MinimumBet // Calculate bet amount based on Fibonacci sequence
			state.FibonacciProfit = -10.0                                // Significant losses to avoid reset

			state.BettingOn = tt.betType
			state.ProcessWin(tt.strategy)

			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}
			if state.FibonacciSequenceIndex != 3 {
				t.Errorf("Fibonacci sequence should move back by 2 after win: got %d, want 3", state.FibonacciSequenceIndex)
			}

			originalBetAmount := float64(GetFibonacciValue(5)) * MinimumBet // Calculate bet amount based on Fibonacci sequence
			expectedPayout := CalculatePayout(tt.betType, originalBetAmount)
			expectedProfitIncrease := expectedPayout / MinimumBet
			expectedProfit := -10.0 + expectedProfitIncrease
			if state.FibonacciProfit != expectedProfit {
				t.Errorf("Fibonacci profit should update: got %.2f, want %.2f", state.FibonacciProfit, expectedProfit)
			}

			expectedBetAmount := float64(GetFibonacciValue(3)) * MinimumBet // Calculate bet amount based on Fibonacci sequence
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should update based on new sequence position: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			state.FibonacciSequenceIndex = 1
			state.BetAmount = float64(GetFibonacciValue(1)) * MinimumBet // Calculate bet amount based on Fibonacci sequence
			state.ProcessWin(tt.strategy)
			if state.FibonacciSequenceIndex != 0 {
				t.Errorf("Fibonacci sequence index should not go below 0: got %d, want 0", state.FibonacciSequenceIndex)
			}

			t.Run("Fibonacci profit resets when reaching +1 wager unit", func(t *testing.T) {
				state := NewSimulatorState()
				state.FibonacciProfit = 0.9
				state.FibonacciSequenceIndex = 2
				state.BetAmount = float64(GetFibonacciValue(2)) * MinimumBet // Calculate bet amount based on Fibonacci sequence
				state.ProcessWin(tt.strategy)

				if state.FibonacciProfit >= 1.0 {
					if state.FibonacciSequenceIndex != 0 {
						t.Errorf("Fibonacci sequence should reset when profit >= 1.0: got %d, want 0", state.FibonacciSequenceIndex)
					}
					if state.FibonacciProfit != 0.0 {
						t.Errorf("Fibonacci profit should reset: got %.2f, want 0.0", state.FibonacciProfit)
					}
					if state.BetAmount != MinimumBet {
						t.Errorf("BetAmount should reset to minimum when profit >= 1.0: got %.2f, want %.2f", state.BetAmount, MinimumBet)
					}
				}
			})
		})
	}
}

func TestSimulatorStateProcessWin_ParoliStrategies(t *testing.T) {
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
			initialWins := state.Wins

			state.BettingOn = tt.betType

			// First Paroli win - start progression
			state.ProcessWin(tt.strategy)
			if !state.IsInParoliProgression {
				t.Error("Paroli progression should start")
			}
			if state.ParoliProgressionLevel != 1 {
				t.Errorf("Paroli progression level should be 1: got %d", state.ParoliProgressionLevel)
			}
			if state.BetAmount != state.BaseBetAmount {
				t.Errorf("BetAmount should be set to base bet: got %.2f, want %.2f", state.BetAmount, state.BaseBetAmount)
			}
			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}

			// Second Paroli win - double the bet
			state.ProcessWin(tt.strategy)
			if !state.IsInParoliProgression {
				t.Error("Paroli progression should be set")
			}
			if state.ParoliProgressionLevel != 2 {
				t.Errorf("Paroli progression level should be 2: got %d", state.ParoliProgressionLevel)
			}
			if state.BetAmount != state.BaseBetAmount*2 {
				t.Errorf("BetAmount should doubled: got %.2f, want %.2f", state.BetAmount, state.BaseBetAmount*2)
			}

			// Third Paroli win - double the bet again
			state.ProcessWin(tt.strategy)
			if !state.IsInParoliProgression {
				t.Error("Paroli progression should be set")
			}
			if state.ParoliProgressionLevel != 3 {
				t.Errorf("Paroli progression level should be 3: got %d", state.ParoliProgressionLevel)
			}
			if state.BetAmount != state.BaseBetAmount*4 {
				t.Errorf("BetAmount should doubled: got %.2f, want %.2f", state.BetAmount, state.BaseBetAmount*4)
			}

			// Fourth Paroli win - goal reached: reset to base bet and end progression
			state.ProcessWin(tt.strategy)
			if state.IsInParoliProgression {
				t.Error("Paroli progression should end after goal reached")
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

func TestSimulatorStateProcessWin_DAlembertStrategies(t *testing.T) {
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
			initialWins := state.Wins

			// Precondition of D'Alembert level
			state.DAlembertLevel = 3
			state.BetAmount = state.DAlembertUnitSize + (float64(state.DAlembertLevel) * state.DAlembertUnitSize)

			state.BettingOn = tt.betType
			state.ProcessWin(tt.strategy)

			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}
			if state.DAlembertLevel != 2 {
				t.Errorf("D'Alembert level should decrease: got %d, want 2", state.DAlembertLevel)
			}

			expectedBetAmount := state.DAlembertUnitSize + (float64(state.DAlembertLevel) * state.DAlembertUnitSize)
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should decrease by one unit: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			t.Run("D'Alembert level and bet amount should not go below minimum", func(t *testing.T) {
				state := NewSimulatorState()

				// Precondition of D'Alembert level = 0
				state.DAlembertLevel = 0
				state.BetAmount = state.DAlembertUnitSize
				state.ProcessWin(tt.strategy)

				if state.DAlembertLevel != 0 {
					t.Errorf("D'Alembert level should not go below 0: got %d", state.DAlembertLevel)
				}
				if state.BetAmount != MinimumBet {
					t.Errorf("BetAmount should stay at minimum bet: got %.2f, want %.2f", state.BetAmount, MinimumBet)
				}
			})
		})
	}
}

func TestSimulatorStateProcessWin_OneThreeTwoSixStrategies(t *testing.T) {
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
			state := NewSimulatorState()
			initialWins := state.Wins

			state.BettingOn = tt.betType

			// First win - progress position 1 (3 units)
			state.ProcessWin(tt.strategy)
			if state.Wins != initialWins+1 {
				t.Errorf("wins should increment: got %d, want %d", state.Wins, initialWins+1)
			}
			if state.LossStreak != 0 {
				t.Errorf("LossStreak should reset: got %d, want 0", state.LossStreak)
			}
			if state.WinsStreak != 1 {
				t.Errorf("WinsStreak should increment: got %d, want 1", state.WinsStreak)
			}
			if state.OneThreeTwoSixSequenceIndex != 1 {
				t.Errorf("Sequence index should be 1: got %d", state.OneThreeTwoSixSequenceIndex)
			}
			expectedBetAmount := float64(GetOneThreeTwoSixValue(1)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should increase: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			// Second win - progress position 2 (2 units)
			state.ProcessWin(tt.strategy)
			if state.OneThreeTwoSixSequenceIndex != 2 {
				t.Errorf("Sequence index should be 2: got %d", state.OneThreeTwoSixSequenceIndex)
			}
			expectedBetAmount = float64(GetOneThreeTwoSixValue(2)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should increase: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			// Third win - progress position 3 (6 units)
			state.ProcessWin(tt.strategy)
			if state.OneThreeTwoSixSequenceIndex != 3 {
				t.Errorf("Sequence index should be 3: got %d", state.OneThreeTwoSixSequenceIndex)
			}
			expectedBetAmount = float64(GetOneThreeTwoSixValue(3)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should increase: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			// Fourth win - progress resets to position 0 (1 unit)
			state.ProcessWin(tt.strategy)
			if state.OneThreeTwoSixSequenceIndex != 0 {
				t.Errorf("Sequence index should reset: got %d, want 0", state.OneThreeTwoSixSequenceIndex)
			}
			expectedBetAmount = float64(GetOneThreeTwoSixValue(0)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should reset: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}

			// Fifth win - progress starts from position 1 (3 units)
			state.ProcessWin(tt.strategy)
			if state.OneThreeTwoSixSequenceIndex != 1 {
				t.Errorf("Sequence index should be 1: got %d", state.OneThreeTwoSixSequenceIndex)
			}
			expectedBetAmount = float64(GetOneThreeTwoSixValue(1)) * MinimumBet
			if state.BetAmount != expectedBetAmount {
				t.Errorf("BetAmount should update: got %.2f, want %.2f", state.BetAmount, expectedBetAmount)
			}
		})
	}
}
