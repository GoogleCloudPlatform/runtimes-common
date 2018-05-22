package hack

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

type Test struct {
	cmd  string
	name string
}

var tests = []Test{
	{
		cmd:  "python ../ftl/integration_tests/ftl_node_integration_tests_yaml.py",
		name: "ftl node",
	},
	{
		cmd:  "python ../ftl/integration_tests/ftl_php_integration_tests_yaml.py",
		name: "ftl php",
	},
	{
		cmd:  "python ../ftl/integration_tests/ftl_python_integration_tests_yaml.py",
		name: "ftl python",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=node-same",
		name: "ftl cached node same",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=node-plus-one",
		name: "ftl cached node plus one",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=php-lock-plus-one",
		name: "ftl cached php plus one",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=python-requirements-same",
		name: "ftl cached python requirements same",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=python-requirements-plus-one",
		name: "ftl cached requirements plus one",
	},

	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=python-pipfile-plus-one",
		name: "ftl cached pipfile plus one",
	},
	{
		cmd:  "python ../ftl/cached/ftl_cached_yaml.py --runtime=python-pipfile-same",
		name: "ftl cached pipfile same",
	},
	// {
	// 	cmd:  "cat ../ftl/integration_tests/ftl_python_error_test.yaml",
	// 	name: "ftl cached python error",
	// },
}

func command(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return "", err
	}
	out, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return string(out), err
	}
	return string(out), nil
}

func TestAll(t *testing.T) {
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			out, err := command(fmt.Sprintf("%s | gcloud container builds submit --config=/dev/stdin ../", tc.cmd))
			if err != nil {
				t.Logf(out)
				t.Fatalf("error running build: %s", err)
			}
		})
	}
}
