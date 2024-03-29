package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"

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
 * @param  string   repo         Must be formatted as {author}/{repo}
 * @param  string   filter_regex Regex filter for the name of the asset to be downloaded
 * @param ...string api_url      Custom API URL if it's not for GitHub
 * @return []*string, error
 */
func getLatestAssets(repo string, filter_regex string, api_url ...string) ([]*string, error) {
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

	if err = decoder.Decode(&response); err != nil {
		return nil, err
	}

	log_add(fmt.Sprintf("* %s latest release: %s\n", repo, response[0].TagName))

	file_paths := []*string{}

	re := regexp.MustCompile(filter_regex)

	for _, asset := range response[0].Assets {
		if re.MatchString(asset.BrowserDownloadUrl) {
			filename, _ := url.QueryUnescape(path.Base(asset.BrowserDownloadUrl))
			file_path := filepath.Join(workdir, filename)

			// Download if not exists
			if _, err := os.Stat(file_path); err == nil {
				log_add(fmt.Sprintf("- %s already exists\n", filename))
			} else {
				log_add(fmt.Sprintf("  Downloading %s... ", filename))
				if err = downloadFile(file_path, asset.BrowserDownloadUrl); err != nil {
					log_add(fmt.Sprintf("\n! Could not download %s: %s\n", filename, err))
					return nil, err
				} else {
					log_add("Done\n")
				}
			}

			file_paths = append(file_paths, &file_path)
		}
	}

	return file_paths, nil
}
