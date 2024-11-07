package protocols

type HttpResponse struct {
	StatusCode int `json:"statusCode" example:"200"`
	Body       any `json:"body"`
}
type RequestCreateFlagger struct {
	Flagger      string   `json:"flagger"`
	ExpirationAt string   `json:"expiration_at"`
	FullRollout  bool     `json:"full_rollout"`
	IDs          []string `json:"ids"`
}

type RequestRelease struct {
	Ids []string `json:"ids"`
}

type RequestFullRollout struct {
	FullRollout bool `json:"full_rollout"`
}
