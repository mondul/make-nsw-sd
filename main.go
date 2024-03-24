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
	// Flags
	do_bootdat := false
	do_lockpick := false
	do_dbi := false

	// Check command line args
	for i, arg := range os.Args {
		if arg == "--with-bootdat" {
			do_bootdat = true
		} else if arg == "--with-lockpick" {
			do_lockpick = true
		} else if arg == "--with-dbi" {
			do_dbi = true
		} else if arg == "--compress" {
			compress(os.Args[i+1])
		}
	}

	// We'll use this folder for all downloaded files
	os.MkdirAll(workdir, os.ModePerm)

	var assets []*string

	// Download latest Atmosphère release
	repo := "Atmosphere-NX/Atmosphere"
	assets, err := getLatestAssets(repo, regexp.MustCompile(`\.zip$`))
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}
	atmosphere_zipfile := assets[0]

	// Download latest Hekate release
	repo = "CTCaer/hekate"
	assets, err = getLatestAssets(repo, regexp.MustCompile(`hekate_ctcaer.+\.zip$`))
	if err != nil {
		fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		os.Exit(1)
	}
	hekate_zipfile := assets[0]

	// Download SX-Gear boot.dat and config to launch Hekate
	var bootdat_zipfile *string = nil
	if do_bootdat {
		bootdat_zipfile, err = getBootDat()
		if err != nil {
			fmt.Printf("! Could not get SX Gear boot files: %s\n", err)
		}
	}

	// Download latest SPs
	sps_zipfile, err := getLatestSPs()
	if err != nil {
		fmt.Printf("! Could not get SPs: %s\n", err)
	}

	// Download latest Lockpick_RCM release
	var lockpick_bin *string = nil
	if do_lockpick {
		repo = "Mirror/Lockpick_RCM"
		assets, err = getLatestAssets(repo, regexp.MustCompile(`\.bin$`), "git.gdm.rocks/api/v1")
		if err != nil {
			fmt.Printf("! Could not get latest %s asset: %s\n", repo, err)
		}
		lockpick_bin = assets[0]
	}

	// Download latest DBI
	var dbi_files []*string
	if do_dbi {
		repo = "rashevskyv/dbi"
		dbi_files, err = getLatestAssets(repo, regexp.MustCompile(`((dbi\.config)|(DBI\.nro))$`))
		if err != nil {
			fmt.Printf("! Could not get latest %s assets: %s\n", repo, err)
			os.Exit(1)
		}
	}

	fmt.Println("-------")

	// Create output dir name
	outdir := fmt.Sprintf("SD_%X", time.Now().Unix())

	// Extract Atmosphère
	fmt.Printf("Extracting %s... ", *atmosphere_zipfile)
	if err = extractZip(*atmosphere_zipfile, outdir); err != nil {
		fmt.Printf("\n! Could not extract %s: %s\n", *atmosphere_zipfile, err)
		os.Exit(1)
	}
	fmt.Println("Done")

	// Extract Hekate
	fmt.Printf("Extracting %s... ", *hekate_zipfile)
	if err = extractZip(*hekate_zipfile, outdir, "hekate_ctcaer"); err != nil {
		fmt.Printf("\n! Could not extract %s: %s\n", *hekate_zipfile, err)
		os.Exit(1)
	}
	fmt.Println("Done")

	// Extract SX Gear boot files
	if do_bootdat && (bootdat_zipfile != nil) {
		fmt.Print("Extracting SX Gear boot files... ")
		if err = extractZip(*bootdat_zipfile, outdir); err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", *bootdat_zipfile, err)
		} else {
			fmt.Println("Done")
		}
	}

	// Extract SPs
	if sps_zipfile != nil {
		fmt.Printf("Extracting %s... ", *sps_zipfile)
		if err = extractZip(*sps_zipfile, outdir); err != nil {
			fmt.Printf("\n! Could not extract %s: %s\n", *sps_zipfile, err)
		} else {
			fmt.Println("Done")
		}
	}

	// Prevent ban
	fmt.Print("Creating ban prevention files... ")
	if err = preventBan(outdir); err != nil {
		fmt.Printf("\n! Could not create files: %s\n", err)
	} else {
		fmt.Println("Done")
	}

	// Move Lockpick_RCM.bin
	if do_lockpick && (lockpick_bin != nil) {
		fmt.Print("Moving Lockpick_RCM to payloads... ")
		if err = os.Rename(
			*lockpick_bin,
			filepath.Join(outdir, "bootloader", "payloads", "Lockpick_RCM.bin"),
		); err != nil {
			fmt.Printf("\n! Could not move Lockpick_RCM: %s\n", err)
		} else {
			fmt.Println("Done")
		}
	}

	// Move DBI files
	if do_dbi && (len(dbi_files) > 0) {
		fmt.Print("Moving DBI files... ")

		dbi_no_errors := true
		dbi_folder := filepath.Join(outdir, "switch", "DBI")
		os.MkdirAll(dbi_folder, os.ModePerm)

		for _, dbi_file := range dbi_files {
			dest_filename := filepath.Base(*dbi_file)

			if err = os.Rename(
				*dbi_file,
				filepath.Join(dbi_folder, dest_filename),
			); err != nil {
				dbi_no_errors = false
				fmt.Printf("\n! Could not move %s: %s\n", dest_filename, err)
			}
		}

		if dbi_no_errors {
			fmt.Println("Done")
		}
	}

	// Extract bootlogo if found
	boot_logo_zip := filepath.Join(workdir, "bootlogo.zip")
	if _, err := os.Stat(boot_logo_zip); err == nil {
		fmt.Print("Extracting custom boot logo... ")
		if err = extractZip(boot_logo_zip, filepath.Join(outdir, "atmosphere", "exefs_patches")); err != nil {
			fmt.Printf("\n! Could not extract boot logo: %s\n", err)
		} else {
			fmt.Println("Done")
		}
	}
}
