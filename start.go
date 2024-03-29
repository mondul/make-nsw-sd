package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

const workdir string = "workdir"

/**
 * Copies a file (why there's no os.Copy ???)
 * @param  src string Source file path
 * @param  dst string Destination file path
 * @return error
 */
func copyFile(src, dst string) error {
	src_file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer src_file.Close()

	_, err = os.Stat(dst)
	if err == nil {
		return errors.New("file already exists")
	}

	dst_file, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dst_file.Close()

	_, err = dst_file.ReadFrom(src_file)

	return err
}

/**
 * Runs the stuff
 * @param dos_type       dos               Processes to follow
 * @param binding.String folder_entry_data Output directory folder entry data
 * @param *widget.Button close_btn         Close button in log window, will be re-enabled after process is finished
 */
var start func(dos_type, binding.String, *widget.Button) = func(dos dos_type, folder_entry_data binding.String, close_btn *widget.Button) {
	// We'll use this folder for all downloaded files
	os.MkdirAll(workdir, os.ModePerm)

	// Download latest Atmosphère release
	var atmosphere_zipfile *string = nil

	if dos.atmosphere {
		repo := "Atmosphere-NX/Atmosphere"
		assets, err := getLatestAssets(repo, `\.zip$`)
		if err != nil {
			log_add(fmt.Sprintf("! Could not get latest %s asset: %s\n", repo, err))
			close_btn.Enable()
			return
		}
		atmosphere_zipfile = assets[0]
	}

	// Download latest Hekate release
	var hekate_zipfile *string = nil
	var bootdat_zipfile *string = nil
	var lockpick_bin *string = nil

	if dos.hekate {
		repo := "CTCaer/hekate"
		assets, err := getLatestAssets(repo, `hekate_ctcaer.+\.zip$`)
		if err != nil {
			log_add(fmt.Sprintf("! Could not get latest %s asset: %s\n", repo, err))
			close_btn.Enable()
			return
		}
		hekate_zipfile = assets[0]

		// Download SX-Gear boot.dat and config to launch Hekate
		if dos.bootdat {
			bootdat_zipfile, err = getBootDat()
			if err != nil {
				log_add(fmt.Sprintf("! Could not get SX Gear boot files: %s\n", err))
			}
		}

		// Download latest Lockpick_RCM release
		if dos.lockpick {
			repo = "Mirror/Lockpick_RCM"
			assets, err = getLatestAssets(repo, `\.bin$`, "git.gdm.rocks/api/v1")
			if err != nil {
				log_add(fmt.Sprintf("! Could not get latest %s asset: %s\n", repo, err))
			}
			lockpick_bin = assets[0]
		}
	}

	// Download latest SPs
	sps_zipfile, err := getLatestSPs()
	if err != nil {
		log_add(fmt.Sprintf("! Could not get SPs: %s\n", err))
	}

	// Download latest DBI
	var dbi_files []*string
	if dos.dbi {
		repo := "rashevskyv/dbi"
		dbi_files, err = getLatestAssets(repo, `((dbi\.config)|(DBI\.nro))$`)
		if err != nil {
			log_add(fmt.Sprintf("! Could not get latest %s assets: %s\n", repo, err))
		}
	}

	outdir, _ := folder_entry_data.Get()
	log_add(fmt.Sprintf("-------\nOutput directory: %s\n-------\n", outdir))

	// If output dir doesn't exist, create it
	os.MkdirAll(outdir, os.ModePerm)

	// Extract Atmosphère
	if dos.atmosphere {
		log_add(fmt.Sprintf("Extracting %s... ", filepath.Base(*atmosphere_zipfile)))
		if err = extractZip(*atmosphere_zipfile, outdir); err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", *atmosphere_zipfile, err))
			close_btn.Enable()
			return
		}
		log_add("Done\n")

		// Prevent ban
		log_add("Creating ban prevention files... ")
		if err = preventBan(outdir); err != nil {
			log_add(fmt.Sprintf("\n! Could not create files: %s\n", err))
		} else {
			log_add("Done\n")
		}

		// Extract bootlogo if found
		boot_logo_zip := filepath.Join(workdir, "bootlogo.zip")
		if _, err := os.Stat(boot_logo_zip); err == nil {
			log_add("Extracting custom boot logo... ")
			if err = extractZip(boot_logo_zip, filepath.Join(outdir, "atmosphere", "exefs_patches")); err != nil {
				log_add(fmt.Sprintf("\n! Could not extract boot logo: %s\n", err))
			} else {
				log_add("Done\n")
			}
		}
	}

	// Extract Hekate
	if dos.hekate {
		log_add(fmt.Sprintf("Extracting %s... ", filepath.Base(*hekate_zipfile)))
		if err = extractZip(*hekate_zipfile, outdir, "hekate_ctcaer"); err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", *hekate_zipfile, err))
			close_btn.Enable()
			return
		}
		log_add("Done\n")

		// Copy hekate payload.bin to output dir
		if dos.payload {
			log_add("Copying Hekate payload.bin... ")
			if err = copyFile(
				filepath.Join(outdir, "bootloader", "update.bin"),
				filepath.Join(outdir, "payload.bin"),
			); err != nil {
				log_add(fmt.Sprintf("\n! Could not create payload.bin: %s\n", err))
			} else {
				log_add("Done\n")
			}
		} else if dos.bootdat && bootdat_zipfile != nil { // Extract SX Gear boot files
			log_add("Extracting SX Gear boot files... ")
			if err = extractZip(*bootdat_zipfile, outdir); err != nil {
				log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", *bootdat_zipfile, err))
			} else {
				log_add("Done\n")
			}
		}

		// Move Lockpick_RCM.bin
		if dos.lockpick && lockpick_bin != nil {
			log_add("Moving Lockpick_RCM to payloads... ")
			if err = os.Rename(
				*lockpick_bin,
				filepath.Join(outdir, "bootloader", "payloads", "Lockpick_RCM.bin"),
			); err != nil {
				log_add(fmt.Sprintf("\n! Could not move Lockpick_RCM: %s\n", err))
			} else {
				log_add("Done\n")
			}
		}
	}

	// Extract SPs
	if dos.sps && sps_zipfile != nil {
		log_add(fmt.Sprintf("Extracting %s... ", filepath.Base(*sps_zipfile)))
		if err = extractZip(*sps_zipfile, outdir); err != nil {
			log_add(fmt.Sprintf("\n! Could not extract %s: %s\n", *sps_zipfile, err))
		} else {
			log_add("Done\n")
		}
	}

	// Move DBI files
	if dos.dbi && len(dbi_files) > 0 {
		log_add("Moving DBI files... ")

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
				log_add(fmt.Sprintf("\n! Could not move %s: %s\n", dest_filename, err))
			}
		}

		if dbi_no_errors {
			log_add("Done\n")
		}
	}

	// Set new output directory just in case
	folder_entry_data.Set(newOutdir())

	close_btn.Enable()
}
