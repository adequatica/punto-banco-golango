package simulator

import (
	"fmt"
	"math/rand"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

var (
	Budget     = 1000.0
	MinimumBet = 10.0

	MaxIntValue         = int(^uint(0) >> 1)
	ProfitableThreshold = Budget * 1.01 // 101% of default budget
)

type StrategyType string

const (
	// Flat betting strategies
	BetOnPunto    StrategyType = "Bet on Punto (player)"
	BetOnBanco    StrategyType = "Bet on Banco (banker)"
	BetOnEgalite  StrategyType = "Bet on Égalité (tie)"
	BetOnLastHand StrategyType = "Bet on last hand"
	BetOnRandom   StrategyType = "Bet on random"
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

	return state
}

type MultipleSimulationsStats struct {
	TotalSimulations int

	AvgRoundsPlayed float64
	MinRoundsPlayed int
	MaxRoundsPlayed int

	AvgWins           float64
	MinWins           int
	MaxWins           int
	WinRate           float64
	GamesWithZeroWins int
	ZeroWinsRate      float64

	AvgMaxConsecutiveWins   float64
	MaxConsecutiveWins      int
	AvgMaxConsecutiveLosses float64
	MaxConsecutiveLosses    int

	AvgMaxBudgetReached       float64
	MaxBudgetReacorded        float64
	GamesWithProfitableBudget int
	ProfitableBudgetRate      float64
}

func NewMultipleSimulationsStats(numSimulations int) MultipleSimulationsStats {
	if numSimulations <= 0 {
		numSimulations = 1
	}

	return MultipleSimulationsStats{
		TotalSimulations: numSimulations,

		AvgRoundsPlayed: 0.0,
		MinRoundsPlayed: MaxIntValue,
		MaxRoundsPlayed: 0,

		AvgWins:           0.0,
		MinWins:           MaxIntValue,
		MaxWins:           0,
		WinRate:           0.0,
		GamesWithZeroWins: 0,
		ZeroWinsRate:      0.0,

		AvgMaxConsecutiveWins:   0.0,
		MaxConsecutiveWins:      0,
		AvgMaxConsecutiveLosses: 0.0,
		MaxConsecutiveLosses:    0,

		AvgMaxBudgetReached:       0.0,
		MaxBudgetReacorded:        Budget,
		GamesWithProfitableBudget: 0,
		ProfitableBudgetRate:      0.0,
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
	totalMaxConsecutiveWins := 0
	totalMaxConsecutiveLosses := 0
	totalMaxBudgetReached := 0.0

	for i := 0; i < numSimulations; i++ {
		state := RunSimulator(strategy)

		// Track games played stats
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

		// Track games with zero wins
		if state.Wins == 0 {
			stats.GamesWithZeroWins++
		}

		// Track win rate
		if state.RoundsPlayed > 0 {
			winRate := float64(state.Wins) / float64(state.RoundsPlayed) * 100
			totalWinRate += winRate
		}

		// Track consecutive wins stats
		totalMaxConsecutiveWins += state.MaxConsecutiveWins
		if state.MaxConsecutiveWins > stats.MaxConsecutiveWins {
			stats.MaxConsecutiveWins = state.MaxConsecutiveWins
		}

		// Track consecutive losses stats
		totalMaxConsecutiveLosses += state.MaxConsecutiveLosses
		if state.MaxConsecutiveLosses > stats.MaxConsecutiveLosses {
			stats.MaxConsecutiveLosses = state.MaxConsecutiveLosses
		}

		// Track max budget reached stats
		totalMaxBudgetReached += state.MaxBudgetReached
		if state.MaxBudgetReached > stats.MaxBudgetReacorded {
			stats.MaxBudgetReacorded = state.MaxBudgetReached
		}

		// Track games with a profit budget during the game
		if state.MaxBudgetReached > ProfitableThreshold {
			stats.GamesWithProfitableBudget++
		}
	}

	// Calculate averages
	stats.AvgRoundsPlayed = float64(totalRoundsPlayed) / float64(numSimulations)
	stats.AvgWins = float64(totalWins) / float64(numSimulations)
	stats.WinRate = totalWinRate / float64(numSimulations)
	stats.ZeroWinsRate = float64(stats.GamesWithZeroWins) / float64(numSimulations) * 100
	stats.AvgMaxConsecutiveWins = float64(totalMaxConsecutiveWins) / float64(numSimulations)
	stats.AvgMaxConsecutiveLosses = float64(totalMaxConsecutiveLosses) / float64(numSimulations)
	stats.AvgMaxBudgetReached = totalMaxBudgetReached / float64(numSimulations)
	stats.ProfitableBudgetRate = float64(stats.GamesWithProfitableBudget) / float64(numSimulations) * 100

	return stats
}
