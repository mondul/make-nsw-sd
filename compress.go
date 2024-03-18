package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"os"
)

func compress(file_path string) {
	file_contents, err := os.ReadFile(file_path)
	if err != nil {
		fmt.Printf("! Could not load file '%s': %s\n", file_path, err)
		os.Exit(1)
	}

	var buf bytes.Buffer
	flate_writer, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		fmt.Printf("! Could not create compressor: %s\n", err)
		os.Exit(1)
	}

	flate_writer.Write(file_contents)
	flate_writer.Close()

	for i, byte := range buf.Bytes() {
		fmt.Printf(" 0x%02X,", byte)
		if ((i + 1) % 8) == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")

	os.Exit(0)
}
