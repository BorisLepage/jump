package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"src/facture/internal/receipts"
	"src/helper_api"
)

func TestAPI_GetUser(t *testing.T) {
	for name, c := range map[string]struct {
		userID         string
		fixtures       TestFixtures
		expectedStatus int
		expectedError  string
	}{
		"nominal case user exists": {
			userID: "1",
			fixtures: TestFixtures{
				Users: []receipts.User{
					{
						ID:        1,
						FirstName: "Bob",
						LastName:  "Loco",
						Balance:   241817,
					},
				},
			},
			expectedStatus: http.StatusOK,
		},
	} {
		t.Run(name, func(t *testing.T) {
			server := setupTest(t, c.fixtures)
			defer server.Close()

			client := server.Client()

			url := fmt.Sprintf("%s/user/%s", server.URL, c.userID)
			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, c.expectedStatus, resp.StatusCode)

			var user receipts.User
			err = json.NewDecoder(resp.Body).Decode(&user)
			require.NoError(t, err)

			assert.Equal(t, c.userID, strconv.Itoa(user.ID))

		})
	}

	for name, c := range map[string]struct {
		userID         string
		fixtures       TestFixtures
		expectedStatus int
		expectedError  string
	}{
		"user not found": {
			userID: "999",
			fixtures: TestFixtures{
				Users: []receipts.User{
					{
						ID:        1,
						FirstName: "Bob",
						LastName:  "Loco",
						Balance:   241817,
					},
				},
			},
			expectedStatus: http.StatusNotFound,
			expectedError:  "user not found",
		},
		"invalid user ID": {
			userID:         "abc",
			fixtures:       TestFixtures{},
			expectedStatus: http.StatusBadRequest,
		},
	} {
		t.Run(name, func(t *testing.T) {
			server := setupTest(t, c.fixtures)
			defer server.Close()

			client := server.Client()

			url := fmt.Sprintf("%s/user/%s", server.URL, c.userID)
			resp, err := client.Get(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, c.expectedStatus, resp.StatusCode)

			var errorResp helper_api.ErrorResponse
			err = json.NewDecoder(resp.Body).Decode(&errorResp)
			require.NoError(t, err)

			if c.expectedError != "" {
				assert.Contains(t, errorResp.Message, c.expectedError)
			}
		})
	}
}

type TestFixtures struct {
	Users []receipts.User
}

var mockUsers map[int]receipts.User

func mockGetUser(userID int) (receipts.User, error) {
	user, exists := mockUsers[userID]
	if !exists {
		return receipts.User{}, receipts.UserNotFound
	}
	return user, nil
}

func setupTest(t *testing.T, fixtures TestFixtures) *httptest.Server {
	t.Helper()

	mockUsers = make(map[int]receipts.User)
	for _, user := range fixtures.Users {
		mockUsers[user.ID] = user
	}

	r := chi.NewRouter()

	r.Get("/user/{id}", func(w http.ResponseWriter, r *http.Request) {
		userID := chi.URLParam(r, "id")

		userIDInt, err := strconv.Atoi(userID)
		if err != nil {
			helper_api.SendErrorResponse(w, "bad_request", err.Error(), http.StatusBadRequest)
			return
		}

		user, err := mockGetUser(userIDInt)
		if errors.Is(err, receipts.UserNotFound) {
			helper_api.SendErrorResponse(w, "not_found", err.Error(), http.StatusNotFound)
			return
		}
		if err != nil {
			helper_api.SendErrorResponse(w, "internal_server_error", err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	})

	server := httptest.NewServer(r)

	return server
}
