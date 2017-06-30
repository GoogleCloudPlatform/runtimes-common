package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/system"
)

// ImageToDir converts an image to an unpacked tar and creates a representation of that directory.
func ImageToDir(img string) (string, string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", "", err
	}
	tarPath, err := ImageToTar(cli, img)
	if err != nil {
		return "", "", err
	}
	err = ExtractTar(tarPath)
	if err != nil {
		return "", "", err
	}
	os.Remove(tarPath)
	path := strings.TrimSuffix(tarPath, filepath.Ext(tarPath))
	jsonPath := path + ".json"
	err = DirToJSON(path, jsonPath)
	if err != nil {
		return "", "", err
	}
	return jsonPath, path, nil
}

// ImageToTar writes an image to a .tar file
func ImageToTar(cli client.APIClient, image string) (string, error) {
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer imgBytes.Close()
	newpath := image + ".tar"
	return newpath, copyToFile(newpath, imgBytes)
}

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
