package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

type PR struct {
	Identifier string `json:"identifier"`
	Title      string `json:"title"`
	URL        string `json:"url"`
}

func main() {
	userName, err := getUserName(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.DefaultGraphQLClient()
	if err != nil {
		log.Fatal(err)
	}

	var query struct {
		Search struct {
			Nodes []struct {
				PullRequest struct {
					Title      string
					Url        string
					Number     int
					Repository struct {
						Name string
						Url  string
					}
				} `graphql:"... on PullRequest"`
			}
			IssueCount int
		} `graphql:"search(first: 100, type: ISSUE, query: $query)"`
	}
	variables := map[string]interface{}{
		"query": graphql.String(fmt.Sprintf("is:open is:pr archived:false user:%s", userName)),
	}
	if err := client.Query("search", &query, variables); err != nil {
		log.Fatal(err)
	}
	prs := make([]PR, 0, len(query.Search.Nodes))
	for _, node := range query.Search.Nodes {
		identifier := fmt.Sprintf("%s#%d", node.PullRequest.Repository.Name, node.PullRequest.Number)
		prs = append(prs, PR{
			Identifier: identifier,
			Title:      node.PullRequest.Title,
			URL:        node.PullRequest.Repository.Url,
		})
	}
	if err := json.NewEncoder(os.Stdout).Encode(prs); err != nil {
		log.Fatal(err)
	}
}

func getUserName(ctx context.Context) (string, error) {
	client, err := api.DefaultRESTClient()
	if err != nil {
		return "", err
	}
	response := struct{ Login string }{}
	err = client.Get("user", &response)
	if err != nil {
		return "", err
	}
	return response.Login, nil
}
