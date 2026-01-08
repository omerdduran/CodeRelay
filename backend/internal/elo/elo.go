package elo

import (
	"math"
)

const (
	// KFactor determines how much ratings change after each race
	// Higher K = more volatile ratings
	KFactor = 32

	// DefaultRating is the starting ELO rating for new users
	DefaultRating = 1200
)

// Player represents a player in a race with their rating and rank.
type Player struct {
	UserID    int64
	Rating    int
	Rank      int // 1 = 1st place, 2 = 2nd place, etc.
	NewRating int
}

// CalculateRatings calculates new ELO ratings for all players in a race.
// Players should be sorted by rank (1st place first).
func CalculateRatings(players []Player) []Player {
	if len(players) == 0 {
		return players
	}

	// For single player, no rating change
	if len(players) == 1 {
		players[0].NewRating = players[0].Rating
		return players
	}

	// Calculate expected scores and new ratings
	for i := range players {
		expectedScore := calculateExpectedScore(players[i].Rating, players, i)
		actualScore := calculateActualScore(players[i].Rank, len(players))
		
		ratingChange := int(math.Round(KFactor * (actualScore - expectedScore)))
		players[i].NewRating = players[i].Rating + ratingChange
		
		// Ensure rating doesn't go below a minimum threshold
		if players[i].NewRating < 100 {
			players[i].NewRating = 100
		}
	}

	return players
}

// calculateExpectedScore calculates the expected score for a player
// based on their rating relative to all other players.
func calculateExpectedScore(playerRating int, allPlayers []Player, playerIndex int) float64 {
	var expectedScore float64
	
	for i, opponent := range allPlayers {
		if i == playerIndex {
			continue
		}
		
		expectedScore += expectedWinProbability(playerRating, opponent.Rating)
	}
	
	// Normalize to 0-1 range
	if len(allPlayers) > 1 {
		expectedScore = expectedScore / float64(len(allPlayers)-1)
	}
	
	return expectedScore
}

// expectedWinProbability calculates the expected probability of player A
// winning against player B based on their ratings.
func expectedWinProbability(ratingA, ratingB int) float64 {
	return 1.0 / (1.0 + math.Pow(10.0, float64(ratingB-ratingA)/400.0))
}

// calculateActualScore converts a rank to a normalized score.
// 1st place = 1.0, last place = 0.0, middle ranks are interpolated.
func calculateActualScore(rank, totalPlayers int) float64 {
	if totalPlayers <= 1 {
		return 0.5
	}
	
	// Linear interpolation: 1st place = 1.0, last place = 0.0
	return float64(totalPlayers-rank) / float64(totalPlayers-1)
}

// CalculateRatingChange is a helper function to calculate rating change
// for a single player against another player (head-to-head).
func CalculateRatingChange(playerRating, opponentRating int, won bool) int {
	expected := expectedWinProbability(playerRating, opponentRating)
	
	var actual float64
	if won {
		actual = 1.0
	} else {
		actual = 0.0
	}
	
	return int(math.Round(KFactor * (actual - expected)))
}
