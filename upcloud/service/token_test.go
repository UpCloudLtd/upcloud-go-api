package service

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"

	"github.com/stretchr/testify/require"
)

func TestToken(t *testing.T) {
	expires := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	tokenRequests := []request.CreateTokenRequest{
		{
			Name:               "my_1st_token",
			ExpiresAt:          expires,
			AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
			CanCreateSubTokens: true,
		},
		{
			Name:               "my_2nd_token",
			ExpiresAt:          expires,
			AllowedIPRanges:    []string{"0.0.0.0/1", "::/0"},
			CanCreateSubTokens: false,
		},
	}

	record(t, "token", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		// Create some tokens
		ids := make([]string, len(tokenRequests))

		for i, req := range tokenRequests {
			token, err := svc.CreateToken(ctx, &req)
			require.NoError(t, err)
			t.Cleanup(cleanupTokenFunc(t, svc, token.ID))

			ids[i] = token.ID
			assert.True(t, strings.HasPrefix(token.APIToken, "ucat_"))
			assert.Equal(t, req.Name, token.Name)
			assert.Equal(t, req.AllowedIPRanges, token.AllowedIPRanges)
			assert.Equal(t, req.ExpiresAt.Format(time.RFC3339), token.ExpiresAt.Format(time.RFC3339))
			assert.Equal(t, req.CanCreateSubTokens, token.CanCreateSubTokens)
		}

		// Get one token
		token, err := svc.GetTokenDetails(ctx, &request.GetTokenDetailsRequest{ID: ids[0]})
		require.NoError(t, err)

		assert.Equal(t, "my_1st_token", token.Name)
		assert.Equal(t, []string{"0.0.0.0/0", "::/0"}, token.AllowedIPRanges)
		assert.Equal(t, expires.Format(time.RFC3339), token.ExpiresAt.Format(time.RFC3339))
		assert.Equal(t, true, token.CanCreateSubTokens)

		// List tokens
		tokens, err := svc.GetTokens(ctx, &request.GetTokensRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(*tokens), len(tokenRequests))

		// Create a token and delete it immediately
		deleteThis, err := svc.CreateToken(ctx, &tokenRequests[0])
		require.NoError(t, err)
		require.NoError(t, svc.DeleteToken(ctx, &request.DeleteTokenRequest{ID: deleteThis.ID}))
	})
}

// TestClientWithToken tests that a client can be created with a token and used to make authenticated API requests
func TestClientWithToken(t *testing.T) {
	if os.Getenv("UPCLOUD_GO_SDK_TEST_NO_CREDENTIALS") == "yes" || testing.Short() {
		t.Skip("Skipping TestGetAccount...")
	}
	expires := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	tokenRequest := request.CreateTokenRequest{
		Name:               "my_1st_token",
		ExpiresAt:          expires,
		AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		CanCreateSubTokens: true,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	// Get the initial token
	user, password := getCredentials()
	svc := New(client.New(user, password))
	token, err := svc.CreateToken(ctx, &tokenRequest)
	require.NoError(t, err, "Failed to create token")
	require.NotNil(t, token, "Token must not be nil")

	// Make sure that we cleanup the token
	t.Cleanup(cleanupTokenFunc(t, svc, token.ID))

	// Create a new client with the token
	svcWithToken := New(client.New("", "", client.WithBearerAuth(token.APIToken)))

	// Make an authenticated API request
	server, err := svcWithToken.GetServers(ctx)
	require.NotEmpty(t, server, "Failed to get servers. This points to a problem with token auth")
	require.NoError(t, err, "Error getting the servers. This points to a problem with token auth")

	// Delete the token
	err = svcWithToken.DeleteToken(ctx, &request.DeleteTokenRequest{ID: token.ID})
	require.NoError(t, err, "Token deletion should not fail")

	// Make sure the token is deleted
	server, err = svcWithToken.GetServers(ctx)
	require.Error(t, err, "Getting servers with deleted token should fail")
	require.Empty(t, server, "Getting servers with deleted token should return empty list")
}

func cleanupTokenFunc(t *testing.T, svc *Service, id string) func() {
	return func() {
		if err := svc.DeleteToken(context.Background(), &request.DeleteTokenRequest{ID: id}); err != nil {
			t.Log(err, "This might not be a problem if the test deleted the token already")
		}
	}
}
