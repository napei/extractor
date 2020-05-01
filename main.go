package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/mholt/archiver/v3"
)

var (
	inputarg     string = ""
	outputarg    string = ""
	versionarg   bool   = false
	dryrunarg    bool   = false
	verbosearg   bool   = false
	overwritearg bool   = false
)

const appVersion string = "v0.7.1"

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
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString, fmt.Sprint(args...))
	}
}

func searchForArchives(inputpath string, verbose bool) (out []string, err error) {
	var outputmessage = "[Looking for Archives]"
	var partRegex = regexp.MustCompile("^.*(part[0-9]+\\.rar)$")
	if dryrunarg {
		outputmessage += " - Dry Run"
	}
	fmt.Println(yellow(outputmessage))

	err = filepath.Walk(inputpath, func(path string, info os.FileInfo, e error) (err error) {
		if e != nil {
			log.Println("Error: ", e)
			return e
		}

		// Skip if the current file a directory
		if info.IsDir() {
			return nil
		}

		// Current filename
		filename := info.Name()

		// Find explicit part01 file
		containsPart01 := strings.Contains(filename, "part01.rar")
		// Check if the file is a rar part file in general
		isAnyPartFile := partRegex.MatchString(filename)
		// Check if file has .rar extension
		isRarFile := filepath.Ext(filename) == ".rar"

		isZipFile := filepath.Ext(filename) == ".zip"

		if containsPart01 || (!isAnyPartFile && isRarFile) {
			p := filepath.Clean(path)
			if verbose {
				fmt.Println(green("[Found RAR Archive]: ") + p)
			}
			out = append(out, p)
		} else if isZipFile {
			p := filepath.Clean(path)
			if verbose {
				fmt.Println(green("[Found ZIP Archive]: ") + p)
			}
			out = append(out, p)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	fmt.Println(teal("Found " + strconv.Itoa(len(out)) + " archives"))

	return out, nil
}

func extractArchives(files []string, outputpath string, verbose bool) (err error) {
	fmt.Println(yellow("[Extracting Archives]"))

	var outPath string
	var currentFile string

	for i := range files {
		currentFile = strconv.Itoa(i+1) + "\\" + strconv.Itoa(len(files))

		if outputarg != "" {
			outPath = outputarg
		} else {
			outPath = filepath.Dir(files[i])
		}
		var outputMessage = green("[Extracting]") + " - " + teal(currentFile)
		if verbose {
			outputMessage += white(" - " + filepath.Base(files[i]))
		}
		fmt.Println(outputMessage)

		arc, err := archiver.ByExtension(files[i])
		if err != nil {
			return err
		}

		switch arc.(type) {
		case *archiver.Rar:
			arc.(*archiver.Rar).OverwriteExisting = overwritearg
			arc.(*archiver.Rar).ContinueOnError = true
			break

		case *archiver.Zip:
			arc.(*archiver.Zip).OverwriteExisting = overwritearg
			arc.(*archiver.Zip).ContinueOnError = true
			break
		}

		err = arc.(archiver.Unarchiver).Unarchive(files[i], outPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func processDirectory(inputpath string, outputpath string, dry bool, verbose bool) (err error) {
	if inputpath == "" {
		fail()
	}

	files, err := searchForArchives(inputpath, verbose)
	if err != nil {
		return err
	}

	if dry {
		fmt.Println(yellow("Dry run complete. No archives extracted"))
	} else {
		err = extractArchives(files, outputpath, verbose)
		if err != nil {
			return err
		}
	}

	return nil
}

func fail() {
	fmt.Println(red("ERROR: input path not specified. Call the program as: ") + filepath.Base(os.Args[0]) + " -input=\"Directory\" <flags>")
	fmt.Println("For help, use the '-h' flag")
	os.Exit(1)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -input=\"Directory\" [flags]\n", filepath.Base(os.Args[0]))

		flag.PrintDefaults()
	}

	flag.StringVar(&inputarg, "input", "", "Specify input directory for scanning in the form -input Directory")
	flag.StringVar(&outputarg, "output", "", "Specify an alternate output directory for all located archives in the form -output Directory.\nBy default, this program will output archives in the same folder.")
	flag.BoolVar(&versionarg, "version", false, "Output the version of the program and exit")
	flag.BoolVar(&dryrunarg, "dryrun", false, "Don't extract archives, only list them. Default: false")
	flag.BoolVar(&verbosearg, "verbose", false, "List archive names. Default: false")
	flag.BoolVar(&overwritearg, "overwrite", false, "Overwrite existing files. Default: false")
	flag.Parse()
}

func main() {

	if len(os.Args) < 1 {
		fail()
	}

	if versionarg {
		fmt.Println("Extractor - " + appVersion)
		os.Exit(0)
	}

	if inputarg != "" {
		inputarg = filepath.Clean(inputarg)
	} else {
		inputarg = ""
	}
	if outputarg != "" {
		outputarg = filepath.Clean(outputarg)
	} else {
		outputarg = ""
	}

	archiver.DefaultRar.OverwriteExisting = overwritearg
	archiver.DefaultRar.ContinueOnError = true

	err := processDirectory(inputarg, outputarg, dryrunarg, verbosearg)
	if err != nil {
		log.Fatalln("An error occured: ", err)
	}
}
