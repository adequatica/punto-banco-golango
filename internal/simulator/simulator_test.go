package simulator

import (
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestGetStrategyOptions(t *testing.T) {
	want := GetStrategyOptions()

	if len(want) != 16 {
		t.Errorf("Strategy options of length %d should be 16", len(want))
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

func TestGetOnlyPuntoBanco(t *testing.T) {
	tests := []struct {
		name            string
		lastWinningHand puntobanco.BetType
		want            puntobanco.BetType
	}{
		{
			name:            "if last winning hand is PuntoPlayer, return PuntoPlayer",
			lastWinningHand: puntobanco.PuntoPlayer,
			want:            puntobanco.PuntoPlayer,
		},
		{
			name:            "if last winning hand is BancoBanker, return BancoBanker",
			lastWinningHand: puntobanco.BancoBanker,
			want:            puntobanco.BancoBanker,
		},
		{
			name:            "if last winning hand is EgaliteTie, return BancoBanker",
			lastWinningHand: puntobanco.EgaliteTie,
			want:            puntobanco.BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetOnlyPuntoBanco(tt.lastWinningHand)
			if got != tt.want {
				t.Errorf("GetOnlyPuntoBanco(%v) = %v, want %v", tt.lastWinningHand, got, tt.want)
			}
		})
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
			name:     "Bet on last hand PB returns last winning hand with minimum bet",
			strategy: BetOnLastHand,
			state: &SimulatorState{
				LastWinningHand: puntobanco.PuntoPlayer,
			},
			wantBetType:   puntobanco.PuntoPlayer,
			wantBetAmount: MinimumBet,
		},
		{
			name:     "Bet on last hand PB returns BancoBanker with minimum bet in case of EgaliteTie",
			strategy: BetOnLastHandPB,
			state: &SimulatorState{
				LastWinningHand: puntobanco.EgaliteTie,
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
	result := RunSimulator(BetOnPunto, nil)
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
			if result.AvgRoundsPerGame != 0.0 {
				t.Errorf("AvgRoundsPerGame = %.2f should be 0.0",
					result.AvgRoundsPerGame)
			}
			if result.MinRoundsPlayed != MaxIntValue {
				t.Errorf("MinRoundsPlayed = %d should be %d",
					result.MinRoundsPlayed, MaxIntValue)
			}
			if result.MaxRoundsPlayed != 0 {
				t.Errorf("MaxRoundsPlayed = %d should be 0",
					result.MaxRoundsPlayed)
			}
			if result.MaxLossStreak != 0 {
				t.Errorf("MaxLossStreak = %d should be 0",
					result.MaxLossStreak)
			}
			if result.MaxWinsStreak != 0 {
				t.Errorf("MaxWinsStreak = %d should be 0",
					result.MaxWinsStreak)
			}
			if result.AvgWinsPerGames != 0.0 {
				t.Errorf("AvgWinsPerGames = %.2f should be 0.0",
					result.AvgWinsPerGames)
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
			if result.AvgMaxWinsStreak != 0.0 {
				t.Errorf("AvgMaxWinsStreak = %.2f should be 0.0",
					result.AvgMaxWinsStreak)
			}
			if result.MaxWinsStreak != 0 {
				t.Errorf("MaxWinsStreak = %d should be 0",
					result.MaxWinsStreak)
			}
			if result.AvgMaxLossStreak != 0.0 {
				t.Errorf("AvgMaxLossStreak = %.2f should be 0.0",
					result.AvgMaxLossStreak)
			}
			if result.MaxLossStreak != 0 {
				t.Errorf("MaxLossStreak = %d should be 0",
					result.MaxLossStreak)
			}
			if result.AvgMaxBankrollReached != 0.0 {
				t.Errorf("AvgMaxBankrollReached = %.2f should be 0.0",
					result.AvgMaxBankrollReached)
			}
			if result.MaxBankrollReacorded != Bankroll {
				t.Errorf("MaxBankrollReacorded = %.2f should be  %.2f",
					result.MaxBankrollReacorded, Bankroll)
			}
			if result.GamesWithProfitableBankroll != 0 {
				t.Errorf("GamesWithProfitableBankroll = %d should be 0",
					result.GamesWithProfitableBankroll)
			}
			if result.ProfitableBankrollRate != 0.0 {
				t.Errorf("ProfitableBankrollRate = %.2f should be 0.0",
					result.ProfitableBankrollRate)
			}
			if result.GamesWithProfitableEnd != 0 {
				t.Errorf("GamesWithProfitableEnd = %d should be 0",
					result.GamesWithProfitableEnd)
			}
			if result.ProfitableEndGamesRate != 0.0 {
				t.Errorf("ProfitableEndGamesRate = %.2f should be 0.0",
					result.ProfitableEndGamesRate)
			}
		})
	}
}

func TestRunMultipleSimulations(t *testing.T) {
	numberOfTestSimulations := 10
	result := RunMultipleSimulations(BetOnPunto, numberOfTestSimulations, false)
	if result.TotalSimulations != numberOfTestSimulations {
		t.Fatal("should run multiple simulations")
	}
	if result.AvgRoundsPerGame <= 0 {
		t.Fatal("multiple simulations should play average at least one round")
	}
	if result.MinRoundsPlayed <= 0 {
		t.Fatal("multiple simulations should play min at least one round")
	}
	if result.MaxRoundsPlayed <= 0 {
		t.Fatal("multiple simulations should play max at least one round")
	}
}
