package utils

import (
	"io/fs"
	// "runtime"
	"syscall"
	"time"
)

func GetTimes(info fs.FileInfo) (created, modified time.Time) {
	modified = info.ModTime()

	switch stat := info.Sys().(type) {
	// case *syscall.Stat_t:
	// 	// macOS / BSD: Birthtimespec is available
	// 	if runtime.GOOS == "darwin" || runtime.GOOS == "freebsd" {
	// 		created = time.Unix(stat.Birthtimespec.Sec, stat.Birthtimespec.Nsec)
	// 	} else {
	// 		// Linux: No true creation time, fallback to ctime (metadata change)
	// 		created = time.Unix(stat.Ctim.Sec, stat.Ctim.Nsec)
	// 	}
	case *syscall.Win32FileAttributeData:
		// Windows: CreationTime is available
		created = time.Unix(0, stat.CreationTime.Nanoseconds())
	default:
		created = time.Time{} // Unknown
	}

	return
}