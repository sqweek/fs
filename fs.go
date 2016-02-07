package fs

import (
	"os"
	"path/filepath"
)

type Dirs struct {
	Cache string
	Data string
	Config string
	Docs string
	ResourcePath []string
}

// home is not exported as it is inappropriate to create files here on some platforms (windows)
var home string

func UserDirs() (*Dirs, error) {
	return userDirs()
}

func AppDirs(fqdn, name string) (*Dirs, error) {
	return appDirs(fqdn, name)
}

func (d Dirs) Locate(filename string) (string, error) {
	paths := make([]string, 0, len(d.ResourcePath) + 1)
	paths = append(paths, d.Data)
	paths = append(paths, d.ResourcePath...)
	return Locate(filename, paths...)
}

func Locate(filename string, paths... string) (string, error) {
	for _, path := range paths {
		f := filepath.Join(path, filename)
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			return f, nil
		}
	}
	return "", &os.PathError{Op: "locate", Path: filename, Err: os.ErrNotExist}
}

func ReplaceFile(src, dst string) error {
	return replaceFile(src, dst)
}

func Create(name string) (*os.File, error) {
	if f, err := os.Create(name); err == nil {
		return f, nil
	} else if os.IsNotExist(err) {
		dir, _ := filepath.Split(name)
		if err = os.MkdirAll(dir, 0777); err != nil {
			return nil, err
		}
		return os.Create(name)
	} else {
		return nil, err
	}
}

func CreateIn(dir, file string) (*os.File, error) {
	return Create(filepath.Join(dir, file))
}



func envs(tpls... string) string {
	for _, tpl := range tpls {
		if s := os.ExpandEnv(tpl); s != "" {
			return s
		}
	}
	return ""
}
