package gitGithub

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pandaci-com/pandaci/pkg/utils/env"
)

func (c *graphqlClient) Query(ctx context.Context, dest any, query string, variables *map[string]interface{}) error {
	graphqlAPI := fmt.Sprintf("%s/graphql", env.GetGithubAPIEndpoint())

	type GraphQLRequest struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}

	reqBody := GraphQLRequest{
		Query: query,
	}

	reqBodyJSON, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(context.Background(), "POST", graphqlAPI, bytes.NewBuffer(reqBodyJSON))
	if err != nil {
		return err
	}

	token, err := c.itr.Token(ctx)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return json.Unmarshal(respBody, dest)
}
