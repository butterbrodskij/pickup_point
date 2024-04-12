package kafka

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/mer_marat/homework/internal/model"
)

func TestSendMessage(t *testing.T) {
	t.Parallel()
	type args struct {
		topic string
		msg   model.RequestMessage
	}
	type fields struct {
		producer func(m *Mockproducer)
	}
	testCases := []struct {
		name        string
		args        args
		fields      fields
		expectError bool
		wantError   string
	}{
		{
			name: "success",
			args: args{
				topic: "valid",
				msg:   model.RequestMessage{Request: httptest.NewRequest("", "/", strings.NewReader(""))},
			},
			fields: fields{
				producer: func(m *Mockproducer) {
					m.EXPECT().SendSyncMessage(gomock.Any()).Return(int32(0), int64(0), nil)
				},
			},
			expectError: false,
			wantError:   "",
		},
		{
			name: "empty request",
			args: args{
				topic: "valid",
				msg:   model.RequestMessage{},
			},
			fields: fields{
				producer: func(m *Mockproducer) {},
			},
			expectError: true,
			wantError:   "empty request",
		},
		{
			name: "empty body request",
			args: args{
				topic: "valid",
				msg:   model.RequestMessage{Request: &http.Request{}},
			},
			fields: fields{
				producer: func(m *Mockproducer) {},
			},
			expectError: true,
			wantError:   "empty body request",
		},
		{
			name: "error",
			args: args{
				topic: "valid",
				msg:   model.RequestMessage{Request: httptest.NewRequest("", "/", strings.NewReader(""))},
			},
			fields: fields{
				producer: func(m *Mockproducer) {
					m.EXPECT().SendSyncMessage(gomock.Any()).Return(int32(0), int64(0), assert.AnError)
				},
			},
			expectError: true,
			wantError:   "assert.AnError general error for testing",
		},
	}
	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := setUp(t, tt.args.topic)
			defer s.tearDown()
			tt.fields.producer(s.mockProd)

			err := s.sender.SendMessage(tt.args.msg)

			if tt.expectError {
				assert.EqualError(t, err, tt.wantError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
