// entry point into our application
package main

import (
	"bookie/router"

	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load() // loads env variables
}

func main() {
	router.Router()
}
