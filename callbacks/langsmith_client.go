package callbacks

// LangSmithClient communicates with the LangSmith API.
// Full implementation requires the langsmith-go SDK.
type LangSmithClient struct {
	apiKey  string
	baseURL string
}

func NewLangSmithClient(apiKey, baseURL string) *LangSmithClient {
	if baseURL == "" {
		baseURL = "https://api.smith.langchain.com"
	}
	return &LangSmithClient{apiKey: apiKey, baseURL: baseURL}
}

// SendRun sends a run to LangSmith.
func (c *LangSmithClient) SendRun(runName, runType, input, output string, err error) {
	// Stub: requires langsmith-go SDK
}
