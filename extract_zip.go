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
 * @param  string  filename
 * @param  string  outdir
 * @param  *string prefix_to_skip
 * @return error
 */
func extractZip(filename string, outdir string, prefix_to_skip *string) error {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	defer archive.Close()

	for _, file := range archive.File {
		if prefix_to_skip != nil && strings.HasPrefix(file.Name, *prefix_to_skip) {
			continue
		}

		extract_path := filepath.Join(outdir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(extract_path, os.ModePerm)
			continue
		}

		src_file, err := file.Open()
		if err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", file.Name, err)
			continue
		}
		defer src_file.Close()

		dst_file, err := os.Create(extract_path)
		if err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", file.Name, err)
			continue
		}
		defer dst_file.Close()

		_, err = dst_file.ReadFrom(src_file)
		if err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", file.Name, err)
		}
	}

	return nil
}
