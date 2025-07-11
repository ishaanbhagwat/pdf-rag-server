package main

import(
	"cmp"
	"context"
	"fmt"
	"os"

	"github.com/weaviate/weaviate-go-client/v5/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

func initWeaviate(ctx context.Context)(*weaviate.Client, error){
		client, err := weaviate.NewClient(weaviate.Config{
		Host:   "localhost:" + cmp.Or(os.Getenv("WVPORT"), "9035"),
		Scheme: "http",
	})

	if err != nil{
		return nil, fmt.Errorf("initializing weaviate: %w", err)
	}

	cls := &models.Class{
		Class: "Document",
		Vectorizer: "none",
	}

	exists, err := client.Schema().ClassExistenceChecker().WithClassName(cls.Class).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("weaviate error: %w", err)
	}
	if !exists {
		err = client.Schema().ClassCreator().WithClass(cls).Do(ctx)
		if err != nil {
			return nil, fmt.Errorf("weaviate error: %w", err)
		}
	}

	return client, nil
}

func combinedWeaviateError(result *models.GraphQLResponse, err error) error {
	if err != nil {
		return err
	}
	if len(result.Errors) != 0 {
		var ss []string
		for _, e := range result.Errors {
			ss = append(ss, e.Message)
		}
		return fmt.Errorf("weaviate error: %v", ss)
	}
	return nil
}