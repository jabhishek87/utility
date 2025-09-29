package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"gopkg.in/yaml.v3"
)

type SpamEmails struct {
	SpamEmails []string `yaml:"spam_emails"`
}

func loadSpamEmails(path string) (*SpamEmails, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var spamEmails SpamEmails
	err = yaml.Unmarshal(data, &spamEmails)
	return &spamEmails, err
}

func trashSpamMessages(client *http.Client, config *Config) error {
	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %v", err)
	}

	spamEmails, err := loadSpamEmails(config.Google.SpamEmailsFile)
	if err != nil {
		return fmt.Errorf("unable to load spam emails: %v", err)
	}

	user := "me"
	totalTrashed := 0
	for _, email := range spamEmails.SpamEmails {
		query := fmt.Sprintf("is:unread from:%s", email)

		r, err := srv.Users.Messages.List(user).Q(query).Do()
		if err != nil {
			fmt.Printf("Error searching messages from %s: %v\n", email, err)
			continue
		}

		if len(r.Messages) == 0 {
			fmt.Printf("No unread messages from %s\n", email)
			continue
		}

		trashed := 0
		for _, msg := range r.Messages {
			_, err = srv.Users.Messages.Trash(user, msg.Id).Do()
			if err != nil {
				fmt.Printf("Error trashing message %s from %s: %v\n", msg.Id, email, err)
				continue
			}
			trashed++
			totalTrashed++
		}

		fmt.Printf("Trashed %d unread messages from %s\n", trashed, email)
	}
	fmt.Printf("Total Trashed %d \n", totalTrashed)
	return nil
}
