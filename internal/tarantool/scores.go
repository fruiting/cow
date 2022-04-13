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
