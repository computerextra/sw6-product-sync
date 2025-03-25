package api

type Config struct {
	AdminApiUrl      *string
	StorefrontApiUrl *string
	//  which grant type to use - can be either 'user_credentials'- or 'resource_owner'
	GrantType string
	AdminApi  AdminApi
}

type AdminApi struct {
	UserCredentials *UserCredentials
	ResourceOwner   *ResourceOwner
}

type UserCredentials struct {
	// With Refresh token
	// we recommend to only use this grant flow for client applications that should
	// perform administrative actions and require a user-based authentication
	Username string
	Password string
}

type ResourceOwner struct {
	// Without Refresh token
	// should be used for automated services
	ClientId     string
	ClientSecret string
}
