package main

import (
	"bytes"
	"compress/flate"
	"fmt"
	"io"
	"os"
)

func compress(file_path string) {
	file, err := os.Open(file_path)
	if err != nil {
		fmt.Printf("! Could not load file '%s': %s\n", file_path, err)
		os.Exit(1)
	}
	defer file.Close()

	var buf bytes.Buffer
	flate_writer, err := flate.NewWriter(&buf, flate.BestCompression)
	if err != nil {
		fmt.Printf("! Could not create compressor: %s\n", err)
		os.Exit(1)
	}
	defer flate_writer.Close()

	io.Copy(flate_writer, file)
	flate_writer.Flush()

	for i, byte := range buf.Bytes() {
		fmt.Printf(" 0x%02X,", byte)
		if ((i + 1) % 8) == 0 {
			fmt.Printf("\n")
		}
	}
	fmt.Printf("\n")

	os.Exit(0)
}
