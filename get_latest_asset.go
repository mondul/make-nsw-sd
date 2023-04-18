package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type GitHubAsset struct {
	BrowserDownloadUrl string `json:"browser_download_url"`
}

type GitHubResponse struct {
	TagName string `json:"tag_name"`
	Assets  []GitHubAsset
}

/**
 * Gets info on a GitHub's repo latest release
 * @param  string repo            Must be formatted as {author}/{repo}
 * @return *string, error
 */
func getLatestAsset(repo string, prefix string) (*string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/"+repo+"/releases/latest", nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept":               {"application/vnd.github+json"},
		"X-GitHub-Api-Version": {"2022-11-28"},
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	decoder := jsoniter.NewDecoder(res.Body)

	var response GitHubResponse

	err = decoder.Decode(&response)

	if err != nil {
		return nil, err
	}

	fmt.Printf("* %s latest release: %s\n", repo, response.TagName)

	var asset_url string

	for _, asset := range response.Assets {
		if strings.Contains(asset.BrowserDownloadUrl, prefix) {
			asset_url = asset.BrowserDownloadUrl
			break
		}
	}

	filename := path.Base(asset_url)
	filename, _ = url.QueryUnescape(filename)

	// Download if not exists
	if _, err := os.Stat(filename); err == nil {
		fmt.Printf("- %s already exists\n", filename)
	} else {
		fmt.Printf("Downloading %s... ", filename)
		if err = downloadFile(filename, asset_url); err != nil {
			fmt.Printf("\n! Could not download %s: %s\n", filename, err)
		} else {
			fmt.Println("Done")
		}
	}

	return &filename, nil
}
