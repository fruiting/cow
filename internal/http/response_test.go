package http

import (
	"bytes"
	"cow/internal/mock"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestWrite(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	serverMock := mock.NewMockResponseWriter(mockCtrl)

	type args struct {
		toMarshal interface{}
	}

	tests := []struct {
		name string
		args *args
		want func() (int, int)
	}{
		{
			name: "success",
			args: &args{
				toMarshal: 1,
			},
			want: func() (int, int) {
				serverMock.EXPECT().Write(gomock.Any()).Return(0, nil)
				return 0, 0
			},
		},
		{
			name: "marshal error",
			args: &args{
				toMarshal: map[string]interface{}{
					"foo": make(chan int),
				},
			},
			want: func() (int, int) {
				return 1, 0
			},
		},
		{
			name: "write error",
			args: &args{
				toMarshal: 1,
			},
			want: func() (int, int) {
				serverMock.EXPECT().Write(gomock.Any()).Return(0, errors.New("mock"))
				return 0, 1
			},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			logger, logs := initLogger()
			wantJsonLogsLen, wantWriteLogsLen := test.want()
			req, err := http.NewRequest("POST", "/v1/add", bytes.NewBuffer([]byte("")))
			assert.NoError(t, err)

			resp := NewResponse(logger)
			resp.write(serverMock, req, test.args.toMarshal)

			assert.Equal(t, wantJsonLogsLen, logs.FilterMessage("can't marshal http response").Len())
			assert.Equal(t, wantWriteLogsLen, logs.FilterMessage("can't write http response").Len())
		})
	}
}
