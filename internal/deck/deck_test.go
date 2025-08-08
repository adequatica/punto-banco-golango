package deck

import (
	"fmt"
	"reflect"
	"testing"
)

var DefaultNumberOfCards = 52

func TestMakeNewDeck(t *testing.T) {
	get := MakeNewDeck(Cards, Suits)
	want := DefaultNumberOfCards

	if len(get) != want {
		t.Errorf("deck length of %d should have %d cards", len(get), want)
	}
}

func TestMultiplyDeck(t *testing.T) {
	originalDeck := MakeNewDeck(Cards, Suits)
	multipleDecks := MultiplyDeck(originalDeck, NumberOfDecks)
	want := DefaultNumberOfCards * NumberOfDecks

	if len(multipleDecks) != want {
		t.Errorf("multiplied deck length of %d should have %d cards", len(multipleDecks), want)
	}
}

func TestShuffleDeck(t *testing.T) {
	originalDeck := MakeNewDeck(Cards, Suits)
	originalLength := len(originalDeck)

	shuffledDeck := ShuffleDeck(MakeNewDeck(Cards, Suits))
	shuffledLength := len(shuffledDeck)

	if shuffledLength != originalLength {
		t.Errorf("shuffled deck should have the same length as original deck")
	}

	if reflect.DeepEqual(shuffledDeck, originalDeck) {
		t.Errorf("shuffled deck should not be equal to original deck")
	}

	for i, card := range shuffledDeck {
		if card.Card == "" {
			t.Errorf("card at index %d should not have an empty Card field", i)
		}
		if card.Suit == "" {
			t.Errorf("card at index %d should not have an empty Suit field", i)
		}
		if card.Value < 0 || card.Value > 10 {
			t.Errorf("card at index %d should have valid Value", i)
		}
	}
}

func TestCutDeck(t *testing.T) {
	originalDeck := MakeNewDeck(Cards, Suits)
	cuttedDecks := CutDeck(originalDeck)

	if len(cuttedDecks) != len(originalDeck) {
		t.Errorf("cutted deck length of %d should have %d cards as original deck", len(cuttedDecks), len(originalDeck))
	}

	if reflect.DeepEqual(cuttedDecks[0], originalDeck[0]) {
		t.Errorf("first cards of cutted deck should be different from original deck's first element")
	}

	if reflect.DeepEqual(cuttedDecks, originalDeck) {
		t.Errorf("cutted deck should not be equal to original deck")
	}
}

func TestBurnCards(t *testing.T) {
	originalDeck := MakeNewDeck(Cards, Suits)
	burnedDeck := BurnCards(originalDeck)

	if len(burnedDeck) == len(originalDeck) {
		t.Errorf("burned deck length of %d should not have the same length of %d as original", len(burnedDeck), len(originalDeck))
	}

	wantMinLength := len(originalDeck) - 11
	wantMaxLength := len(originalDeck) - 1

	if len(burnedDeck) < wantMinLength || len(burnedDeck) > wantMaxLength {
		t.Errorf("burned deck length of %d should be between %d and %d cards", len(burnedDeck), wantMinLength, wantMaxLength)
	}

	if reflect.DeepEqual(originalDeck, burnedDeck) {
		t.Errorf("burned deck should not be equal to original deck")
	}

	cardsRemoved := len(originalDeck) - len(burnedDeck)
	pseudoOriginalDeck := originalDeck[cardsRemoved:]

	if !reflect.DeepEqual(burnedDeck, pseudoOriginalDeck) {
		t.Errorf("burned deck should be equal to original deck except removed cards")
	}
}

func TestMakeNewShow(t *testing.T) {
	get := MakeNewShoe()
	want := DefaultNumberOfCards * NumberOfDecks

	wantMinLength := want - 11
	wantMaxLength := want - 1

	if len(get) < wantMinLength || len(get) > wantMaxLength {
		t.Errorf("shoe length of %d should be between %d and %d cards after burning", len(get), wantMinLength, wantMaxLength)
	}

	originalDeck := MakeNewDeck(Cards, Suits)
	multipleDecks := MultiplyDeck(originalDeck, NumberOfDecks)
	cardsRemoved := len(multipleDecks) - len(get)
	pseudoMultipleDecks := multipleDecks[cardsRemoved:]

	if reflect.DeepEqual(get, pseudoMultipleDecks) {
		t.Errorf("shoe should not be equal to unshuffled decks")
	}
}

func TestGetRemainingRounds(t *testing.T) {
	originalDeck := MakeNewDeck(Cards, Suits)
	multipleDecks := MultiplyDeck(originalDeck, NumberOfDecks)

	minIdealRounds := len(multipleDecks) / 6
	maxIdealRounds := len(multipleDecks) / 4

	tests := []struct {
		name string
		deck []Card
		want string
	}{
		{
			name: "nil deck",
			deck: nil,
			want: "0",
		},
		{
			name: "zero cards",
			deck: []Card{},
			want: "0",
		},
		{
			name: "cards that would result in same min and max",
			deck: make([]Card, 7),
			want: "1",
		},
		{
			name: "cards that would result in decimal values",
			deck: make([]Card, 27),
			want: "4–6",
		},
		{
			name: "cards that would result in exact division",
			deck: multipleDecks,
			want: fmt.Sprintf("%d–%d", minIdealRounds, maxIdealRounds),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetRemainingRounds(tt.deck)
			if result != tt.want {
				t.Errorf("GetRemainingRounds() = %v should be %v", result, tt.want)
			}
		})
	}
}
