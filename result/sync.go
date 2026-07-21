package result

type Sync struct {
	Profile     *Profile      `json:"profile"`
	Folders     []*Folder     `json:"folders"`
	Collections []*Collection `json:"collections"`
	Items       []*Item       `json:"ciphers"`

	// policies
	// sends
	// domains
	// policiesNew
	// userDecryption
}
