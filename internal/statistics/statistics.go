package statistics

import (
	puntobanco "github.com/adequatica/punto-banco-golango/internal/punto_banco"
)

type SessionStatistics struct {
	TotalRounds int
	PuntoWins   int
	BancoWins   int
	Ties        int
	UserWins    int
	UserBets    map[puntobanco.BetType]int
}

func NewSessionStatistics() SessionStatistics {
	return SessionStatistics{
		TotalRounds: 0,
		PuntoWins:   0,
		BancoWins:   0,
		Ties:        0,
		UserWins:    0,
		UserBets:    make(map[puntobanco.BetType]int),
	}
}

func (s *SessionStatistics) UpdateStatistics(gameResult puntobanco.BetType, userBet puntobanco.BetType) {
	switch gameResult {
	case puntobanco.PuntoPlayer:
		s.PuntoWins++
	case puntobanco.BancoBanker:
		s.BancoWins++
	case puntobanco.EgaliteTie:
		s.Ties++
	default:
		// Stop execution (don't update anything else) because of incorrect game result
		return
	}

	if gameResult == userBet {
		s.UserWins++
	}

	s.TotalRounds++
	s.UserBets[userBet]++
}

func (s *SessionStatistics) GetPuntoWinsPercentage() float64 {
	if s.TotalRounds == 0 {
		return 0.0
	}

	return float64(s.PuntoWins) / float64(s.TotalRounds) * 100.0
}

func (s *SessionStatistics) GetBancoWinsPercentage() float64 {
	if s.TotalRounds == 0 {
		return 0.0
	}

	return float64(s.BancoWins) / float64(s.TotalRounds) * 100.0
}

func (s *SessionStatistics) GetTiesPercentage() float64 {
	if s.TotalRounds == 0 {
		return 0.0
	}

	return float64(s.Ties) / float64(s.TotalRounds) * 100.0
}

func (s *SessionStatistics) GetUserWinsPercentage() float64 {
	if s.TotalRounds == 0 {
		return 0.0
	}

	return float64(s.UserWins) / float64(s.TotalRounds) * 100.0
}

func (s *SessionStatistics) GetUserBetsDistribution() map[puntobanco.BetType]float64 {
	distribution := make(map[puntobanco.BetType]float64)

	if s.TotalRounds == 0 {
		return distribution
	}

	for betType, count := range s.UserBets {
		distribution[betType] = float64(count) / float64(s.TotalRounds) * 100.0
	}

	return distribution
}

func (s *SessionStatistics) ResetStatistics() {
	*s = NewSessionStatistics()
}
