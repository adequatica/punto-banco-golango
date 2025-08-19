package main

import (
	"reflect"
	"testing"

	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
	"github.com/adequatica/punto-banco-golango/internal/statistics"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
)

func TestInitialModel(t *testing.T) {
	s := spinner.New()
	s.Spinner = spinner.Dot

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
		spinner:           s,
	}

	actualModel := initialModel()

	// Compare stateUI
	if actualModel.stateUI != expectedModel.stateUI {
		t.Errorf("stateUI mismatch: got %v, want %v", actualModel.stateUI, expectedModel.stateUI)
	}

	// Compare stateGame
	if actualModel.stateGame.GetResult() != expectedModel.stateGame.GetResult() {
		t.Errorf("stateGame.Result mismatch: got %v, want %v", actualModel.stateGame.GetResult(), expectedModel.stateGame.GetResult())
	}

	// Compare statistics
	if !reflect.DeepEqual(actualModel.statistics, expectedModel.statistics) {
		t.Errorf("statistics mismatch: got %v, want %v", actualModel.statistics, expectedModel.statistics)
	}

	// Compare statistics show
	if actualModel.showStatistics != expectedModel.showStatistics {
		t.Errorf("statistics show: got %v, want %v", actualModel.statistics, expectedModel.statistics)
	}

	// Compare cursor
	if actualModel.cursor != expectedModel.cursor {
		t.Errorf("cursor mismatch: got %d, want %d", actualModel.cursor, expectedModel.cursor)
	}

	// Compare number of betting options
	if len(actualModel.bettingOptions) != len(expectedModel.bettingOptions) {
		t.Errorf("bettingOptions length mismatch: got %d, want %d", len(actualModel.bettingOptions), len(expectedModel.bettingOptions))
	}

	// Compare after round options
	if !reflect.DeepEqual(actualModel.afterRoundOptions, expectedModel.afterRoundOptions) {
		t.Errorf("afterRoundOptions mismatch: got %v, want %v", actualModel.afterRoundOptions, expectedModel.afterRoundOptions)
	}

	// Compare selected option
	if actualModel.selectedOption != expectedModel.selectedOption {
		t.Errorf("selectedOption mismatch: got %v, want %v", actualModel.selectedOption, expectedModel.selectedOption)
	}

	// Compare keys
	if !reflect.DeepEqual(actualModel.keys, expectedModel.keys) {
		t.Errorf("keys mismatch: got %v, want %v", actualModel.keys, expectedModel.keys)
	}

	// Compare help
	if !reflect.DeepEqual(actualModel.help, expectedModel.help) {
		t.Errorf("help mismatch: got %v, want %v", actualModel.help, expectedModel.help)
	}

	// Compare spinner
	if reflect.DeepEqual(actualModel.spinner, expectedModel.spinner) {
		t.Errorf("spinner should have different instance")
	}
}
