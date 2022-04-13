package internal

//go:generate mockgen -destination=mock/scores_storage.go -package=mock -source=scores.go ScoresStorage

type ScoresStorage interface {
	// Replace saves scores
	Replace(score *Score) error
}

// Score player score
type Score struct {
	GameId    string // 1
	Name      string // 2
	Score     int8   // 3
	ExpiresAt int64  // 4
}
