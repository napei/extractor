package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/mholt/archiver/v3"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var files []string
var outputdir string
var appVersion string = "v0.4"

func searchForArchives(rootpath string) error {
	var partRegex = regexp.MustCompile("^.*(part[0-9]+\\.rar)$")

	fmt.Printf(color.YellowString("[Looking for Archives]\n"))
	err := filepath.Walk(rootpath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.Contains(info.Name(), "part01.rar") || (!partRegex.MatchString(info.Name()) && strings.Contains(info.Name(), ".rar"))) {
			fmt.Printf(color.GreenString("[Found Archive]: ") + path + "\n")
			files = append(files, path)
		}
		return nil
	})

	return err
}

func extractArchives() {
	fmt.Printf(color.YellowString("\n[Extracting Archives]\n"))
	for i := range files {
		outputpath := filepath.Dir(files[i])
		fmt.Printf(color.GreenString("[Extracting Archive]: ") + files[i] + color.BlueString(" - to directory ["+filepath.Dir(files[i])+"]\n"))
		archiver.Unarchive(files[i], outputpath)
	}
}

func main() {
	args := os.Args[1:]
	var version bool
	// Handle Args
	if !(len(os.Args) > 1) {
		fmt.Printf(color.RedString("ERROR: input path not specified. Call the program as: ") + filepath.Base(os.Args[0]) + " [input path] <flags>\n")
		fmt.Printf("For help, use the '-h' flag\n")
		os.Exit(1)
	}
	flag.StringVar(&outputdir, "output", "", "Specify an alternate output directory for all located archives. By default, this program will output archives in the same folder.")
	flag.BoolVar(&version, "version", false, "Output the version of the program")
	flag.Parse()

	if version {
		fmt.Printf("Extractor - " + appVersion)
		os.Exit(0)
	}

	searchForArchives(args[0])
	extractArchives()
}
