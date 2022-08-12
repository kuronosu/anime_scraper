package scrape

import (
	"net/http"

	cloudflarebp "github.com/DaRealFreak/cloudflare-bp-go"
)

func GetCloudFlareRoundTripper() http.RoundTripper {
	client := &http.Client{}
	client.Transport = cloudflarebp.AddCloudFlareByPass(client.Transport)
	return client.Transport
}
