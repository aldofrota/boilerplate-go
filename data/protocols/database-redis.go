package protocols

type FlagPayload struct {
	Flagger      string   `json:"flagger"`
	ExpirationAt string   `json:"expiration_at"`
	FullRollout  bool     `json:"full_rollout"`
	Description  string   `json:"description"`
	IDs          []string `json:"ids"`
}

type Redis interface {
	FindAll() ([]FlagPayload, error)
	FindFlagger(flagger string) (FlagPayload, error)
	FindFlaggersCompany(id string) ([]string, error)
	CreateFlagger(payload FlagPayload) error
	DeleteFlagger(flagger string) error
	AddCompanyInFlagger(flagger string, id []string) error
	RemoveCompanyInFlagger(flagger string, ids []string) error
	UpdateFullRollout(flagger string, status bool) error
}
