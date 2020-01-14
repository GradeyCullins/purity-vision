package main

import (
	"flag"

	"github.com/GradeyCullins/GoogleVisionFilter/src"
)

var portFlag int

func main() {
	flag.IntVar(&portFlag, "port", 8080, "port to run the service on")
	flag.Parse()

	src.InitWebServer(portFlag)

	// flag.Usage = func() {
	// 	fmt.Fprintf(os.Stderr, "Usage: %s <path-to-image>\n", filepath.Base(os.Args[0]))
	// }
	// flag.Parse()

	// args := flag.Args()
	// if len(args) == 0 {
	// 	flag.Usage()
	// 	os.Exit(1)
	// }

}
