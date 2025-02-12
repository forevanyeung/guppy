package internal

var verbose bool
var desktop bool

func IsDesktop() bool {
	return desktop
}

func IsVerbose() bool {
	return verbose
}

func SetVerbose(v bool) {
	verbose = v
}

func SetDesktop(d bool) {
	desktop = d
}
