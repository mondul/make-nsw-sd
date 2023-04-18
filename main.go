package main

import (
	"fmt"
	"os"
	"time"
)

/**
 * Program entry point
 */
func main() {
	// Download latest Atmosphère release
	repo := "Atmosphere-NX/Atmosphere"
	atmosphere_zipfile, err := getLatestAsset(repo, "atmosphere")
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}

	// Download latest Hekate release
	repo = "CTCaer/hekate"
	prefix := "hekate_ctcaer"
	hekate_zipfile, err := getLatestAsset(repo, prefix)
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}

	// Download latest SPs
	sps_zipfile, err := getLatestSPs()
	if err != nil {
		fmt.Printf("! Could not get SPs: %s\n", err)
		os.Exit(1)
	}

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
	fmt.Printf("Extracting %s... ", *sps_zipfile)
	if err = extractZip(*sps_zipfile, outdir, nil); err != nil {
		fmt.Printf("\n! Could not extract %s: %s\n", *sps_zipfile, err)
		os.Exit(1)
	}
	fmt.Println("Done")

	// Prevent ban
	fmt.Print("Creating ban prevention files... ")
	if err = preventBan(outdir); err != nil {
		fmt.Printf("\n! Could not create files: %s\n", err)
	} else {
		fmt.Println("Done")
	}
}
