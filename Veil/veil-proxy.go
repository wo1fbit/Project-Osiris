package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

const defaultPort = "8081"

func main() {
	var listenPort string
	if len(os.Args) > 1 {
		listenPort = os.Args[1]
	} else {
		listenPort = defaultPort
	}

	ln, err := net.Listen("tcp", ":"+listenPort)
	if err != nil {
		fmt.Println("Error starting proxy: ", err)
		os.Exit(1)
	}
	fmt.Println("Proxy listening on port ", listenPort)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		fmt.Println("Error reading request: ", err)
		return
	}

	// Open the request in the default text editor
	file, err := ioutil.TempFile("", "request")
	if err != nil {
		fmt.Println("Error creating temp file: ", err)
		return
	}
	defer os.Remove(file.Name())

	// Write the request headers and data to the file
	writeRequestToFile(req, file)
	file.Close()

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("start", file.Name())
	} else {
		cmd = exec.Command("open", file.Name())
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("Error opening request in text editor: ", err)
		return
	}

	// Wait for the file to be closed by the text editor
	cmd.Wait()

	// Read the updated request from the file
	file, err = os.Open(file.Name())
	if err != nil {
		fmt.Println("Error opening temp file: ", err)
		return
	}
	defer file.Close()

	req, err = http.ReadRequest(bufio.NewReader(file))
	if err != nil {
		fmt.Println("Error reading updated request: ", err)
		return
	}

	// Send the updated request to the remote server
	transport := &http.Transport{
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending updated request: ", err)
		return
	}
	defer res.Body.Close()

	// Copy the response back to the client
	res.Write(conn)
}

func writeRequestToFile(req *http.Request, file *os.File) {
	req.Write(file)
	file.WriteString("\r\n")
	if req.Body != nil {
		if req.ContentLength > 0 {
			file.WriteString("Content-Length: " + strconv.FormatInt(req.ContentLength, 10) + "\r\n")
		}
		file.WriteString("\r\n")
		io.Copy(file, req.Body)
		req.Body.Close()

		file.Seek(0, io.SeekStart)
		body, _ := ioutil.ReadAll(file)
		req.Body = ioutil.NopCloser(bytes.NewReader(body))
		req.ContentLength = int64(len(body))
	}
}
