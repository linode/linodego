package linodego

import "context"

// LinodeQuota represents a Linode-related quota information on your account.
type LinodeQuota struct {
	QuotaID        int    `json:"quota_id"`
	QuotaName      string `json:"quota_name"`
	Description    string `json:"description"`
	QuotaLimit     int    `json:"quota_limit"`
	ResourceMetric string `json:"resource_metric"`
	RegionApplied  string `json:"region_applied"`
}

// LinodeQuotaUsage is the usage data for a specific Linode-related quota on your account.
type LinodeQuotaUsage struct {
	QuotaLimit int  `json:"quota_limit"`
	Usage      *int `json:"usage"`
}

// ListLinodeQuotas lists the active Linode-related quotas applied to your account.
// Linode Quota related features are under v4beta and may not currently be available to all users.
func (c *Client) ListLinodeQuotas(ctx context.Context, opts *ListOptions) ([]LinodeQuota, error) {
	return getPaginatedResults[LinodeQuota](ctx, c, formatAPIPath("linode/quotas"), opts)
}

// GetLinodeQuota gets information about a specific Linode-related quota on your account.
// The operation includes any quota overrides in the response.
// Linode Quota related features are under v4beta and may not currently be available to all users.
func (c *Client) GetLinodeQuota(ctx context.Context, quotaID int) (*LinodeQuota, error) {
	e := formatAPIPath("linode/quotas/%d", quotaID)
	return doGETRequest[LinodeQuota](ctx, c, e)
}

// GetLinodeQuotaUsage gets usage data for a specific Linode Quota resource you can have on your account and the current usage for that resource.
// Linode Quota related features are under v4beta and may not currently be available to all users.
func (c *Client) GetLinodeQuotaUsage(ctx context.Context, quotaID int) (*LinodeQuotaUsage, error) {
	e := formatAPIPath("linode/quotas/%d/usage", quotaID)
	return doGETRequest[LinodeQuotaUsage](ctx, c, e)
}
