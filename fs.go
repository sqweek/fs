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

func (d Dirs) DataPath() []string {
	paths := make([]string, 0, len(d.ResourcePath) + 1)
	paths = append(paths, d.Data)
	paths = append(paths, d.ResourcePath...)
	return paths
}

func (d Dirs) Locate(filename string) (string, error) {
	return Locate(filename, d.DataPath()...)
}

func (d Dirs) Search(pattern string) (string, error) {
	return Search(pattern, d.DataPath()...)
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

func Search(pattern string, paths... string) (string, error) {
	for _, dir := range paths {
		if f, err := os.Open(dir); err == nil {
			defer f.Close()
			names, _ := f.Readdirnames(64)
			for len(names) > 0 {
				for _, name := range(names) {
					if match, _ := filepath.Match(pattern, name); match {
						return filepath.Join(dir, name), nil
					}
				}
				names, err = f.Readdirnames(64)
			}
		}
	}
	return "", &os.PathError{Op: "search", Path: pattern, Err: os.ErrNotExist}
}

func ReplaceFile(src, dst string) error {
	return replaceFile(src, dst)
}

/* CreateVia attempts to create a file by calling fn. If fn fails due to an
 * IsNotExist kind of error, the specified directory is created and the
 * function attempted once more. */
func CreateVia(dir string, fn func() error) (err error) {
	if err = fn(); err != nil && os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0777); err == nil {
			err = fn() // directory created, try again
		}
	}
	return
}

/* CreateIn creates a file in the given directory, creating the directory if necessary. */
func CreateIn(dir, filename string) (file *os.File, err error) {
	path := filepath.Join(dir, filename)
	err = CreateVia(dir, func() (e error) {
		file, e = os.Create(path)
		return e
	})
	return
}

/* Returns a path appropriate for storing a single config file */
func SingleConfigPath(usr, app *Dirs, name string) string {
	return filepath.Join(singleConfigDir(usr, app), name)
}

func envs(tpls... string) string {
	for _, tpl := range tpls {
		if s := os.ExpandEnv(tpl); s != "" {
			return s
		}
	}
	return ""
}
