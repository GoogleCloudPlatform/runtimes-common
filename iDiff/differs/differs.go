package differs

import "errors"

var diffs = map[string]func(string, string) (string, error){
	"hist": History,
	"dir":  Package,
}

func Diff(arg1, arg2, differ string) (string, error) {
	if f, exists := diffs[differ]; exists {
		return callDiffer(arg1, arg2, f)
	} else {
		return "", errors.New("Unknown differ.")
	}
}

func callDiffer(arg1, arg2 string, differ func(string, string) (string, error)) (string, error) {
	return differ(arg1, arg2)
}
