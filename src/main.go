package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var versionNum = "1.0.4"
var saveUrlContents = flag.Bool("keepweb", false, "Stores the contents of http requests in individual files")

func main() {
	fmt.Printf("Running tests, results are being written to result.txt")
	flag.Parse()
	path, err := os.Getwd()
	if err != nil {
		fmt.Printf("Fatal error: could not get working directory\r\n")
		panic(err)
	}
	outfile, err := os.Create(path + "/result.txt")
	if err != nil {
		fmt.Printf("Fatal error: could not create results file\r\n")
		panic(err)
	}
	defer outfile.Close()
	os.Stdout = outfile
	fmt.Printf("CCP connection test tool version " + versionNum + "\r\n")
	if *saveUrlContents == false {
		fmt.Printf("successful web requests will not be stored for examination, specify -keepweb=true to store them\r\n")
	}
	fmt.Printf("begin tests\r\n")
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
	fmt.Printf("======PORT FORWARDING TEST======\r\n")
	urlStr := "http://tuq.in/tools/port.txt?port=" + strconv.FormatUint(port, 10)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Printf("port open check failed, could not get address\r\n")
		fmt.Println(err, "\r\n")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("port open check failed, could not read response\r\n")
		fmt.Println(err, "\r\n")
		return
	}
	bodyString := string(body[:])
	fmt.Printf("The scan of port " + strconv.FormatUint(port, 10) + " returned " + bodyString + "\r\n")
}

func testLauncherURL(url string) {
	fmt.Printf("======HTTP TEST======\r\n")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("failed to connect to url " + url + " with error:\r\n")
		fmt.Println(err, "\r\n")
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Connected, but failed to read contents of URL " + url + " with error:\r\n")
		fmt.Println(err, "\r\n")
		return
	}
	if *saveUrlContents {
		filename := "./" + cleanURL(url) + ".txt"
		outfile, err := os.Create(filename)
		if err != nil {
			fmt.Printf("could not create results file\r\n")
			panic(err)
		}
		defer outfile.Close()
		bodyString := string(body[:])
		outfile.WriteString(bodyString)
	}
	fmt.Printf("connected and read contents of " + url + " successfully\r\n")

}

func testPing() {
	fmt.Printf("======PING TEST======\r\n")
	if runtime.GOOS == "windows" {
		cmd := exec.Command("ping", "87.237.38.200")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Printf("OS does not have ping utility\r\n")
			fmt.Println(err, "\r\n")
			return
		}
	} else if runtime.GOOS == "darwin" {
		cmd := exec.Command("ping", "-c 5", "87.237.38.200")
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			fmt.Printf("OS does not have ping utility\r\n")
			fmt.Println(err, "\r\n")
			return
		}
	} else {
		fmt.Printf("unsupported OS, no ping\r\n")
	}
}

func tcpConnect(port uint64) {
	fmt.Printf("======TCP TEST======\r\n")
	fmt.Printf("tcpConnect on port " + strconv.FormatUint(port, 10) + "\r\n")
	conStr := "87.237.38.200:" + strconv.FormatUint(port, 10)
	conn, err := net.Dial("tcp", conStr)
	if err != nil {
		fmt.Printf("Error connecting:\r\n")
		fmt.Println(err, "\r\n")
		return
	}
	defer conn.Close()
	buffer := make([]byte, 0, 4096)
	_, err = conn.Read(buffer)
	if err != nil {
		if err != io.EOF {
			fmt.Printf("error reading buffer\r\n")
			fmt.Println(err, "\r\n")
			return
		}
	}
	fmt.Printf(string(buffer))
	fmt.Printf("connection successful\r\n")
}

func cleanURL(url string) string {
	holder := strings.Replace(url, ":", "-", -1)
	holder = strings.Replace(holder, "/", "", -1)
	holder = strings.Replace(holder, ".", "_", -1)
	return holder
}
