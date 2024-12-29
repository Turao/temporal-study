package createproject

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	mockservice "github.com/turao/temporal-study/mocks/src/service"
	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/service"
	"go.temporal.io/sdk/testsuite"
	"go.uber.org/mock/gomock"
)

func TestExecute(t *testing.T) {
	var (
		mockProjectName = "project-name"
		mockOwnerID     = uuid.Must(uuid.NewV4())
	)

	tests := map[string]struct {
		Context   context.Context
		Request   Request
		Assertion func(t *testing.T, response *Response, err error)
		// mocks
		MockProjectService func(t *testing.T) service.ProjectService
	}{
		"success": {
			Context: context.Background(),
			Request: Request{
				ProjectName: mockProjectName,
				OwnerID:     mockOwnerID.String(),
			},
			Assertion: func(t *testing.T, response *Response, err error) {
				assert.NotNil(t, response)
				assert.Equal(t, mockProjectName, response.Entity.Name)
				assert.Equal(t, mockOwnerID, response.Entity.OwnerID)
			},
			MockProjectService: func(t *testing.T) service.ProjectService {
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
				ProjectName: mockProjectName,
				OwnerID:     mockOwnerID.String(),
			},
			Assertion: func(t *testing.T, response *Response, err error) {
				assert.ErrorContains(t, err, assert.AnError.Error())
			},
			MockProjectService: func(t *testing.T) service.ProjectService {
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
			suite := testsuite.WorkflowTestSuite{}
			env := suite.NewTestActivityEnvironment()

			activity := Activity{
				ProjectService: test.MockProjectService(t),
			}
			env.RegisterActivity(activity.Execute)
			val, err := env.ExecuteActivity(activity.Execute, test.Request)

			var response *Response
			if val != nil {
				val.Get(&response)
			}

			test.Assertion(t, response, err)
		})
	}
}
