package logutil

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

//SendLog sends the log message over tcp and throws an error if the log message is an error.
//If a log file is given it will write the data to the log file
func SendLog(address string, isErr bool, data []string, logFile *os.File, id string) {
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range data {
		logString := fmt.Sprintf("%s [%s] %s", id, time.Now().Format("Jan 2 2006 15:04:05 MST"), v)
		conn.Write([]byte(logString))
		if logFile != nil {
			WriteToLog([]string{logString}, logFile)
		}
	}
	if isErr {
		log.Fatalln(data)
	}
}

//LogServerRequest creates a summary of the http connection information and send it to the connected logger.
//if a logfile is provided it will also write the log messages to a log file
func LogServerRequest(w http.ResponseWriter, r *http.Request, loggerAddr string, logFile *os.File, id string) {
	method := r.Method
	url := r.URL
	httpVer := r.Proto
	host := r.Host
	address := r.RemoteAddr
	reqData := fmt.Sprintf("%s %s %s %s %s", address, method, url, httpVer, host)
	SendLog(loggerAddr, false, []string{reqData}, logFile, id)
}

//WriteToLog writes the data passed into data to the given file
func WriteToLog(data []string, file *os.File) {
	for _, v := range data {
		//prepare string for writing to file
		v = strings.ToValidUTF8(v, "")

		_, err := file.Write([]byte(v))
		if err != nil {
			log.Fatalln(err)
		}
	}
}

//CreateLogServerAndListen runs a tcp server at address:port
func CreateLogServerAndListen(address string, port string, logFile *os.File) {
	//make sure the port number is in the format ":####"
	if !strings.ContainsAny(port, ":") {
		port = ":" + port
	}
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Logger is listening at %s\n", address)
	defer l.Close()

	//wait for a connection
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		buffer := make([]byte, 1024)
		conn.Read(buffer)

		//write the contents of buffer to the log file
		bufferText := string(buffer) + "\n"
		fmt.Print(bufferText)
		WriteToLog([]string{bufferText}, logFile)

		go func(c net.Conn) {
			c.Write(buffer)
			c.Close()
		}(conn)
	}
}

//OpenLogFile opens the log file stored in path.
//If the file doesn't exist it is created
func OpenLogFile(path string) *os.File {
	//if path does not end in a slash, add it
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	date := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("%s%s.txt", path, date)
	//opens file with options to append string on write, and open in write only mode
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	return logFile
}
