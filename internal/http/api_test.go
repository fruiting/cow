package http

import (
	"bytes"
	"cow/internal/mock"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
	"net/http"
	"net/http/httptest"
	"testing"
)

// initLogger инициализирует логгер
func initLogger() (*zap.Logger, *observer.ObservedLogs) {
	core, logs := observer.New(zap.InfoLevel)
	return zap.New(core), logs
}

func TestApi_SetScore(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	scoresStorage := mock.NewMockScoresStorage(mockCtrl)

	type args struct {
		params string
	}

	tests := []struct {
		name string
		args *args
		want func() (int, int, int, int, int)
	}{
		{
			"success",
			&args{
				"gameId=123&name=Roman&score=50&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				scoresStorage.EXPECT().Replace(gomock.Any()).Return(nil)
				return 0, 0, 0, 0, 0
			},
		},
		{
			"game_id is required",
			&args{
				"name=Roman&score=50&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				return 1, 0, 0, 0, 0
			},
		},
		{
			"name is required",
			&args{
				"gameId=123&score=50&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				return 0, 1, 0, 0, 0
			},
		},
		{
			"score is required",
			&args{
				"gameId=123&name=Roman&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				return 0, 0, 1, 0, 0
			},
		},
		{
			"can't convert score",
			&args{
				"gameId=123&name=Roman&score=mock&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				return 0, 0, 0, 1, 0
			},
		},
		{
			"can't replace",
			&args{
				"gameId=123&name=Roman&score=50&expiresAt=1649686891",
			},
			func() (int, int, int, int, int) {
				scoresStorage.EXPECT().Replace(gomock.Any()).Return(errors.New("mock"))
				return 0, 0, 0, 0, 1
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			logger, logs := initLogger()
			resp := NewResponse(logger)

			wantGameIdLogsLen, wantNameLogsLen, wantScoreLogsLen, wantConvertLogsLen, wantReplaceLogsLen := test.want()

			api := NewApi(scoresStorage, logger, resp)
			req, err := http.NewRequest(
				"POST",
				fmt.Sprintf("/v1/add?%s", test.args.params),
				bytes.NewBuffer([]byte("")),
			)
			assert.NoError(t, err)

			api.SetScore(httptest.NewRecorder(), req)

			assert.Equal(t, wantGameIdLogsLen, logs.FilterMessage("gameId is required").Len())
			assert.Equal(t, wantNameLogsLen, logs.FilterMessage("name is required").Len())
			assert.Equal(t, wantScoreLogsLen, logs.FilterMessage("score is required").Len())
			assert.Equal(t, wantConvertLogsLen, logs.FilterMessage("can't convert score").Len())
			assert.Equal(t, wantReplaceLogsLen, logs.FilterMessage("can't replace score").Len())
		})
	}
}
