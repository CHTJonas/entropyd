package entropy

type EntropyClient struct {
	serverURL string
	minBits   int
	maxBits   int
}

func NewClient(serverURL string, minBits int, maxBits int) *EntropyClient {
	return &EntropyClient{
		serverURL: serverURL,
		minBits:   minBits,
		maxBits:   maxBits,
	}
}
