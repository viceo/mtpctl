package util

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/viceo/tplibcmd/dto"
)

func PanicIfError(e error) {
	if e != nil {
		panic(e)
	}
}

// Error Trap
func ErrorTrap() {
	if r := recover(); r != nil {
		var stackTrace string
		stackTrace = string(debug.Stack())
		stackTrace = strings.ReplaceAll(stackTrace, "\n", "")
		stackTrace = strings.ReplaceAll(stackTrace, "\t", " ")
		stackTrace = strconv.Quote(stackTrace)
		errorResponse := dto.CommandResponse{
			HasError:   true,
			Error:      fmt.Sprintf("%s", r),
			StackTrace: stackTrace,
		}
		errorResponse.Print()
	}
}
