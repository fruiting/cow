package tarantool

import (
	"cow/internal"
	"cow/internal/mock"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoresStorage_Replace(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	tntClient := mock.NewMockConnector(mockCtrl)

	type args struct {
		score *internal.Score
	}

	tests := []struct {
		name string
		args *args
		want func() error
	}{
		{
			"succeess",
			&args{
				&internal.Score{
					"123",
					"Roman",
					50,
					1649674495,
				},
			},
			func() error {
				tntClient.EXPECT().Call17(gomock.Any(), gomock.Any()).Return(nil, nil)

				return nil
			},
		},
		{
			"error",
			&args{
				&internal.Score{
					"123",
					"Roman",
					50,
					1649674495,
				},
			},
			func() error {
				err := errors.New("mock")
				tntClient.EXPECT().Call17(gomock.Any(), gomock.Any()).Return(nil, err)

				return fmt.Errorf("can't save score: %w", err)
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			wantErr := test.want()

			storage := NewScoresStorage(tntClient)
			err := storage.Replace(test.args.score)

			assert.Equal(t, wantErr, err)
		})
	}
}
