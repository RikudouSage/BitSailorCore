package result

import (
	"net"
	"time"

	"github.com/google/uuid"
)

type AuthRequest struct {
	ID                      uuid.UUID  `json:"Id"`
	PublicKey               []byte     `json:"PublicKey"`
	RequestDeviceType       string     `json:"RequestDeviceType"`
	RequestDeviceTypeValue  int        `json:"RequestDeviceTypeValue"`
	RequestDeviceIdentifier string     `json:"RequestDeviceIdentifier"`
	RequestIPAddress        net.IP     `json:"RequestIpAddress"`
	RequestCountryName      string     `json:"RequestCountryName"`
	Key                     *string    `json:"Key"`
	CreationDate            time.Time  `json:"CreationDate"`
	RequestApproved         *bool      `json:"RequestApproved"`
	ResponseDate            *time.Time `json:"ResponseDate"`
	RequestDeviceID         *uuid.UUID `json:"RequestDeviceId"`
}
