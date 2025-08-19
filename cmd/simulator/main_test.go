package main

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/adequatica/punto-banco-golango/internal/simulator"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
)

func TestInitialModel(t *testing.T) {
	s := spinner.New()
	s.Spinner = spinner.Dot

	ti := textinput.New()
	ti.Placeholder = fmt.Sprintf("%d", defaultNumberOfSimulations)

	expectedModel := model{
		stateUI:         stateSelectStrategy,
		cursor:          0,
		strategyOptions: simulator.GetStrategyOptions(),
		textInput:       ti,
		numSimulations:  0,
		keys:            defaultKeys,
		help:            help.New(),
		spinner:         s,
	}

	actualModel := InitialModel()

	// Compare stateUI
	if actualModel.stateUI != expectedModel.stateUI {
		t.Errorf("stateUI mismatch: got %v, want %v", actualModel.stateUI, expectedModel.stateUI)
	}

	// Compare cursor
	if actualModel.cursor != expectedModel.cursor {
		t.Errorf("cursor mismatch: got %d, want %d", actualModel.cursor, expectedModel.cursor)
	}

	// Compare strategy options
	if !reflect.DeepEqual(actualModel.strategyOptions, expectedModel.strategyOptions) {
		t.Errorf("strategyOptions mismatch: got %v, want %v", actualModel.strategyOptions, expectedModel.strategyOptions)
	}

	// Compare text input properties
	if actualModel.textInput.Placeholder != expectedModel.textInput.Placeholder {
		t.Errorf("textInput.Placeholder mismatch: got %v, want %v", actualModel.textInput.Placeholder, expectedModel.textInput.Placeholder)
	}

	// Compare number of simulations
	if actualModel.numSimulations != expectedModel.numSimulations {
		t.Errorf("numSimulations mismatch: got %d, want %d", actualModel.numSimulations, expectedModel.numSimulations)
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
