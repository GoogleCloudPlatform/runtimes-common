package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/golang/glog"
)

// ValidDockerVersion determines if there is a Docker client of the correct version locally installed.
func ValidDockerVersion() (bool, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return false, fmt.Errorf("Docker client error: %s", err)
	}
	version := cli.ClientVersion()
	if version == "1.31" {
		return true, nil
	}
	return false, nil
}

func GetImageHistory(img string) ([]image.HistoryResponseItem, error) {
	imageID := img
	var err error
	var history []image.HistoryResponseItem
	if !CheckImageID(img) {
		imageID, _, err = pullImageCmd(img)
		if err != nil {
			return history, err
		}
	}
	histArgs := []string{"history", "--no-trunc", imageID}
	dockerHistCmd := exec.Command("docker", histArgs...)
	var response bytes.Buffer
	dockerHistCmd.Stdout = &response
	if err := dockerHistCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() > 0 {
			}
		} else {
			return history, err
		}
	}
	history, err = processHistOutput(response)
	if err != nil {
		return history, err
	}
	return history, nil

}

func processHistOutput(response bytes.Buffer) ([]image.HistoryResponseItem, error) {
	respReader := bytes.NewReader(response.Bytes())
	reader := bufio.NewReader(respReader)
	var history []image.HistoryResponseItem
	var CreatedByIndex int
	var SizeIndex int
	for {
		var event image.HistoryResponseItem
		text, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return history, err
		}

		line := string(text)
		if CreatedByIndex == 0 {
			CreatedByIndex = strings.Index(line, "CREATED BY")
			SizeIndex = strings.Index(line, "SIZE")
			continue
		}
		event.CreatedBy = line[CreatedByIndex:SizeIndex]
		history = append(history, event)
	}
	return history, nil
}

func processPullCmdOutput(image string, response bytes.Buffer) (string, string, error) {
	respReader := bytes.NewReader(response.Bytes())
	reader := bufio.NewReader(respReader)

	var events []Event
	for {
		var event Event
		text, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", "", err
		}
		event.Status = string(text)
		events = append(events, event)
	}

	imageDigest, err := getImagePullResponse(image, events)
	if err != nil {
		return "", "", err
	}

	URLPattern := regexp.MustCompile("^.+/(.+(:.+){0,1})$")
	URLMatch := URLPattern.FindStringSubmatch(image)
	imageName := strings.Split(URLMatch[1], ":")[0]
	imageURL := strings.TrimSuffix(image, URLMatch[2])
	imageID := imageURL + "@" + imageDigest

	return imageID, imageName, nil
}

func pullImageCmd(image string) (string, string, error) {
	pullArgs := []string{"pull", image}
	dockerPullCmd := exec.Command("docker", pullArgs...)
	var response bytes.Buffer
	dockerPullCmd.Stdout = &response
	if err := dockerPullCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() > 0 {
			}
		} else {
			return "", "", err
		}
	}
	return processPullCmdOutput(image, response)
}

func imageToTarCmd(image string) (string, error) {
	imageName := image
	imageID := image
	var err error
	// If not an already existing image ID, assuming URL, have to pull it from a repo before saving it
	if !CheckImageID(image) {
		imageID, imageName, err = pullImageCmd(image)
		if err != nil {
			return "", err
		}
	}

	// Convert the image into a tar
	cmdArgs := []string{"save", imageID}
	dockerSaveCmd := exec.Command("docker", cmdArgs...)
	var out bytes.Buffer
	dockerSaveCmd.Stdout = &out
	if err := dockerSaveCmd.Run(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok && status.ExitStatus() > 0 {
				glog.Error("Docker Save Command Exit Status: ", status.ExitStatus())
			}
		} else {
			return "", err
		}
	}
	imageTarPath := imageName + ".tar"
	reader := bytes.NewReader(out.Bytes())
	err = copyToFile(imageTarPath, reader)
	if err != nil {
		return "", err
	}
	return imageTarPath, nil
}
