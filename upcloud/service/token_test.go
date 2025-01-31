package service

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/client"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/stretchr/testify/assert"

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
		// TODO: obfuscate real tokes from fixtures. Currently committed tokens in token.yaml are from local env
		//  with the url changed to prod host. rec.AddFilter() for the win.
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

func TestClientWithToken(t *testing.T) {
	expires := time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC)
	tokenRequest := request.CreateTokenRequest{
		Name:               "my_1st_token",
		ExpiresAt:          expires,
		AllowedIPRanges:    []string{"0.0.0.0/0", "::/0"},
		CanCreateSubTokens: true,
	}

	// Create client that retries the initial token
	user, password := getCredentials()
	clt := client.New(user, password)
	svc := New(clt)

	// Get the initial token
	token, err := svc.CreateToken(context.Background(), &tokenRequest)
	require.NoError(t, err)

	// Create a new client with the initial token
	authCfg := client.WithBearerAuth(token.APIToken)
	cltWithToken := client.New("", "", authCfg)

	// Create a new service with the client with the initial token
	svcWithToken := New(cltWithToken)

	account, err := svcWithToken.GetAccount(context.Background())
	require.NoError(t, err)

	// Print account details
	t.Logf("Account: %+v", account)

	if account.UserName != user {
		t.Errorf("TestGetAccount expected %s, got %s", user, account.UserName)
	}

	assert.NotZero(t, account.ResourceLimits.Cores)
	assert.NotZero(t, account.ResourceLimits.Memory)
	assert.NotZero(t, account.ResourceLimits.Networks)
	assert.NotZero(t, account.ResourceLimits.PublicIPv6)
	assert.NotZero(t, account.ResourceLimits.StorageHDD)
	assert.NotZero(t, account.ResourceLimits.StorageSSD)

}

func cleanupTokenFunc(t *testing.T, svc *Service, id string) func() {
	return func() {
		if err := svc.DeleteToken(context.Background(), &request.DeleteTokenRequest{ID: id}); err != nil {
			t.Log(err)
		}
	}
}
