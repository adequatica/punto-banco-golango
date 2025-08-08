package puntobanco

import (
	"fmt"

	"github.com/adequatica/punto-banco-golango/internal/deck"
)

type BetType string

const (
	PuntoPlayer BetType = "Punto (player)"
	BancoBanker BetType = "Banco (banker)"
	EgaliteTie  BetType = "Égalité (tie)"
)

func GetBettingOptions() []string {
	return []string{
		string(PuntoPlayer),
		string(BancoBanker),
		string(EgaliteTie),
	}
}

func CountInitialDeal(firstCard deck.Card, secondCard deck.Card) int {
	firstCardValue := firstCard.Value
	secondCardValue := secondCard.Value

	sum := firstCardValue + secondCardValue
	// Only the units digit of the sum is counted, as valuing hands modulo 10
	return sum % 10
}

func CountThirdCard(initialDeal int, thirdCard deck.Card) int {
	thirdCardValue := thirdCard.Value

	sum := initialDeal + thirdCardValue
	// Only the units digit of the sum is counted, as valuing hands modulo 10
	return sum % 10
}

func IsNatural(puntoPoints int, bancoPoints int) bool {
	return puntoPoints >= 8 || bancoPoints >= 8
}

func PlayPuntoBanco(shoe []deck.Card) (GameResultState, error) {
	// A cut-card is usully placed in front of the seventh from last card to indicate the last round of the shoe
	if len(shoe) < 8 {
		shoe = deck.MakeNewShoe()
	}

	// Create a copy of the shoe to avoid modifying the original
	gameShoe := make([]deck.Card, len(shoe))
	copy(gameShoe, shoe)

	puntoState := PlayerState{}
	bancoState := PlayerState{}

	// Deal first four cards
	// Punto (player) gets 1st and 3rd cards, Banco (banker) gets 2nd and 4th cards
	puntoState.FirstCard = &deck.Card{Card: gameShoe[0].Card, Value: gameShoe[0].Value, Suit: gameShoe[0].Suit}
	puntoState.SecondCard = &deck.Card{Card: gameShoe[2].Card, Value: gameShoe[2].Value, Suit: gameShoe[2].Suit}
	bancoState.FirstCard = &deck.Card{Card: gameShoe[1].Card, Value: gameShoe[1].Value, Suit: gameShoe[1].Suit}
	bancoState.SecondCard = &deck.Card{Card: gameShoe[3].Card, Value: gameShoe[3].Value, Suit: gameShoe[3].Suit}

	// Remove played cards from shoe
	gameShoe = gameShoe[4:]

	// Count points of initial deal
	puntoState.Points = CountInitialDeal(*puntoState.FirstCard, *puntoState.SecondCard)
	bancoState.Points = CountInitialDeal(*bancoState.FirstCard, *bancoState.SecondCard)

	// Check for 'natural' (8 or 9)
	if IsNatural(puntoState.Points, bancoState.Points) {
		return DetermineGameResultState(puntoState, bancoState, gameShoe), nil
	}

	// Player's rule for third card
	if puntoState.Points <= 5 {
		if len(gameShoe) == 0 {
			return GameResultState{}, fmt.Errorf("insufficient cards to draw a third card for Punto (player)")
		}

		puntoState.ThirdCard = &gameShoe[0]
		puntoState.Points = CountThirdCard(puntoState.Points, *puntoState.ThirdCard)
		// Remove drawn card from shoe
		gameShoe = gameShoe[1:]
	}

	// Banker's rule for third card
	shouldBancoDraw := DrawThirdCardBanco(bancoState.Points, puntoState.ThirdCard)
	if shouldBancoDraw {
		if len(gameShoe) == 0 {
			return GameResultState{}, fmt.Errorf("insufficient cards to draw a third card for Banco (banker)")
		}

		bancoState.ThirdCard = &gameShoe[0]
		bancoState.Points = CountThirdCard(bancoState.Points, *bancoState.ThirdCard)
		// Remove drawn card from shoe
		gameShoe = gameShoe[1:]
	}

	return DetermineGameResultState(puntoState, bancoState, gameShoe), nil
}

func DrawThirdCardBanco(bancoPoints int, puntoThirdCard *deck.Card) bool {
	// Banker stands on 7 or higher (does not draw a third card)
	if bancoPoints >= 7 {
		return false
	}

	// If the Banker's total is 0, 1, or 2, the Banker draws a third card regardless of the Player's third card
	if bancoPoints <= 2 {
		return true
	}

	// Banker draws a third card if player didn't draw a third card (stood on 6 or 7),
	// and the Banker's initial total is 0, 1, 2, 3, 4, or 5
	if puntoThirdCard == nil {
		return bancoPoints <= 5
	}

	puntoThirdCardValue := puntoThirdCard.Value

	switch bancoPoints {
	case 3:
		// Banker draws a third card unless the Player's third card is an 8
		// (draws when the Player's third card is a 9 or a 10/face card = 0)
		return puntoThirdCardValue != 8
	case 4:
		// If the Banker's total is 4, the Banker draws a third card,
		// if the Player's third card is a 2, 3, 4, 5, 6, or 7
		return puntoThirdCardValue >= 2 && puntoThirdCardValue <= 7
	case 5:
		// If the Banker's total is 5, the Banker draws a third card,
		// if the Player's third card is a 4, 5, 6, or 7
		return puntoThirdCardValue >= 4 && puntoThirdCardValue <= 7
	case 6:
		// If the Banker's total is 6, the Banker draws a third card
		// if the Player's third card is a 6 or 7
		return puntoThirdCardValue == 6 || puntoThirdCardValue == 7
	default:
		return false
	}
}

func DetermineResult(puntoPoints int, bancoPoints int) BetType {
	if puntoPoints > bancoPoints {
		return PuntoPlayer
	} else if bancoPoints > puntoPoints {
		return BancoBanker
	} else {
		return EgaliteTie
	}
}

func DetermineGameResultState(puntoState PlayerState, bancoState PlayerState, remainingShoe []deck.Card) GameResultState {
	winner := DetermineResult(puntoState.Points, bancoState.Points)

	return GameResultState{
		Result:        &winner,
		PuntoState:    &puntoState,
		BancoState:    &bancoState,
		RemainingShoe: remainingShoe,
	}
}
