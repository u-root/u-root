// +build darwin freebsd linux
// +build cgo

package group

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"syscall"
	"unsafe"
)

/*
#include <unistd.h>
#include <sys/types.h>
#include <grp.h>
#include <stdlib.h>

static int mygetgrgid_r(int gid, struct group *grp,
	char *buf, size_t buflen, struct group **result) {
 return getgrgid_r(gid, grp, buf, buflen, result);
}

int array_len(char **array) {
	int len = 0;
	for (; *array != NULL; array++)
		len++;
	return len;
}
*/
import "C"

// Current returns the primary group of the caller.
func Current() (*Group, error) {
	return lookup(syscall.Getgid(), "", false)
}

// Lookup looks up a group by name. If the group cannot be found,
// the returned error is of type UnknownGroupError.
func Lookup(groupname string) (*Group, error) {
	return lookup(-1, groupname, true)
}

// LookupId looks up a group by groupid. If the group cannot be found,
// the returned error is of type UnknownGroupIdError.
func LookupId(gid string) (*Group, error) {
	i, e := strconv.Atoi(gid)
	if e != nil {
		return nil, e
	}
	return lookup(i, "", false)
}

func lookup(gid int, groupname string, lookupByName bool) (*Group, error) {
	var (
		grp    C.struct_group
		result *C.struct_group
	)
	var bufSize C.long
	if runtime.GOOS == "freebsd" {
		// FreeBSD doesn't have _SC_GETPW_R_SIZE_MAX
		// and just returns -1.  So just use the same
		// size that Linux returns
		bufSize = 1024
	} else {
		bufSize = C.sysconf(C._SC_GETPW_R_SIZE_MAX)
		if bufSize <= 0 || bufSize > 1<<20 {
			return nil, fmt.Errorf(
				"user: unreasonable _SC_GETPW_R_SIZE_MAX of %d", bufSize)
		}
	}
	buf := C.malloc(C.size_t(bufSize))
	defer C.free(buf)
	var rv C.int
	if lookupByName {
		nameC := C.CString(groupname)
		defer C.free(unsafe.Pointer(nameC))
		rv = C.getgrnam_r(nameC,
			&grp,
			(*C.char)(buf),
			C.size_t(bufSize),
			&result)
		if rv != 0 {
			return nil, fmt.Errorf(
				"group: lookup groupname %s: %s", groupname, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupError(groupname)
		}

	} else {
		rv = C.mygetgrgid_r(C.int(gid),
			&grp,
			(*C.char)(buf),
			C.size_t(bufSize),
			&result)
		if rv != 0 {
			return nil, fmt.Errorf("group: lookup groupid %d: %s", gid, syscall.Errno(rv))
		}
		if result == nil {
			return nil, UnknownGroupIdError(gid)
		}
	}
	g := &Group{
		Gid:     strconv.Itoa(int(grp.gr_gid)),
		Name:    C.GoString(grp.gr_name),
		Members: getMembers(grp),
	}

	return g, nil
}

func getMembers(grp C.struct_group) []string {
	// We need to count the members before we can create a slice.
	nmembers := int(C.array_len(grp.gr_mem))

	members := make([]string, nmembers)

	// Create a slice over the C grp.gr_mem char* array
	var raw_gr_mem []*C.char
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&raw_gr_mem)))
	sliceHeader.Cap = nmembers
	sliceHeader.Len = nmembers
	sliceHeader.Data = uintptr(unsafe.Pointer(grp.gr_mem))

	for idx, m := range raw_gr_mem {
		members[idx] = C.GoString(m)
	}

	return members
}
