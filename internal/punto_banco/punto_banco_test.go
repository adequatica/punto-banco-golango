package puntobanco

import (
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/deck"
)

func TestGetBettingOptions(t *testing.T) {
	want := GetBettingOptions()

	if len(want) != 3 {
		t.Errorf("Betting options of length %d should be 3", len(want))
	}

	if len(want) > 0 && want[0] != "Punto (player)" {
		t.Errorf("First betting option of '%s' should be 'Punto (player)'", want[0])
	}

	if len(want) > 0 && want[1] != "Banco (banker)" {
		t.Errorf("Second betting option of '%s' should be 'Banco (banker)'", want[1])
	}

	if len(want) > 0 && want[2] != "Égalité (tie)" {
		t.Errorf("Third betting option of '%s' should be 'Égalité (tie)'", want[2])
	}
}

func TestCountInitialDeal(t *testing.T) {
	tests := []struct {
		name       string
		firstCard  deck.Card
		secondCard deck.Card
		want       int
	}{
		{
			name:       "Ace + Ace = 2",
			firstCard:  deck.Card{Card: "A", Value: 1, Suit: "Spades"},
			secondCard: deck.Card{Card: "A", Value: 1, Suit: "Hearts"},
			want:       2,
		},
		{
			name:       "Ace + 9 = 10 => 0",
			firstCard:  deck.Card{Card: "A", Value: 1, Suit: "Spades"},
			secondCard: deck.Card{Card: "9", Value: 9, Suit: "Hearts"},
			want:       0,
		},
		{
			name:       "9 + 9 = 18 => 8",
			firstCard:  deck.Card{Card: "9", Value: 9, Suit: "Spades"},
			secondCard: deck.Card{Card: "9", Value: 9, Suit: "Hearts"},
			want:       8,
		},
		{
			name:       "5 + 5 = 10 => 0",
			firstCard:  deck.Card{Card: "5", Value: 5, Suit: "Spades"},
			secondCard: deck.Card{Card: "5", Value: 5, Suit: "Hearts"},
			want:       0,
		},
		{
			name:       "King + Queen = 0",
			firstCard:  deck.Card{Card: "K", Value: 0, Suit: "Spades"},
			secondCard: deck.Card{Card: "Q", Value: 0, Suit: "Hearts"},
			want:       0,
		},
		{
			name:       "Ace + King = 1",
			firstCard:  deck.Card{Card: "A", Value: 1, Suit: "Spades"},
			secondCard: deck.Card{Card: "K", Value: 0, Suit: "Hearts"},
			want:       1,
		},
		{
			name:       "Jack + 8 = 8",
			firstCard:  deck.Card{Card: "J", Value: 0, Suit: "Spades"},
			secondCard: deck.Card{Card: "8", Value: 8, Suit: "Hearts"},
			want:       8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountInitialDeal(tt.firstCard, tt.secondCard)
			if got != tt.want {
				t.Errorf("Initial deal of %v and %v cards should have value %v, but got %v", tt.firstCard, tt.secondCard, got, tt.want)
			}
		})
	}
}

func TestCountThirdCard(t *testing.T) {
	tests := []struct {
		name        string
		initialDeal int
		thirdCard   deck.Card
		want        int
	}{
		{
			name:        "Initial deal 8 + Ace = 9",
			initialDeal: 8,
			thirdCard:   deck.Card{Card: "A", Value: 1, Suit: "Spades"},
			want:        9,
		},
		{
			name:        "Initial deal 8 + 2 = 10 => 0",
			initialDeal: 8,
			thirdCard:   deck.Card{Card: "2", Value: 2, Suit: "Hearts"},
			want:        0,
		},
		{
			name:        "Initial deal 5 + 5 = 10 => 0",
			initialDeal: 5,
			thirdCard:   deck.Card{Card: "5", Value: 5, Suit: "Diamonds"},
			want:        0,
		},
		{
			name:        "Initial deal 9 + King = 9",
			initialDeal: 9,
			thirdCard:   deck.Card{Card: "K", Value: 0, Suit: "Hearts"},
			want:        9,
		},
		{
			name:        "Initial deal 0 + Queen = 0",
			initialDeal: 0,
			thirdCard:   deck.Card{Card: "Q", Value: 0, Suit: "Diamonds"},
			want:        0,
		},
		{
			name:        "Initial deal 0 + 10 = 10 => 0",
			initialDeal: 0,
			thirdCard:   deck.Card{Card: "10", Value: 0, Suit: "Diamonds"},
			want:        0,
		},
		{
			name:        "Initial deal 0 + 2 = 2",
			initialDeal: 0,
			thirdCard:   deck.Card{Card: "2", Value: 2, Suit: "Diamonds"},
			want:        2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CountThirdCard(tt.initialDeal, tt.thirdCard)
			if got != tt.want {
				t.Errorf("third card %v with initial deal %v should have value %v, but got %v", tt.thirdCard, tt.initialDeal, tt.want, got)
			}
		})
	}
}

func TestIsNatural(t *testing.T) {
	tests := []struct {
		name        string
		puntoPoints int
		bancoPoints int
		want        bool
	}{
		{
			name:        "Punto 8, Banco 6 = natural",
			puntoPoints: 8,
			bancoPoints: 6,
			want:        true,
		},
		{
			name:        "Punto 6, Banco 8 = natural",
			puntoPoints: 6,
			bancoPoints: 8,
			want:        true,
		},
		{
			name:        "Punto 9, Banco 9 = natural",
			puntoPoints: 9,
			bancoPoints: 9,
			want:        true,
		},
		{
			name:        "Punto 7, Banco 5 = not natural",
			puntoPoints: 7,
			bancoPoints: 5,
			want:        false,
		},
		{
			name:        "Punto 0, Banco 0 = not natural",
			puntoPoints: 0,
			bancoPoints: 0,
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsNatural(tt.puntoPoints, tt.bancoPoints)
			if got != tt.want {
				t.Errorf("punto %v and banco %v should count as %v natural, but got %v", tt.puntoPoints, tt.bancoPoints, tt.want, got)
			}
		})
	}

}

func TestPlayPuntoBanco(t *testing.T) {
	t.Run("new shoe", func(t *testing.T) {
		initialShoe := deck.MakeNewShoe()

		got, err := PlayPuntoBanco(initialShoe)

		if err != nil {
			t.Errorf("should not have error playing the game: %v\n", err)
		}

		if got.Result == nil {
			t.Errorf("result should not be nil")
			return
		}

		if got.PuntoState == nil {
			t.Errorf("PuntoState should not be nil")
		}

		if got.BancoState == nil {
			t.Errorf("BancoState should not be nil")
		}

		if len(got.RemainingShoe) >= len(initialShoe) {
			t.Errorf("remaining shoe length of %d should be less that initial %d", len(got.RemainingShoe), len(initialShoe))
		}
	})

	t.Run("end of shoe", func(t *testing.T) {
		initialShoe := []deck.Card{
			{Card: "A", Value: 1, Suit: "Spades"},
			{Card: "K", Value: 0, Suit: "Hearts"},
			{Card: "Q", Value: 0, Suit: "Diamonds"},
			{Card: "J", Value: 0, Suit: "Clubs"},
			{Card: "10", Value: 0, Suit: "Spades"},
			{Card: "9", Value: 9, Suit: "Hearts"},
			{Card: "8", Value: 8, Suit: "Diamonds"},
		}

		got, err := PlayPuntoBanco(initialShoe)

		if err != nil {
			t.Errorf("should not have error playing the game: %v\n", err)
		}

		if got.Result == nil {
			t.Errorf("result should not be nil")
			return
		}

		if got.PuntoState == nil {
			t.Errorf("PuntoState should not be nil")
		}

		if got.BancoState == nil {
			t.Errorf("BancoState should not be nil")
		}

		if len(got.RemainingShoe) <= len(initialShoe) {
			// A new shoe (game) should start if initial shoe is less than 8 cards
			t.Errorf("remaining shoe length of %d should be more that initial %d", len(got.RemainingShoe), len(initialShoe))
		}
	})

	t.Run("multiple shoes game", func(t *testing.T) {
		currentShoe := deck.MakeNewShoe()
		// Magic number for test iteration: 312 cards is a full shoe, it should be at least 4 games with full shoe usage
		inerations := 312

		for i := 0; i < inerations; i++ {
			got, err := PlayPuntoBanco(currentShoe)

			if err != nil {
				t.Errorf("should not have error playing the game %d: %v\n", i+1, err)
			}

			if got.Result == nil {
				t.Errorf("game %d: result should not be nil", i+1)
				return
			}

			if got.PuntoState == nil {
				t.Errorf("game %d: PuntoState should not be nil", i+1)
				return
			}

			if got.BancoState == nil {
				t.Errorf("game %d: BancoState should not be nil", i+1)
				return
			}

			if len(got.RemainingShoe) == 0 {
				t.Errorf("game %d: remaining shoe should not be empty", i+1)
				return
			}
		}
	})
}

func TestDrawThirdCardBanco(t *testing.T) {
	// Tableau for punto banco (Banker's total / Player's third card value):
	// |     | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 |
	// |-----|---|---|---|---|---|---|---|---|---|---|
	// | 0-2 | H | H | H | H | H | H | H | H | H | H |
	// | 3   | H | H | H | H | H | H | H | H | S | H |
	// | 4   | S | S | H | H | H | H | H | H | S | S |
	// | 5   | S | S | S | S | H | H | H | H | S | S |
	// | 6   | S | S | S | S | S | S | H | H | S | S |
	// | 7+  | S | S | S | S | S | S | S | S | S | S |
	// Want results: H (Hit, draw another card) = true, S (Stand, no more cards) = false

	tests := []struct {
		name           string
		bancoPoints    int
		puntoThirdCard *deck.Card
		wantShouldHit  bool
	}{
		// Banker total 0-2: always hit (H) regardless of player's third card
		{"Banco 0, Player 0", 0, &deck.Card{Value: 0}, true},
		{"Banco 0, Player 1", 0, &deck.Card{Value: 1}, true},
		{"Banco 0, Player 2", 0, &deck.Card{Value: 2}, true},
		{"Banco 0, Player 3", 0, &deck.Card{Value: 3}, true},
		{"Banco 0, Player 4", 0, &deck.Card{Value: 4}, true},
		{"Banco 0, Player 5", 0, &deck.Card{Value: 5}, true},
		{"Banco 0, Player 6", 0, &deck.Card{Value: 6}, true},
		{"Banco 0, Player 7", 0, &deck.Card{Value: 7}, true},
		{"Banco 0, Player 8", 0, &deck.Card{Value: 8}, true},
		{"Banco 0, Player 9", 0, &deck.Card{Value: 9}, true},

		{"Banco 1, Player 0", 1, &deck.Card{Value: 0}, true},
		{"Banco 1, Player 1", 1, &deck.Card{Value: 1}, true},
		{"Banco 1, Player 2", 1, &deck.Card{Value: 2}, true},
		{"Banco 1, Player 3", 1, &deck.Card{Value: 3}, true},
		{"Banco 1, Player 4", 1, &deck.Card{Value: 4}, true},
		{"Banco 1, Player 5", 1, &deck.Card{Value: 5}, true},
		{"Banco 1, Player 6", 1, &deck.Card{Value: 6}, true},
		{"Banco 1, Player 7", 1, &deck.Card{Value: 7}, true},
		{"Banco 1, Player 8", 1, &deck.Card{Value: 8}, true},
		{"Banco 1, Player 9", 1, &deck.Card{Value: 9}, true},

		{"Banco 2, Player 0", 2, &deck.Card{Value: 0}, true},
		{"Banco 2, Player 1", 2, &deck.Card{Value: 1}, true},
		{"Banco 2, Player 2", 2, &deck.Card{Value: 2}, true},
		{"Banco 2, Player 3", 2, &deck.Card{Value: 3}, true},
		{"Banco 2, Player 4", 2, &deck.Card{Value: 4}, true},
		{"Banco 2, Player 5", 2, &deck.Card{Value: 5}, true},
		{"Banco 2, Player 6", 2, &deck.Card{Value: 6}, true},
		{"Banco 2, Player 7", 2, &deck.Card{Value: 7}, true},
		{"Banco 2, Player 8", 2, &deck.Card{Value: 8}, true},
		{"Banco 2, Player 9", 2, &deck.Card{Value: 9}, true},

		// Banker total 3: hit on all except player's 8
		{"Banco 3, Player 0", 3, &deck.Card{Value: 0}, true},
		{"Banco 3, Player 1", 3, &deck.Card{Value: 1}, true},
		{"Banco 3, Player 2", 3, &deck.Card{Value: 2}, true},
		{"Banco 3, Player 3", 3, &deck.Card{Value: 3}, true},
		{"Banco 3, Player 4", 3, &deck.Card{Value: 4}, true},
		{"Banco 3, Player 5", 3, &deck.Card{Value: 5}, true},
		{"Banco 3, Player 6", 3, &deck.Card{Value: 6}, true},
		{"Banco 3, Player 7", 3, &deck.Card{Value: 7}, true},
		{"Banco 3, Player 8", 3, &deck.Card{Value: 8}, false}, // Stand on 8
		{"Banco 3, Player 9", 3, &deck.Card{Value: 9}, true},

		// Banker total 4: stand on player's 0,1 and 8,9; hit on 2-7
		{"Banco 4, Player 0", 4, &deck.Card{Value: 0}, false}, // Stand on 0
		{"Banco 4, Player 1", 4, &deck.Card{Value: 1}, false}, // Stand on 1
		{"Banco 4, Player 2", 4, &deck.Card{Value: 2}, true},
		{"Banco 4, Player 3", 4, &deck.Card{Value: 3}, true},
		{"Banco 4, Player 4", 4, &deck.Card{Value: 4}, true},
		{"Banco 4, Player 5", 4, &deck.Card{Value: 5}, true},
		{"Banco 4, Player 6", 4, &deck.Card{Value: 6}, true},
		{"Banco 4, Player 7", 4, &deck.Card{Value: 7}, true},
		{"Banco 4, Player 8", 4, &deck.Card{Value: 8}, false}, // Stand on 8
		{"Banco 4, Player 9", 4, &deck.Card{Value: 9}, false}, // Stand on 9

		// Banker total 5: stand on player's 0-3 and 8,9; hit on 4-7
		{"Banco 5, Player 0", 5, &deck.Card{Value: 0}, false}, // Stand on 0
		{"Banco 5, Player 1", 5, &deck.Card{Value: 1}, false}, // Stand on 1
		{"Banco 5, Player 2", 5, &deck.Card{Value: 2}, false}, // Stand on 2
		{"Banco 5, Player 3", 5, &deck.Card{Value: 3}, false}, // Stand on 3
		{"Banco 5, Player 4", 5, &deck.Card{Value: 4}, true},
		{"Banco 5, Player 5", 5, &deck.Card{Value: 5}, true},
		{"Banco 5, Player 6", 5, &deck.Card{Value: 6}, true},
		{"Banco 5, Player 7", 5, &deck.Card{Value: 7}, true},
		{"Banco 5, Player 8", 5, &deck.Card{Value: 8}, false}, // Stand on 8
		{"Banco 5, Player 9", 5, &deck.Card{Value: 9}, false}, // Stand on 9

		// Banker total 6: stand on player's 0-5 and 8,9; hit on 6-7
		{"Banco 6, Player 0", 6, &deck.Card{Value: 0}, false}, // Stand on 0
		{"Banco 6, Player 1", 6, &deck.Card{Value: 1}, false}, // Stand on 1
		{"Banco 6, Player 2", 6, &deck.Card{Value: 2}, false}, // Stand on 2
		{"Banco 6, Player 3", 6, &deck.Card{Value: 3}, false}, // Stand on 3
		{"Banco 6, Player 4", 6, &deck.Card{Value: 4}, false}, // Stand on 4
		{"Banco 6, Player 5", 6, &deck.Card{Value: 5}, false}, // Stand on 5
		{"Banco 6, Player 6", 6, &deck.Card{Value: 6}, true},
		{"Banco 6, Player 7", 6, &deck.Card{Value: 7}, true},
		{"Banco 6, Player 8", 6, &deck.Card{Value: 8}, false}, // Stand on 8
		{"Banco 6, Player 9", 6, &deck.Card{Value: 9}, false}, // Stand on 9

		// Banker total 7: always stand (no more cards)
		{"Banco 7, Player 0", 7, &deck.Card{Value: 0}, false},
		{"Banco 7, Player 1", 7, &deck.Card{Value: 1}, false},
		{"Banco 7, Player 2", 7, &deck.Card{Value: 2}, false},
		{"Banco 7, Player 3", 7, &deck.Card{Value: 3}, false},
		{"Banco 7, Player 4", 7, &deck.Card{Value: 4}, false},
		{"Banco 7, Player 5", 7, &deck.Card{Value: 5}, false},
		{"Banco 7, Player 6", 7, &deck.Card{Value: 6}, false},
		{"Banco 7, Player 7", 7, &deck.Card{Value: 7}, false},
		{"Banco 7, Player 8", 7, &deck.Card{Value: 8}, false},
		{"Banco 7, Player 9", 7, &deck.Card{Value: 9}, false},

		// Banker total 8: always stand (no more cards)
		{"Banco 8, Player 0", 8, &deck.Card{Value: 0}, false},
		{"Banco 8, Player 1", 8, &deck.Card{Value: 1}, false},
		{"Banco 8, Player 2", 8, &deck.Card{Value: 2}, false},
		{"Banco 8, Player 3", 8, &deck.Card{Value: 3}, false},
		{"Banco 8, Player 4", 8, &deck.Card{Value: 4}, false},
		{"Banco 8, Player 5", 8, &deck.Card{Value: 5}, false},
		{"Banco 8, Player 6", 8, &deck.Card{Value: 6}, false},
		{"Banco 8, Player 7", 8, &deck.Card{Value: 7}, false},
		{"Banco 8, Player 8", 8, &deck.Card{Value: 8}, false},
		{"Banco 8, Player 9", 8, &deck.Card{Value: 9}, false},

		// Banker total 9: always stand (no more cards)
		{"Banco 9, Player 0", 9, &deck.Card{Value: 0}, false},
		{"Banco 9, Player 1", 9, &deck.Card{Value: 1}, false},
		{"Banco 9, Player 2", 9, &deck.Card{Value: 2}, false},
		{"Banco 9, Player 3", 9, &deck.Card{Value: 3}, false},
		{"Banco 9, Player 4", 9, &deck.Card{Value: 4}, false},
		{"Banco 9, Player 5", 9, &deck.Card{Value: 5}, false},
		{"Banco 9, Player 6", 9, &deck.Card{Value: 6}, false},
		{"Banco 9, Player 7", 9, &deck.Card{Value: 7}, false},
		{"Banco 9, Player 8", 9, &deck.Card{Value: 8}, false},
		{"Banco 9, Player 9", 9, &deck.Card{Value: 9}, false},

		// Edge case: player doesn't draw the third card (nil)
		{"Banco 3, Player no third card", 3, nil, true},
		{"Banco 4, Player no third card", 4, nil, true},
		{"Banco 5, Player no third card", 5, nil, true},
		{"Banco 6, Player no third card", 6, nil, false},
		{"Banco 7, Player no third card", 7, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DrawThirdCardBanco(tt.bancoPoints, tt.puntoThirdCard)
			if got != tt.wantShouldHit {
				action := "hit"
				if !tt.wantShouldHit {
					action = "stand"
				}
				t.Errorf("banco draws the third card (%d, %v) = %v, got %v (should %s)",
					tt.bancoPoints, tt.puntoThirdCard, got, tt.wantShouldHit, action)
			}
		})
	}
}

func TestDetermineResult(t *testing.T) {
	tests := []struct {
		name        string
		puntoPoints int
		bancoPoints int
		want        BetType
	}{
		{
			name:        "Punto 9, Banco 6 = Punto wins",
			puntoPoints: 9,
			bancoPoints: 6,
			want:        PuntoPlayer,
		},
		{
			name:        "Punto 8, Banco 7 = Punto wins",
			puntoPoints: 8,
			bancoPoints: 7,
			want:        PuntoPlayer,
		},
		{
			name:        "Punto 6, Banco 9 = Banco wins",
			puntoPoints: 6,
			bancoPoints: 9,
			want:        BancoBanker,
		},
		{
			name:        "Punto 5, Banco 8 = Banco wins",
			puntoPoints: 5,
			bancoPoints: 8,
			want:        BancoBanker,
		},
		{
			name:        "Punto 9, Banco 9 = Tie",
			puntoPoints: 9,
			bancoPoints: 9,
			want:        EgaliteTie,
		},
		{
			name:        "Punto 8, Banco 8 = Tie",
			puntoPoints: 8,
			bancoPoints: 8,
			want:        EgaliteTie,
		},
		{
			name:        "Punto 0, Banco 0 = Tie",
			puntoPoints: 0,
			bancoPoints: 0,
			want:        EgaliteTie,
		},
		{
			name:        "Punto 6, Banco 5 = Punto wins",
			puntoPoints: 6,
			bancoPoints: 5,
			want:        PuntoPlayer,
		},
		{
			name:        "Punto 4, Banco 7 = Banco wins",
			puntoPoints: 4,
			bancoPoints: 7,
			want:        BancoBanker,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineResult(tt.puntoPoints, tt.bancoPoints)
			if got != tt.want {
				t.Errorf("punto %v and banco %v should result as %v, but got %v", tt.puntoPoints, tt.bancoPoints, tt.want, got)
			}
		})
	}
}

func TestDetermineGameResultState(t *testing.T) {
	shoe := deck.MakeNewShoe()

	tests := []struct {
		name          string
		puntoState    PlayerState
		bancoState    PlayerState
		remainingShoe []deck.Card
		wantResult    BetType
		wantShoeLen   int
	}{
		{
			name: "Punto wins with 9 vs Banco 6",
			puntoState: PlayerState{
				FirstCard:  &deck.Card{Card: "9", Value: 9, Suit: "Spades"},
				SecondCard: &deck.Card{Card: "K", Value: 0, Suit: "Hearts"},
				ThirdCard:  nil,
				Points:     9,
			},
			bancoState: PlayerState{
				FirstCard:  &deck.Card{Card: "6", Value: 6, Suit: "Diamonds"},
				SecondCard: &deck.Card{Card: "Q", Value: 0, Suit: "Clubs"},
				ThirdCard:  nil,
				Points:     6,
			},
			remainingShoe: []deck.Card{
				{Card: "A", Value: 1, Suit: "Spades"},
			},
			wantResult:  PuntoPlayer,
			wantShoeLen: 1,
		},
		{
			name: "Banco wins with 8 vs Punto 0",
			puntoState: PlayerState{
				FirstCard:  &deck.Card{Card: "10", Value: 0, Suit: "Spades"},
				SecondCard: &deck.Card{Card: "J", Value: 0, Suit: "Hearts"},
				ThirdCard:  nil,
				Points:     0,
			},
			bancoState: PlayerState{
				FirstCard:  &deck.Card{Card: "8", Value: 8, Suit: "Diamonds"},
				SecondCard: &deck.Card{Card: "K", Value: 0, Suit: "Clubs"},
				ThirdCard:  nil,
				Points:     8,
			},
			remainingShoe: []deck.Card{
				{Card: "A", Value: 1, Suit: "Spades"},
				{Card: "1", Value: 1, Suit: "Hearts"},
			},
			wantResult:  BancoBanker,
			wantShoeLen: 2,
		},
		{
			name: "Tie with both having 7 and empty shoe",
			puntoState: PlayerState{
				FirstCard:  &deck.Card{Card: "7", Value: 7, Suit: "Spades"},
				SecondCard: &deck.Card{Card: "K", Value: 0, Suit: "Hearts"},
				ThirdCard:  nil,
				Points:     7,
			},
			bancoState: PlayerState{
				FirstCard:  &deck.Card{Card: "7", Value: 7, Suit: "Diamonds"},
				SecondCard: &deck.Card{Card: "Q", Value: 0, Suit: "Clubs"},
				ThirdCard:  nil,
				Points:     7,
			},
			remainingShoe: []deck.Card{},
			wantResult:    EgaliteTie,
			wantShoeLen:   0,
		},
		{
			name: "Tie with both having 0 and full shoe",
			puntoState: PlayerState{
				FirstCard:  &deck.Card{Card: "J", Value: 0, Suit: "Spades"},
				SecondCard: &deck.Card{Card: "K", Value: 0, Suit: "Hearts"},
				ThirdCard:  &deck.Card{Card: "Q", Value: 0, Suit: "Diamonds"},
				Points:     6,
			},
			bancoState: PlayerState{
				FirstCard:  &deck.Card{Card: "J", Value: 0, Suit: "Clubs"},
				SecondCard: &deck.Card{Card: "K", Value: 0, Suit: "Spades"},
				ThirdCard:  &deck.Card{Card: "Q", Value: 0, Suit: "Hearts"},
				Points:     9,
			},
			remainingShoe: shoe,
			wantResult:    BancoBanker,
			wantShoeLen:   len(shoe),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetermineGameResultState(tt.puntoState, tt.bancoState, tt.remainingShoe)

			if got.Result == nil {
				t.Errorf("result should not be nil")
				return
			}
			if *got.Result != tt.wantResult {
				t.Errorf("result should be %v, but got %v", tt.wantResult, *got.Result)
			}

			if got.PuntoState == nil {
				t.Errorf("PuntoState should not be nil")
			}

			if got.BancoState == nil {
				t.Errorf("BancoState should not be nil")
			}

			if len(got.RemainingShoe) != tt.wantShoeLen {
				t.Errorf("remaining shoe length of %d should be %d", len(got.RemainingShoe), tt.wantShoeLen)
			}
		})
	}
}
