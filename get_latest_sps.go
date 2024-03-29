package main

import (
	"bufio"
	"bytes"
	"compress/flate"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var compressed_forum_url = []byte{
	0x04, 0xC0, 0x01, 0x0E, 0x82, 0x30, 0x0C, 0x05,
	0xD0, 0x13, 0x95, 0x6A, 0x90, 0x98, 0x78, 0x9B,
	0x8A, 0x7F, 0xAB, 0x21, 0xB8, 0x66, 0xFF, 0x7B,
	0x7F, 0x5E, 0x4A, 0xC5, 0x97, 0x7B, 0x7F, 0x87,
	0x70, 0xD6, 0xF2, 0x83, 0x5C, 0x39, 0x11, 0x1F,
	0x3A, 0xBF, 0xBD, 0x42, 0x7B, 0x82, 0xD6, 0xC6,
	0xB4, 0xD0, 0x39, 0x58, 0x89, 0x09, 0x4B, 0x1C,
	0x21, 0x58, 0x23, 0x6F, 0xD6, 0xFE, 0x04, 0xAC,
	0x62, 0x3F, 0xA2, 0x63, 0x5D, 0xB6, 0xE7, 0x7D,
	0x7B, 0xAC, 0x7E, 0x05, 0x00, 0x00, 0xFF, 0xFF,
}

func getLatestSPs() (*string, error) {
	var b bytes.Buffer
	r := flate.NewReader(bytes.NewReader(compressed_forum_url))
	b.ReadFrom(r)
	r.Close()

	forum_url := b.String()

	// Load forum post
	res, err := http.Get(forum_url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %s", res.Status)
	}

	scanner := bufio.NewScanner(res.Body)

	var download_url string
	var sps_filename string

	for scanner.Scan() {
		current_line := scanner.Text()

		text_start := strings.Index(current_line, "/attachments/hekate-ams")
		if text_start > 0 {
			text_end := strings.Index(current_line, "\" target")
			download_url = forum_url[:19] + current_line[text_start:text_end]
			continue
		}

		text_start = strings.Index(current_line, "Hekate+AMS")
		if text_start > 0 {
			text_end := strings.Index(current_line, "\">")
			sps_filename = current_line[text_start:text_end]
			break
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, err
	}

	sps_file_path := filepath.Join(workdir, sps_filename)

	// Download if not exists
	if _, err := os.Stat(sps_file_path); err == nil {
		log_add(fmt.Sprintf("* %s already exists\n", sps_filename))
	} else {
		log_add(fmt.Sprintf("* Downloading %s... ", sps_filename))
		if err = downloadFile(sps_file_path, download_url); err != nil {
			return nil, err
		} else {
			log_add("Done\n")
		}
	}

	return &sps_file_path, nil
}
