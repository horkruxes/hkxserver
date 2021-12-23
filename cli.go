package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func help() {
	fmt.Println("Usage : hkxserver help / version / update / run")
}

func version() {
	fmt.Println("hkxserver v0.7.7")
}

func downloadAndSaveLastVersion() error {
	fmt.Println("Downloading latest version")
	resp, err := http.Get("https://github.com/horkruxes/hkxserver/releases/latest/download/hkxserver_linux_amd64.tar.gz")
	if err != nil {
		fmt.Println("Can't fetch latest upgrade online")
		return err
	}
	defer resp.Body.Close()

	err = Untar(resp.Body)
	if err != nil {
		fmt.Println("Can't untar the downloaded doc")
		return err
	}
	fmt.Println("Downloaded latest version. Please restart server.")

	return nil

}

// Untar takes a destination path and a reader; a tar reader loops over the tarfile
// creating the file structure at 'dst' along the way, and writing any files
func Untar(r io.Reader) error {

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Clean(header.Name)
		fmt.Println(target)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0750); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			//#nosec G110 -- Copying trusted content (downloading the new executable from the github CI)
			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			if err := f.Close(); err != nil {
				return err
			}
		}
	}
}
