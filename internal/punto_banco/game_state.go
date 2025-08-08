package puntobanco

import (
	"fmt"
	"strings"

	"github.com/adequatica/punto-banco-golango/internal/deck"
)

type PlayerState struct {
	FirstCard  *deck.Card
	SecondCard *deck.Card
	ThirdCard  *deck.Card
	Points     int
}

type GameResultState struct {
	Result        *BetType
	PuntoState    *PlayerState
	BancoState    *PlayerState
	RemainingShoe []deck.Card
}

func GetNewGameResultState() GameResultState {
	return GameResultState{
		Result:        nil,
		PuntoState:    nil,
		BancoState:    nil,
		RemainingShoe: deck.MakeNewShoe(),
	}
}

func (g *GameResultState) GetResult() *BetType {
	return g.Result
}

func (g *GameResultState) SetResult(result *BetType) {
	g.Result = result
}

func (g *GameResultState) GetShoe() []deck.Card {
	return g.RemainingShoe
}

func (g *GameResultState) SetShoe(shoe []deck.Card) error {
	if len(shoe) == 0 {
		return fmt.Errorf("empty shoe is not allowed for the next round")
	}

	g.RemainingShoe = shoe
	return nil
}

func (g *GameResultState) Render() string {
	var result strings.Builder

	// Render Punto state
	result.WriteString("Punto: ")
	if g.PuntoState != nil {
		result.WriteString(g.renderPlayerState(g.PuntoState))
	} else {
		result.WriteString("no cards")
	}
	result.WriteString("\n")

	// Render Banco state
	result.WriteString("Banco: ")
	if g.BancoState != nil {
		result.WriteString(g.renderPlayerState(g.BancoState))
	} else {
		result.WriteString("no cards")
	}
	result.WriteString("\n")

	return result.String()
}

func (g *GameResultState) renderPlayerState(state *PlayerState) string {
	var cards []string

	if state.FirstCard != nil {
		cards = append(cards, fmt.Sprintf("%s of %s", state.FirstCard.Card, state.FirstCard.Suit))
	}
	if state.SecondCard != nil {
		cards = append(cards, fmt.Sprintf("%s of %s", state.SecondCard.Card, state.SecondCard.Suit))
	}
	if state.ThirdCard != nil {
		cards = append(cards, fmt.Sprintf("%s of %s", state.ThirdCard.Card, state.ThirdCard.Suit))
	}

	if len(cards) == 0 {
		return "no cards"
	}

	cardString := strings.Join(cards, " ")
	return fmt.Sprintf("%s = %d", cardString, state.Points)
}
