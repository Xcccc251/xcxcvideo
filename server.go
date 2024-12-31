package main

import "XcxcVideo/router"

func main() {
	r := router.Router()
	r.Run(":7070")
}
