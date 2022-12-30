package scrape

var _VERBOSE = false

func SetVerbose(v bool) {
	_VERBOSE = v
}

func Verbose() bool {
	return _VERBOSE
}
