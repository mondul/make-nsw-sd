package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"slices"

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
 * Gets files from a GitHub's repo latest release according to a regex filter
 * @param  string         repo         Must be formatted as {author}/{repo}
 * @param  *regexp.Regexp filter_regex Regex filter for the name of the asset to be downloaded
 * @param ...string       api_url      Custom API URL if it's not for GitHub
 * @return *string, error
 */
func getLatestAssets(repo string, filter_regex *regexp.Regexp, api_url ...string) ([]*string, error) {
	base_url := "api.github.com"
	no_gh := len(api_url) > 0

	if no_gh {
		base_url = api_url[0]
	}

	req, err := http.NewRequest(http.MethodGet, "https://"+base_url+"/repos/"+repo+"/releases?per_page=1", nil)
	if err != nil {
		return nil, err
	}

	req.Header = http.Header{
		"Accept":               {"application/vnd.github+json"},
		"X-GitHub-Api-Version": {"2022-11-28"},
	}

	if no_gh {
		req.Header = http.Header{
			"Accept": {"application/json"},
		}
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

	var response []GitHubResponse

	err = decoder.Decode(&response)

	if err != nil {
		return nil, err
	}

	fmt.Printf("* %s latest release: %s\n", repo, response[0].TagName)

	var file_paths []*string

	for _, asset := range response[0].Assets {
		if filter_regex.MatchString(asset.BrowserDownloadUrl) {
			filename, _ := url.QueryUnescape(path.Base(asset.BrowserDownloadUrl))
			file_path := filepath.Join(workdir, filename)

			// Download if not exists
			if _, err := os.Stat(file_path); err == nil {
				fmt.Printf("- %s already exists\n", filename)
			} else {
				fmt.Printf("Downloading %s... ", filename)
				if err = downloadFile(file_path, asset.BrowserDownloadUrl); err != nil {
					fmt.Printf("\n! Could not download %s: %s\n", filename, err)
					return nil, err
				} else {
					fmt.Println("Done")
				}
			}

			file_paths = slices.Insert(file_paths, 0, &file_path)
		}
	}

	return file_paths, nil
}
