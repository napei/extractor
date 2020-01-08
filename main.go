package main

import (
	"flag"
	"fmt"
	"github.com/mholt/archiver/v3"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var files []string
var inputarg string = ""
var outputarg string = ""
var versionarg bool = false
var dryrunarg bool = false
var verbosearg bool = false

var appVersion string = "v0.5"

var (
	black   = consoleColor("\033[1;30m%s\033[0m")
	red     = consoleColor("\033[1;31m%s\033[0m")
	green   = consoleColor("\033[1;32m%s\033[0m")
	yellow  = consoleColor("\033[1;33m%s\033[0m")
	purple  = consoleColor("\033[1;34m%s\033[0m")
	magenta = consoleColor("\033[1;35m%s\033[0m")
	teal    = consoleColor("\033[1;36m%s\033[0m")
	white   = consoleColor("\033[1;37m%s\033[0m")
)

func consoleColor(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func searchForArchives() error {
	var outputmessage = "[Looking for Archives]"
	var partRegex = regexp.MustCompile("^.*(part[0-9]+\\.rar)$")
	if dryrunarg {
		outputmessage += " - Dry Run"
	}
	fmt.Println(yellow(outputmessage))

	err := filepath.Walk(inputarg, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && (strings.Contains(info.Name(), "part01.rar") || (!partRegex.MatchString(info.Name()) && strings.Contains(info.Name(), ".rar"))) {
			if verbosearg {
				fmt.Println(green("[Found Archive]: ") + filepath.Clean(path))
			}
			files = append(files, path)
		}
		return nil
	})
	fmt.Println(teal("Found " + strconv.Itoa(len(files)) + " archives"))
	return err
}

func extractArchives() {
	fmt.Println(yellow("[Extracting Archives]"))
	var outputpath string
	var currentItem string
	for i := range files {
		currentItem = strconv.Itoa(i+1) + "\\" + strconv.Itoa(len(files))

		if outputarg != "" {
			outputpath = outputarg
		} else {
			outputpath = filepath.Dir(files[i])
		}
		var outputstring = green("[Extracting]") + " - " + teal(currentItem)
		if verbosearg {
			outputstring += white(" - " + filepath.Base(files[i]))
		}
		fmt.Println(outputstring)
		archiver.Unarchive(files[i], outputpath)
	}
}
func fail() {
	fmt.Println(red("ERROR: input path not specified. Call the program as: ") + filepath.Base(os.Args[0]) + " -input=\"Directory\" <flags>")
	fmt.Println("For help, use the '-h' flag")
	os.Exit(1)
}

func init() {
	flag.StringVar(&inputarg, "input", "", "Specify input directory for scanning in the form -input Directory")
	flag.StringVar(&outputarg, "output", "", "Specify an alternate output directory for all located archives in the form -output Directory. By default, this program will output archives in the same folder.")
	flag.BoolVar(&versionarg, "version", false, "Output the version of the program and exit")
	flag.BoolVar(&dryrunarg, "dryrun", false, "Don't extract archives, only list them")
	flag.BoolVar(&verbosearg, "verbose", false, "List every archive individually")
	flag.Parse()
}

func main() {

	if !(len(os.Args) > 1) {
		fail()
	}

	if versionarg {
		fmt.Println("Extractor - " + appVersion)
		os.Exit(0)
	}

	if inputarg != "" {
		inputarg = filepath.Clean(inputarg)
	} else {
		fail()
	}
	if outputarg != "" {
		outputarg = filepath.Clean(outputarg)
	}

	searchForArchives()

	if dryrunarg {
		fmt.Println(yellow("Dry run complete. No archives extracted"))
	} else {
		extractArchives()
	}
}
