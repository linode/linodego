package integration

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSecurityQuestions_List(t *testing.T) {
	warnSensitiveTest(t)

	client, teardown := createTestClient(t, "fixtures/TestSecurityQuestions_List")
	defer teardown()

	questions, err := client.SecurityQuestionsList(context.Background())
	require.NoError(t, err, "Error getting security questions, expected struct")

	require.NotEmpty(t, questions.SecurityQuestions, "Expected to see security questions returned")

	require.Equal(
		t,
		"What was the name of your first pet?",
		questions.SecurityQuestions[0].Question,
		"Expected question 'What was the name of your first pet?'",
	)
}
