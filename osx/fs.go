package osx

import (
	"bytes"
	"unsafe"
)

//#cgo darwin LDFLAGS: -framework Foundation
//#include <sys/syslimits.h>
//#include "fs.h"
import "C"

const (
	UserDomain = 1 // NSUserDomainMask
	LocalDomain = 2 // NSLocalDomainMask

	DocumentDir = 9 // NSDocumentDirectory
	CacheDir = 13 // NSCachesDirectory
	SupportDir = 14 // NSApplicationSupportDirectory
)

func Dir(dir, domain int) string {
	buf := make([]byte, int(C.PATH_MAX))
	C.get_dir(C.int(dir), C.int(domain), (*C.char)(unsafe.Pointer(&buf[0])), C.PATH_MAX)
	return string(buf[:bytes.Index(buf, []byte{0})])
}