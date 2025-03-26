package api

type Config struct {
	AdminApiUrl      *string // the api url, like : 'https://shop.yourdomain.com/api'
	StorefrontApiUrl *string // the storefront api url, like : 'https://shop.yourdomain.com/store-api'
	GrantType        string  // which grant type to use - can be either 'user_credentials'- or 'resource_owner'
	AdminApi         AdminApi
	AccessKey        *string // sw-access-key set in Administration/Sales Channels/API
}

type AdminApi struct {
	UserCredentials *UserCredentials // with refresh token
	ResourceOwner   *ResourceOwner   // no refresh token
}

type UserCredentials struct {
	// we recommend to only use this grant flow for client applications that should
	// perform administrative actions and require a user-based authentication
	Username string
	Password string
}

type ResourceOwner struct {
	// should be used for automated services
	ClientId     string // the client ID, setup at Web Administration Interface > settings > system > integration > access_id
	ClientSecret string // the client secret, setup at Web Administration Interface > settings > system > integration > access_secret
}
