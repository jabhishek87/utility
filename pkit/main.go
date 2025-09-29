package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	version    = "dev"
	configPath string
	config     *Config
)

type Config struct {
	App struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"app"`
	Settings struct {
		Debug        bool   `yaml:"debug"`
		OutputFormat string `yaml:"output_format"`
	} `yaml:"settings"`
	Google struct {
		Auth struct {
			Scopes          []string `yaml:"scopes"`
			CredentialsFile string   `yaml:"credentials_file"`
			TokenFile       string   `yaml:"token_file"`
		} `yaml:"auth"`
		SpamEmailsFile string `yaml:"spam_emails_file"`
	} `yaml:"google"`
}

func loadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	return &config, err
}

var rootCmd = &cobra.Command{
	Use:     "pkit",
	Short:   "Personal CLI toolkit",
	Version: version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		var err error
		config, err = loadConfig(configPath)
		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			os.Exit(1)
		}
	},
}

var deleteSpamCmd = &cobra.Command{
	Use:   "delete-spam",
	Short: "Delete unread messages from spam emails",
	Run: func(cmd *cobra.Command, args []string) {
		authService := NewAuthService(config)
		client, err := authService.GetClient(context.Background())
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			os.Exit(1)
		}
		err = trashSpamMessages(client, config)
		if err != nil {
			fmt.Printf("Error trashing spam messages: %v\n", err)
			os.Exit(1)
		}
	},
}

var downloadDriveCmd = &cobra.Command{
	Use:   "download-drive [folder-link]",
	Short: "Download files from Google Drive folder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authService := NewAuthService(config)
		client, err := authService.GetClient(context.Background())
		if err != nil {
			fmt.Printf("Authentication failed: %v\n", err)
			os.Exit(1)
		}
		err = downloadDriveFolder(client, args[0])
		if err != nil {
			fmt.Printf("Error downloading folder: %v\n", err)
			os.Exit(1)
		}
	},
}

var createTokenCmd = &cobra.Command{
	Use:   "create_token [filename]",
	Short: "Create Google API token file with specified filename",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		authService := NewAuthService(config)
		
		tok, err := authService.tokenFromFile()
		if err != nil || !tok.Valid() {
			_, err = authService.GetClient(context.Background())
			if err != nil {
				fmt.Printf("Authentication failed: %v\n", err)
				os.Exit(1)
			}
			tok, _ = authService.tokenFromFile()
		}
		
		err = authService.saveCredentialsWithToken(tok, args[0])
		if err != nil {
			fmt.Printf("Error creating credentials file: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Credentials file created: %s\n", args[0])
	},
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "path to config file")
	rootCmd.AddCommand(deleteSpamCmd)
	rootCmd.AddCommand(downloadDriveCmd)
	rootCmd.AddCommand(createTokenCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}