package simulator

import (
	"regexp"
	"strings"
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestFormatCard(t *testing.T) {
	tests := []struct {
		name     string
		card     *deck.Card
		expected string
	}{
		{
			name:     "nil card returns empty string",
			card:     nil,
			expected: "",
		},
		{
			name: "Spades card",
			card: &deck.Card{
				Card:  "A",
				Value: 1,
				Suit:  "Spades",
			},
			expected: "AS",
		},
		{
			name: "Clubs card",
			card: &deck.Card{
				Card:  "K",
				Value: 0,
				Suit:  "Clubs",
			},
			expected: "KC",
		},
		{
			name: "Hearts card",
			card: &deck.Card{
				Card:  "Q",
				Value: 0,
				Suit:  "Hearts",
			},
			expected: "QH",
		},
		{
			name: "Diamonds card",
			card: &deck.Card{
				Card:  "10",
				Value: 0,
				Suit:  "Diamonds",
			},
			expected: "10D",
		},
		{
			name: "Unknown suit returns question mark",
			card: &deck.Card{
				Card:  "2",
				Value: 2,
				Suit:  "Unknown",
			},
			expected: "2?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatCard(tt.card)
			if result != tt.expected {
				t.Errorf("FormatCard() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFormatBetAndResultType(t *testing.T) {
	tests := []struct {
		name     string
		betType  puntobanco.BetType
		expected string
	}{
		{
			name:     "PuntoPlayer returns punto",
			betType:  puntobanco.PuntoPlayer,
			expected: "punto",
		},
		{
			name:     "BancoBanker returns banko",
			betType:  puntobanco.BancoBanker,
			expected: "banko",
		},
		{
			name:     "EgaliteTie returns egalite",
			betType:  puntobanco.EgaliteTie,
			expected: "egalite",
		},
		{
			name:     "Unknown bet type defaults to punto",
			betType:  puntobanco.BetType("Unknown"),
			expected: "punto",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBetAndResultType(tt.betType)
			if result != tt.expected {
				t.Errorf("FormatBetAndResultType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewDataCollector(t *testing.T) {
	tests := []struct {
		name                string
		strategy            StrategyType
		decksInShoe         int
		startingBankroll    float64
		standardBet         float64
		numberOfSimulations int
		expectedStrategy    string
		expectedDecksInShoe int
		expectedBankroll    float64
		expectedBet         float64
		expectedSimulations int
	}{
		{
			name:                "Create data collector with valid parameters",
			strategy:            BetOnPunto,
			decksInShoe:         6,
			startingBankroll:    1000.0,
			standardBet:         10.0,
			numberOfSimulations: 100,
			expectedStrategy:    string(BetOnPunto),
			expectedDecksInShoe: 6,
			expectedBankroll:    1000.0,
			expectedBet:         10.0,
			expectedSimulations: 100,
		},
		{
			name:                "Create data collector with different parameters",
			strategy:            MartingaleOnPunto,
			decksInShoe:         8,
			startingBankroll:    2000.0,
			standardBet:         20.0,
			numberOfSimulations: 10000,
			expectedStrategy:    string(MartingaleOnPunto),
			expectedDecksInShoe: 8,
			expectedBankroll:    2000.0,
			expectedBet:         20.0,
			expectedSimulations: 10000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := NewDataCollector(
				tt.strategy,
				tt.decksInShoe,
				tt.startingBankroll,
				tt.standardBet,
				tt.numberOfSimulations,
			)

			if dc == nil {
				t.Fatal("NewDataCollector() returned nil")
			}

			if dc.data == nil {
				t.Fatal("DataCollector.data should not be nil")
			}

			if dc.data.Strategy != tt.expectedStrategy {
				t.Errorf("Strategy = %v, want %v", dc.data.Strategy, tt.expectedStrategy)
			}

			if dc.data.DecksInShoe != tt.expectedDecksInShoe {
				t.Errorf("DecksInShoe = %v, want %v", dc.data.DecksInShoe, tt.expectedDecksInShoe)
			}

			if dc.data.StartingBankroll != tt.expectedBankroll {
				t.Errorf("StartingBankroll = %v, want %v", dc.data.StartingBankroll, tt.expectedBankroll)
			}

			if dc.data.StandardBet != tt.expectedBet {
				t.Errorf("StandardBet = %v, want %v", dc.data.StandardBet, tt.expectedBet)
			}

			if dc.data.NumberOfSimulations != tt.expectedSimulations {
				t.Errorf("NumberOfSimulations = %v, want %v", dc.data.NumberOfSimulations, tt.expectedSimulations)
			}

			if dc.data.Games == nil {
				t.Fatal("Games slice should not be nil")
			}

			if len(dc.data.Games) != 0 {
				t.Errorf("Games should be empty initially, got length %d", len(dc.data.Games))
			}

			if dc.currentGameID != 0 {
				t.Errorf("currentGameID should be 0 initially, got %d", dc.currentGameID)
			}

			if dc.currentHandID != 0 {
				t.Errorf("currentHandID should be 0 initially, got %d", dc.currentHandID)
			}

			if dc.currentShoeNum != 0 {
				t.Errorf("currentShoeNum should be 0 initially, got %d", dc.currentShoeNum)
			}
		})
	}
}

func TestDataCollector_StartNewGame(t *testing.T) {
	dc := NewDataCollector(BetOnPunto, 6, 1000.0, 10.0, 100)
	shoe := deck.MakeNewShoe()

	// Start first game
	dc.StartNewGame(shoe)

	if dc.currentGameID != 1 {
		t.Errorf("currentGameID = %d, want 1", dc.currentGameID)
	}

	if dc.currentHandID != 0 {
		t.Errorf("currentHandID should be reset to 0, got %d", dc.currentHandID)
	}

	if dc.currentShoeNum != 1 {
		t.Errorf("currentShoeNum = %d, want 1", dc.currentShoeNum)
	}

	if dc.previousShoeLen != len(shoe) {
		t.Errorf("previousShoeLen = %d, want %d", dc.previousShoeLen, len(shoe))
	}

	if len(dc.data.Games) != 1 {
		t.Errorf("Games length = %d, want 1", len(dc.data.Games))
	}

	if len(dc.data.Games[0]) != 0 {
		t.Errorf("First game should have no hands initially, got %d", len(dc.data.Games[0]))
	}

	// Start second game
	dc.StartNewGame(shoe)

	if dc.currentGameID != 2 {
		t.Errorf("currentGameID = %d, want 2", dc.currentGameID)
	}

	if dc.currentShoeNum != 1 {
		t.Errorf("currentShoeNum = %d, want 1 (should reset to 1 for new game)", dc.currentShoeNum)
	}

	if len(dc.data.Games) != 2 {
		t.Errorf("Games length = %d, want 2", len(dc.data.Games))
	}
}

func TestDataCollector_CollectHandData(t *testing.T) {
	dc := NewDataCollector(BetOnPunto, 6, 1000.0, 10.0, 100)
	shoe := deck.MakeNewShoe()
	dc.StartNewGame(shoe)

	// Create test state
	state := &SimulatorState{
		CurrentBankroll: 1000.0,
		BettingOn:       puntobanco.PuntoPlayer,
		BetAmount:       10.0,
	}

	// Create test game result
	puntoCard1 := &deck.Card{Card: "A", Value: 1, Suit: "Spades"}
	puntoCard2 := &deck.Card{Card: "2", Value: 2, Suit: "Clubs"}
	bancoCard1 := &deck.Card{Card: "K", Value: 0, Suit: "Hearts"}
	bancoCard2 := &deck.Card{Card: "Q", Value: 0, Suit: "Diamonds"}

	result := puntobanco.PuntoPlayer
	gameResult := &puntobanco.GameResultState{
		Result: &result,
		PuntoState: &puntobanco.PlayerState{
			FirstCard:  puntoCard1,
			SecondCard: puntoCard2,
			ThirdCard:  nil,
			Points:     3,
		},
		BancoState: &puntobanco.PlayerState{
			FirstCard:  bancoCard1,
			SecondCard: bancoCard2,
			ThirdCard:  nil,
			Points:     0,
		},
		RemainingShoe: shoe[4:],
	}

	// Collect hand data
	dc.CollectHandData(state, gameResult, len(shoe))

	// Verify data was collected
	if len(dc.data.Games) != 1 {
		t.Fatalf("Games length = %d, want 1", len(dc.data.Games))
	}

	if len(dc.data.Games[0]) != 1 {
		t.Fatalf("First game should have 1 hand, got %d", len(dc.data.Games[0]))
	}

	hand := dc.data.Games[0][0]

	if hand.GameID != 1 {
		t.Errorf("GameID = %d, want 1", hand.GameID)
	}

	if hand.HandID != 1 {
		t.Errorf("HandID = %d, want 1", hand.HandID)
	}

	if hand.ShoeNumber != 1 {
		t.Errorf("ShoeNumber = %d, want 1", hand.ShoeNumber)
	}

	if len(hand.PuntoHand) != 2 {
		t.Errorf("PuntoHand length = %d, want 2", len(hand.PuntoHand))
	}

	if hand.PuntoHand[0] != "AS" {
		t.Errorf("PuntoHand[0] = %s, want AS", hand.PuntoHand[0])
	}

	if hand.PuntoHand[1] != "2C" {
		t.Errorf("PuntoHand[1] = %s, want 2C", hand.PuntoHand[1])
	}

	if len(hand.BankoHand) != 2 {
		t.Errorf("BankoHand length = %d, want 2", len(hand.BankoHand))
	}

	if hand.BankoHand[0] != "KH" {
		t.Errorf("BankoHand[0] = %s, want KH", hand.BankoHand[0])
	}

	if hand.BankoHand[1] != "QD" {
		t.Errorf("BankoHand[1] = %s, want QD", hand.BankoHand[1])
	}

	if hand.PuntoTotal != 3 {
		t.Errorf("PuntoTotal = %d, want 3", hand.PuntoTotal)
	}

	if hand.BankoTotal != 0 {
		t.Errorf("BankoTotal = %d, want 0", hand.BankoTotal)
	}

	if hand.Result != "punto" {
		t.Errorf("Result = %s, want punto", hand.Result)
	}

	if hand.Bet.BetOn != "punto" {
		t.Errorf("Bet.BetOn = %s, want punto", hand.Bet.BetOn)
	}

	if !hand.Bet.IsWin {
		t.Errorf("Bet.IsWin = %v, want true", hand.Bet.IsWin)
	}

	if hand.Bet.BetAmount != 10.0 {
		t.Errorf("Bet.BetAmount = %f, want 10.0", hand.Bet.BetAmount)
	}

	if hand.Bet.FinalBankroll != 1000.0 {
		t.Errorf("Bet.FinalBankroll = %f, want 1000.0", hand.Bet.FinalBankroll)
	}
}

func TestDataCollector_CollectHandData_WithThirdCard(t *testing.T) {
	dc := NewDataCollector(BetOnBanco, 6, 1000.0, 10.0, 100)
	shoe := deck.MakeNewShoe()
	dc.StartNewGame(shoe)

	state := &SimulatorState{
		CurrentBankroll: 1000.0,
		BettingOn:       puntobanco.BancoBanker,
		BetAmount:       10.0,
	}

	puntoCard1 := &deck.Card{Card: "A", Value: 1, Suit: "Spades"}
	puntoCard2 := &deck.Card{Card: "2", Value: 2, Suit: "Clubs"}
	puntoCard3 := &deck.Card{Card: "3", Value: 3, Suit: "Hearts"}
	bancoCard1 := &deck.Card{Card: "K", Value: 0, Suit: "Diamonds"}
	bancoCard2 := &deck.Card{Card: "Q", Value: 0, Suit: "Spades"}
	bancoCard3 := &deck.Card{Card: "J", Value: 0, Suit: "Clubs"}

	result := puntobanco.BancoBanker
	gameResult := &puntobanco.GameResultState{
		Result: &result,
		PuntoState: &puntobanco.PlayerState{
			FirstCard:  puntoCard1,
			SecondCard: puntoCard2,
			ThirdCard:  puntoCard3,
			Points:     6,
		},
		BancoState: &puntobanco.PlayerState{
			FirstCard:  bancoCard1,
			SecondCard: bancoCard2,
			ThirdCard:  bancoCard3,
			Points:     0,
		},
		RemainingShoe: shoe[6:],
	}

	dc.CollectHandData(state, gameResult, len(shoe))

	hand := dc.data.Games[0][0]

	if len(hand.PuntoHand) != 3 {
		t.Errorf("PuntoHand length = %d, want 3", len(hand.PuntoHand))
	}

	if hand.PuntoHand[2] != "3H" {
		t.Errorf("PuntoHand[2] = %s, want 3H", hand.PuntoHand[2])
	}

	if len(hand.BankoHand) != 3 {
		t.Errorf("BankoHand length = %d, want 3", len(hand.BankoHand))
	}

	if hand.BankoHand[2] != "JC" {
		t.Errorf("BankoHand[2] = %s, want JC", hand.BankoHand[2])
	}
}

func TestDataCollector_CollectHandData_NewShoeDetection(t *testing.T) {
	dc := NewDataCollector(BetOnPunto, 6, 1000.0, 10.0, 100)
	shoe := deck.MakeNewShoe()
	dc.StartNewGame(shoe)

	state := &SimulatorState{
		CurrentBankroll: 1000.0,
		BettingOn:       puntobanco.PuntoPlayer,
		BetAmount:       10.0,
	}

	// Simulate a new shoe being created (previous length < 8, current >= 8)
	result := puntobanco.PuntoPlayer
	gameResult := &puntobanco.GameResultState{
		Result:        &result,
		RemainingShoe: shoe, // Large shoe (>= 8 cards)
	}

	dc.previousShoeLen = 5 // Previous shoe was small
	dc.CollectHandData(state, gameResult, 5)

	if dc.currentShoeNum != 2 {
		t.Errorf("currentShoeNum = %d, want 2 (new shoe detected)", dc.currentShoeNum)
	}
}

func TestDataCollector_GetData(t *testing.T) {
	dc := NewDataCollector(BetOnPunto, 6, 1000.0, 10.0, 100)

	data := dc.GetSimulationData()

	if data == nil {
		t.Fatal("GetSimulationData() returned nil")
	}

	if data.Strategy != string(BetOnPunto) {
		t.Errorf("Strategy = %s, want %s", data.Strategy, string(BetOnPunto))
	}

	// Modify data through collector
	shoe := deck.MakeNewShoe()
	dc.StartNewGame(shoe)

	// Get data again and verify it's the same reference
	data2 := dc.GetSimulationData()

	if data != data2 {
		t.Error("GetSimulationData() should return the same reference")
	}

	if len(data2.Games) != 1 {
		t.Errorf("Games length = %d, want 1", len(data2.Games))
	}
}

func TestCreateSimulationDataFilename(t *testing.T) {
	tests := []struct {
		name                string
		strategy            string
		numberOfSimulations int
		useGzip             bool
		expectedPattern     string
		expectedExtension   string
		expectedSanitized   string
	}{
		{
			name:                "Simple strategy name without compression",
			strategy:            "Bet on Punto",
			numberOfSimulations: 50,
			useGzip:             false,
			expectedExtension:   ".json",
			expectedSanitized:   "Bet_on_Punto",
		},
		{
			name:                "Simple strategy name with compression",
			strategy:            "Bet on Banco",
			numberOfSimulations: 101,
			useGzip:             true,
			expectedExtension:   ".json.gz",
			expectedSanitized:   "Bet_on_Banco",
		},
		{
			name:                "Strategy with parentheses",
			strategy:            "Bet on Punto (player)",
			numberOfSimulations: 100,
			useGzip:             false,
			expectedExtension:   ".json",
			expectedSanitized:   "Bet_on_Punto_player",
		},
		{
			name:                "Strategy with special characters é and É",
			strategy:            "Égalité",
			numberOfSimulations: 150,
			useGzip:             true,
			expectedExtension:   ".json.gz",
			expectedSanitized:   "Egalite",
		},
		{
			name:                "Strategy with multiple special characters",
			strategy:            "Test (Strategy) with 'quotes'",
			numberOfSimulations: 1,
			useGzip:             false,
			expectedExtension:   ".json",
			expectedSanitized:   "Test_Strategy_with_quotes",
		},
		{
			name:                "Empty strategy name",
			strategy:            "",
			numberOfSimulations: 10000,
			useGzip:             true,
			expectedExtension:   ".json.gz",
			expectedSanitized:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := CreateSimulationDataFilename(tt.strategy, tt.numberOfSimulations, tt.useGzip)

			// Verify filename is not empty
			if filename == "" {
				t.Fatal("Filename should not be empty")
			}

			// Verify extension
			if !strings.HasSuffix(filename, tt.expectedExtension) {
				t.Errorf("Filename should end with %s, got: %s", tt.expectedExtension, filename)
			}

			// Verify sanitized strategy name
			if tt.expectedSanitized != "" {
				if !strings.Contains(filename, tt.expectedSanitized) {
					t.Errorf("Filename should contain sanitized strategy '%s', got: %s", tt.expectedSanitized, filename)
				}
			}

			// Verify spaces are replaced with underscores (the filename should not contain spaces)
			if strings.Contains(filename, " ") {
				t.Errorf("Filename should not contain spaces, got: %s", filename)
			}

			// Verify number of simulations
			// The format is: {sanitized}_{number}_{datetime}.{ext}
			numberPattern := regexp.MustCompile(`_\d+_`)
			if !numberPattern.MatchString(filename) {
				t.Errorf("Filename should contain number of simulations between underscores, got: %s", filename)
			}

			// Verify date and time format (YYYY-MM-DD_HH.MM.SS)
			dateTimePattern := regexp.MustCompile(`\d{4}-\d{2}-\d{2}_\d{2}\.\d{2}\.\d{2}`)
			if !dateTimePattern.MatchString(filename) {
				t.Errorf("Filename should contain date/time in format YYYY-MM-DD_HH.MM.SS, got: %s", filename)
			}
		})
	}
}
