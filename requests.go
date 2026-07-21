package bitwarden

type preLoginRequest struct {
	Email string `json:"email"`
}

type passwordLoginRequest struct {
	GrantType string `url:"grant_type"`

	Username string `url:"username"`
	Password string `url:"password"`

	Scope    string `url:"scope"`
	ClientID string `url:"client_id"`

	DeviceType       int    `url:"deviceType"`
	DeviceIdentifier string `url:"deviceIdentifier"`
	DeviceName       string `url:"deviceName"`

	TwoFactorProvider       *int    `url:"twoFactorProvider,omitempty"`
	TwoFactorToken          *string `url:"twoFactorToken,omitempty"`
	TwoFactorRemember       *int    `url:"twoFactorRemember,omitempty"`
	SsoEmail2faSessionToken *string `url:"ssoEmail2faSessionToken,omitempty"`
}

type apiKeyLoginRequest struct {
	GrantType    string `url:"grant_type"`
	ClientID     string `url:"client_id"`
	ClientSecret string `url:"client_secret"`
	Scope        string `url:"scope"`

	DeviceType       int    `url:"deviceType"`
	DeviceIdentifier string `url:"deviceIdentifier"`
	DeviceName       string `url:"deviceName"`
}
