package notifyprojectowner

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
	"go.uber.org/mock/gomock"

	servicemock "github.com/turao/temporal-study/mocks/src/service"
)

func TestExecute(t *testing.T) {
	const (
		mockEntityID       = "entity-id"
		mockNotificationID = "notification-id"
	)

	tests := map[string]struct {
		Context                 context.Context
		Request                 Request
		ExpectedResponse        *Response
		ExpectedError           error
		MockNotificationService func(t *testing.T) service.NotificationService
	}{
		"success": {
			Context: context.Background(),
			Request: Request{
				Request: &api.NotifyRequest{
					EntityID: mockEntityID,
				},
			},
			ExpectedResponse: &Response{
				Response: &api.NotifyResponse{
					NotificationID: mockNotificationID,
				},
			},
			ExpectedError: nil,
			MockNotificationService: func(t *testing.T) service.NotificationService {
				controller := gomock.NewController(t)
				mock := servicemock.NewMockNotificationService(controller)
				mock.EXPECT().
					Notify(gomock.Any(), gomock.Any()).
					Return(
						&api.NotifyResponse{
							NotificationID: mockNotificationID,
						},
						nil,
					)
				return mock
			},
		},
		"error - notification service error": {
			Context: context.Background(),
			Request: Request{
				Request: &api.NotifyRequest{
					EntityID: mockEntityID,
				},
			},
			ExpectedResponse: nil,
			ExpectedError:    assert.AnError,
			MockNotificationService: func(t *testing.T) service.NotificationService {
				controller := gomock.NewController(t)
				mock := servicemock.NewMockNotificationService(controller)
				mock.EXPECT().
					Notify(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
				return mock
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			activity := Activity{
				NotificationService: test.MockNotificationService(t),
			}

			res, err := activity.Execute(test.Context, test.Request)
			assert.Equal(t, test.ExpectedError, err)
			assert.Equal(t, test.ExpectedResponse, res)
		})
	}
}
