package integration

import (
	"context"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego"
)

func TestSecurityQuestions_List(t *testing.T) {
	client := createMockClient(t)

	desiredResponse := linodego.SecurityQuestionsListResponse{
		SecurityQuestions: []linodego.SecurityQuestion{
			{
				ID:       1,
				Question: "Really cool question",
				Response: "uhhhh",
			},
		},
	}

	httpmock.RegisterRegexpResponder("GET", mockRequestURL(t, "/profile/security-questions"),
		httpmock.NewJsonResponderOrPanic(200, &desiredResponse))

	questions, err := client.SecurityQuestionsList(context.Background())
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(*questions, desiredResponse) {
		t.Fatalf("actual response does not equal desired response: %s", cmp.Diff(questions, desiredResponse))
	}
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
