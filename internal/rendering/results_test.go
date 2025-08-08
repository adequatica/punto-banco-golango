package rendering

import (
	"strings"
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

func TestConvertSuitToSymbol(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Spades suit",
			input: "Spades",
			want:  "♠",
		},
		{
			name:  "Clubs suit",
			input: "Clubs",
			want:  "♣",
		},
		{
			name:  "Hearts suit",
			input: "Hearts",
			want:  "♥",
		},
		{
			name:  "Diamonds suit",
			input: "Diamonds",
			want:  "♦",
		},
		{
			name:  "unknown suit",
			input: "Unknown",
			want:  "Unknown",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConvertSuitToSymbol(tt.input)
			if got != tt.want {
				t.Errorf("ConvertSuitToSymbol() = %v should be %v", got, tt.want)
			}
		})
	}
}

func TestRenderPlayingCard(t *testing.T) {
	t.Run("valid cards", func(t *testing.T) {
		tests := []struct {
			name string
			card *deck.Card
			want string
		}{
			{
				name: "Ace of Spades",
				card: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				want: "A♠",
			},
			{
				name: "King of Hearts",
				card: &deck.Card{
					Card:  "K",
					Value: 0,
					Suit:  "Hearts",
				},
				want: "K♥",
			},
			{
				name: "10 of Diamonds",
				card: &deck.Card{
					Card:  "10",
					Value: 0,
					Suit:  "Diamonds",
				},
				want: "10♦",
			},
			{
				name: "2 of Clubs",
				card: &deck.Card{
					Card:  "2",
					Value: 2,
					Suit:  "Clubs",
				},
				want: "2♣",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := RenderPlayingCard(tt.card)
				if got != tt.want {
					t.Errorf("RenderPlayingCard() = %v should be %v", got, tt.want)
				}
			})
		}
	})

	t.Run("nil card", func(t *testing.T) {
		result := RenderPlayingCard(nil)
		if result != "" {
			t.Errorf("RenderPlayingCard(nil) = %v should be empty", result)
		}
	})
}
func TestConvertStringToBetType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    puntobanco.BetType
		wantErr bool
	}{
		{
			name:    "valid Punto player",
			input:   string(puntobanco.PuntoPlayer),
			want:    puntobanco.PuntoPlayer,
			wantErr: false,
		},
		{
			name:    "valid Banco banker",
			input:   string(puntobanco.BancoBanker),
			want:    puntobanco.BancoBanker,
			wantErr: false,
		},
		{
			name:    "valid Égalité tie",
			input:   string(puntobanco.EgaliteTie),
			want:    puntobanco.EgaliteTie,
			wantErr: false,
		},
		{
			name:    "invalid bet type",
			input:   "invalid",
			want:    "",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertStringToBetType(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertStringToBetType() error = %v should be %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertStringToBetType() = %v should be %v", got, tt.want)
			}
		})
	}
}

func TestRenderBetResult(t *testing.T) {
	t.Run("valid game results", func(t *testing.T) {
		tests := []struct {
			name       string
			betString  string
			gameResult puntobanco.BetType
			wantWin    bool // true if result contains "win", false if contains "lose" or "invalid"
			wantError  bool // true if result contains "invalid"
		}{
			{
				name:       "winning bet - Punto",
				betString:  string(puntobanco.PuntoPlayer),
				gameResult: puntobanco.PuntoPlayer,
				wantWin:    true,
				wantError:  false,
			},
			{
				name:       "winning bet - Banco",
				betString:  string(puntobanco.BancoBanker),
				gameResult: puntobanco.BancoBanker,
				wantWin:    true,
				wantError:  false,
			},
			{
				name:       "winning bet - Égalité",
				betString:  string(puntobanco.EgaliteTie),
				gameResult: puntobanco.EgaliteTie,
				wantWin:    true,
				wantError:  false,
			},
			{
				name:       "losing bet - Punto vs Banco",
				betString:  string(puntobanco.PuntoPlayer),
				gameResult: puntobanco.BancoBanker,
				wantWin:    false,
				wantError:  false,
			},
			{
				name:       "losing bet - Banco vs Punto",
				betString:  string(puntobanco.BancoBanker),
				gameResult: puntobanco.PuntoPlayer,
				wantWin:    false,
				wantError:  false,
			},
			{
				name:       "losing bet - Égalité vs Punto",
				betString:  string(puntobanco.EgaliteTie),
				gameResult: puntobanco.PuntoPlayer,
				wantWin:    false,
				wantError:  false,
			},
			{
				name:       "invalid bet string",
				betString:  "invalid",
				gameResult: puntobanco.PuntoPlayer,
				wantWin:    false,
				wantError:  true,
			},
			{
				name:       "empty bet string",
				betString:  "",
				gameResult: puntobanco.PuntoPlayer,
				wantWin:    false,
				wantError:  true,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result := RenderBetResult(&tt.gameResult, tt.betString)

				containsWin := strings.Contains(result, "won")
				containsLose := strings.Contains(result, "lost")
				containsInvalid := strings.Contains(result, "Invalid")

				if tt.wantError {
					if !containsInvalid {
						t.Errorf("RenderBetResult() should contain 'Invalid' for invalid input, got: %s", result)
					}
				} else if tt.wantWin {
					if !containsWin {
						t.Errorf("RenderBetResult() should contain 'won' for winning bet, got: %s", result)
					}
				} else {
					if !containsLose {
						t.Errorf("RenderBetResult() should contain 'lost' for losing bet, got: %s", result)
					}
				}
			})
		}
	})

	t.Run("nil game results", func(t *testing.T) {
		result := RenderBetResult(nil, string(puntobanco.PuntoPlayer))

		if result == "" {
			t.Errorf("RenderBetResult() should not return empty string for nil gameResult")
		}

		if !strings.Contains(result, "not available") {
			t.Errorf("RenderBetResult() should contain 'not available' for nil gameResult")
		}
	})
}

func TestRenderDrawnCards(t *testing.T) {
	t.Run("valid player state with all cards", func(t *testing.T) {
		state := &puntobanco.PlayerState{
			FirstCard: &deck.Card{
				Card:  "A",
				Value: 1,
				Suit:  "Spades",
			},
			SecondCard: &deck.Card{
				Card:  "K",
				Value: 0,
				Suit:  "Hearts",
			},
			ThirdCard: &deck.Card{
				Card:  "Q",
				Value: 0,
				Suit:  "Diamonds",
			},
			Points: 1,
		}

		result := RenderDrawnCards(state)
		want := "A♠ K♥ Q♦ = 1"

		if result != want {
			t.Errorf("RenderDrawnCards() = %v should be %v", result, want)
		}
	})

	t.Run("player state with two cards", func(t *testing.T) {
		state := &puntobanco.PlayerState{
			FirstCard: &deck.Card{
				Card:  "J",
				Value: 0,
				Suit:  "Clubs",
			},
			SecondCard: &deck.Card{
				Card:  "9",
				Value: 9,
				Suit:  "Spades",
			},
			ThirdCard: nil,
			Points:    9,
		}

		result := RenderDrawnCards(state)
		want := "J♣ 9♠ = 9"

		if result != want {
			t.Errorf("RenderDrawnCards() = %v should be %v", result, want)
		}
	})

	t.Run("player state with one card", func(t *testing.T) {
		state := &puntobanco.PlayerState{
			FirstCard: &deck.Card{
				Card:  "10",
				Value: 0,
				Suit:  "Hearts",
			},
			SecondCard: nil,
			ThirdCard:  nil,
			Points:     0,
		}

		result := RenderDrawnCards(state)
		want := "10♥ = 0"

		if result != want {
			t.Errorf("RenderDrawnCards() = %v shold be %v", result, want)
		}
	})

	t.Run("player state with no cards", func(t *testing.T) {
		state := &puntobanco.PlayerState{
			FirstCard:  nil,
			SecondCard: nil,
			ThirdCard:  nil,
			Points:     0,
		}

		result := RenderDrawnCards(state)
		want := "no cards"

		if result != want {
			t.Errorf("RenderDrawnCards() = %v should be %v", result, want)
		}
	})

	t.Run("nil player state", func(t *testing.T) {
		result := RenderDrawnCards(nil)
		want := "no cards"

		if result != want {
			t.Errorf("RenderDrawnCards(nil) = %v should be %v", result, want)
		}
	})
}

func TestRenderGameResultState(t *testing.T) {
	t.Run("complete game state", func(t *testing.T) {
		puntoResult := puntobanco.PuntoPlayer
		gameState := &puntobanco.GameResultState{
			Result: &puntoResult,
			PuntoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				SecondCard: &deck.Card{
					Card:  "K",
					Value: 0,
					Suit:  "Hearts",
				},
				ThirdCard: nil,
				Points:    1,
			},
			BancoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "10",
					Value: 0,
					Suit:  "Clubs",
				},
				SecondCard: &deck.Card{
					Card:  "J",
					Value: 0,
					Suit:  "Diamonds",
				},
				ThirdCard: nil,
				Points:    0,
			},
			RemainingShoe: []deck.Card{},
		}

		result := RenderGameResultState(gameState, string(puntobanco.PuntoPlayer))

		if !strings.Contains(result, "Punto:") {
			t.Errorf("RenderGameResultState() should contain 'Punto:'")
		}
		if !strings.Contains(result, "Banco:") {
			t.Errorf("RenderGameResultState() should contain 'Banco:'")
		}
		if !strings.Contains(result, "You won") {
			t.Errorf("RenderGameResultState() should contain 'You won' for winning bet")
		}
	})

	t.Run("game state with nil punto state", func(t *testing.T) {
		bancoResult := puntobanco.BancoBanker
		gameState := &puntobanco.GameResultState{
			Result:     &bancoResult,
			PuntoState: nil,
			BancoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     1,
			},
			RemainingShoe: []deck.Card{},
		}

		result := RenderGameResultState(gameState, string(puntobanco.BancoBanker))

		if !strings.Contains(result, "Punto: no cards") {
			t.Errorf("RenderGameResultState() should show 'no cards' for nil PuntoState")
		}
		if !strings.Contains(result, "Banco:") {
			t.Errorf("RenderGameResultState() should contain 'Banco:'")
		}
		if !strings.Contains(result, "You won") {
			t.Errorf("RenderGameResultState() should contain 'You won' for winning bet")
		}
	})

	t.Run("game state with nil banco state", func(t *testing.T) {
		puntoResult := puntobanco.PuntoPlayer
		gameState := &puntobanco.GameResultState{
			Result: &puntoResult,
			PuntoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     1,
			},
			BancoState:    nil,
			RemainingShoe: []deck.Card{},
		}

		result := RenderGameResultState(gameState, string(puntobanco.EgaliteTie))

		if !strings.Contains(result, "Punto:") {
			t.Errorf("RenderGameResultState() should contain 'Punto:'")
		}
		if !strings.Contains(result, "Banco: no cards") {
			t.Errorf("RenderGameResultState() should show 'no cards' for nil BancoState")
		}
		if !strings.Contains(result, "You lost") {
			t.Errorf("RenderGameResultState() should contain 'You lost' for losing bet")
		}
	})

	t.Run("game state with nil result", func(t *testing.T) {
		gameState := &puntobanco.GameResultState{
			Result: nil,
			PuntoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     1,
			},
			BancoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "K",
					Value: 0,
					Suit:  "Hearts",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     0,
			},
			RemainingShoe: []deck.Card{},
		}

		result := RenderGameResultState(gameState, string(puntobanco.PuntoPlayer))

		if !strings.Contains(result, resultIsNotAvailable) {
			t.Errorf("RenderGameResultState() should show '%s' for nil result", resultIsNotAvailable)
		}
	})

	t.Run("nil game state", func(t *testing.T) {
		result := RenderGameResultState(nil, string(puntobanco.PuntoPlayer))

		if result != "" {
			t.Errorf("RenderGameResultState(nil) = %v should be empty string", result)
		}
	})

	t.Run("invalid bet string", func(t *testing.T) {
		puntoResult := puntobanco.PuntoPlayer
		gameState := &puntobanco.GameResultState{
			Result: &puntoResult,
			PuntoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "A",
					Value: 1,
					Suit:  "Spades",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     1,
			},
			BancoState: &puntobanco.PlayerState{
				FirstCard: &deck.Card{
					Card:  "K",
					Value: 0,
					Suit:  "Hearts",
				},
				SecondCard: nil,
				ThirdCard:  nil,
				Points:     0,
			},
			RemainingShoe: []deck.Card{},
		}

		result := RenderGameResultState(gameState, "invalid")

		if !strings.Contains(result, "Invalid bet:") {
			t.Errorf("RenderGameResultState() should contain 'Invalid bet:' for invalid bet string")
		}
	})
}
