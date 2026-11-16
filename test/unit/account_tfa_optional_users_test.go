package unit

import (
	"context"
	"maps"
	"net/http"
	"strconv"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccountTFAOptionalUsers_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_tfa_optional_users_list")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	t.Cleanup(func() {
		base.TearDown(t)
	})

	base.MockGet("account/tfa-optional-users", fixtureData)

	users, err := base.Client.ListTFAOptionalUsers(context.Background(), nil)
	require.NoError(t, err)
	assert.Equal(t, []linodego.TFAOptionalUser{
		{Username: "user1", URL: "/account/users/user1"},
		{Username: "user2", URL: "/account/users/user2"},
	}, users)
}

func TestAccountTFAOptionalUsers_Update(t *testing.T) {
	for _, testCase := range []struct {
		name      string
		usernames []string
		fixture   string
	}{
		{name: "replace", usernames: []string{"user1", "user2"}, fixture: "account_tfa_optional_users_update"},
		{name: "clear", usernames: []string{}, fixture: "account_tfa_optional_users_update_empty"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			fixtureData, err := fixtures.GetFixture(testCase.fixture)
			require.NoError(t, err)

			var base ClientBaseCase
			base.SetUp(t)
			t.Cleanup(func() {
				base.TearDown(t)
			})

			opts := linodego.TFAOptionalUsersUpdateOptions{Usernames: testCase.usernames}
			httpmock.RegisterResponder("PUT", base.BaseURL+"account/tfa-optional-users", mockRequestBodyValidate(t, opts, fixtureData))

			response, err := base.Client.UpdateTFAOptionalUsers(context.Background(), opts)
			require.NoError(t, err)
			assert.Equal(t, len(testCase.usernames), response.TFAOptionalUsersCount)
			assert.Equal(t, "/account/tfa-optional-users", response.URL)
		})
	}
}

func TestAccountTFAOptionalUsers_ListEmpty(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_tfa_optional_users_list_empty")
	require.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	t.Cleanup(func() {
		base.TearDown(t)
	})

	base.MockGet("account/tfa-optional-users", fixtureData)

	users, err := base.Client.ListTFAOptionalUsers(context.Background(), nil)
	require.NoError(t, err)
	assert.Empty(t, users)
}

func TestAccountTFAOptionalUsers_ListPagination(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("account_tfa_optional_users_list")
	require.NoError(t, err)
	fixturePage := fixtureData.(map[string]interface{})
	fixtureUsers := fixturePage["data"].([]interface{})

	for _, testCase := range []struct {
		name          string
		page          int
		expectedUsers []linodego.TFAOptionalUser
		expectedPages []int
	}{
		{
			name: "all_pages",
			expectedUsers: []linodego.TFAOptionalUser{
				{Username: "user1", URL: "/account/users/user1"},
				{Username: "user2", URL: "/account/users/user2"},
			},
			expectedPages: []int{1, 2},
		},
		{
			name:          "single_page",
			page:          2,
			expectedUsers: []linodego.TFAOptionalUser{{Username: "user2", URL: "/account/users/user2"}},
			expectedPages: []int{2},
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var base ClientBaseCase
			base.SetUp(t)
			t.Cleanup(func() {
				base.TearDown(t)
			})

			var requestedPages []int
			httpmock.RegisterResponder("GET", base.BaseURL+"account/tfa-optional-users", func(req *http.Request) (*http.Response, error) {
				page, err := strconv.Atoi(req.URL.Query().Get("page"))
				require.NoError(t, err)
				require.Contains(t, []int{1, 2}, page)
				assert.Equal(t, "1", req.URL.Query().Get("page_size"))
				requestedPages = append(requestedPages, page)

				response := maps.Clone(fixturePage)
				response["data"] = fixtureUsers[page-1 : page]
				response["page"] = page
				response["pages"] = 2
				return httpmock.NewJsonResponse(http.StatusOK, response)
			})

			opts := linodego.NewListOptions(testCase.page, "")
			opts.PageSize = 1
			users, err := base.Client.ListTFAOptionalUsers(context.Background(), opts)
			require.NoError(t, err)
			assert.Equal(t, testCase.expectedUsers, users)
			assert.Equal(t, testCase.expectedPages, requestedPages)
			assert.Equal(t, 2, opts.Page)
			assert.Equal(t, 2, opts.Pages)
			assert.Equal(t, 2, opts.Results)
		})
	}
}

func TestAccountTFAOptionalUsers_Errors(t *testing.T) {
	for _, testCase := range []struct {
		name   string
		method string
		code   int
		reason string
	}{
		{name: "list_forbidden", method: "GET", code: http.StatusForbidden, reason: "Permission denied"},
		{name: "update_forbidden", method: "PUT", code: http.StatusForbidden, reason: "Permission denied"},
		{name: "invalid_username", method: "PUT", code: http.StatusBadRequest, reason: "Invalid username"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			var base ClientBaseCase
			base.SetUp(t)
			t.Cleanup(func() {
				base.TearDown(t)
			})

			httpmock.RegisterResponder(testCase.method, base.BaseURL+"account/tfa-optional-users",
				httpmock.NewJsonResponderOrPanic(testCase.code, linodego.APIError{
					Errors: []linodego.APIErrorReason{{Reason: testCase.reason}},
				}))

			var err error
			if testCase.method == "GET" {
				_, err = base.Client.ListTFAOptionalUsers(context.Background(), nil)
			} else {
				_, err = base.Client.UpdateTFAOptionalUsers(context.Background(), linodego.TFAOptionalUsersUpdateOptions{
					Usernames: []string{"invalid-user"},
				})
			}

			var apiError *linodego.Error
			require.ErrorAs(t, err, &apiError)
			assert.Equal(t, testCase.code, apiError.Code)
			assert.Contains(t, apiError.Message, testCase.reason)
		})
	}
}
