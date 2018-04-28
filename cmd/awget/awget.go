package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"proj3/ssnet"
)

func checkErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", e.Error())
		os.Exit(1)
	}
}

var (
	chainFileName    string
	defaultChainFile string
	chainFileUsage   string
)

const (
	debug = false
)

func init() {
	defaultChainFile = "chaingang.txt"
	chainFileUsage = "filename containing stepping stones for the request"
}

// exit if no URL specified
func usageExit() {
	fmt.Println("URL not specified")
	os.Exit(1)
}

func validateURL(a string) (string, error) {
	_, err := url.ParseRequestURI(a)
	checkErr(err)

	return a, err
}

func main() {
	// no URL specified, exit
	if len(os.Args) < 2 {
		usageExit()
	} else if len(os.Args) < 3 {
		// chainfile not specified, use the default chainfile
		chainFileName = defaultChainFile
	} else {
		chainFileName = os.Args[2]
	}

	URL := os.Args[1]
	fmt.Printf("Request: %s\n", URL)

	ssList, err := ssnet.ReadChainFile(chainFileName)
	checkErr(err)

	ssList.Print()

	if debug {
		fmt.Println("Read chainfile successfully")
	}

	conn, err := ssnet.SendReqToRandomStone(URL, ssList)
	checkErr(err)
	defer conn.Close()

	// wait for response from ss and read the whole file
	fmt.Println("waiting for file...")
	readFile, err := ioutil.ReadAll(conn)
	checkErr(err)

	wgetFName := ssnet.GetFileNameFromURL(URL)
	if debug {
		fmt.Printf("wgetFName: %s\n", wgetFName)
	}

	err = ioutil.WriteFile(wgetFName, readFile, 0644)
	checkErr(err)

	fmt.Printf("received file: %s\nGoodbye!\n", wgetFName)
}
