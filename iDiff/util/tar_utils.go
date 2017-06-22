package tarUtil

import (
	"archive/tar"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/system"
)

// copyToFile writes the content of the reader to the specified file
func copyToFile(outfile string, r io.Reader) error {
	// We use sequential file access here to avoid depleting the standby list
	// on Windows. On Linux, this is a call directly to ioutil.TempFile
	tmpFile, err := system.TempFileSequential(filepath.Dir(outfile), ".docker_temp_")
	if err != nil {
		return err
	}

	tmpPath := tmpFile.Name()

	_, err = io.Copy(tmpFile, r)
	tmpFile.Close()

	if err != nil {
		os.Remove(tmpPath)
		return err
	}

	if err = os.Rename(tmpPath, outfile); err != nil {
		os.Remove(tmpPath)
		return err
	}

	return nil
}

// ImageToTar writes an image to a .tar file
func ImageToTar(cli client.APIClient, image string) error {
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return err
	}
	defer imgBytes.Close()
	return copyToFile(image+".tar", imgBytes)
}

// Dir stores a representaiton of a file directory.
type Dir struct {
	Root    string
	Content []string
}

// UnTar takes in a path to a tar file and writes the untarred version to the provided target.
// Only untars one level, does not untar nested tars.
func UnTar(filename string, path string) error {
	if _, ok := os.Stat(path); ok != nil {
		os.MkdirAll(path, 0777)

	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	tr := tar.NewReader(file)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		target := filepath.Join(path, header.Name)
		mode := header.FileInfo().Mode()
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, mode); err != nil {
					return err
				}
				continue
			}

		// if it's a file create it
		case tar.TypeReg:

			currFile, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
			if err != nil {
				return err
			}
			defer currFile.Close()
			_, err = io.Copy(currFile, tr)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isTar(path string) bool {
	return filepath.Ext(path) == ".tar"
}

// ExtractTar extracts the tar and any nested tar at the given path.
// After execution the original tar file is removed and the untarred version is in it place.
func ExtractTar(path string) error {
	removeTar := false

	var untarWalkFn func(path string, info os.FileInfo, err error) error

	untarWalkFn = func(path string, info os.FileInfo, err error) error {
		if isTar(path) {
			target := strings.TrimSuffix(path, filepath.Ext(path))
			UnTar(path, target)
			if removeTar {
				os.Remove(path)
			}
			// remove nested tar files that get copied but not the original tar passed
			removeTar = true
			filepath.Walk(target, untarWalkFn)
		}
		return nil
	}

	return filepath.Walk(path, untarWalkFn)
}

// DirToJSON records the directory structure starting at the provided path as in a json file.
func DirToJSON(path string, target string) error {
	var directory Dir
	directory.Root = path

	tarJSONWalkFn := func(currPath string, info os.FileInfo, err error) error {
		newContent := strings.TrimPrefix(currPath, directory.Root)
		directory.Content = append(directory.Content, newContent)
		return nil
	}

	filepath.Walk(path, tarJSONWalkFn)
	data, err := json.Marshal(directory)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(target, data, 0777)
}
