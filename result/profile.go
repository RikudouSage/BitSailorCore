package result

import (
	"time"

	"github.com/google/uuid"
)

type Profile struct {
	ID                      uuid.UUID `json:"id"`
	Email                   string    `json:"email"`
	Premium                 bool      `json:"premium"`
	PrivateKey              string    `json:"privateKey"`
	SecurityStamp           string    `json:"securityStamp"`
	Name                    string    `json:"name"`
	PremiumFromOrganization bool      `json:"premiumFromOrganization"`
	Locale                  string    `json:"culture"`
	Key                     string    `json:"key"`
	EmailVerified           bool      `json:"emailVerified"`
	ForcePasswordReset      bool      `json:"forcePasswordReset"`
	UsesKeyConnector        bool      `json:"usesKeyConnector"`
	CreationDate            time.Time `json:"creationDate"`
	TwoFactorEnabled        bool      `json:"twoFactorEnabled"`
	VerifyDevices           bool      `json:"verifyDevices"`

	//AvatarColor *string `json:"avatarColor"`
	// providerOrganizations
	// accountKeys
	// organizations
}
