package utils

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/system"
	"github.com/golang/glog"
)

// ImageToDir converts an image to an unpacked tar and creates a representation of that directory.
func ImageToDir(img string, eng bool) (string, string, error) {
	var tarName string

	if !CheckTar(img) {
		// If not an image tar already existing in the filesystem, create client to obtain image
		// check client compatibility with Docker API
		valid, err := ValidDockerVersion(eng)
		if err != nil {
			return "", "", err
		}
		var imageTar string
		if !valid {
			glog.Info("Docker version incompatible with api, shelling out to local Docker client.")
			imageTar, err = imageToTarCmd(img)
		} else {
			imageTar, err = saveImageToTar(img)
		}
		if err != nil {
			return "", "", err
		}
		tarName = imageTar
		defer os.Remove(tarName)
	} else {
		tarName = img
	}
	return TarToJSON(tarName)
}

// ImageToTar writes an image to a .tar file
func ImageToTar(cli client.APIClient, image, tarName string) (string, error) {
	imgBytes, err := cli.ImageSave(context.Background(), []string{image})
	if err != nil {
		return "", err
	}
	defer imgBytes.Close()
	newpath := tarName + ".tar"
	return newpath, copyToFile(newpath, imgBytes)
}

func CheckImageID(image string) bool {
	pattern := regexp.MustCompile("[a-z|0-9]{12}")
	if exp := pattern.FindString(image); exp != image {
		return false
	}
	return true
}

func CheckImageURL(image string) bool {
	pattern := regexp.MustCompile("^.+/.+(:.+){0,1}$")
	if exp := pattern.FindString(image); exp != image || CheckTar(image) {
		return false
	}
	return true
}

func saveImageToTar(image string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return "", err
	}

	fromImage := image
	toTar := image
	// If not an already existing image ID, have to pull it from a repo before saving it
	if !CheckImageID(image) {
		imageID, imageName, err := pullImageFromRepo(cli, image)
		if err != nil {
			return "", err
		}
		fromImage = imageID
		toTar = imageName
	}
	// Convert the image into a tar
	imageTarPath, err := ImageToTar(cli, fromImage, toTar)
	if err != nil {
		return "", err
	}
	return imageTarPath, nil
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
