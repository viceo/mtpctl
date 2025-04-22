package util

func PanicIfError(e error) {
	if e != nil {
		panic(e)
	}
}
