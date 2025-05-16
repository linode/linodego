package linodego

import "context"

type UserRolePermissions struct {
	AccountAccess []string     `json:"account_access"`
	EntityAccess  []UserAccess `json:"entity_access"`
}

type UserRolePermissionsUpdateOptions struct {
	AccountAccess []string     `json:"account_access"`
	EntityAccess  []UserAccess `json:"entity_access"`
}

type UserAccess struct {
	ID    int      `json:"id"`
	Type  string   `json:"type"`
	Roles []string `json:"roles"`
}

type AccountRolePermissions struct {
	AccountAccess []AccountAccess `json:"account_access"`
	EntityAccess  []AccountAccess `json:"entity_access"`
}

type AccountAccess struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Roles []Role `json:"roles"`
}

type Role struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Permissions []string `json:"permissions"`
}

func (c *Client) GetUserRolePermissions(ctx context.Context, username string) (*UserRolePermissions, error) {
	return doGETRequest[UserRolePermissions](ctx, c,
		formatAPIPath("iam/users/%s/role-permissions", username),
	)
}

func (c *Client) UpdateUserRolePermissions(ctx context.Context, username string, opts UserRolePermissionsUpdateOptions) (*UserRolePermissions, error) {
	return doPUTRequest[UserRolePermissions](ctx, c,
		formatAPIPath("iam/users/%s/role-permissions", username),
		opts,
	)
}

// TODO: GET UPDATE USER ROLE OPTIONS

func (c *Client) GetAccountRolePermissions(ctx context.Context) (*AccountRolePermissions, error) {
	return doGETRequest[AccountRolePermissions](ctx, c, "iam/role-permissions")
}
