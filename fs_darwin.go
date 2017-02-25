package fs

import (
	"errors"
	"os"
	"path/filepath"
	"github.com/sqweek/fs/osx"
	"unsafe"
)

// #include <mach-o/dyld.h>
import "C"

func userDirs() (*Dirs, error) {
	home = os.Getenv("HOME")
	if home == "" {
		return nil, errors.New("$HOME not set")
	}
	d := Dirs{
		Cache: osx.Dir(osx.CacheDir, osx.UserDomain),
		Data: osx.Dir(osx.SupportDir, osx.UserDomain),
		Docs: osx.Dir(osx.DocumentDir, osx.UserDomain),
	}
	d.Config = d.Data
	d.ResourcePath = []string{}
	return &d, nil
}

func appDirs(rqdn, name string) (*Dirs, error) {
	d, err := userDirs()
	if err != nil {
		return nil, err
	}
	d.Cache += "/" + rqdn
	d.Data += "/" + rqdn
	d.Config += "/" + rqdn
	d.Docs += "/" + name

	exe := exeDir()
	fd, err := os.Open(exe + "/../Resources")
	if fd != nil {
		fd.Close()
	}
	if err == nil || !os.IsNotExist(err) {
		d.ResourcePath = append(d.ResourcePath, exeDir() + "/../Resources")
	} else {
		d.ResourcePath = append(d.ResourcePath, exeDir())
	}
	return d, nil
}

func replaceFile(src, dst string) error {
	return os.Rename(src, dst)
}

func exeFileName() string {
	size := C.uint32_t(0)
	C._NSGetExecutablePath(nil, &size)
	buf := make([]byte, int(size))
	C._NSGetExecutablePath((*C.char)(unsafe.Pointer((&buf[0]))), &size)
	for buf[len(buf)-1] == '\000' {
		buf = buf[:len(buf) - 1]
	}
	return string(buf)
}

func exeDir() string {
	// TODO dereference symlinks from exeFileName?
	return filepath.Dir(exeFileName())
}

func singleConfigDir(usr, app *Dirs) string {
	return app.Config
}
