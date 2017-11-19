package ci

// Webhook represents a struct in payload
type Webhook struct {
	ID                int    `json:"id"`
	Number            string `json:"number"`
	PullRequestNumber int    `json:"pull_request_number"`
	PullRequestTitle  string `json:"pull_request_title"`
	Duration          int    `json:"duration"`
	AuthorName        string `json:"author_name"`
	AuthorEmail       string `json:"author_email"`
	Type              string `json:"type"`
	State             string `json:"state"`
	BuildURL          string `json:"build_url"`
}
