package result

import "time"

type AuthData struct {
	AccessToken  string    `json:"accessToken"`
	ExpiresAt    time.Time `json:"expiresAt"`
	RefreshToken string    `json:"refreshToken"`
	TokenType    string    `json:"tokenType"`
}
