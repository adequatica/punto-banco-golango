package main

import (
	"reflect"
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
	"github.com/adequatica/punto-banco-golango/internal/statistics"
	"github.com/charmbracelet/bubbles/help"
)

func TestInitialModel(t *testing.T) {
	expectedModel := model{
		stateUI:           stateIsBetting,
		stateGame:         puntobanco.GetNewGameResultState(),
		statistics:        statistics.NewSessionStatistics(),
		showStatistics:    false,
		cursor:            0,
		bettingOptions:    puntobanco.GetBettingOptions(),
		afterRoundOptions: defaultAfterRoundOptions,
		selectedOption:    "",
		keys:              defaultKeys,
		help:              help.New(),
		// progress:          progress.New(),
		progressPercent: 0.0,
	}

	actualModel := initialModel()

	// Compare stateUI
	if actualModel.stateUI != expectedModel.stateUI {
		t.Errorf("stateUI mismatch: got %v, got %v", actualModel.stateUI, expectedModel.stateUI)
	}

	// Compare stateGame
	if actualModel.stateGame.GetResult() != expectedModel.stateGame.GetResult() {
		t.Errorf("stateGame.Result mismatch: got %v, got %v", actualModel.stateGame.GetResult(), expectedModel.stateGame.GetResult())
	}

	// Compare statistics
	if !reflect.DeepEqual(actualModel.statistics, expectedModel.statistics) {
		t.Errorf("statistics mismatch: got %v, got %v", actualModel.statistics, expectedModel.statistics)
	}

	// Compare statistics show
	if actualModel.showStatistics != expectedModel.showStatistics {
		t.Errorf("statistics show: got %v, got %v", actualModel.statistics, expectedModel.statistics)
	}

	// Compare cursor
	if actualModel.cursor != expectedModel.cursor {
		t.Errorf("cursor mismatch: got %d, got %d", actualModel.cursor, expectedModel.cursor)
	}

	// Compare number of betting options
	if len(actualModel.bettingOptions) != len(expectedModel.bettingOptions) {
		t.Errorf("bettingOptions length mismatch: got %d, got %d", len(actualModel.bettingOptions), len(expectedModel.bettingOptions))
	}

	// Compare after round options
	if !reflect.DeepEqual(actualModel.afterRoundOptions, expectedModel.afterRoundOptions) {
		t.Errorf("afterRoundOptions mismatch: got %v, got %v", actualModel.afterRoundOptions, expectedModel.afterRoundOptions)
	}

	// Compare selected option
	if actualModel.selectedOption != expectedModel.selectedOption {
		t.Errorf("selectedOption mismatch: got %v, got %v", actualModel.selectedOption, expectedModel.selectedOption)
	}

	// Compare keys
	if !reflect.DeepEqual(actualModel.keys, expectedModel.keys) {
		t.Errorf("keys mismatch: got %v, got %v", actualModel.keys, expectedModel.keys)
	}

	// Compare help
	if !reflect.DeepEqual(actualModel.help, expectedModel.help) {
		t.Errorf("help mismatch: got %v, got %v", actualModel.help, expectedModel.help)
	}

	// Compare progress percent
	if actualModel.progressPercent != expectedModel.progressPercent {
		t.Errorf("progressPercent mismatch: got %v, got %v", actualModel.progressPercent, expectedModel.progressPercent)
	}
}
