package puntobanco

import (
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/deck"
)

func TestGetNewGameResultState(t *testing.T) {
	newGameState := GetNewGameResultState()

	if newGameState.GetShoe() == nil {
		t.Error("nextShoe should not be nil")
	}

	if newGameState.GetResult() != nil {
		t.Error("lastResult should be nil by default")
	}
}

func TestGameState_GetResult(t *testing.T) {
	gameState := GetNewGameResultState()

	if gameState.GetResult() != nil {
		t.Error("GetResult should return nil by default")
	}

	result := PuntoPlayer
	gameState.SetResult(&result)

	if gameState.GetResult() == nil {
		t.Error("GetResult should not be nil after setting a result")
	}

	if *gameState.GetResult() != result {
		t.Errorf("GetResult should return %v after setting %v", *gameState.GetResult(), result)
	}
}

func TestGameState_SetResult(t *testing.T) {
	gameState := GetNewGameResultState()

	result := BancoBanker
	gameState.SetResult(&result)

	if gameState.GetResult() == nil {
		t.Error("GetResult should not be nil after setting a result")
	}

	if *gameState.GetResult() != result {
		t.Errorf("GetResult should return %v after setting %v", *gameState.GetResult(), result)
	}

	gameState.SetResult(nil)

	if gameState.GetResult() != nil {
		t.Error("GetResult should be nil after setting nil")
	}
}

func TestGameState_GetShoe(t *testing.T) {
	gameState := GetNewGameResultState()

	if gameState.GetShoe() == nil {
		t.Error("GetShoe should not return nil by default")
	}

	shoe := gameState.GetShoe()
	if len(shoe) == 0 {
		t.Error("GetShoe should return a non-empty slice")
	}
}

func TestGameState_SetShoe(t *testing.T) {
	t.Run("set valid shoe", func(t *testing.T) {
		gameState := GetNewGameResultState()

		validShoe := deck.MakeNewShoe()
		err := gameState.SetShoe(validShoe)
		if err != nil {
			t.Errorf("SetShoe should not return error for valid shoe")
		}

		setShoe := gameState.GetShoe()
		if len(setShoe) != len(validShoe) {
			t.Errorf("Set shoe length %d should match input shoe length %d", len(setShoe), len(validShoe))
		}
	})

	t.Run("set invalid shoe", func(t *testing.T) {
		gameState := GetNewGameResultState()

		emptyShoe := []deck.Card{}
		err := gameState.SetShoe(emptyShoe)
		if err == nil {
			t.Error("SetShoe should return error for empty shoe")
		}
	})
}
