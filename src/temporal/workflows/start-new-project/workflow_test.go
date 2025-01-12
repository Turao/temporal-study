package startnewproject

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/turao/temporal-study/src/api"
	"github.com/turao/temporal-study/src/entity/project"
	"github.com/turao/temporal-study/src/temporal/activities"
	createproject "github.com/turao/temporal-study/src/temporal/activities/create-project"
	notifyprojectowner "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/testsuite"
)

func TestExecute(t *testing.T) {
	registerActivityStubs := func(env *testsuite.TestWorkflowEnvironment) {
		env.RegisterActivityWithOptions(
			func(ctx context.Context, req createproject.Request) (*createproject.Response, error) {
				return nil, nil
			},
			activity.RegisterOptions{
				Name: activities.ActivityNameCreateProject,
			},
		)

		env.RegisterActivityWithOptions(
			func(ctx context.Context, req notifyprojectowner.Request) (*notifyprojectowner.Response, error) {
				return nil, nil
			},
			activity.RegisterOptions{
				Name: activities.ActivityNameNotifyProjectOwner,
			},
		)
	}

	tests := map[string]struct {
		Request          Request
		ExpectedResponse *Response
		ExpectedError    error
		// mocks
		MockTestWorkflowEnvironment func(t *testing.T) *testsuite.TestWorkflowEnvironment
	}{
		"success": {
			Request: Request{
				ProjectName: "project-name",
				OwnerID:     "00000000-0000-0000-0000-000000000001",
			},
			ExpectedResponse: &Response{},
			ExpectedError:    nil,
			MockTestWorkflowEnvironment: func(t *testing.T) *testsuite.TestWorkflowEnvironment {
				suite := testsuite.WorkflowTestSuite{}
				env := suite.NewTestWorkflowEnvironment()

				registerActivityStubs(env)

				env.OnActivity(activities.ActivityNameCreateProject, mock.Anything, mock.Anything).Return(
					&createproject.Response{
						Entity: &project.Project{
							ID:      uuid.Must(uuid.NewV4()),
							Name:    "project-name",
							OwnerID: uuid.Must(uuid.NewV4()),
						},
					},
					nil,
				)

				env.OnActivity(activities.ActivityNameNotifyProjectOwner, mock.Anything, mock.Anything).Return(
					&notifyprojectowner.Response{
						Response: &api.NotifyResponse{
							NotificationID: uuid.Must(uuid.NewV4()).String(),
						},
					},
					nil,
				)

				return env
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			env := test.MockTestWorkflowEnvironment(t)

			workflow := Workflow{}
			env.ExecuteWorkflow(workflow.Execute, test.Request)
			err := env.GetWorkflowError()
			assert.Equal(t, test.ExpectedError, err)

			var res *Response
			err = env.GetWorkflowResult(&res)
			require.NoError(t, err)

			assert.Equal(t, test.ExpectedResponse, res)
		})
	}
}
