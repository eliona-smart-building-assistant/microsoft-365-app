package msgraph

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"template/conf"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/eliona-smart-building-assistant/go-utils/log"
)

type Proxy struct {
}

type Response struct {
	ConfigID int64       `json:"config_id"`
	Username string      `json:"username"`
	Code     int         `json:"code"`
	Body     interface{} `json:"body"`
}

func (proxy *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	configs, err := conf.GetConfigs(context.Background())
	if err != nil {
		log.Fatal("conf", "Couldn't read configs from DB: %v", err)
		return
	}
	if len(configs) == 0 {
		return
	}

	var responses []Response
	for _, config := range configs {
		graph := NewGraphHelper()
		if config.ClientSecret == nil || config.Username == nil || config.Password == nil {
			log.Error("conf", "Shouldn't happen: some values are nil")
			return
		}
		if err := graph.InitializeGraphForUserAuth(config.ClientId, config.TenantId, *config.ClientSecret, *config.Username, *config.Password); err != nil {
			log.Error("ms-graph", "initializing graph for user auth: %v", err)
			return
		}

		requestURL := "https://graph.microsoft.com/v1.0/" + r.URL.Path
		log.Info("ms-graph", requestURL)

		graphReq, err := http.NewRequest(r.Method, requestURL, r.Body)
		if err != nil {
			msg := fmt.Sprintf("Error creating request to Microsoft Graph API %s: %v", requestURL, err.Error())
			http.Error(w, msg, http.StatusInternalServerError)
			return
		}

		// Copy the headers from the original request.
		for name, values := range r.Header {
			for _, value := range values {
				graphReq.Header.Add(name, value)
			}
		}

		// Refers to all permissions
		scopes := []string{"https://graph.microsoft.com/.default"}
		token, err := graph.credential.GetToken(context.Background(), policy.TokenRequestOptions{Scopes: scopes})
		if err != nil {
			http.Error(w, "Error getting bearer token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		graphReq.Header.Add("Authorization", "Bearer "+token.Token)

		graphRes, err := http.DefaultClient.Do(graphReq)
		if err != nil {
			http.Error(w, "Error sending request to Microsoft Graph API: "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer graphRes.Body.Close()

		body, err := io.ReadAll(graphRes.Body)
		if err != nil {
			http.Error(w, "Error reading body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		var b interface{}
		if err := json.Unmarshal(body, &b); err != nil {
			http.Error(w, "Error parsing body: "+err.Error(), http.StatusInternalServerError)
			return
		}
		responses = append(responses, Response{
			ConfigID: *config.Id,
			Username: *config.Username,
			Code:     graphRes.StatusCode,
			Body:     b,
		})
	}
	w.WriteHeader(http.StatusOK)

	responseJSON, err := json.Marshal(responses)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(responseJSON); err != nil {
		http.Error(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
