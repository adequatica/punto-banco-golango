package simulator

import (
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

var ParoliMaxLevel = 3

type SimulatorState struct {
	CurrentBankroll     float64
	MaxBankrollReached  float64
	LastWinningHand     puntobanco.BetType
	BettingOn           puntobanco.BetType
	RoundsPlayed        int
	Wins                int
	BetAmount           float64
	GameEndedProfitably bool
	// Martingale-specific fields
	BaseBetAmount float64
	LossStreak    int
	WinsStreak    int
	MaxLossStreak int
	MaxWinsStreak int
	// Fibonacci-specific fields
	FibonacciSequenceIndex int     // Current position in Fibonacci sequence (0=based)
	FibonacciProfit        float64 // Current profit in wager units (1 unit = MinimumBet)
	// Paroli-specific fields
	IsInParoliProgression  bool
	ParoliProgressionLevel int // Current level in the Paroli progression (1, 2, 3)
	ParoliGoal             int // Goal to reach before resetting
	// D'Alembert-specific fields
	DAlembertUnitSize float64 // Base unit size for D'Alembert progression
	DAlembertLevel    int     // Current level in D'Alembert progression
	// 1-3-2-6 specific fields
	OneThreeTwoSixSequenceIndex int // Current position in 1-3-2-6 sequence (0=1, 1=3, 2=2, 3=6)
}

func NewSimulatorState() *SimulatorState {
	return &SimulatorState{
		CurrentBankroll:     Bankroll,
		MaxBankrollReached:  Bankroll,
		LastWinningHand:     puntobanco.PuntoPlayer,
		BettingOn:           puntobanco.PuntoPlayer,
		RoundsPlayed:        0,
		Wins:                0,
		BetAmount:           MinimumBet,
		GameEndedProfitably: false,
		// Martingale-specific fields
		BaseBetAmount: MinimumBet,
		LossStreak:    0,
		WinsStreak:    0,
		MaxLossStreak: 0,
		MaxWinsStreak: 0,
		// Fibonacci-specific fields
		FibonacciSequenceIndex: 0,
		FibonacciProfit:        0.0,
		// Paroli-specific fields
		IsInParoliProgression:  false,
		ParoliProgressionLevel: 0,
		ParoliGoal:             ParoliMaxLevel,
		// D'Alembert-specific fields
		DAlembertUnitSize: MinimumBet,
		DAlembertLevel:    0,
		// 1-3-2-6 specific fields
		OneThreeTwoSixSequenceIndex: 0,
	}
}

// Fibonacci sequence 1, 1, 2, 3, 5, 8, 13, 21, 34, 55, 89, 144, 233, 377, 610, 987, ...
func GetFibonacciValue(index int) int {
	if index <= 0 {
		return 1
	}
	if index == 1 {
		return 1
	}
	// Prevent memory exhaustion
	if index > 100000 {
		return 1
	}

	// Calculation for larger indices
	fib := make([]int, index+1)
	fib[0] = 1
	fib[1] = 1

	for i := 2; i <= index; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}

	return fib[index]
}

func GetOneThreeTwoSixValue(index int) int {
	sequence := []int{1, 3, 2, 6}
	if index < 0 || index >= len(sequence) {
		return 1 // Default to 1 unit if index is out of bounds
	}
	return sequence[index]
}

func CalculatePayout(betType puntobanco.BetType, betAmount float64) float64 {
	switch betType {

	case puntobanco.PuntoPlayer:
		// Winning bets on Punto hand pay even money (1:1)
		return betAmount

	case puntobanco.BancoBanker:
		// Winning bets on Banco hand pay 19 to 20 (5% commission)
		return betAmount * 0.95

	case puntobanco.EgaliteTie:
		// Standard payout for tie bet is 8-to-1
		return betAmount * 8.0

	default:
		return 0.0
	}
}

func (s *SimulatorState) CanPlaceBet() bool {
	return s.CurrentBankroll >= s.BetAmount
}

func (s *SimulatorState) PlaceBet() {
	s.CurrentBankroll -= s.BetAmount
}

func (s *SimulatorState) ProcessWin(strategy StrategyType) {
	s.Wins++

	payoutAmount := CalculatePayout(s.BettingOn, s.BetAmount)
	s.CurrentBankroll += s.BetAmount + payoutAmount
	// Track maximum bankroll reached
	if s.CurrentBankroll > s.MaxBankrollReached {
		s.MaxBankrollReached = s.CurrentBankroll
	}

	// Martingale strategy: reset loss streak and return to base bet.
	s.LossStreak = 0
	s.WinsStreak++
	// Track maximum wins streak
	if s.WinsStreak > s.MaxWinsStreak {
		s.MaxWinsStreak = s.WinsStreak
	}
	// If your bet wins, your stake remains the same for the next round, and you return to your original bet amount.
	// This means you continue betting the same amount.
	// Once you win a hand, you return to your original or opening bet unit and start the process over again
	if strategy == MartingaleOnPunto || strategy == MartingaleOnBanco {
		s.BetAmount = s.BaseBetAmount
	}

	// Fibonacci strategy: move back two places in sequence after a win.
	if strategy == FibonacciOnPunto || strategy == FibonacciOnBanco {
		// Calculate profit in wager units (1 unit = MinimumBet)
		payoutUnits := payoutAmount / MinimumBet
		s.FibonacciProfit += payoutUnits

		// Move back two places in the Fibonacci sequence
		if s.FibonacciSequenceIndex >= 2 {
			s.FibonacciSequenceIndex -= 2
		} else {
			s.FibonacciSequenceIndex = 0
		}

		// Calculate new bet amount based on Fibonacci sequence
		fibValue := GetFibonacciValue(s.FibonacciSequenceIndex)
		s.BetAmount = float64(fibValue) * MinimumBet

		// Reset if profit reaches +1 wager unit
		if s.FibonacciProfit >= 1.0 {
			s.FibonacciSequenceIndex = 0
			s.FibonacciProfit = 0.0
			s.BetAmount = MinimumBet
		}
	}

	// Paroli strategy: handle progression after a win.
	if strategy == ParoliOnPunto || strategy == ParoliOnBanco {
		if !s.IsInParoliProgression {
			// Start Paroli progression with base bet
			s.IsInParoliProgression = true
			s.ParoliProgressionLevel = 1
			s.BetAmount = s.BaseBetAmount
		} else {
			// Continue progression
			s.ParoliProgressionLevel++
			if s.ParoliProgressionLevel <= s.ParoliGoal {
				// Double the bet for next round
				s.BetAmount = s.BetAmount * 2
			} else {
				// Goal reached: reset to base bet and end progression
				s.BetAmount = s.BaseBetAmount
				s.IsInParoliProgression = false
				s.ParoliProgressionLevel = 0
			}
		}
	}

	// D'Alembert strategy: decrease bet by 1 unit after a win.
	if strategy == DAlembertOnPunto || strategy == DAlembertOnBanco {
		if s.DAlembertLevel > 0 {
			s.DAlembertLevel--
		}
		// Calculate new bet amount: base + (level * unit size)
		newBetAmount := s.DAlembertUnitSize + (float64(s.DAlembertLevel) * s.DAlembertUnitSize)
		// Ensure bet doesn't go below minimum bet
		if newBetAmount < MinimumBet {
			newBetAmount = MinimumBet
			s.DAlembertLevel = 0
		}
		s.BetAmount = newBetAmount
	}

	// 1-3-2-6 strategy: progress through sequence after a win.
	if strategy == OneThreeTwoSixOnPunto || strategy == OneThreeTwoSixOnBanco {
		// Move to next position in the sequence
		s.OneThreeTwoSixSequenceIndex++

		// Reset to start if the sequence is completed (reached index 3)
		if s.OneThreeTwoSixSequenceIndex >= 4 {
			s.OneThreeTwoSixSequenceIndex = 0
		}

		// Calculate new bet amount based on current position in sequence
		sequenceValue := GetOneThreeTwoSixValue(s.OneThreeTwoSixSequenceIndex)
		s.BetAmount = float64(sequenceValue) * MinimumBet
	}
}

func (s *SimulatorState) ProcessLoss(strategy StrategyType) {
	// Martingale strategy: increment loss streakand double the bet for next round.
	s.LossStreak++
	s.WinsStreak = 0
	// Track maximum loss streak
	if s.LossStreak > s.MaxLossStreak {
		s.MaxLossStreak = s.LossStreak
	}
	// If your bet loses, you double your bet size for the next hand.
	// You continue doubling your bet after each loss until you win a hand.
	if strategy == MartingaleOnPunto || strategy == MartingaleOnBanco {
		s.BetAmount = s.BetAmount * 2
	}

	// Fibonacci strategy: move to next number in sequence after a loss.
	if strategy == FibonacciOnPunto || strategy == FibonacciOnBanco {
		// Calculate loss in wager units (1 unit = MinimumBet)
		lossUnits := s.BetAmount / MinimumBet
		s.FibonacciProfit -= lossUnits

		// Move to next number in the Fibonacci sequence
		s.FibonacciSequenceIndex++

		// Calculate new bet amount based on Fibonacci sequence
		fibValue := GetFibonacciValue(s.FibonacciSequenceIndex)
		s.BetAmount = float64(fibValue) * MinimumBet
	}

	// Paroli strategy: reset to base bet when a loss occurs.
	if strategy == ParoliOnPunto || strategy == ParoliOnBanco {
		s.BetAmount = s.BaseBetAmount
		s.IsInParoliProgression = false
		s.ParoliProgressionLevel = 0
	}

	// D'Alembert strategy: increase bet by 1 unit after a loss.
	if strategy == DAlembertOnPunto || strategy == DAlembertOnBanco {
		s.DAlembertLevel++
		// Calculate new bet amount: base + (level * unit size)
		newBetAmount := s.DAlembertUnitSize + (float64(s.DAlembertLevel) * s.DAlembertUnitSize)
		s.BetAmount = newBetAmount
	}

	// 1-3-2-6 strategy: reset to start of sequence after a loss.
	if strategy == OneThreeTwoSixOnPunto || strategy == OneThreeTwoSixOnBanco {
		// Reset to the beginning of the sequence (1 unit)
		s.OneThreeTwoSixSequenceIndex = 0
		s.BetAmount = MinimumBet
	}
}
