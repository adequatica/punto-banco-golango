package simulator

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/adequatica/punto-banco-golango/internal/deck"
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

// Data structures for saving simulation data
type SimulationData struct {
	Strategy            string    `json:"strategy"`
	DecksInShoe         int       `json:"decksInShoe"`
	StartingBankroll    float64   `json:"startingBankroll"`
	StandardBet         float64   `json:"standardBet"`
	NumberOfSimulations int       `json:"numberOfSimulations"`
	Games               [][]Hands `json:"games"`
}

type Hands struct {
	GameID     int      `json:"gameId"`
	HandID     int      `json:"handId"`
	ShoeNumber int      `json:"shoeNumber"`
	PuntoHand  []string `json:"puntoHand"`
	BankoHand  []string `json:"bankoHand"`
	PuntoTotal int      `json:"puntoTotal"`
	BankoTotal int      `json:"bankoTotal"`
	Result     string   `json:"result"`
	Bet        BetData  `json:"bet"`
}

type BetData struct {
	BetOn         string  `json:"betOn"`
	IsWin         bool    `json:"isWin"`
	BetAmount     float64 `json:"betAmount"`
	Payout        float64 `json:"payout"`
	FinalBankroll float64 `json:"finalBankroll"`
}

func FormatCard(card *deck.Card) string {
	if card == nil {
		return ""
	}

	// Get suit abbreviation
	suitAbbr := ""
	switch strings.ToLower(card.Suit) {
	case "spades":
		suitAbbr = "S"
	case "clubs":
		suitAbbr = "C"
	case "hearts":
		suitAbbr = "H"
	case "diamonds":
		suitAbbr = "D"
	default:
		suitAbbr = "?"
	}

	return card.Card + suitAbbr
}

func FormatBetAndResultType(betType puntobanco.BetType) string {
	switch betType {
	case puntobanco.PuntoPlayer:
		return "punto"
	case puntobanco.BancoBanker:
		return "banko"
	case puntobanco.EgaliteTie:
		return "egalite"
	default:
		return "punto"
	}
}

type DataCollector struct {
	data            *SimulationData
	currentGameID   int
	currentHandID   int
	currentShoeNum  int
	previousShoeLen int
}

// Create a new data collector for simulation data
// Collecting data for a large number of simulations (over 1000) may cause memory exhaustion
func NewDataCollector(
	strategy StrategyType,
	decksInShoe int,
	startingBankroll float64,
	standardBet float64,
	numberOfSimulations int) *DataCollector {
	return &DataCollector{
		data: &SimulationData{
			Strategy:            string(strategy),
			DecksInShoe:         decksInShoe,
			StartingBankroll:    startingBankroll,
			StandardBet:         standardBet,
			NumberOfSimulations: numberOfSimulations,
			Games:               make([][]Hands, 0, numberOfSimulations),
		},
		currentGameID:   0,
		currentHandID:   0,
		currentShoeNum:  0,
		previousShoeLen: 0,
	}
}

// Initialize a new game in the data collector
func (dc *DataCollector) StartNewGame(shoe []deck.Card) {
	dc.currentGameID++
	dc.currentHandID = 0
	dc.currentShoeNum = 1
	dc.previousShoeLen = len(shoe)
	dc.data.Games = append(dc.data.Games, make([]Hands, 0))
}

// Collect data for a single round
func (dc *DataCollector) CollectHandData(
	state *SimulatorState,
	gameResult *puntobanco.GameResultState,
	previousShoeLength int,
) {
	if dc.data == nil {
		return
	}

	// Track shoe number - increment when a new shoe is created
	if gameResult != nil {
		currentShoeLength := len(gameResult.RemainingShoe)
		// A new shoe is created when the remaining shoe has less than 8 cards
		// If previous shoe length was small (< 8) and current is large (>= 8), then a new shoe was created
		if previousShoeLength < 8 && currentShoeLength >= 8 {
			dc.currentShoeNum++
		}
		dc.previousShoeLen = currentShoeLength
	}

	// Increment hand ID each round
	dc.currentHandID++

	// Get punto and banco hands
	var puntoHand []string
	var bankoHand []string
	var puntoTotal int
	var bankoTotal int
	var result string

	if gameResult != nil {
		if gameResult.PuntoState != nil {
			if gameResult.PuntoState.FirstCard != nil {
				puntoHand = append(puntoHand, FormatCard(gameResult.PuntoState.FirstCard))
			}
			if gameResult.PuntoState.SecondCard != nil {
				puntoHand = append(puntoHand, FormatCard(gameResult.PuntoState.SecondCard))
			}
			if gameResult.PuntoState.ThirdCard != nil {
				puntoHand = append(puntoHand, FormatCard(gameResult.PuntoState.ThirdCard))
			}
			puntoTotal = gameResult.PuntoState.Points
		}

		if gameResult.BancoState != nil {
			if gameResult.BancoState.FirstCard != nil {
				bankoHand = append(bankoHand, FormatCard(gameResult.BancoState.FirstCard))
			}
			if gameResult.BancoState.SecondCard != nil {
				bankoHand = append(bankoHand, FormatCard(gameResult.BancoState.SecondCard))
			}
			if gameResult.BancoState.ThirdCard != nil {
				bankoHand = append(bankoHand, FormatCard(gameResult.BancoState.ThirdCard))
			}
			bankoTotal = gameResult.BancoState.Points
		}

		if gameResult.Result != nil {
			result = FormatBetAndResultType(*gameResult.Result)
		}
	}

	// Get bet information
	var betOn string
	var isWin bool
	var betAmount float64
	var payout float64
	var finalBankroll float64

	if state != nil {
		betOn = FormatBetAndResultType(state.BettingOn)
		betAmount = state.BetAmount
		finalBankroll = state.CurrentBankroll

		// Determine if this hand was a win
		if gameResult != nil && gameResult.Result != nil && *gameResult.Result == state.BettingOn {
			isWin = true
			payout = CalculatePayout(state.BettingOn, state.BetAmount)
		} else {
			isWin = false
			payout = 0.0
		}
	}

	// Create hand data
	handData := Hands{
		GameID:     dc.currentGameID,
		HandID:     dc.currentHandID,
		ShoeNumber: dc.currentShoeNum,
		PuntoHand:  puntoHand,
		BankoHand:  bankoHand,
		PuntoTotal: puntoTotal,
		BankoTotal: bankoTotal,
		Result:     result,
		Bet: BetData{
			BetOn:         betOn,
			IsWin:         isWin,
			BetAmount:     betAmount,
			Payout:        payout,
			FinalBankroll: finalBankroll,
		},
	}

	// This should not happen in normal flow, but handle gracefully
	// Initialize a new game if somehow StartNewGame was missed
	if len(dc.data.Games) == 0 {
		dc.data.Games = append(dc.data.Games, make([]Hands, 0))
		dc.currentGameID = 1
	}
	gameIndex := len(dc.data.Games) - 1
	dc.data.Games[gameIndex] = append(dc.data.Games[gameIndex], handData)
}

// Return the collected simulation data
func (dc *DataCollector) GetSimulationData() *SimulationData {
	return dc.data
}

func CreateSimulationDataFilename(strategy string, numberOfSimulations int, useGzip bool) string {
	// Sanitize strategy name for filename (remove special characters)
	replacer := strings.NewReplacer(
		" ", "_",
		"(", "",
		")", "",
		"'", "",
		"é", "e",
		"É", "E",
	)
	sanitizedStrategy := replacer.Replace(strategy)

	numberOfSimulationsStr := fmt.Sprintf("%d", numberOfSimulations)

	// Format date and time: YYYY-MM-DD_HH.MM.SS
	now := time.Now()
	dateTimeStr := now.Format("2006-01-02_15.04.05")

	// Filename format: strategy name + date + time.json (or .json.gz if compressed)
	fileExtension := ".json"
	if useGzip {
		fileExtension = ".json.gz"
	}

	return fmt.Sprintf("%s_%s_%s%s", sanitizedStrategy, numberOfSimulationsStr, dateTimeStr, fileExtension)
}

// Save simulation data to a JSON file
func SaveSimulationData(data *SimulationData) error {
	if data == nil {
		return fmt.Errorf("Simulation data is nil")
	}

	// Determine the need of gzip compression
	// 100 simulations can create a .json file larger than 43 MB with over 92K hands
	// 1000 simulations can create a .json file larger than 440 MB with over 930K hands
	// 10000 simulations can create a .json file larger than 4.4 GB with over 9.3M hands
	useGzip := data.NumberOfSimulations > 100

	filename := CreateSimulationDataFilename(data.Strategy, data.NumberOfSimulations, useGzip)

	// Create /datasets directory if it doesn't exist
	datasetsDir := "datasets"
	err := os.MkdirAll(datasetsDir, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create datasets directory: %w", err)
	}

	// Create full file path
	filepath := fmt.Sprintf("%s/%s", datasetsDir, filename)

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("Failed to marshal simulation data: %w", err)
	}

	// Write to file (with or without gzip compression)
	if useGzip {
		// Write compressed file
		file, err := os.Create(filepath)
		if err != nil {
			return fmt.Errorf("Failed to create gzip file: %w", err)
		}
		defer file.Close()

		gzipWriter := gzip.NewWriter(file)
		defer gzipWriter.Close()

		_, err = gzipWriter.Write(jsonData)
		if err != nil {
			return fmt.Errorf("Failed to write compressed data: %w", err)
		}
	} else {
		// Write uncompressed file
		err = os.WriteFile(filepath, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("Failed to write simulation data to file: %w", err)
		}
	}

	return nil
}
