package types

type Claim struct {
	UserID     string   `json:"user_id"`
	Username   string   `json:"username"`
	ConsumerID string   `json:"consumer_id"`
	Roles      []string `json:"roles"`
	Modules    []string `json:"modules"`
}
