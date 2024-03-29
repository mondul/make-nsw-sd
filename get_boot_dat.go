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
		log_add(fmt.Sprintf("* %s already exists\n", filename))
	} else {
		log_add(fmt.Sprintf("* Downloading %s... ", filename))
		if err = downloadFile(file_path, "https://raw.githubusercontent.com/mondul/MakeNSWSD-GUI/main/"+filename); err != nil {
			log_add(fmt.Sprintf("\n! Could not download %s: %s\n", filename, err))
			return nil, err
		} else {
			log_add("Done\n")
		}
	}

	return &file_path, nil
}
