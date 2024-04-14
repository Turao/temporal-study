package upsertproject

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	mockservice "github.com/turao/temporal-study/mocks/src/service"
	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
	"go.uber.org/mock/gomock"
)

func TestExecute(t *testing.T) {
	tests := map[string]struct {
		Context          context.Context
		Request          Request
		ExpectedResponse *Response
		ExpectedError    error
		// mocks
		ProjectService func(t *testing.T) service.ProjectService
	}{
		"success": {
			Context: context.Background(),
			Request: Request{
				Request: &api.UpsertProjectRequest{},
			},
			ExpectedResponse: &Response{
				Response: &api.UpsertProjectResponse{},
			},
			ExpectedError: nil,
			ProjectService: func(t *testing.T) service.ProjectService {
				controller := gomock.NewController(t)
				mock := mockservice.NewMockProjectService(controller)
				mock.EXPECT().
					UpsertProject(gomock.Any(), gomock.Any()).
					Return(&api.UpsertProjectResponse{}, nil)
				return mock
			},
		},
		"error - upsert project failure": {
			Context: context.Background(),
			Request: Request{
				Request: &api.UpsertProjectRequest{},
			},
			ExpectedResponse: nil,
			ExpectedError:    assert.AnError,
			ProjectService: func(t *testing.T) service.ProjectService {
				controller := gomock.NewController(t)
				mock := mockservice.NewMockProjectService(controller)
				mock.EXPECT().
					UpsertProject(gomock.Any(), gomock.Any()).
					Return(nil, assert.AnError)
				return mock
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			activity := Activity{
				ProjectService: test.ProjectService(t),
			}

			res, err := activity.Execute(test.Context, test.Request)
			assert.Equal(t, test.ExpectedError, err)
			assert.Equal(t, test.ExpectedResponse, res)
		})
	}
}
