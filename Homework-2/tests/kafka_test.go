//go:build integration

package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	handler "gitlab.ozon.dev/mer_marat/homework/internal/api/handlers/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/middleware"
	"gitlab.ozon.dev/mer_marat/homework/internal/api/router"
	"gitlab.ozon.dev/mer_marat/homework/internal/pkg/kafka"
	"gitlab.ozon.dev/mer_marat/homework/internal/service/pickpoint"
	"gitlab.ozon.dev/mer_marat/homework/internal/storage/postgres"
	"gitlab.ozon.dev/mer_marat/homework/tests/dummy"
)

func TestKafka(t *testing.T) {
	var (
		ctx               = context.Background()
		handl             = dummy.NewHandler()
		consumer          = kafka.NewConsumerGroup(map[string]kafka.Handler{cfg.Kafka.Topic: handl}, cfg.Kafka.Topic)
		receiver, errRec  = kafka.NewReceiverGroup(consumer, cfg.Kafka.Brokers)
		producer, errProd = kafka.NewProducer(cfg.Kafka.Brokers)
		sender            = kafka.NewKafkaSender(producer, cfg.Kafka.Topic)
		middleware        = middleware.NewMiddleware(cfg, sender)
		router            = router.MakeRouter(handler.NewHandler(pickpoint.NewService(postgres.NewRepo(db.DB))), middleware, cfg)
	)
	require.NoError(t, errRec)
	require.NoError(t, errProd)
	err := receiver.Subscribe(cfg.Kafka.Topic)
	require.NoError(t, err)
	defer receiver.Close()
	type args struct {
		body     string
		login    string
		password string
		method   string
		url      string
	}
	testCases := []struct {
		name        string
		args        args
		expectError bool
		wantError   string
		wantBody    string
	}{
		{
			name: "success",
			args: args{
				body:     `{"name":"Chertanovo", "address":"Chertanovskaya street, 13", "contacts":"+7(999)888-77-66"}`,
				login:    "admin",
				password: "password",
				method:   "POST",
				url:      "/pickpoint",
			},
			expectError: false,
			wantError:   "",
			wantBody:    "New Request:\n\tMethod: POST\tPath: /pickpoint\tlogin: admin\tBody: {\"name\":\"Chertanovo\", \"address\":\"Chertanovskaya street, 13\", \"contacts\":\"+7(999)888-77-66\"}",
		},
		{
			name: "bad request handling",
			args: args{
				body:     `bad request`,
				login:    "admin",
				password: "password",
				method:   "POST",
				url:      "/pickpoint",
			},
			expectError: false,
			wantError:   "",
			wantBody:    "New Request:\n\tMethod: POST\tPath: /pickpoint\tlogin: admin\tBody: bad request",
		},
		{
			name: "another example",
			args: args{
				body:     `body`,
				login:    "pirate",
				password: "wrong-password",
				method:   "GET",
				url:      "/pickpoint/10",
			},
			expectError: false,
			wantError:   "",
			wantBody:    "New Request:\n\tMethod: GET\tPath: /pickpoint/10\tlogin: pirate\tBody: body",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			db.SetUp(t, "pickpoints")
			defer db.TearDown(t, "pickpoints")
			req, _ := http.NewRequestWithContext(ctx, tt.args.method, tt.args.url, strings.NewReader(tt.args.body))
			req.SetBasicAuth(tt.args.login, tt.args.password)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)
			<-handl.Wait() // wait for handler to handle message

			if tt.expectError {
				require.EqualError(t, err, tt.wantError)
			} else {
				require.NoError(t, handl.Err)
			}
			assert.Equal(t, tt.wantBody, handl.Result)
		})
	}
}
