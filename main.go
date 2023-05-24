package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

const workdir string = "workdir"

/**
 * Program entry point
 */
func main() {
	// We'll use this folder for all downloaded files
	os.MkdirAll(workdir, os.ModePerm)

	// Download latest Atmosphère release
	repo := "Atmosphere-NX/Atmosphere"
	atmosphere_zipfile, err := getLatestAssets(repo, regexp.MustCompile(`\.zip$`))
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}

	// Download latest Hekate release
	repo = "CTCaer/hekate"
	prefix := "hekate_ctcaer"
	hekate_zipfile, err := getLatestAssets(repo, regexp.MustCompile(`hekate_ctcaer.+\.zip$`))
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}

	// Download latest SPs
	sps_zipfile, err := getLatestSPs()
	if err != nil {
		fmt.Printf("! Could not get SPs: %s\n", err)
	}

	// Download latest Lockpick_RCM release
	repo = "shchmue/Lockpick_RCM"
	lockpick_bin, err := getLatestAssets(repo, regexp.MustCompile(`\.bin$`))
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
	}

	fmt.Println("-------")

	// Create output dir name
	outdir := fmt.Sprintf("SD_%X", time.Now().Unix())

	// Extract Atmosphère
	fmt.Printf("Extracting %s... ", *atmosphere_zipfile)
	if err = extractZip(*atmosphere_zipfile, outdir, nil); err != nil {
		fmt.Printf("\n! Could not extract %s: %s\n", *atmosphere_zipfile, err)
		os.Exit(1)
	}
	fmt.Println("Done")

	// Extract Hekate
	fmt.Printf("Extracting %s... ", *hekate_zipfile)
	if err = extractZip(*hekate_zipfile, outdir, &prefix); err != nil {
		fmt.Printf("\n! Could not extract %s: %s\n", *hekate_zipfile, err)
		os.Exit(1)
	}
	fmt.Println("Done")

	// Extract SPs
	if sps_zipfile != nil {
		fmt.Printf("Extracting %s... ", *sps_zipfile)
		if err = extractZip(*sps_zipfile, outdir, nil); err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", *sps_zipfile, err)
		}
		fmt.Println("Done")
	}

	// Prevent ban
	fmt.Print("Creating ban prevention files... ")
	if err = preventBan(outdir); err != nil {
		fmt.Printf("\n! Could not create files: %s\n", err)
	} else {
		fmt.Println("Done")
	}

	// Move Lockpick_RCM.bin
	if lockpick_bin != nil {
		fmt.Print("Moving Lockpick_RCM to payloads... ")
		if err = os.Rename(
			*lockpick_bin,
			filepath.Join(outdir, "bootloader/payloads/", *lockpick_bin),
		); err != nil {
			fmt.Printf("\n! Could not move Lockpick_RCM: %s\n", err)
		} else {
			fmt.Println("Done")
		}
	}
}
