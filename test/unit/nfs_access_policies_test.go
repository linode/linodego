package unit

import (
	"context"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/linode/linodego/v2"
	"github.com/stretchr/testify/assert"
)

func TestNFSSpaceAccessPolicy_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_space_access_policy")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/access-policy", fixtureData)

	policy, err := base.Client.GetNFSSpaceAccessPolicy(context.Background(), 123)
	assert.NoError(t, err)
	assertNFSSpaceAccessPolicy(t, policy)
}

func TestNFSSpaceAccessPolicy_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_space_access_policy")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSSpaceAccessPolicyUpdateOptions{
		Label:      linodego.Pointer("production-policy"),
		Enabled:    linodego.Pointer(true),
		VPCs:       linodego.Pointer([]linodego.NFSSpaceAccessPolicyVPCOptions{{ID: 789, Subnets: []int{101}}}),
		MTLSMode:   linodego.Pointer(linodego.NFSMTLSModeOptional),
		MTLSCACert: linodego.Pointer("certificate"),
	}

	httpmock.RegisterRegexpResponder(
		"PUT",
		mockRequestURL(t, "nfs/spaces/123/access-policy"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	policy, err := client.UpdateNFSSpaceAccessPolicy(context.Background(), 123, opts)
	assert.NoError(t, err)
	assertNFSSpaceAccessPolicy(t, policy)
}

func TestNFSFilesystemAccessPolicy_Get(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem_access_policy")
	assert.NoError(t, err)

	var base ClientBaseCase
	base.SetUp(t)
	defer base.TearDown(t)

	base.MockGet("nfs/spaces/123/filesystems/456/access-policy", fixtureData)

	policy, err := base.Client.GetNFSFilesystemAccessPolicy(context.Background(), 123, 456)
	assert.NoError(t, err)
	assertNFSFilesystemAccessPolicy(t, policy)
}

func TestNFSFilesystemAccessPolicy_Update(t *testing.T) {
	fixtureData, err := fixtures.GetFixture("nfs_filesystem_access_policy")
	assert.NoError(t, err)

	client := createMockClient(t)
	opts := linodego.NFSFilesystemAccessPolicyUpdateOptions{
		Label:        linodego.Pointer("filesystem-policy"),
		Enabled:      linodego.Pointer(true),
		LinodeIDs:    linodego.Pointer([]int{987}),
		SquashPolicy: linodego.Pointer(linodego.NFSSquashPolicyRootSquash),
		Protocols:    linodego.Pointer([]linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}),
	}

	httpmock.RegisterRegexpResponder(
		"PUT",
		mockRequestURL(t, "nfs/spaces/123/filesystems/456/access-policy"),
		mockRequestBodyValidate(t, opts, fixtureData),
	)

	policy, err := client.UpdateNFSFilesystemAccessPolicy(context.Background(), 123, 456, opts)
	assert.NoError(t, err)
	assertNFSFilesystemAccessPolicy(t, policy)
}

func assertNFSSpaceAccessPolicy(t *testing.T, policy *linodego.NFSSpaceAccessPolicy) {
	t.Helper()
	if !assert.NotNil(t, policy) {
		return
	}

	assert.Equal(t, 123, policy.SpaceID)
	assert.Equal(t, "production-policy", policy.Label)
	assert.True(t, policy.Enabled)
	assert.Equal(t, linodego.NFSMTLSModeOptional, policy.MTLSMode)
	assert.Equal(t, linodego.NFSAccessPolicyStatusActive, policy.Status)
	if assert.Len(t, policy.VPCACL, 1) {
		assert.Equal(t, 789, policy.VPCACL[0].ID)
		assert.Len(t, policy.VPCACL[0].Subnets, 1)
	}
	assert.Equal(t, "2026-05-01T14:00:00Z", policy.Created.Format(time.RFC3339))
	assert.Equal(t, "2026-05-10T10:00:00Z", policy.Updated.Format(time.RFC3339))
}

func assertNFSFilesystemAccessPolicy(t *testing.T, policy *linodego.NFSFilesystemAccessPolicy) {
	t.Helper()
	if !assert.NotNil(t, policy) {
		return
	}

	assert.Equal(t, 456, policy.FilesystemID)
	assert.Equal(t, "filesystem-policy", policy.Label)
	assert.True(t, policy.Enabled)
	assert.Equal(t, linodego.NFSSquashPolicyRootSquash, policy.SquashPolicy)
	assert.Equal(t, []linodego.NFSProtocolVersion{linodego.NFSProtocolVersionV4}, policy.Protocols)
	assert.Equal(t, linodego.NFSAccessPolicyStatusActive, policy.Status)
	if assert.Len(t, policy.LinodeACL, 1) {
		assert.Equal(t, 987, policy.LinodeACL[0].ID)
	}
	assert.Equal(t, "2026-05-01T14:00:00Z", policy.Created.Format(time.RFC3339))
	assert.Equal(t, "2026-05-10T10:00:00Z", policy.Updated.Format(time.RFC3339))
}
