package domain

type UserOutlook struct {
	AvatarURL string `json:"image"`
	Name      string `json:"name"`
}

type UserOutlookAPI interface {
	GenerateAvatarAndName(id int) (*UserOutlook, error)
}
