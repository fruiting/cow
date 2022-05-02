package tarantool

import (
	"cow/internal"
	"fmt"
	"github.com/tarantool/go-tarantool"
)

// ScoresStorage tarantool storage of a scores
type ScoresStorage struct {
	client tarantool.Connector
}

// NewScoresStorage inits and returns ScoresStorage struct
func NewScoresStorage(client tarantool.Connector) *ScoresStorage {
	return &ScoresStorage{client: client}
}

// Replace saves scores
func (s *ScoresStorage) Replace(score *internal.Score) error {
	_, err := s.client.Call17(
		"API.scores.set.v1",
		[]interface{}{score.GameId, score.Name, score.Score, score.ExpiresAt},
	)
	if err != nil {
		return fmt.Errorf("can't save score: %w", err)
	}

	return nil
}

// Find returns scores
func (s *ScoresStorage) Find(gameId, name string) (*internal.Score, error) {
	var scores [][]*internal.Score
	err := s.client.Call17Typed(
		"API.scores.get.v1",
		[]interface{}{gameId, name},
		&scores,
	)
	if err != nil {
		return nil, fmt.Errorf("can't find score: %w", err)
	}

	if len(scores) == 0 {
		return nil, nil
	}

	return scores[0][0], nil
}
