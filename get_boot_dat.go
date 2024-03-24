package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func getBootDat() (*string, error) {
	filename := "sxgearboot.zip"
	file_path := filepath.Join(workdir, filename)

	// Download if not exists
	if _, err := os.Stat(file_path); err == nil {
		fmt.Printf("* %s already exists\n", filename)
	} else {
		fmt.Printf("* Downloading %s... ", filename)
		if err = downloadFile(file_path, "https://raw.githubusercontent.com/mondul/MakeNSWSD-GUI/main/"+filename); err != nil {
			fmt.Printf("\n! Could not download %s: %s\n", filename, err)
			return nil, err
		} else {
			fmt.Println("Done")
		}
	}

	return &file_path, nil
}
