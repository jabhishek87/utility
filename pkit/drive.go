package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func getFolderIDFromLink(driveLink string) string {
	re := regexp.MustCompile(`(?:id=|folders/)([a-zA-Z0-9_-]+)`)
	matches := re.FindStringSubmatch(driveLink)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func downloadDriveFolder(client *http.Client, folderLink string) error {
	ctx := context.Background()
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	folderID := getFolderIDFromLink(folderLink)
	if folderID == "" {
		return fmt.Errorf("invalid Google Drive folder link")
	}

	folder, err := srv.Files.Get(folderID).Do()
	if err != nil {
		return fmt.Errorf("unable to get folder info: %v", err)
	}

	fmt.Printf("Downloading folder: %s\n", folder.Name)
	return downloadFolderRecursive(srv, folderID, folder.Name)
}

func downloadFolderRecursive(srv *drive.Service, folderID, localPath string) error {
	err := os.MkdirAll(localPath, 0755)
	if err != nil {
		return fmt.Errorf("unable to create folder: %v", err)
	}

	query := fmt.Sprintf("'%s' in parents", folderID)
	r, err := srv.Files.List().Q(query).Do()
	if err != nil {
		return fmt.Errorf("unable to retrieve files: %v", err)
	}

	for _, file := range r.Files {
		if file.MimeType == "application/vnd.google-apps.folder" {
			subfolderPath := filepath.Join(localPath, file.Name)
			fmt.Printf("Entering folder: %s\n", file.Name)
			err = downloadFolderRecursive(srv, file.Id, subfolderPath)
			if err != nil {
				fmt.Printf("Error downloading folder %s: %v\n", file.Name, err)
			}
		} else {
			err = downloadFile(srv, file, localPath)
			if err != nil {
				fmt.Printf("Error downloading %s: %v\n", file.Name, err)
			}
		}
	}
	return nil
}

func downloadFile(srv *drive.Service, file *drive.File, folderPath string) error {
	filePath := filepath.Join(folderPath, file.Name)
	
	if _, err := os.Stat(filePath); err == nil {
		fmt.Printf("Skipping existing file: %s\n", file.Name)
		return nil
	}

	resp, err := srv.Files.Get(file.Id).Download()
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err == nil {
		fmt.Printf("Downloaded: %s\n", file.Name)
	}
	return err
}