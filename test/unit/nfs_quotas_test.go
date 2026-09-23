package unit

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestNFSQuota_List(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_quotas_list")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/filesystems/456/quotas", fixtureData)

	quotas, err := base.Client.ListNFSQuotas(context.Background(), 123, 456, &linodego.ListOptions{})
	assert.NoError(t, err)
	if assert.Len(t, quotas, 1) {
		assertNFSQuota(t, &quotas[0])
		assert.NotNil(t, quotas[0].CollectedAt)
	}
}

func TestNFSQuota_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_quota")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/filesystems/456/quotas/789", fixtureData)

	quota, err := base.Client.GetNFSQuota(context.Background(), 123, 456, 789)
	assert.NoError(t, err)
	assertNFSQuota(t, quota)
	assert.Nil(t, quota.UsedCapacityBytes)
	assert.Nil(t, quota.UsedFileCount)
	assert.Nil(t, quota.CollectedAt)
}

func TestNFSQuota_Create(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_quota")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSQuotaCreateOptions{
		Path:             "/training-data-10932/data/project-x",
		MaxCapacityBytes: 107374182400,
		MaxFileCount:     1000000,
		UserGroupConfig: &linodego.NFSUserGroupConfigUpdateOptions{
			DefaultUserLimit: &linodego.NFSCapacityLimitUpdateOptions{
				MaxCapacityBytes: linodego.Pointer[int64](107374182400),
				MaxFileCount:     linodego.Pointer[int64](1000000),
			},
			DefaultGroupLimit: &linodego.NFSCapacityLimitUpdateOptions{
				MaxCapacityBytes: linodego.Pointer[int64](107374182400),
				MaxFileCount:     linodego.Pointer[int64](1000000),
			},
			UserLimits: &[]linodego.NFSIdentifiedLimitUpdateOptions{
				{
					IdentifierType:   linodego.NFSQuotaRuleIdentifierTypeUID,
					Identifier:       "1001",
					MaxCapacityBytes: linodego.Pointer[int64](107374182400),
					MaxFileCount:     linodego.Pointer[int64](1000000),
				},
			},
			GroupLimits: &[]linodego.NFSIdentifiedLimitUpdateOptions{
				{
					IdentifierType:   linodego.NFSQuotaRuleIdentifierTypeUID,
					Identifier:       "1001",
					MaxCapacityBytes: linodego.Pointer[int64](107374182400),
					MaxFileCount:     linodego.Pointer[int64](1000000),
				},
			},
		},
	}

	httpmock.RegisterRegexpResponder(
		"POST",
		mockRequestURL(t, "nfs/spaces/123/filesystems/456/quotas"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	quota, err := client.CreateNFSQuota(context.Background(), 123, 456, opts)
	assert.NoError(t, err)
	assertNFSQuota(t, quota)
}

func TestNFSQuota_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_quota")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSQuotaUpdateOptions{
		MaxCapacityBytes: linodego.Pointer[int64](107374182400),
		MaxFileCount:     linodego.Pointer[int64](1000000),
		UserGroupConfig: &linodego.NFSUserGroupConfigUpdateOptions{
			DefaultUserLimit: &linodego.NFSCapacityLimitUpdateOptions{
				MaxCapacityBytes: linodego.Pointer[int64](107374182400),
				MaxFileCount:     linodego.Pointer[int64](1000000),
			},
			DefaultGroupLimit: &linodego.NFSCapacityLimitUpdateOptions{
				MaxCapacityBytes: linodego.Pointer[int64](107374182400),
				MaxFileCount:     linodego.Pointer[int64](1000000),
			},
			UserLimits: &[]linodego.NFSIdentifiedLimitUpdateOptions{
				{
					IdentifierType:   linodego.NFSQuotaRuleIdentifierTypeUID,
					Identifier:       "1001",
					MaxCapacityBytes: linodego.Pointer[int64](107374182400),
					MaxFileCount:     linodego.Pointer[int64](1000000),
				},
			},
			GroupLimits: &[]linodego.NFSIdentifiedLimitUpdateOptions{
				{
					IdentifierType:   linodego.NFSQuotaRuleIdentifierTypeUID,
					Identifier:       "1001",
					MaxCapacityBytes: linodego.Pointer[int64](107374182400),
					MaxFileCount:     linodego.Pointer[int64](1000000),
				},
			},
		},
	}

	httpmock.RegisterRegexpResponder(
		"PUT",
		mockRequestURL(t, "nfs/spaces/123/filesystems/456/quotas/789"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	quota, err := client.UpdateNFSQuota(context.Background(), 123, 456, 789, opts)
	assert.NoError(t, err)
	assertNFSQuota(t, quota)
}

func TestNFSQuota_Delete(t *testing.T) {
	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockDelete("nfs/spaces/123/filesystems/456/quotas/789", nil)

	assert.NoError(t, base.Client.DeleteNFSQuota(context.Background(), 123, 456, 789))
}

func assertNFSQuota(t *testing.T, quota *linodego.NFSQuota) {
	t.Helper()
	if !assert.NotNil(t, quota) {
		return
	}

	assert.Equal(t, 789, quota.ID)
	assert.Equal(t, 456, quota.FilesystemID)
	assert.Equal(t, "/training-data-10932/data/project-x", quota.Path)
	if assert.NotNil(t, quota.MaxCapacityBytes) {
		assert.Equal(t, int64(107374182400), *quota.MaxCapacityBytes)
	}
	if assert.NotNil(t, quota.MaxFileCount) {
		assert.Equal(t, int64(1000000), *quota.MaxFileCount)
	}
	assert.Equal(t, linodego.NFSQuotaStatusActive, quota.Status)
	assert.Equal(t, "2026-05-01T14:00:00Z", quota.Created.Format(time.RFC3339))

	config := quota.UserGroupConfig

	if assert.NotNil(t, config.DefaultUserLimit) {
		assert.NotNil(t, config.DefaultUserLimit.MaxCapacityBytes)
		assert.Equal(t, int64(21474836480), *config.DefaultUserLimit.MaxCapacityBytes)
		assert.NotNil(t, config.DefaultUserLimit.MaxFileCount)
		assert.Equal(t, int64(1000000), *config.DefaultUserLimit.MaxFileCount)
	}

	if assert.NotNil(t, config.DefaultGroupLimit) {
		assert.NotNil(t, config.DefaultGroupLimit.MaxCapacityBytes)
		assert.Equal(t, int64(21474836480), *config.DefaultGroupLimit.MaxCapacityBytes)
		assert.NotNil(t, config.DefaultGroupLimit.MaxFileCount)
		assert.Equal(t, int64(1000000), *config.DefaultGroupLimit.MaxFileCount)
	}

	if assert.Len(t, config.UserLimits, 1) {
		userLimit := config.UserLimits[0]
		assert.Equal(t, linodego.NFSQuotaRuleIdentifierTypeUID, userLimit.IdentifierType)
		assert.Equal(t, "1001", userLimit.Identifier)
		assert.NotNil(t, userLimit.MaxCapacityBytes)
		assert.Equal(t, int64(21474836480), *userLimit.MaxCapacityBytes)
		assert.NotNil(t, userLimit.MaxFileCount)
		assert.Equal(t, int64(1000000), *userLimit.MaxFileCount)
	}

	if assert.Len(t, config.GroupLimits, 1) {
		groupLimit := config.GroupLimits[0]
		assert.Equal(t, linodego.NFSQuotaRuleIdentifierTypeUID, groupLimit.IdentifierType)
		assert.Equal(t, "1001", groupLimit.Identifier)
		assert.NotNil(t, groupLimit.MaxCapacityBytes)
		assert.Equal(t, int64(21474836480), *groupLimit.MaxCapacityBytes)
		assert.NotNil(t, groupLimit.MaxFileCount)
		assert.Equal(t, int64(1000000), *groupLimit.MaxFileCount)
	}
}
