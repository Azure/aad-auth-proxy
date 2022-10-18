package utils

import "log"

// Handle Panic or unexpected function completion - log and exit process
// would only happen due to some out of memory/corruption so letting host restart process safer than attempting to recover
func HandlePanic(function string) {
	if err := recover(); err != nil {
		log.Fatal("Unexpected error in "+function+" process exiting: ", err)
	}
	log.Fatal(function + " unexpectedly finished, process exiting")
}

// Handle Panic for Goroutines that will exit - only fatal if an err
func HandlePanicFunctionExits(function string) {
	if err := recover(); err != nil {
		log.Fatal("Unexpected error in "+function+" process exiting: ", err)
	}
}
