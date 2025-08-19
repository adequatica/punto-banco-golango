package simulator

import (
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestGetStrategyOptions(t *testing.T) {
	want := GetStrategyOptions()

	if len(want) != 15 {
		t.Errorf("Strategy options of length %d should be 15", len(want))
	}

	if len(want) > 0 && want[0] != "Bet on Punto (player)" {
		t.Errorf("First betting option of '%s' should be 'Bet on Punto (player)'", want[0])
	}
}

func TestGetRandomBetType(t *testing.T) {
	result := GetRandomBetType()

	if result != puntobanco.PuntoPlayer && result != puntobanco.BancoBanker {
		t.Errorf("GetRandomBetType() returned unexpected value: %s", result)
	}
}

func TestMakeStrategy(t *testing.T) {
	tests := []struct {
		name          string
		strategy      StrategyType
		state         *SimulatorState
		wantBetType   puntobanco.BetType
		wantBetAmount float64
	}{
		{
			name:          "Bet on Punto returns PuntoPlayer with minimum bet",
			strategy:      BetOnPunto,
			state:         NewSimulatorState(),
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: MinimumBet,
		},
		{
			name:          "Bet on Banco returns BancoBanker with minimum bet",
			strategy:      BetOnBanco,
			state:         NewSimulatorState(),
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: MinimumBet,
		},
		{
			name:          "Bet on Égalité returns EgaliteTie with minimum bet",
			strategy:      BetOnEgalite,
			state:         NewSimulatorState(),
			wantBetType:   puntobanco.EgaliteTie,
			wantBetAmount: MinimumBet,
		},
		{
			name:     "Bet on last hand returns last winning hand with minimum bet",
			strategy: BetOnLastHand,
			state: &SimulatorState{
				LastWinningHand: puntobanco.BancoBanker,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: MinimumBet,
		},
		{
			name:          "Bet on random returns Punto or Banco hand with minimum bet",
			strategy:      BetOnRandom,
			state:         NewSimulatorState(),
			wantBetType:   GetRandomBetType(),
			wantBetAmount: MinimumBet,
		},
		{
			name:     "Martingale on Punto returns PuntoPlayer with current bet amount",
			strategy: MartingaleOnPunto,
			state: &SimulatorState{
				BetAmount: 10.0,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: 10.0,
		},
		{
			name:     "Martingale on Banco returns BancoBanker with current bet amount",
			strategy: MartingaleOnBanco,
			state: &SimulatorState{
				BetAmount: 20.0,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: 20.0,
		},
		{
			name:     "Fibonacci on Punto returns PuntoPlayer with current bet amount",
			strategy: FibonacciOnPunto,
			state: &SimulatorState{
				BetAmount: 30.0,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: 30.0,
		},
		{
			name:     "Fibonacci on Banco returns BancoBanker with current bet amount",
			strategy: FibonacciOnBanco,
			state: &SimulatorState{
				BetAmount: 40.0,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: 40.0,
		},
		{
			name:     "Paroli on Punto returns PuntoPlayer with current bet amount",
			strategy: ParoliOnPunto,
			state: &SimulatorState{
				BetAmount: 50.0,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: 50.0,
		},
		{
			name:     "Paroli on Banco returns BancoBanker with current bet amount",
			strategy: ParoliOnBanco,
			state: &SimulatorState{
				BetAmount: 60.0,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: 60.0,
		},
		{
			name:     "D'Alembert on Punto returns PuntoPlayer with current bet amount",
			strategy: DAlembertOnPunto,
			state: &SimulatorState{
				BetAmount: 70.0,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: 70.0,
		},
		{
			name:     "D'Alembert on Banco returns BancoBanker with current bet amount",
			strategy: DAlembertOnBanco,
			state: &SimulatorState{
				BetAmount: 80.0,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: 80.0,
		},
		{
			name:     "1-3-2-6 on Punto returns PuntoPlayer with current bet amount",
			strategy: OneThreeTwoSixOnPunto,
			state: &SimulatorState{
				BetAmount: 90.0,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: 90.0,
		},
		{
			name:     "1-3-2-6 on Punto returns BancoBanker with current bet amount",
			strategy: OneThreeTwoSixOnBanco,
			state: &SimulatorState{
				BetAmount: 100.0,
			},
			wantBetType:   puntobanco.BancoBanker,
			wantBetAmount: 100.0,
		},
		{
			name:          "Unknown strategy defaults to PuntoPlayer with minimum bet",
			strategy:      StrategyType("Unknown"),
			state:         NewSimulatorState(),
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: MinimumBet,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			betType, betAmount := MakeStrategy(tt.strategy, tt.state)

			// Don't check for random bet type due to randomness
			if tt.strategy != BetOnRandom {
				if betType != tt.wantBetType {
					t.Errorf("MakeStrategy() betType = %v, want %v", betType, tt.wantBetType)
				}
			}

			if betAmount != tt.wantBetAmount {
				t.Errorf("MakeStrategy() betAmount = %.2f, want %.2f", betAmount, tt.wantBetAmount)
			}
		})
	}
}

func TestRunSimulator(t *testing.T) {
	result := RunSimulator(BetOnPunto)
	if result == nil {
		t.Fatal("simulator's result should not be nil")
	}
	if result.RoundsPlayed == 0 {
		t.Fatal("simulator should play at least one round")
	}
}

func TestNewMultipleSimulationsStats(t *testing.T) {
	tests := []struct {
		name                 string
		numSimulations       int
		wantTotalSimulations int
	}{
		{
			name:                 "new stats with 100 simulations",
			numSimulations:       100,
			wantTotalSimulations: 100,
		},
		{
			name:                 "new stats with 1 simulation",
			numSimulations:       1,
			wantTotalSimulations: 1,
		},
		{
			name:                 "new stats with zero simulations",
			numSimulations:       0,
			wantTotalSimulations: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewMultipleSimulationsStats(tt.numSimulations)

			if result.TotalSimulations != tt.wantTotalSimulations {
				t.Errorf("TotalSimulations = %d should be %d",
					result.TotalSimulations, tt.wantTotalSimulations)
			}

			// Default values
			if result.AvgRoundsPlayed != 0.0 {
				t.Errorf("AvgRoundsPlayed = %.2f should be 0.0",
					result.AvgRoundsPlayed)
			}
			if result.MinRoundsPlayed != MaxIntValue {
				t.Errorf("MinRoundsPlayed = %d should be %d",
					result.MinRoundsPlayed, MaxIntValue)
			}
			if result.MaxRoundsPlayed != 0 {
				t.Errorf("MaxRoundsPlayed = %d should be 0",
					result.MaxRoundsPlayed)
			}
			if result.MaxConsecutiveLosses != 0 {
				t.Errorf("MaxConsecutiveLosses = %d should be 0",
					result.MaxConsecutiveLosses)
			}
			if result.MaxConsecutiveWins != 0 {
				t.Errorf("MaxConsecutiveWins = %d should be 0",
					result.MaxConsecutiveWins)
			}
			if result.AvgWins != 0.0 {
				t.Errorf("AvgWins = %.2f should be 0.0",
					result.AvgWins)
			}
			if result.MinWins != MaxIntValue {
				t.Errorf("MinWins = %d should be %d",
					result.MinWins, MaxIntValue)
			}
			if result.MaxWins != 0 {
				t.Errorf("MaxWins = %d should be 0",
					result.MaxWins)
			}
			if result.WinRate != 0.0 {
				t.Errorf("WinRate = %.2f should be 0.0",
					result.WinRate)
			}
			if result.GamesWithZeroWins != 0 {
				t.Errorf("GamesWithZeroWins = %d should be 0",
					result.GamesWithZeroWins)
			}
			if result.ZeroWinsRate != 0 {
				t.Errorf("ZeroWinsRate = %.2f should be 0.0",
					result.ZeroWinsRate)
			}
			if result.AvgMaxConsecutiveWins != 0.0 {
				t.Errorf("AvgMaxConsecutiveWins = %.2f should be 0.0",
					result.AvgMaxConsecutiveWins)
			}
			if result.MaxConsecutiveWins != 0 {
				t.Errorf("MaxConsecutiveWins = %d should be 0",
					result.MaxConsecutiveWins)
			}
			if result.AvgMaxConsecutiveLosses != 0.0 {
				t.Errorf("AvgMaxConsecutiveLosses = %.2f should be 0.0",
					result.AvgMaxConsecutiveLosses)
			}
			if result.MaxConsecutiveLosses != 0 {
				t.Errorf("MaxConsecutiveLosses = %d should be 0",
					result.MaxConsecutiveLosses)
			}
			if result.AvgMaxBudgetReached != 0.0 {
				t.Errorf("AvgMaxBudgetReached = %.2f should be 0.0",
					result.AvgMaxBudgetReached)
			}
			if result.MaxBudgetReacorded != Budget {
				t.Errorf("MaxBudgetReacorded = %.2f should be  %.2f",
					result.MaxBudgetReacorded, Budget)
			}
			if result.GamesWithProfitableBudget != 0 {
				t.Errorf("GamesWithProfitableBudget = %d should be 0",
					result.GamesWithProfitableBudget)
			}
			if result.ProfitableBudgetRate != 0.0 {
				t.Errorf("ProfitableBudgetRate = %.2f should be 0.0",
					result.ProfitableBudgetRate)
			}
		})
	}
}

func TestRunMultipleSimulations(t *testing.T) {
	numberOfTestSimulations := 10
	result := RunMultipleSimulations(BetOnPunto, numberOfTestSimulations)
	if result.TotalSimulations != numberOfTestSimulations {
		t.Fatal("should run multiple simulations")
	}
	if result.AvgRoundsPlayed <= 0 {
		t.Fatal("multiple simulations should play average at least one round")
	}
	if result.MinRoundsPlayed <= 0 {
		t.Fatal("multiple simulations should play min at least one round")
	}
	if result.MaxRoundsPlayed <= 0 {
		t.Fatal("multiple simulations should play max at least one round")
	}
}
