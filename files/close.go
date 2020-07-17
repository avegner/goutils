package files

import "os"

func Close(fd *os.File, err *error) {
	cerr := fd.Close()
	if *err == nil {
		*err = cerr
	}
}
