package fs

import (
	"errors"
	"os"
	"os/exec"
	"path"
	"strings"
)

func userDirs() (*Dirs, error) {
	home = os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("$HOME not set")
	}
	d := Dirs{
		Cache: envs("$XDG_CACHE_HOME", home + "/.cache"),
		Data: envs("$XDG_DATA_HOME", home + "/.local/share"),
		Config: envs("$XDG_CONFIG_HOME", home + "/.config"),
		Docs: envs("$XDG_DOCUMENTS_DIR", home + "/Documents"),
	}
	for _, dir := range strings.Split(envs("$XDG_DATA_DIRS", "/usr/local/share:/usr/share:/opt"), ":") {
		d.ResourcePath = append(d.ResourcePath, dir)
	}
	return &d, nil
}

func appDirs(fqdn, name string) (*Dirs, error) {
	d, err := userDirs()
	if err != nil {
		return nil, err
	}
	d.Cache += "/" + name
	d.Data += "/" + name
	d.Config += "/" + name
	d.Docs += "/" + name
	for i, _ := range d.ResourcePath {
		d.ResourcePath[i] += "/" + name
	}
	sysPath := d.ResourcePath
	e := exeDir()
	if dir, _ := path.Split(e); !strings.HasSuffix(dir, "/bin") {
		d.ResourcePath = []string{dir}
	} else {
		d.ResourcePath = []string{dir + "/../share/" + name}
	}
	d.ResourcePath = append(d.ResourcePath, sysPath...)
	return d, nil
}

func exeDir() string {
	var exe string
	if strings.Contains(os.Args[0], "/") {
		exe = os.Args[0]
	} else {
		p, err := exec.LookPath(os.Args[0])
		if err != nil {
			return "." // couldn't find self in path
		}
		exe = p
	}
	// absolute or relative path; just drop the executable name
	dir, _ := path.Split(exe)
	return dir
}

func replaceFile(src, dst string) error {
	return os.Rename(src, dst)
}
