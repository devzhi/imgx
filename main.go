package main

import (
	"flag"
	"fmt"
)

var (
	imageName = flag.String("name", "", "name of the image")
	tag       = flag.String("tag", "latest", "tag of the image")
	arch      = flag.String("arch", "amd64", "architecture of the image")
)

func main() {
	flag.Parse()
	if *imageName == "" {
		fmt.Println("Image name is required")
		return
	}
	fmt.Println("Pulling image", *imageName, "with tag", *tag, "for architecture", *arch)
}
