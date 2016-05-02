package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var tranquility = "87.237.34.202"
var versionNum = "1.0.11"
var saveURLContents = false
var unlimitedPing = false

var writer io.Writer

type CarriageReturnReplacer struct {
	w io.Writer
}

func (c CarriageReturnReplacer) Write(p []byte) (int, error) {
	_, err := c.w.Write(bytes.Replace(p, []byte("\n"), []byte("\r\n"), -1))

	if err != nil {
		return 0, err
	}

	return len(p), nil
}

func main() {
	fmt.Println("Running tests, results are being written to result.txt")

	saveChoice := false
	for saveChoice == false {
		fmt.Println("Do you wish to save the results of http requests for examination? y/n (default: n)")
		reader := bufio.NewReader(os.Stdin)
		inputChar, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Error: could not read user input, selecting default: n:")
			saveURLContents = false
			saveChoice = true
		} else {
			if inputChar == 89 || inputChar == 121 {
				saveURLContents = true
				saveChoice = true
			} else if inputChar == 78 || inputChar == 110 || inputChar == 13 || inputChar == 10 {
				saveURLContents = false
				saveChoice = true
			} else {
				fmt.Println("Choice invalid")
				saveChoice = false
			}
		}
	}
    
    pingChoice := false
	for pingChoice == false {
		fmt.Println("Do you wish to run the ping command against tranquility indefinitely? y/n (default: n)")
		reader := bufio.NewReader(os.Stdin)
		inputChar, _, err := reader.ReadRune()
		if err != nil {
			fmt.Println("Error: could not read user input, selecting default: n:")
			unlimitedPing = false
			pingChoice = true
		} else {
			if inputChar == 89 || inputChar == 121 {
				unlimitedPing = true
				pingChoice = true
			} else if inputChar == 78 || inputChar == 110 || inputChar == 13 || inputChar == 10 {
				unlimitedPing = false
				pingChoice = true
			} else {
				fmt.Println("Choice invalid")
				pingChoice = false
			}
		}
	}

	outfile, err := os.Create("result.txt")

	if err != nil {
		fmt.Println("Fatal error: could not create results file")
		panic(err)
	}

	defer outfile.Close()

	writer = io.MultiWriter(os.Stdout, CarriageReturnReplacer{outfile})

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(writer)

	log.Println("CCP connection test tool version ", versionNum)

	log.Printf("begin tests")
	runTests()
}

func runTests() {
	testPing("8.8.8.8")
	testPing(tranquility)
	tcpConnect(26000)
	tcpConnect(3724)
	testLauncherURL("http://client.eveonline.com/patches/premium_patchinfoTQ_inc.txt")
	testLauncherURL("http://web.ccpgamescdn.com/launcher/tranquility/selfupdates.htm")
    if unlimitedPing {
        unlimitedPingTest(tranquility)
    }
}

func testLauncherURL(url string) {
	log.Println("======HTTP TEST======")

	resp, err := http.Get(url)

	if err != nil {
		log.Println("failed to connect to url " + url + " with error:")
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("Connected, but failed to read contents of URL ", url, " with error:")
		log.Println(err)
		return
	}

	if saveURLContents {
		filename := "./" + cleanURL(url) + ".txt"
		outfile, err := os.Create(filename)

		if err != nil {
			log.Println("could not create results file")
			log.Panicln(err)
		}

		defer outfile.Close()
		bodyString := string(body[:])
		outfile.WriteString(bodyString)
	}

	log.Println("connected and read contents of ", url, " successfully")
}

func testPing(pingTarget string) {
	log.Println("======PING TEST======")

	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", pingTarget)
	} else {
		cmd = exec.Command("ping", "-c 5", pingTarget)
	}

	cmd.Stdout = writer
	err := cmd.Run()

	if err != nil {
		log.Println("ping test failure")
		log.Println(err)
	}
}

func unlimitedPingTest(pingTarget string) {
	log.Println("======UNLIMITED PING TEST======")
    fmt.Println("You are now in unlimited ping mode. The connection tester will continue to ping tranquility and log the result until it errors, or until you close this window.")
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("ping", "-t", pingTarget)
	} else {
		cmd = exec.Command("ping", pingTarget)
	}

	cmd.Stdout = writer
	err := cmd.Run()

	if err != nil {
		log.Println("ping test failure")
		log.Println(err)
	}
}

func tcpConnect(port uint64) {
	log.Println("======TCP TEST======")

	log.Println("tcpConnect on port ", port)
	conStr := tranquility + ":" + strconv.FormatUint(port, 10)
	conn, err := net.Dial("tcp", conStr)

	if err != nil {
		log.Println("Error connecting:")
		log.Println(err)
		return
	}

	defer conn.Close()
	buffer := make([]byte, 0, 4096)
	_, err = conn.Read(buffer)

	if err != nil {
		if err != io.EOF {
			log.Println("error reading buffer")
			log.Println(err)
			return
		}
	}

	log.Println(string(buffer))
	log.Println("connection successful")
}

func cleanURL(url string) string {
	holder := strings.Replace(url, ":", "-", -1)
	holder = strings.Replace(holder, "/", "", -1)
	holder = strings.Replace(holder, ".", "_", -1)
	return holder
}
