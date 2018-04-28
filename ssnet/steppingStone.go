package ssnet

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func checkErr(e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", e.Error())
		os.Exit(1)
	}
}

var (
	debug = false
)

type steppingStone struct {
	Addr net.IP
	Port uint16
}

// SteppingStones - list of stepping stones
type SteppingStones []steppingStone

// returns string representation of stepping stone (ip:port)
func (s steppingStone) String() string {
	return fmt.Sprintf("%s:%d", s.Addr, s.Port)
}

// Print - prints the stepping stones in the list
func (ssList SteppingStones) Print() {
	if len(ssList) == 0 {
		fmt.Println("chainlist is empty")
		return
	}
	fmt.Println("chainlist is:")
	for _, ss := range ssList {
		fmt.Println(ss.String())
	}
}

// GetFileNameFromURL - Retrieves the filename from the URL
func GetFileNameFromURL(URL string) (fname string) {
	// path part of the url contains the filename
	parsedURL, err := url.Parse(URL)
	checkErr(err)
	URLPathParts := strings.Split(parsedURL.Path, "/")

	for _, part := range URLPathParts {
		if strings.Contains(part, ".") {
			fname = part
		}
	}
	if fname == "" {
		fname = "index.html"
	}

	return
}

// SendReqToRandomStone - Sends an awget request to a random stepping stone
func SendReqToRandomStone(urlStr string, sStones SteppingStones) (conn net.Conn, err error) {
	// choose ss at random, seed with time so generated sequence is not always the same
	rand.Seed(time.Now().Unix())
	numSS := len(sStones)
	randSS := sStones[rand.Intn(numSS)].String()
	//randSS := sStones[0].String()
	tcpAddr, err := net.ResolveTCPAddr("tcp", randSS)
	checkErr(err)

	// create a connection with the random stepping stone
	conn, err = net.DialTCP("tcp", nil, tcpAddr)
	checkErr(err)

	// create a request struct with the URL and stepping stones
	request := NewGetRequest(urlStr, sStones)

	// create a json encoder that will write the request to the conection writer
	err = json.NewEncoder(conn).Encode(request)
	checkErr(err)

	fmt.Printf("sending to next random ss: %s\n", randSS)

	return
}

// HandleReqFromConn - Forward the req if not the last stepping stone in the ssList.
// Otherwise, extract the URL and send it back to the previous ss
func HandleReqFromConn(conn net.Conn) (err error) {
	var request GetRequest
	// decode the JSON and store it in value pointed to by request
	err = json.NewDecoder(conn).Decode(&request)
	checkErr(err)

	// retrieve the URL and stepping stones
	URL := request.URL
	ssList := request.SSlist

	hostname, err := os.Hostname()
	checkErr(err)
	myIPs, err := net.LookupIP(hostname)
	myIP := myIPs[0]

	// remove this stepping stone from the ssList
	for idx, sStone := range ssList {
		if bytes.Equal(myIP, sStone.Addr) {
			// remove this ss from the ssList
			ssList = append(ssList[:idx], ssList[idx+1:]...)
			//fmt.Printf("Removed myself from list, new ssList: %q\n", ssList)
		}
	}

	fmt.Printf("Request: %s\n", URL)
	ssList.Print()

	// last stepping stone in the chain, extract the URL
	// and send back to previous stepping stone
	if len(ssList) == 0 {
		//wget := exec.Command("wget", URL)
		//err = wget.Run()
		response, err := http.Get(URL)
		checkErr(err)
		defer response.Body.Close()

		// create file to write contents to
		wgetFile, err := ioutil.TempFile("", "temp")
		checkErr(err)
		// delete the temp file after using it to hold contents of wget
		defer wgetFile.Close()

		bytesWritten, err := io.Copy(wgetFile, response.Body)
		checkErr(err)
		fmt.Printf("issuing wget for file: %s\n", GetFileNameFromURL(URL))

		if debug {
			fmt.Printf("Downloaded file size: %d\n", bytesWritten)
		}

		// rewind file ptr back to beginning so we can write file contents
		wgetFile.Seek(0, 0)

		// write file back to prev ss
		bytesSent, err := io.Copy(conn, wgetFile)
		checkErr(err)
		fmt.Println("relaying file")

		if debug {
			fmt.Printf("Sent %d bytes\n", bytesSent)
		}

		// close the connection, awget job is done
		fmt.Println("Goodbye!")
		conn.Close()
	} else {
		// not the last stepping stone in the chain, forward request to the next ss
		nextSSConn, err := SendReqToRandomStone(URL, ssList)
		checkErr(err)

		// write file back to prev ss when response received
		fmt.Println("relaying file...")
		bytesWritten, err := io.Copy(conn, nextSSConn)
		checkErr(err)

		if debug {
			fmt.Printf("Received %d bytes from next ss, sending to prev ss\n", bytesWritten)
			fmt.Printf("Closing connection to prev ss: %s\n", conn.RemoteAddr())
		}

		conn.Close()
	}

	return
}

func newSteppingStone(ipStr, portStr string) (ss steppingStone, err error) {
	//ip := net.ParseIP(ipStr)
	IPs, err := net.LookupIP(ipStr)
	checkErr(err)
	ip := IPs[0]
	port64, err := strconv.ParseUint(portStr, 10, 16)
	port16 := uint16(port64)
	checkErr(err)
	ss = steppingStone{ip, port16}
	return
}

// ReadChainFile - reads the chainfile and creates a stepping stone list
func ReadChainFile(fname string) (ssList SteppingStones, err error) {
	file, err := os.Open(fname)
	if err != nil {
		fmt.Printf("Unable to locate chainfile: %s\nExiting\n", fname)
		os.Exit(1)
	}

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		lineSplit := strings.Fields(line)

		// number of SS, skip this line
		if len(lineSplit) < 2 {
			continue
		} else { // ip portNums
			steppingStone, err := newSteppingStone(lineSplit[0], lineSplit[1])
			checkErr(err)
			ssList = append(ssList, steppingStone)
		}
	}
	return
}
