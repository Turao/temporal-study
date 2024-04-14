package createproject

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/turao/temporal-study/src/temporal/activities"
	createprojectentity "github.com/turao/temporal-study/src/temporal/activities/create-project-entity"
	notifyprojectowner "github.com/turao/temporal-study/src/temporal/activities/notify-project-owner"
	upsertproject "github.com/turao/temporal-study/src/temporal/activities/upsert-project"
	"go.temporal.io/sdk/testsuite"
	"go.uber.org/mock/gomock"
)

func TestExecute(t *testing.T) {
	tests := map[string]struct {
		Request          Request
		ExpectedResponse *Response
		ExpectedError    error
		// mocks
	}{
		"success": {
			Request:          Request{},
			ExpectedResponse: &Response{},
			ExpectedError:    nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			suite := testsuite.WorkflowTestSuite{}
			env := suite.NewTestWorkflowEnvironment()

			// todo: this shit won't work because the environment does not know how to resolve a string into a workflow mock
			env.OnActivity(activities.ActivityNameCreateProjectEntity, gomock.Any()).Return(&createprojectentity.Response{}, nil)
			env.OnActivity(activities.ActivityNameUpsertProject, gomock.Any()).Return(&upsertproject.Response{}, nil)
			env.OnActivity(activities.ActivityNameNotifyProjectOwner, gomock.Any()).Return(&notifyprojectowner.Response{}, nil)

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
