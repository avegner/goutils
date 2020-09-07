package errs

func Finally(f func() error, err *error) {
	ferr := f()
	if *err == nil {
		*err = ferr
	}
}
