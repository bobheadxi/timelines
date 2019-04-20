package gh

// Event denotes basic event metadata
type Event struct {
	GitHubEvent    string `json:"X-GitHub-Event"`
	GitHubDelivery string `json:"X-GitHub-Delivery"`
	HubSignature   string `json:"X-Hub-Signature"`
}
