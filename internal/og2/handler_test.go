package og2_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"hunter.io/og2/internal/og2"
	"hunter.io/og2/internal/og2/game"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostUser(t *testing.T) {
	testCases := []struct {
		name     string
		input    game.User
		expected game.User
	}{
		{
			name: "test create user",
			input: game.User{
				Name: "john doe",
			},
			expected: game.User{
				Name: "john doe",
			},
		},
	}

	for _, testCase := range testCases {
		require := require.New(t)
		assert := assert.New(t)

		db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
		require.NoError(err)
		defer db.Close()

		sessions, err := og2.NewSessions(db)
		require.NoError(err)

		input := og2.UserRequest{
			User: testCase.input,
		}

		b1, err := json.Marshal(input)
		require.NoError(err)

		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(b1))
		w := httptest.NewRecorder()

		h := og2.NewHandler(sessions).HandleUser()
		h.ServeHTTP(w, req)

		resp := w.Result()
		b2, err := ioutil.ReadAll(resp.Body)
		require.NoError(err)

		actual, err := game.Unmarshal(b2)
		require.NoError(err)

		assert.Equal(testCase.expected, actual.User)
	}
}
