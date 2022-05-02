package http

import (
	"cow/internal"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

// Api http api
type Api struct {
	scoresStorage internal.ScoresStorage
	logger        *zap.Logger
	resp          *response
}

// NewApi inits and returns Api struct
func NewApi(scoresStorage internal.ScoresStorage, logger *zap.Logger, resp *response) *Api {
	return &Api{
		scoresStorage: scoresStorage,
		logger:        logger,
		resp:          resp,
	}
}

// Ping func to ping api
func (a *Api) Ping(w http.ResponseWriter, req *http.Request) {
	a.resp.createResponse(w, req, "Pong")
}

// SetScore saves score in storage
func (a *Api) SetScore(w http.ResponseWriter, req *http.Request) {
	gameId := req.URL.Query().Get("gameId")
	if gameId == "" {
		a.resp.createBadRequestResponse(w, req, errors.New("gameId is required"))
		a.logger.Info("gameId is required", zap.Any("query", req.URL.Query()))
		return
	}

	name := req.URL.Query().Get("name")
	if name == "" {
		a.resp.createBadRequestResponse(w, req, errors.New("name is required"))
		a.logger.Info("name is required", zap.Any("query", req.URL.Query()))
		return
	}

	score := req.URL.Query().Get("score")
	if score == "" {
		a.resp.createBadRequestResponse(w, req, errors.New("score is required"))
		a.logger.Info("score is required", zap.Any("query", req.URL.Query()))
		return
	}

	scoreInt, err := strconv.ParseInt(score, 10, 8)
	if err != nil {
		a.resp.createInternalErrorResponse(w, req, fmt.Errorf("can't convert score: %w", err))
		a.logger.Error("can't convert score", zap.Any("query", req.URL.Query()), zap.Error(err))
		return
	}

	err = a.scoresStorage.Replace(&internal.Score{
		GameId:    gameId,
		Name:      name,
		Score:     int8(scoreInt),
		ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
	})
	if err != nil {
		a.resp.createInternalErrorResponse(w, req, fmt.Errorf("can't replace: %w", err))
		a.logger.Error("can't replace score", zap.Error(err))
		return
	}

	a.resp.createResponse(w, req, "Ok")
}

// FindScore returns score from a storage
func (a *Api) FindScore(w http.ResponseWriter, req *http.Request) {
	gameId := req.URL.Query().Get("gameId")
	if gameId == "" {
		a.resp.createBadRequestResponse(w, req, errors.New("gameId is required"))
		a.logger.Info("gameId is required", zap.Any("query", req.URL.Query()))
		return
	}

	name := req.URL.Query().Get("name")
	if name == "" {
		a.resp.createBadRequestResponse(w, req, errors.New("name is required"))
		a.logger.Info("name is required", zap.Any("query", req.URL.Query()))
		return
	}

	score, err := a.scoresStorage.Find(gameId, name)
	if err != nil {
		a.resp.createInternalErrorResponse(w, req, fmt.Errorf("error find score: %w", err))
		a.logger.Error("can't find score", zap.Error(err))
	}

	a.resp.createResponse(w, req, score)
}
