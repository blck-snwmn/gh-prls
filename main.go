package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/cli/go-gh/v2/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
)

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
	var builder strings.Builder
	for _, node := range query.Search.Nodes {
		identifier := fmt.Sprintf("%s#%d", node.PullRequest.Repository.Name, node.PullRequest.Number)
		builder.WriteString(fmt.Sprintf("- %s: %s\n", identifier, node.PullRequest.Title))
		builder.WriteString(fmt.Sprintf("    - %s \n", node.PullRequest.Url))
	}
	renderer, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(100))
	out, err := renderer.Render(builder.String())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(out)
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
