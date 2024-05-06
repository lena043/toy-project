package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/identitystore"
)

func main() {
	ctx := context.Background()
	// Load the Shared AWS Configuration (~/.aws/config) - sso profile used
	// ref) https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/
	// how to use github secret
	ConfigProfile := os.Getenv("SSO_PROFILE")
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithSharedConfigProfile(ConfigProfile))
	if err != nil {
		log.Fatal(err)
	}

	// Initialize IdentityStore client with the loaded configuration
	client := identitystore.NewFromConfig(cfg)

	// valid IdentityStoreId(SSO Directory ID)
	identityStoreId := os.Getenv("IDENTITY_STORE_ID")

	// Create a paginator to list users in the identity store
	paginator := identitystore.NewListUsersPaginator(client, &identitystore.ListUsersInput{
		IdentityStoreId: &identityStoreId,
	})

	// user list print
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatal(err)
		}

		for _, user := range output.Users {
			fmt.Println(*user.UserName)
		}
	}
}
