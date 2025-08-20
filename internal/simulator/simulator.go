package simulator

import (
	"fmt"
	"math/rand"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

var (
	Bankroll   = 1000.0
	MinimumBet = 10.0

	MaxIntValue         = int(^uint(0) >> 1)
	ProfitableThreshold = Bankroll * 1.01 // 101% of default bankroll
)

type StrategyType string

const (
	// Flat betting strategies
	BetOnPunto      StrategyType = "Bet on Punto (player)"
	BetOnBanco      StrategyType = "Bet on Banco (banker)"
	BetOnEgalite    StrategyType = "Bet on Égalité (tie)"
	BetOnLastHand   StrategyType = "Bet on last hand"
	BetOnLastHandPB StrategyType = "Bet on last hand PB"
	BetOnRandom     StrategyType = "Bet on random"
	// Progressive betting strategies
	MartingaleOnPunto     StrategyType = "Martingale on Punto"
	MartingaleOnBanco     StrategyType = "Martingale on Banco"
	FibonacciOnPunto      StrategyType = "Fibonacci on Punto"
	FibonacciOnBanco      StrategyType = "Fibonacci on Banco"
	ParoliOnPunto         StrategyType = "Paroli on Punto"
	ParoliOnBanco         StrategyType = "Paroli on Banco"
	DAlembertOnPunto      StrategyType = "D'Alembert on Punto"
	DAlembertOnBanco      StrategyType = "D'Alembert on Banco"
	OneThreeTwoSixOnPunto StrategyType = "1-3-2-6 on Punto"
	OneThreeTwoSixOnBanco StrategyType = "1-3-2-6 on Banco"
)

func GetStrategyOptions() []string {
	return []string{
		// Flat betting strategies
		string(BetOnPunto),
		string(BetOnBanco),
		string(BetOnEgalite),
		string(BetOnLastHand),
		string(BetOnLastHandPB),
		string(BetOnRandom),
		// Progressive betting strategies
		string(MartingaleOnPunto),
		string(MartingaleOnBanco),
		string(FibonacciOnPunto),
		string(FibonacciOnBanco),
		string(ParoliOnPunto),
		string(ParoliOnBanco),
		string(DAlembertOnPunto),
		string(DAlembertOnBanco),
		string(OneThreeTwoSixOnPunto),
		string(OneThreeTwoSixOnBanco),
	}
}

func GetRandomBetType() puntobanco.BetType {
	if rand.Intn(2) == 0 {
		return puntobanco.PuntoPlayer
	} else {
		return puntobanco.BancoBanker
	}
}

func GetOnlyPuntoBanco(lastWinningHand puntobanco.BetType) puntobanco.BetType {
	if lastWinningHand == puntobanco.EgaliteTie {
		// Bet on Banco to maximize chance of winning
		return puntobanco.BancoBanker
	} else {
		return lastWinningHand
	}
}

func MakeStrategy(strategy StrategyType, state *SimulatorState) (puntobanco.BetType, float64) {
	switch strategy {
	case BetOnPunto:
		return puntobanco.PuntoPlayer, MinimumBet
	case BetOnBanco:
		return puntobanco.BancoBanker, MinimumBet
	case BetOnEgalite:
		return puntobanco.EgaliteTie, MinimumBet
	case BetOnLastHand:
		return state.LastWinningHand, MinimumBet
	case BetOnLastHandPB:
		return GetOnlyPuntoBanco(state.LastWinningHand), MinimumBet
	case BetOnRandom:
		return GetRandomBetType(), MinimumBet

	case MartingaleOnPunto:
		return puntobanco.PuntoPlayer, state.BetAmount
	case MartingaleOnBanco:
		return puntobanco.BancoBanker, state.BetAmount

	case FibonacciOnPunto:
		return puntobanco.PuntoPlayer, state.BetAmount
	case FibonacciOnBanco:
		return puntobanco.BancoBanker, state.BetAmount

	case ParoliOnPunto:
		return puntobanco.PuntoPlayer, state.BetAmount
	case ParoliOnBanco:
		return puntobanco.BancoBanker, state.BetAmount

	case DAlembertOnPunto:
		return puntobanco.PuntoPlayer, state.BetAmount
	case DAlembertOnBanco:
		return puntobanco.BancoBanker, state.BetAmount

	case OneThreeTwoSixOnPunto:
		return puntobanco.PuntoPlayer, state.BetAmount
	case OneThreeTwoSixOnBanco:
		return puntobanco.BancoBanker, state.BetAmount

	default:
		return puntobanco.PuntoPlayer, MinimumBet
	}
}

func RunSimulator(strategy StrategyType) *SimulatorState {
	state := NewSimulatorState()
	shoe := deck.MakeNewShoe()

	// Run simulation until player cannot bet anymore
	for state.CanPlaceBet() {
		betType, betAmount := MakeStrategy(strategy, state)
		state.BettingOn = betType
		state.BetAmount = betAmount

		state.PlaceBet()

		// Play the game
		gameResult, err := puntobanco.PlayPuntoBanco(shoe)
		if err != nil {
			fmt.Printf("Error playing game: %v\n", err)
			break
		}

		// Update shoe for next game
		shoe = gameResult.RemainingShoe

		// Track the last winning hand
		if gameResult.Result != nil {
			state.LastWinningHand = *gameResult.Result
		}

		if gameResult.Result != nil && *gameResult.Result == state.BettingOn {
			state.ProcessWin(strategy)
		} else {
			state.ProcessLoss(strategy)
		}

		state.RoundsPlayed++
	}

	// Track if the game ended profitably (when player can't bet any longer)
	if state.CurrentBankroll > Bankroll {
		state.GameEndedProfitably = true
	}

	return state
}

type MultipleSimulationsStats struct {
	TotalSimulations int

	AvgRoundsPerGame float64
	MinRoundsPlayed  int
	MaxRoundsPlayed  int

	AvgWinsPerGames   float64
	MinWins           int
	MaxWins           int
	WinRate           float64
	GamesWithZeroWins int
	ZeroWinsRate      float64

	AvgMaxWinsStreak float64
	MaxWinsStreak    int
	AvgMaxLossStreak float64
	MaxLossStreak    int

	AvgMaxBankrollReached       float64
	MaxBankrollReacorded        float64
	GamesWithProfitableBankroll int
	ProfitableBankrollRate      float64
	GamesWithProfitableEnd      int
	ProfitableEndGamesRate      float64
}

func NewMultipleSimulationsStats(numSimulations int) MultipleSimulationsStats {
	if numSimulations <= 0 {
		numSimulations = 1
	}

	return MultipleSimulationsStats{
		TotalSimulations: numSimulations,

		AvgRoundsPerGame: 0.0,
		MinRoundsPlayed:  MaxIntValue,
		MaxRoundsPlayed:  0,

		AvgWinsPerGames:   0.0,
		MinWins:           MaxIntValue,
		MaxWins:           0,
		WinRate:           0.0,
		GamesWithZeroWins: 0,
		ZeroWinsRate:      0.0,

		AvgMaxWinsStreak: 0.0,
		MaxWinsStreak:    0,
		AvgMaxLossStreak: 0.0,
		MaxLossStreak:    0,

		AvgMaxBankrollReached:       0.0,
		MaxBankrollReacorded:        Bankroll,
		GamesWithProfitableBankroll: 0,
		ProfitableBankrollRate:      0.0,
		GamesWithProfitableEnd:      0,
		ProfitableEndGamesRate:      0.0,
	}
}

func RunMultipleSimulations(strategy StrategyType, numSimulations int) MultipleSimulationsStats {
	if numSimulations <= 0 {
		numSimulations = 1
	}

	stats := NewMultipleSimulationsStats(numSimulations)

	totalRoundsPlayed := 0
	totalWins := 0
	totalWinRate := 0.0
	totalMaxWinsStreak := 0
	totalMaxLossStreak := 0
	totalMaxBankrollReached := 0.0

	for i := 0; i < numSimulations; i++ {
		state := RunSimulator(strategy)

		// Track played games stats
		totalRoundsPlayed += state.RoundsPlayed
		if state.RoundsPlayed < stats.MinRoundsPlayed {
			stats.MinRoundsPlayed = state.RoundsPlayed
		}
		if state.RoundsPlayed > stats.MaxRoundsPlayed {
			stats.MaxRoundsPlayed = state.RoundsPlayed
		}

		// Track wins stats
		totalWins += state.Wins
		if state.Wins < stats.MinWins {
			stats.MinWins = state.Wins
		}
		if state.Wins > stats.MaxWins {
			stats.MaxWins = state.Wins
		}

		// Track win rate
		if state.RoundsPlayed > 0 {
			winRate := float64(state.Wins) / float64(state.RoundsPlayed) * 100
			totalWinRate += winRate
		}

		// Track zero-wins games
		if state.Wins == 0 {
			stats.GamesWithZeroWins++
		}

		// Track wins streak stats
		totalMaxWinsStreak += state.MaxWinsStreak
		if state.MaxWinsStreak > stats.MaxWinsStreak {
			stats.MaxWinsStreak = state.MaxWinsStreak
		}

		// Track loss streak stats
		totalMaxLossStreak += state.MaxLossStreak
		if state.MaxLossStreak > stats.MaxLossStreak {
			stats.MaxLossStreak = state.MaxLossStreak
		}

		// Track max bankroll reached stats
		totalMaxBankrollReached += state.MaxBankrollReached
		if state.MaxBankrollReached > stats.MaxBankrollReacorded {
			stats.MaxBankrollReacorded = state.MaxBankrollReached
		}

		// Track games with a profitable bankroll during the game
		if state.MaxBankrollReached > ProfitableThreshold {
			stats.GamesWithProfitableBankroll++
		}

		// Track games with profitable end (when player couldn't bet anymore)
		if state.GameEndedProfitably {
			stats.GamesWithProfitableEnd++
		}
	}

	// Calculate averages
	stats.AvgRoundsPerGame = float64(totalRoundsPlayed) / float64(numSimulations)
	stats.AvgWinsPerGames = float64(totalWins) / float64(numSimulations)
	stats.WinRate = totalWinRate / float64(numSimulations)
	stats.ZeroWinsRate = float64(stats.GamesWithZeroWins) / float64(numSimulations) * 100
	stats.AvgMaxWinsStreak = float64(totalMaxWinsStreak) / float64(numSimulations)
	stats.AvgMaxLossStreak = float64(totalMaxLossStreak) / float64(numSimulations)
	stats.AvgMaxBankrollReached = totalMaxBankrollReached / float64(numSimulations)
	stats.ProfitableBankrollRate = float64(stats.GamesWithProfitableBankroll) / float64(numSimulations) * 100
	stats.ProfitableEndGamesRate = float64(stats.GamesWithProfitableEnd) / float64(numSimulations) * 100

	return stats
}
