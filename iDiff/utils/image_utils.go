package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/system"
	"github.com/golang/glog"
)

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

type Event struct {
	Status         string `json:"status"`
	Error          string `json:"error"`
	Progress       string `json:"progress"`
	ProgressDetail struct {
		Current int `json:"current"`
		Total   int `json:"total"`
	} `json:"progressDetail"`
}

func getImagePullResponse(image string, response []Event) (string, error) {
	var imageDigest string
	for _, event := range response {
		if event.Error != "" {
			err := fmt.Errorf("Error pulling image %s: %s", image, event.Error)
			return "", err
		}
		digestPattern := regexp.MustCompile("^Digest: (sha256:[a-z|0-9]{64})$")
		digestMatch := digestPattern.FindStringSubmatch(event.Status)
		if len(digestMatch) != 0 {
			imageDigest = digestMatch[1]
			return imageDigest, nil
		}
	}
	err := fmt.Errorf("Could not pull image %s", image)
	return "", err
}

func processImagePullEvents(image string, events []Event) (string, string, error) {
	imageDigest, err := getImagePullResponse(image, events)
	if err != nil {
		return "", "", err
	}

	URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	URLMatch := URLPattern.FindStringSubmatch(image)
	imageName := strings.Replace(URLMatch[1], ":", "", -1)
	imageURL := strings.TrimSuffix(image, URLMatch[2])
	imageID := imageURL + "@" + imageDigest

	return imageID, imageName, nil
}

func pullImageFromRepo(cli client.APIClient, image string) (string, string, error) {
	response, err := cli.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return "", "", err
	}
	defer response.Close()

	d := json.NewDecoder(response)

	var events []Event
	for {
		var event Event
		if err := d.Decode(&event); err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		events = append(events, event)
	}
	return processImagePullEvents(image, events)
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
