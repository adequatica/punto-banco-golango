package rendering

import (
	"fmt"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
	"github.com/charmbracelet/lipgloss"
)

var (
	greenStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // Green
	redStyle             = lipgloss.NewStyle().Foreground(lipgloss.Color("1")) // Red
	blackStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // Gray
	resultIsNotAvailable = "Game result is not available"
)

func ConvertSuitToSymbol(suit string) string {
	switch suit {
	case "Spades":
		return "♠"
	case "Clubs":
		return "♣"
	case "Hearts":
		return "♥"
	case "Diamonds":
		return "♦"
	default:
		return suit
	}
}

func RenderPlayingCard(card *deck.Card) string {
	if card == nil {
		return ""
	}

	suitSymbol := ConvertSuitToSymbol(card.Suit)
	playingCard := fmt.Sprintf("%s%s", card.Card, suitSymbol)

	if card.Suit == "Hearts" || card.Suit == "Diamonds" {
		return redStyle.Render(playingCard)
	} else {
		return blackStyle.Render(playingCard)
	}
}

func ConvertStringToBetType(s string) (puntobanco.BetType, error) {
	switch s {
	case string(puntobanco.PuntoPlayer):
		return puntobanco.PuntoPlayer, nil
	case string(puntobanco.BancoBanker):
		return puntobanco.BancoBanker, nil
	case string(puntobanco.EgaliteTie):
		return puntobanco.EgaliteTie, nil
	default:
		return "", fmt.Errorf("invalid bet type: %s", s)
	}
}

func RenderBetResult(gameResult *puntobanco.BetType, betString string) string {
	if betString == "" && gameResult == nil {
		return "\n"
	}

	betType, err := ConvertStringToBetType(betString)
	if err != nil {
		return fmt.Sprintf("Invalid bet: %s\n\n", err.Error())
	}

	if gameResult == nil {
		return fmt.Sprintf("%s\n\n", resultIsNotAvailable)
	}

	if betType == *gameResult {
		return fmt.Sprintf("You %s\n\n", greenStyle.Bold(true).Render("won"))
	} else {
		return fmt.Sprintf("You %s\n\n", redStyle.Bold(true).Render("lost"))
	}
}

func RenderDrawnCards(state *puntobanco.PlayerState) string {
	if state == nil {
		return "no cards"
	}

	var cards []string

	if state.FirstCard != nil {
		cards = append(cards, RenderPlayingCard(state.FirstCard))
	}
	if state.SecondCard != nil {
		cards = append(cards, RenderPlayingCard(state.SecondCard))
	}
	if state.ThirdCard != nil {
		cards = append(cards, RenderPlayingCard(state.ThirdCard))
	}

	if len(cards) == 0 {
		return "no cards"
	}

	cardsString := ""
	for i, card := range cards {
		if i > 0 {
			cardsString += " "
		}
		cardsString += card
	}

	return fmt.Sprintf("%s = %d", cardsString, state.Points)
}

func RenderGameResultState(gameState *puntobanco.GameResultState, betString string) string {
	if gameState == nil {
		return ""
	}

	var result = "\n"
	// Render Punto state
	result += "\nPunto: "
	if gameState.PuntoState != nil {
		result += RenderDrawnCards(gameState.PuntoState)
	} else {
		result += "no cards"
	}

	// Render Banco state
	result += "\nBanco: "
	if gameState.BancoState != nil {
		result += RenderDrawnCards(gameState.BancoState)
	} else {
		result += "no cards"
	}

	result += "\n"
	result += RenderBetResult(gameState.Result, betString)

	return result
}
