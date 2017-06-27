package differs

import "errors"

var diffs = map[string]func(string, string) (string, error){
	"hist": History,
	"dir":  Package,
}

func Diff(arg1, arg2, differ string) (string, error) {
	if f, exists := diffs[differ]; exists {
		return f(arg1, arg2)
	} else {
		return "", errors.New("Unknown differ.")
	}
}
