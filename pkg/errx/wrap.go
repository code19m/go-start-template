package errx

import (
	"fmt"
	"runtime"
)

// Wrap wraps an error with an ErrorX. If the error is nil, Wrap returns nil.
// If error is not an ErrorX, it is wrapped with an ErrorX by utility functions like fromPG and fromGRPC.
// If after wrapping the error is still not an ErrorX, it is considered that the error is not properly handled
// and is wrapped with an ErrorX with Internal type and the error message as a detail.
// The function also appends the caller information as a stacktrace of the function that invoked the function
func Wrap(err error) error {
	if err == nil {
		return nil
	}

	err = fromPG(err)
	err = fromGRPC(err)

	e, ok := err.(*ErrorX)
	if !ok {
		e = ErrInternal.WithDetail("error", err.Error())
	}

	e.addTrace()
	return e
}

// addTrace appends the caller information of the function that invoked the function
// that called addTrace to the error's trace field. This helps in tracking the chain
// of function calls leading to the error, providing a detailed trace for debugging.
func (e *ErrorX) addTrace() {
	// Skip 0 to get the current function, 1 to get the caller of the current function, etc.
	// Here, we skip 2 to get the caller of the function that invoked addTrace.
	pc, filepath, line, ok := runtime.Caller(2)
	if !ok {
		panic("could not get runtime.Caller")
	}

	fn := runtime.FuncForPC(pc)
	if fn == nil {
		panic("could not get runtime.FuncForPC")
	}

	_, filename := pathSplit(filepath)

	callerInfo := fmt.Sprintf("%s:%d|%s", filename, line, fn.Name())

	if e.trace == "" {
		e.trace = callerInfo
	} else {
		e.trace = fmt.Sprintf("%s ➡️ %s", e.trace, callerInfo)
	}
}

// pathSplit splits a path into the directory and the file name
func pathSplit(path string) (string, string) {
	for i := len(path) - 1; i > 0; i-- {
		if path[i] == '/' {
			return path[:i], path[i+1:]
		}
	}
	return "", path
}
