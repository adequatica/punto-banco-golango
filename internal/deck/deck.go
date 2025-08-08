package deck

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

type BlankCard struct {
	Card  string `json:"card"`
	Value int    `json:"value"`
}

type Card struct {
	Card  string `json:"card"`
	Value int    `json:"value"`
	Suit  string `json:"suit"`
}

var Cards = []BlankCard{
	{Card: "A", Value: 1},
	{Card: "2", Value: 2},
	{Card: "3", Value: 3},
	{Card: "4", Value: 4},
	{Card: "5", Value: 5},
	{Card: "6", Value: 6},
	{Card: "7", Value: 7},
	{Card: "8", Value: 8},
	{Card: "9", Value: 9},
	{Card: "10", Value: 0},
	{Card: "J", Value: 0},
	{Card: "Q", Value: 0},
	{Card: "K", Value: 0},
}

var Suits = []string{"Spades", "Clubs", "Hearts", "Diamonds"}

var NumberOfDecks = 6

// It takes seven shuffles to randomize a deck of card
var NumberOfShuffles = 7

func MakeNewDeck(cards []BlankCard, suits []string) []Card {
	total := len(cards) * len(suits)
	deck := make([]Card, total)

	for i := 0; i < total; i++ {
		cardIdx := i / len(suits)
		suitIdx := i % len(suits)
		deck[i] = Card{
			Card:  cards[cardIdx].Card,
			Value: cards[cardIdx].Value,
			Suit:  suits[suitIdx],
		}
	}

	return deck
}

func MultiplyDeck(deck []Card, multiplier int) []Card {
	severalDecks := make([]Card, 0, len(deck)*multiplier)

	for i := 0; i < multiplier; i++ {
		severalDecks = append(severalDecks, deck...)
	}

	return severalDecks
}

func ShuffleDeck(deck []Card) []Card {
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < NumberOfShuffles; i++ {
		random.Shuffle(len(deck), func(i, j int) {
			deck[i], deck[j] = deck[j], deck[i]
		})
	}
	return deck
}

func CutDeck(deck []Card) []Card {
	if len(deck) <= 1 {
		return deck
	}

	var random = rand.New(rand.NewSource(time.Now().UnixNano()))

	cutPoint := random.Intn(len(deck)-1) + 1
	cutDeck := make([]Card, len(deck))

	copy(cutDeck, deck[cutPoint:])
	copy(cutDeck[len(deck)-cutPoint:], deck[:cutPoint])

	return cutDeck
}

// This procedure is implemented in casinos to prevent card counting
func BurnCards(deck []Card) []Card {
	if len(deck) == 0 {
		return deck
	}

	firstCard := deck[0]
	cardsToRemove := firstCard.Value

	// If the burn card value is 0 (10, Jack, Queen, King), then it is 10 for burning
	if cardsToRemove == 0 {
		cardsToRemove = 10
	}

	// Remove burning cards together with the first opened
	totalToRemove := cardsToRemove + 1
	if totalToRemove >= len(deck) {
		return []Card{}
	}

	remainingDeck := deck[totalToRemove:]

	return remainingDeck
}

func MakeNewShoe() []Card {
	deck := MakeNewDeck(Cards, Suits)
	shoe := MultiplyDeck(deck, NumberOfDecks)
	ShuffleDeck(shoe)
	shoe = CutDeck(shoe)
	shoe = BurnCards(shoe)

	return shoe
}

func GetRemainingRounds(deck []Card) string {
	if len(deck) == 0 {
		return "0"
	}

	minRounds := int(math.Floor(float64(len(deck)) / 6.0))
	maxRounds := int(math.Floor(float64(len(deck)) / 4.0))

	if minRounds == 0 && maxRounds == 0 {
		return "0"
	} else if minRounds != maxRounds {
		return fmt.Sprintf("%d–%d", minRounds, maxRounds)
	} else {
		return fmt.Sprintf("%d", minRounds)
	}
}
