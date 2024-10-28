package unit

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestSecurityQuestions_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("profile_security_question_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("profile/security-questions", fixtureData)

	securityQuestions, err := base.Client.SecurityQuestionsList(context.Background())
	assert.NoError(t, err)

	assert.Len(t, securityQuestions.SecurityQuestions, 1)
	assert.Equal(t, 1, securityQuestions.SecurityQuestions[0].ID)
	assert.Equal(t, "In what city were you born?", securityQuestions.SecurityQuestions[0].Question)
	assert.Equal(t, "Gotham City", securityQuestions.SecurityQuestions[0].Response)
}

func TestSecurityQuestions_Answer(t *testing.T) {
	client := createMockClient(t)

	requestData := linodego.SecurityQuestionsAnswerOptions{
		SecurityQuestions: []linodego.SecurityQuestionsAnswerQuestion{
			{
				QuestionID: 1,
				Response:   "cool",
			},
		},
	}

	httpmock.RegisterRegexpResponder("POST", mockRequestURL(t, "/profile/security-questions"),
		mockRequestBodyValidate(t, requestData, nil))

	if err := client.SecurityQuestionsAnswer(context.Background(), requestData); err != nil {
		t.Fatal(err)
	}
}
