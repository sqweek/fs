package fs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

var (
	kernel = syscall.MustLoadDLL("kernel32.dll")
	getModuleFileName = kernel.MustFindProc("GetModuleFileNameW")
	shell32 = syscall.MustLoadDLL("shell32.dll")
	getFolderPath = shell32.MustFindProc("SHGetFolderPathW")
)

const (
	CSIDL_APPDATA = 0x1a
	CSIDL_LOCAL_APPDATA = 0x1c
	CSIDL_PROGRAM_FILES = 0x26
	CSIDL_PERSONAL = 0x05
)

type Path [syscall.MAX_PATH]uint16

func userDirs() (*Dirs, error){
	home = os.Getenv("USERPROFILE")
	if home == "" {
		return nil, errors.New("%USERPROFILE% not set")
	}
	docs := "Documents"
	if strings.Contains(home, "Documents and Settings") {
		docs = "My Documents"
	}
	var path Path
	d := Dirs{
		Cache: path.folder(CSIDL_LOCAL_APPDATA, "$LOCALAPPDATA"),
		Data: path.folder(CSIDL_APPDATA, "$APPDATA"),
		Docs: path.folder(CSIDL_PERSONAL, home + "\\" + docs),
	}
	d.Config = d.Data
	d.ResourcePath = []string{}
	return &d, nil
}

func appDirs(fqdn string, subdir... string) (*Dirs, error) {
	d, err := userDirs()
	if err != nil {
		return nil, err
	}
	name := filepath.Join(subdir...)
	d.Cache += "\\" + name
	d.Data += "\\" + name
	d.Docs += "\\" + name
	d.Config += "\\" + name
	syspath := d.ResourcePath
	d.ResourcePath = []string{exeDir()}
	d.ResourcePath = append(d.ResourcePath, syspath...)

	return d, nil
}

func exeFileName() (string, error) {
	buf := make([]uint16, syscall.MAX_PATH)
	r, _, err := getModuleFileName.Call(0, uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if r == 0 {
		return "", err
	}
	return string(utf16.Decode(buf[0:uint32(r)])), nil
}

func exeDir() string {
	if exe, err := exeFileName(); err != nil {
		fmt.Println(os.Stderr, "warning: GetModuleFileNameW:", err)
		return "."
	} else {
		return filepath.Dir(exe)
	}
}

func folderPath(csidl int, buf [syscall.MAX_PATH]uint16) (string, error) {
	r, _, err := getFolderPath.Call(0, uintptr(csidl), 0, 0, uintptr(unsafe.Pointer(&buf[0])))
	if r == 0 {
		return "", err
	}
	return string(utf16.Decode(buf[0:uint32(r)])), nil
}

func (p *Path) folder(csidl int, fallbacks... string) string {
	if s, _ := folderPath(csidl, ([syscall.MAX_PATH]uint16)(*p)); s != "" {
		return s
	}
	return envs(fallbacks...)
}

func replaceFile(src, dst string) error {
	bak := dst + ".bak"
	for {
		err := os.Rename(dst, bak)
		if err == nil || os.IsNotExist(err) {
			/* rename succeeded or 'dst' doesn't exist yet; we can proceed */
			break
		} else if os.IsExist(err) {
			if err = os.Remove(bak); err == nil {
				continue /* old backup is gone, rename should work now */
			}
		}
		return err
	}
	if err := os.Rename(src, dst); err == nil {
		os.Remove(bak)
		return nil
	} else {
		return err
	}
}

func singleConfigDir(usr, app *Dirs) string {
	return app.Config
}
