package game

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sort"
)

type HighScore struct {
	Name  string `json:"name"`
	Score int    `json:"score"`
}

type HighScores []HighScore

const (
	maxHighScores = 10 // Keep top 10 scores
)

// GetHighScoreFilePath returns the path to the high scores file
func GetHighScoreFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".go-dinorun-scores")
}

// LoadHighScores loads high scores from file
func LoadHighScores() (HighScores, error) {
	filePath := GetHighScoreFilePath()
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(HighScores, 0), nil
		}
		return nil, err
	}

	// Unmarshal JSON
	var scores HighScores
	if err := json.Unmarshal(data, &scores); err != nil {
		return nil, errors.New("file corrupted")
	}

	return scores, nil
}

// SaveHighScores saves high scores to file
func SaveHighScores(scores HighScores) error {
	// Sort scores in descending order
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	// Keep only top scores
	if len(scores) > maxHighScores {
		scores = scores[:maxHighScores]
	}

	// Marshal to JSON
	data, err := json.Marshal(scores)
	if err != nil {
		return err
	}

	// Save to file
	return os.WriteFile(GetHighScoreFilePath(), data, 0600)
}

// IsHighScore checks if a score qualifies for the high score list
func IsHighScore(score int) bool {
	scores, err := LoadHighScores()
	if err != nil {
		// If there's an error loading scores, treat it as a high score
		return true
	}

	if len(scores) < maxHighScores {
		return true
	}

	// Check if score is higher than the lowest high score
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].Score > scores[j].Score
	})

	return score > scores[len(scores)-1].Score
}

// AddHighScore adds a new high score to the list
func AddHighScore(name string, score int) error {
	scores, err := LoadHighScores()
	if err != nil {
		// If file is corrupted or doesn't exist, start fresh
		scores = make(HighScores, 0)
	}

	scores = append(scores, HighScore{
		Name:  name,
		Score: score,
	})

	return SaveHighScores(scores)
}
