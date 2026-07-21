package result

import "time"

type AuthData struct {
	AccessToken  string
	ExpiresAt    time.Time
	RefreshToken string
	TokenType    string
}
