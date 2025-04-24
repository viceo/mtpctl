package util

import (
	"fmt"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

func PanicIfError(e error) {
	if e != nil {
		panic(e)
	}
}

// Recover panics and log them in JSON format ending process.
// Manually JSON formatting because marshaling can return error
func ErrorHandler() {
	if r := recover(); r != nil {
		var stackTrace string
		stackTrace = string(debug.Stack())
		stackTrace = strings.ReplaceAll(stackTrace, "\n", "")
		stackTrace = strings.ReplaceAll(stackTrace, "\t", " ")
		stackTrace = strconv.Quote(stackTrace)
		fmt.Printf("{\"hasError\":true, \"error\": \"%s\", \"trace\": %s}\n", r, stackTrace)
		os.Exit(2)
	}
}
