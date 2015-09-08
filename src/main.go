package main

import (
	"flag"
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

var versionNum = "1.0.4"
var saveURLContents = flag.Bool("keepweb", false, "Stores the contents of http requests in individual files")

var writer io.Writer

func main() {
	fmt.Println("Running tests, results are being written to result.txt")
	flag.Parse()

	outfile, err := os.Create("result.txt")

	if err != nil {
		fmt.Println("Fatal error: could not create results file")
		panic(err)
	}

	defer outfile.Close()

	writer = io.MultiWriter(os.Stdout, outfile)

	log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	log.SetOutput(writer)

	log.Println("CCP connection test tool version ", versionNum)

	if *saveURLContents == false {
		log.Println("successful web requests will not be stored for examination, specify -keepweb=true to store them")
	}

	log.Printf("begin tests")
	runTests()
}

func runTests() {
	testPing()
	tcpConnect(26000)
	tcpConnect(3724)
	testPortOpen(26000)
	testPortOpen(3724)
	testLauncherURL("http://client.eveonline.com/patches/premium_patchinfoTQ_inc.txt")
	testLauncherURL("http://web.ccpgamescdn.com/launcher/tranquility/selfupdates.htm")
}

func testPortOpen(port uint64) {
	log.Println("======PORT FORWARDING TEST======")

	urlStr := "http://tuq.in/tools/port.txt?port=" + strconv.FormatUint(port, 10)
	resp, err := http.Get(urlStr)

	if err != nil {
		log.Println("port open check failed, could not get address")
		log.Println(err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println("port open check failed, could not read response")
		log.Println(err)
		return
	}

	bodyString := string(body[:])
	log.Println("The scan of port ", port, " returned ", bodyString)
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

	if *saveURLContents {
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

func testPing() {
	log.Println("======PING TEST======")

	if runtime.GOOS == "windows" {
		cmd := exec.Command("ping", "87.237.38.200")
		cmd.Stdout = writer
		err := cmd.Run()

		if err != nil {
			log.Println("OS does not have ping utility")
			log.Println(err)
			return
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("ping", "-c 5", "87.237.38.200")
		cmd.Stdout = writer
		err := cmd.Run()

		if err != nil {
			log.Println("OS does not have ping utility")
			log.Println(err)
			return
		}
	} else {
		log.Println("unsupported OS, no ping")
	}
}

func tcpConnect(port uint64) {
	log.Println("======TCP TEST======")

	log.Println("tcpConnect on port ", port)
	conStr := "87.237.38.200:" + strconv.FormatUint(port, 10)
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
