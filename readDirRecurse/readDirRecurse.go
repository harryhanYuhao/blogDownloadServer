package readDirRecurse

import "os"

var ret = make([]string, 0)

func recurseUtil(path string) error {
	files, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	if path == "./" {
		path = ""
	}

	for _, file := range files {
		// ignore hidden files
		if file.Name()[0] == '.' {
			continue
		}
		if file.IsDir() {
			subpath := path + file.Name() + "/"
			err := recurseUtil(subpath)
			if err != nil {
				return err
			}
		} else {
			// ignore readme files
			if file.Name() == "readme.md" {
				continue
			}
			ret = append(ret, path+file.Name())
		}
	}
	return nil
}

func ReadDirRecurse(path string) ([]string, error) {
	ret = make([]string, 0)
	err := recurseUtil(path)
	if err != nil {
		return nil, err
	}
	return ret, nil
}
