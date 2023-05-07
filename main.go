package main

import (
	"context"
	"fmt"
	"log"

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
	fmt.Println(query.Search.IssueCount)
	for _, node := range query.Search.Nodes {
		fmt.Println(node.PullRequest.Title)
		fmt.Println(node.PullRequest.Url)
		fmt.Println(node.PullRequest.Repository.Name)
		fmt.Println(node.PullRequest.Repository.Url)
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
