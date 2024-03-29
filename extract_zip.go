package main

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

/**
 * Extracts a zip file
 * @param  string    filename
 * @param  string    outdir
 * @param  ...string prefix Prefix to skip
 * @return error
 */
func extractZip(filename string, outdir string, prefix ...string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer archive.Close()

	check_prefix := len(prefix) > 0

	for _, file := range archive.File {
		if check_prefix && strings.HasPrefix(file.Name, prefix[0]) {
			continue
		}

		extract_path := filepath.Join(outdir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extract_path, os.ModePerm)
			continue
		}

		src_file, err := file.Open()
		if err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", file.Name, err))
			continue
		}
		defer src_file.Close()

		dst_file, err := os.Create(extract_path)
		if err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", file.Name, err))
			continue
		}
		defer dst_file.Close()

		_, err = dst_file.ReadFrom(src_file)
		if err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", file.Name, err))
		}
	}

	return nil
}
