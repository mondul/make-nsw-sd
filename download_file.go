package main

import (
	"fmt"
	"net/http"
	"os"
)

/**
 * Downloads a file into the current directory
 * @param  string filename
 * @param  string url
 * @return error
 */
func downloadFile(filename string, url string) error {
	// Get the data
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	// Create the file
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	// Writer the body to file
	_, err = out.ReadFrom(res.Body)
	if err != nil {
		return err
	}

	return nil
}
