package upcloud

import (
	"encoding/xml"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestUnmarshalAccount tests that Account objects unmarshal correctly
func TestUnmarshalAccount(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<account>
  <credits>22465.536</credits>
  <resource_limits>
    <cores>100</cores>
    <memory>307200</memory>
    <networks>100</networks>
    <public_ipv4>20</public_ipv4>
    <public_ipv6>100</public_ipv6>
    <storage_hdd>10240</storage_hdd>
    <storage_ssd>10240</storage_ssd>
  </resource_limits>
  <username>foobar</username>
</account>`

	account := Account{}
	err := xml.Unmarshal([]byte(originalXML), &account)
	assert.Nil(t, err)
	assert.Equal(t, 22465.536, account.Credits)
	assert.Equal(t, "foobar", account.UserName)
	assert.Equal(t, 100, account.ResourceLimits.Cores)
	assert.Equal(t, 307200, account.ResourceLimits.Memory)
	assert.Equal(t, 100, account.ResourceLimits.Networks)
	assert.Equal(t, 20, account.ResourceLimits.PublicIpv4)
	assert.Equal(t, 100, account.ResourceLimits.PublicIpv6)
	assert.Equal(t, 10240, account.ResourceLimits.StorageHdd)
	assert.Equal(t, 10240, account.ResourceLimits.StorageSsd)
}

func TestUnmarshalAccountList(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<accounts>
  <account>
    <roles>
      <role>billing</role>
      <role>technical</role>
    </roles>
    <type>main</type>
    <username>foobar</username>
  </account>
</accounts>`

	accountList := AccountList{}
	err := xml.Unmarshal([]byte(originalXML), &accountList)
	assert.Nil(t, err)
	assert.Equal(t, "foobar", accountList.Accounts[0].UserName)
	assert.Equal(t, "main", accountList.Accounts[0].Type)
	assert.Equal(t, "billing", accountList.Accounts[0].Roles[0])
	assert.Equal(t, "technical", accountList.Accounts[0].Roles[1])
}

func TestUnmarshalAccountDetails(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<account>
  <address>Address Name St</address>
  <allow_api>yes</allow_api>
  <campaigns>
  </campaigns>
  <city>Indianapolis</city>
  <company>FooBar Inc.</company>
  <country>USA</country>
  <currency>USD</currency>
  <email>user@email.com</email>
  <enable_3rd_party_services>yes</enable_3rd_party_services>
  <first_name>John</first_name>
  <ip_filters>
  </ip_filters>
  <language>en</language>
  <last_name>Doe</last_name>
  <phone>+1.3173333333</phone>
  <postal_code>46250</postal_code>
  <roles>
    <role>billing</role>
    <role>technical</role>
  </roles>
  <simple_backup>no</simple_backup>
  <state>Indiana</state>
  <timezone>UTC</timezone>
  <type>main</type>
  <username>foobar</username>
  <vat_number>FI24315605</vat_number>
</account>`

	accountDetails := AccountDetails{}
	err := xml.Unmarshal([]byte(originalXML), &accountDetails)
	assert.Nil(t, err)
	assert.Equal(t, "John", accountDetails.FirstName)
	assert.Equal(t, "Doe", accountDetails.LastName)
	assert.Equal(t, "+1.3173333333", accountDetails.Phone)
	assert.Equal(t, "user@email.com", accountDetails.Email)
	assert.Equal(t, "FooBar Inc.", accountDetails.Company)
	assert.Equal(t, "Address Name St", accountDetails.Address)
	assert.Equal(t, "Indianapolis", accountDetails.City)
	assert.Equal(t, "Indiana", accountDetails.State)
	assert.Equal(t, 46250, accountDetails.PostalCode)
	assert.Equal(t, "USA", accountDetails.Country)
	assert.Equal(t, "UTC", accountDetails.Timezone)
	assert.Equal(t, "USD", accountDetails.Currency)
	assert.Equal(t, "yes", accountDetails.AllowApi)
	assert.Equal(t, "yes", accountDetails.Enable3rdPartyServices)
	assert.Equal(t, "en", accountDetails.Language)
	assert.Equal(t, "billing", accountDetails.Roles[0])
	assert.Equal(t, "no", accountDetails.SimpleBackup)
	assert.Equal(t, "main", accountDetails.Type)
	assert.Equal(t, "foobar", accountDetails.UserName)
	assert.Equal(t, "FI24315605", accountDetails.VatNumber)
}
