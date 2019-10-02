package upcloud

import (
	"encoding/xml"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestUnmarshalStorage tests that Storages and Storage struct are unmarshaled correctly
func TestUnmarshalStorage(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<storages>
    <storage>
        <access>public</access>
        <license>0</license>
        <size>1</size>
        <state>online</state>
        <title>Windows Server 2003 R2 Standard (CD 1)</title>
        <type>cdrom</type>
        <uuid>01000000-0000-4000-8000-000010010101</uuid>
    </storage>
    <storage>
        <access>public</access>
        <license>0</license>
        <size>1</size>
        <state>online</state>
        <title>Windows Server 2003 R2 Standard (CD 2)</title>
        <type>cdrom</type>
        <uuid>01000000-0000-4000-8000-000010010102</uuid>
    </storage>
    <storage>
        <access>public</access>
        <license>0</license>
        <size>1</size>
        <state>online</state>
        <title>Windows Server 2003 R2 Standard (CD 1)</title>
        <type>cdrom</type>
        <uuid>01000000-0000-4000-8000-000010010201</uuid>
    </storage>
</storages>`

	storages := Storages{}
	err := xml.Unmarshal([]byte(originalXML), &storages)

	assert.Nil(t, err)
	assert.Len(t, storages.Storages, 3)

	firstStorage := storages.Storages[0]
	assert.Equal(t, "public", firstStorage.Access)
	assert.Equal(t, 0.0, firstStorage.License)
	assert.Equal(t, 1, firstStorage.Size)
	assert.Equal(t, "Windows Server 2003 R2 Standard (CD 1)", firstStorage.Title)
	assert.Equal(t, StorageTypeCDROM, firstStorage.Type)
	assert.Equal(t, "01000000-0000-4000-8000-000010010101", firstStorage.UUID)
}

// TestUnmarshalStorage tests that Storages and Storage struct are unmarshaled correctly for private and backup storages
func TestUnmarshalStoragesPrivateAndBackup(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<storages>
    <storage>
        <access>private</access>
        <license>0</license>
        <size>10</size>
        <state>online</state>
	<tier>hdd</tier>
        <title>Operating system disk</title>
        <type>normal</type>
        <uuid>01eff7ad-168e-413e-83b0-054f6a28fa23</uuid>
        <zone>uk-lon1</zone>
    </storage>
    <storage>
        <access>private</access>
	<created>2019-09-17T14:34:43Z</created>
        <license>0</license>
	<origin>01eff7ad-168e-413e-83b0-054f6a28fa23</origin>
        <size>10</size>
        <state>online</state>
        <title>On demand backup</title>
        <type>backup</type>
        <uuid>01287ad1-496c-4b5f-bb67-0fc2e3494740</uuid>
        <zone>uk-lon1</zone>
    </storage>
    <storage>
        <access>private</access>
        <license>0</license>
        <part_of_plan>yes</part_of_plan>
        <size>50</size>
        <state>online</state>
        <tier>maxiops</tier>
        <title>Databases</title>
        <type>normal</type>
        <uuid>01f3286c-a5ea-4670-8121-d0b9767d625b</uuid>
        <zone>fi-hel1</zone>
    </storage>
</storages>`

	storages := Storages{}
	err := xml.Unmarshal([]byte(originalXML), &storages)

	assert.Nil(t, err)
	assert.Len(t, storages.Storages, 3)

	firstStorage := storages.Storages[0]
	assert.Equal(t, "private", firstStorage.Access)
	assert.Equal(t, 0.0, firstStorage.License)
	assert.Equal(t, 10, firstStorage.Size)
	assert.Equal(t, "online", firstStorage.State)
	assert.Equal(t, "hdd", firstStorage.Tier)
	assert.Equal(t, "Operating system disk", firstStorage.Title)
	assert.Equal(t, "normal", firstStorage.Type)
	assert.Equal(t, "01eff7ad-168e-413e-83b0-054f6a28fa23", firstStorage.UUID)
	assert.Equal(t, "uk-lon1", firstStorage.Zone)
	assert.Equal(t, time.Time{}, firstStorage.Created)
	assert.Equal(t, "", firstStorage.Origin)

	secondStorage := storages.Storages[1]
	assert.Equal(t, "private", secondStorage.Access)
	created, err := time.Parse(time.RFC3339, "2019-09-17T14:34:43Z")
	assert.Nil(t, err)
	assert.Equal(t, created, secondStorage.Created)
	assert.Equal(t, 0.0, secondStorage.License)
	assert.Equal(t, "01eff7ad-168e-413e-83b0-054f6a28fa23", secondStorage.Origin)
	assert.Equal(t, 10, secondStorage.Size)
	assert.Equal(t, "online", secondStorage.State)
	assert.Equal(t, "", secondStorage.Tier)
	assert.Equal(t, "On demand backup", secondStorage.Title)
	assert.Equal(t, "backup", secondStorage.Type)
	assert.Equal(t, "01287ad1-496c-4b5f-bb67-0fc2e3494740", secondStorage.UUID)
	assert.Equal(t, "uk-lon1", secondStorage.Zone)
}

// TestUnmarshalStorageDetails tests that StorageDetails struct is unmarshaled correctly
func TestUnmarshalStorageDetails(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<storage>
    <access>private</access>
    <backup_rule>
        <interval>daily</interval>
        <retention>1</retention>
        <time>0400</time>
    </backup_rule>
    <backups>
        <backup>37c96670-9c02-4d5d-8f60-291d38f9a80c</backup>
        <backup>ecfda9f2-e071-4bbb-b38f-079ed26eb32a</backup>
    </backups>
    <license>0</license>
    <servers>
        <server>33850294-50f4-4712-8463-aeb7b42de42f</server>
    </servers>
    <size>500</size>
    <state>online</state>
    <tier>maxiops</tier>
    <title>Debian server (Disk 1)</title>
    <type>normal</type>
    <uuid>bf3da6c2-66c4-4e70-9640-5b4896aacd5c</uuid>
    <zone>fi-hel1</zone>
</storage>`
	storageDeviceDetails := StorageDetails{}
	err := xml.Unmarshal([]byte(originalXML), &storageDeviceDetails)

	assert.Nil(t, err)
	assert.Equal(t, "private", storageDeviceDetails.Access)
	assert.Equal(t, 0.0, storageDeviceDetails.License)
	assert.Equal(t, 500, storageDeviceDetails.Size)
	assert.Equal(t, "online", storageDeviceDetails.State)
	assert.Equal(t, "maxiops", storageDeviceDetails.Tier)
	assert.Equal(t, "Debian server (Disk 1)", storageDeviceDetails.Title)
	assert.Equal(t, StorageTypeNormal, storageDeviceDetails.Type)
	assert.Equal(t, "bf3da6c2-66c4-4e70-9640-5b4896aacd5c", storageDeviceDetails.UUID)
	assert.Equal(t, "fi-hel1", storageDeviceDetails.Zone)

	assert.Equal(t, BackupRuleIntervalDaily, storageDeviceDetails.BackupRule.Interval)
	assert.Equal(t, 1, storageDeviceDetails.BackupRule.Retention)
	assert.Equal(t, "0400", storageDeviceDetails.BackupRule.Time)

	assert.Equal(t, 2, len(storageDeviceDetails.BackupUUIDs))
	assert.Equal(t, "37c96670-9c02-4d5d-8f60-291d38f9a80c", storageDeviceDetails.BackupUUIDs[0])
	assert.Equal(t, "ecfda9f2-e071-4bbb-b38f-079ed26eb32a", storageDeviceDetails.BackupUUIDs[1])

	assert.Equal(t, 1, len(storageDeviceDetails.ServerUUIDs))
	assert.Equal(t, "33850294-50f4-4712-8463-aeb7b42de42f", storageDeviceDetails.ServerUUIDs[0])
}

// TestUnmarshalServerStorageDevice tests that ServerStorageDevice objects are properly unmarshaled
func TestUnmarshalServerStorageDevice(t *testing.T) {
	originalXML := `<?xml version="1.0" encoding="utf-8"?>
<storage_device>
    <address>virtio:0</address>
    <part_of_plan>yes</part_of_plan>
    <storage>01c8df16-d1c6-4223-9bfc-d3c06b208c88</storage>
    <storage_size>30</storage_size>
    <storage_title>test-disk0</storage_title>
    <type>disk</type>
</storage_device>`

	storageDevice := ServerStorageDevice{}
	err := xml.Unmarshal([]byte(originalXML), &storageDevice)

	assert.Nil(t, err)
	assert.Equal(t, "virtio:0", storageDevice.Address)
	assert.Equal(t, "yes", storageDevice.PartOfPlan)
	assert.Equal(t, "01c8df16-d1c6-4223-9bfc-d3c06b208c88", storageDevice.UUID)
	assert.Equal(t, 30, storageDevice.Size)
	assert.Equal(t, "test-disk0", storageDevice.Title)
	assert.Equal(t, StorageTypeDisk, storageDevice.Type)
}

// TestMarshalCreateStorageDevice tests that CreateStorageDevice objects are correctly marshaled. We don't need to
// test unmarshaling because these data structures are never returned from the API.
func TestMarshalCreateStorageDevice(t *testing.T) {
	storage := CreateServerStorageDevice{
		Action:  CreateServerStorageDeviceActionClone,
		Storage: "01000000-0000-4000-8000-000030060200",
		Title:   "disk1",
		Size:    30,
		Tier:    StorageTierMaxIOPS,
	}

	expectedXML := "<storage_device><action>clone</action><storage>01000000-0000-4000-8000-000030060200</storage><title>disk1</title><size>30</size><tier>maxiops</tier></storage_device>"

	actualXML, err := xml.Marshal(storage)
	assert.Nil(t, err)
	assert.Equal(t, expectedXML, string(actualXML))
}

// TestMarshalBackupRule tests that BackupRule objects are properly marshaled
func TestMarshalBackupRule(t *testing.T) {
	backupRule := BackupRule{
		Interval:  BackupRuleIntervalDaily,
		Time:      "0430",
		Retention: 30,
	}

	ruleXML, err := xml.Marshal(backupRule)
	assert.Nil(t, err)
	assert.Equal(t, "<backup_rule><interval>daily</interval><time>0430</time><retention>30</retention></backup_rule>", string(ruleXML))
}

// TestUnmarshalBackupRule tests that BackupRule objects are properly unmarshaled
func TestUnmarshalBackupRule(t *testing.T) {
	originalXML := "<backup_rule><interval>daily</interval><time>0430</time><retention>30</retention></backup_rule>"

	backupRule := BackupRule{}
	err := xml.Unmarshal([]byte(originalXML), &backupRule)
	assert.Nil(t, err)
	assert.Equal(t, BackupRuleIntervalDaily, backupRule.Interval)
	assert.Equal(t, "0430", backupRule.Time)
	assert.Equal(t, 30, backupRule.Retention)
}
