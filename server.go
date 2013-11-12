package main

import (
	"github.com/hoisie/web"
)

func collections(val string) string {
	return "<html><body><h1>collections <i>" + val + "</i></h1></body></html>"
}

func main() {
	web.Get("/(.*)", collections)
	web.Run("0.0.0.0:8080")
}
