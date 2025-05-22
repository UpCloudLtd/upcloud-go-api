package service

import (
	"context"
	"testing"

	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud"
	"github.com/UpCloudLtd/upcloud-go-api/v8/upcloud/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

func TestCreatePartnerAccount(t *testing.T) {
	record(t, "createpartneraccount", func(ctx context.Context, t *testing.T, rec *recorder.Recorder, svc *Service) {
		username := "sdk_test_partner_account"
		expectedAccount := upcloud.PartnerAccount{
			Username:  username,
			FirstName: "Test",
			LastName:  "User",
			Country:   "FIN",
			Email:     "test@myhost.mydomain",
			Phone:     "+358.31245434",
		}

		account, err := svc.CreatePartnerAccount(ctx, &request.CreatePartnerAccountRequest{
			Username: username,
			Password: "superSecret123",
			ContactDetails: &request.CreatePartnerAccountContactDetails{
				FirstName: expectedAccount.FirstName,
				LastName:  expectedAccount.LastName,
				Country:   expectedAccount.Country,
				Email:     expectedAccount.Email,
				Phone:     expectedAccount.Phone,
			},
		})
		require.NoError(t, err)
		assert.Equal(t, expectedAccount, *account)

		accounts, err := svc.GetPartnerAccounts(ctx)
		require.NoError(t, err)
		accountFound := false
		for _, a := range accounts {
			if a.Username == username {
				accountFound = true
				assert.Equal(t, expectedAccount, a)
			}
		}
		assert.True(t, accountFound)
	})
}
