package elo

import (
	"testing"
)

func TestCalculateRatings_TwoPlayers(t *testing.T) {
	players := []Player{
		{UserID: 1, Rating: 1200, Rank: 1}, // Winner
		{UserID: 2, Rating: 1200, Rank: 2}, // Loser
	}

	result := CalculateRatings(players)

	// Winner should gain rating, loser should lose rating
	if result[0].NewRating <= 1200 {
		t.Errorf("Winner should gain rating, got %d", result[0].NewRating)
	}
	if result[1].NewRating >= 1200 {
		t.Errorf("Loser should lose rating, got %d", result[1].NewRating)
	}

	// Total rating change should be balanced
	totalChange := (result[0].NewRating - result[0].Rating) + (result[1].NewRating - result[1].Rating)
	if totalChange != 0 {
		t.Errorf("Total rating change should be 0, got %d", totalChange)
	}
}

func TestCalculateRatings_UpsetsGiveMorePoints(t *testing.T) {
	// Low rated player beats high rated player
	players := []Player{
		{UserID: 1, Rating: 1000, Rank: 1}, // Low rated winner
		{UserID: 2, Rating: 1400, Rank: 2}, // High rated loser
	}

	result := CalculateRatings(players)

	// Winner should gain more points in an upset
	winnerGain := result[0].NewRating - result[0].Rating
	loserLoss := result[1].Rating - result[1].NewRating

	if winnerGain <= 16 {
		t.Errorf("Upset winner should gain significant points, got %d", winnerGain)
	}
	if loserLoss <= 16 {
		t.Errorf("Upset loser should lose significant points, got %d", loserLoss)
	}
}

func TestCalculateRatings_ExpectedWinGivesFewerPoints(t *testing.T) {
	// High rated player beats low rated player (expected outcome)
	players := []Player{
		{UserID: 1, Rating: 1400, Rank: 1}, // High rated winner
		{UserID: 2, Rating: 1000, Rank: 2}, // Low rated loser
	}

	result := CalculateRatings(players)

	// Winner should gain fewer points for expected win
	winnerGain := result[0].NewRating - result[0].Rating

	if winnerGain >= 16 {
		t.Errorf("Expected winner should gain fewer points, got %d", winnerGain)
	}
}

func TestCalculateRatings_ThreePlayers(t *testing.T) {
	players := []Player{
		{UserID: 1, Rating: 1200, Rank: 1}, // 1st place
		{UserID: 2, Rating: 1200, Rank: 2}, // 2nd place
		{UserID: 3, Rating: 1200, Rank: 3}, // 3rd place
	}

	result := CalculateRatings(players)

	// 1st should gain, 3rd should lose, 2nd should be roughly neutral
	if result[0].NewRating <= 1200 {
		t.Errorf("1st place should gain rating")
	}
	if result[2].NewRating >= 1200 {
		t.Errorf("3rd place should lose rating")
	}
}

func TestCalculateRatings_SinglePlayer(t *testing.T) {
	players := []Player{
		{UserID: 1, Rating: 1200, Rank: 1},
	}

	result := CalculateRatings(players)

	// Single player should not change rating
	if result[0].NewRating != 1200 {
		t.Errorf("Single player rating should not change, got %d", result[0].NewRating)
	}
}

func TestCalculateRatings_MinimumRating(t *testing.T) {
	// Player with very low rating loses
	players := []Player{
		{UserID: 1, Rating: 1400, Rank: 1},
		{UserID: 2, Rating: 150, Rank: 2}, // Very low rating
	}

	result := CalculateRatings(players)

	// Rating should not go below minimum
	if result[1].NewRating < 100 {
		t.Errorf("Rating should not go below 100, got %d", result[1].NewRating)
	}
}

func TestExpectedWinProbability(t *testing.T) {
	tests := []struct {
		ratingA  int
		ratingB  int
		expected float64
		name     string
	}{
		{1200, 1200, 0.5, "equal ratings"},
		{1400, 1200, 0.76, "200 point advantage"},
		{1200, 1400, 0.24, "200 point disadvantage"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expectedWinProbability(tt.ratingA, tt.ratingB)
			// Allow 0.01 tolerance for floating point
			if result < tt.expected-0.01 || result > tt.expected+0.01 {
				t.Errorf("expectedWinProbability(%d, %d) = %f, want approximately %f",
					tt.ratingA, tt.ratingB, result, tt.expected)
			}
		})
	}
}

func TestCalculateActualScore(t *testing.T) {
	tests := []struct {
		rank         int
		totalPlayers int
		expected     float64
		name         string
	}{
		{1, 2, 1.0, "1st of 2"},
		{2, 2, 0.0, "2nd of 2"},
		{1, 3, 1.0, "1st of 3"},
		{2, 3, 0.5, "2nd of 3"},
		{3, 3, 0.0, "3rd of 3"},
		{1, 4, 1.0, "1st of 4"},
		{2, 4, 0.666, "2nd of 4"},
		{3, 4, 0.333, "3rd of 4"},
		{4, 4, 0.0, "4th of 4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateActualScore(tt.rank, tt.totalPlayers)
			// Allow 0.01 tolerance
			if result < tt.expected-0.01 || result > tt.expected+0.01 {
				t.Errorf("calculateActualScore(%d, %d) = %f, want approximately %f",
					tt.rank, tt.totalPlayers, result, tt.expected)
			}
		})
	}
}
