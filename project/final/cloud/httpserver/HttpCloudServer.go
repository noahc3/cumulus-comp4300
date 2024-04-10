package httpserver

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type HttpCloudServer struct{}

type ConnectionWrapper struct {
	connection *websocket.Conn
	context    context.Context
	new        bool
}

type ProcessWrapper struct {
	id          string
	cmd         *exec.Cmd
	stdin       *io.WriteCloser
	stdout      *io.ReadCloser
	stderr      *io.ReadCloser
	bufioStdout *bufio.Reader
	bufioStderr *bufio.Reader
}

type WsMessage struct {
	Command string  `json:"command"`
	Output  string  `json:"output"`
	Value   float64 `json:"value"`
	Used    int64   `json:"used"`
	Total   int64   `json:"total"`
}

var activeProcess *ProcessWrapper = nil

var outputHistoryMutex sync.Mutex
var outputHistory string = ""
var newOutputHistory string = ""

var activeConnections []*ConnectionWrapper = nil

var lastCpuStats map[int]*WsMessage = make(map[int]*WsMessage)
var lastMemStats *WsMessage = &WsMessage{}
var lastDiskStats *WsMessage = &WsMessage{}

func RandomUUID() (uuid string) {
	newUUID, err := exec.Command("uuidgen").Output()
	if err != nil {
		log.Fatal(err)
	}

	uuid = string(newUUID)
	return
}

func StartProcess(command string, directory string) error {
	outputHistoryMutex.Lock()
	newOutputHistory += "Starting server...\n"
	outputHistoryMutex.Unlock()

	split := strings.Split(command, " ")

	cmd := exec.Command(split[0], split[1:]...)
	cmd.Dir = directory
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	pw := &ProcessWrapper{
		id:          RandomUUID(),
		cmd:         cmd,
		stdin:       &stdin,
		stdout:      &stdout,
		stderr:      &stderr,
		bufioStdout: bufio.NewReader(stdout),
		bufioStderr: bufio.NewReader(stderr),
	}

	activeProcess = pw

	go ReadOutput()

	return nil
}

func StopProcess() {
	if activeProcess != nil {
		outputHistoryMutex.Lock()
		newOutputHistory += "Stopping server...\n"
		outputHistoryMutex.Unlock()

		(*activeProcess.stdin).Close()
		(*activeProcess.stdout).Close()
		(*activeProcess.stderr).Close()
		activeProcess.cmd.Process.Kill()
		activeProcess = nil
	}
}

func StartLoop() {
	go OutputLoop()
	go CPUStats(1)
	go CPUStats(10)
	go CPUStats(30)
	go CPUStats(60)
	go MemStats()
	go DiskStats()
}

func DiskStats() {
	for {
		cmd := exec.Command("df", "/", "-m")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		buf, err := io.ReadAll(stdout)

		if err != nil {
			fmt.Println("Disk error:", err)
		} else if err == nil {
			pattern := regexp.MustCompile(`(?mi)/dev[^\s]*\s*(\d*)\s*(\d*)`)
			match := pattern.FindStringSubmatch(string(buf))

			if len(match) > 2 {
				totalStr := match[1]
				usedStr := match[2]

				total, err1 := strconv.ParseInt(totalStr, 10, 64)
				used, err2 := strconv.ParseInt(usedStr, 10, 64)

				if err1 == nil && err2 == nil {
					msg := WsMessage{
						Command: "disk",
						Used:    used,
						Total:   total,
					}

					lastDiskStats = &msg

					for _, conn := range activeConnections {
						wsjson.Write(conn.context, conn.connection, msg)
					}
				} else {
					fmt.Println("Disk error:", err1, err2)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func MemStats() {
	for {
		cmd := exec.Command("vmstat", "-s", "-S", "m")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		buf, err := io.ReadAll(stdout)

		if err != nil {
			fmt.Println("Mem error:", err)
		} else {
			pattern := regexp.MustCompile(`\s*(\d*).*\s*(\d*)\sm used memory`)
			match := pattern.FindStringSubmatch(string(buf))

			if len(match) > 2 {
				totalStr := match[1]
				usedStr := match[2]

				total, err1 := strconv.ParseInt(totalStr, 10, 64)
				used, err2 := strconv.ParseInt(usedStr, 10, 64)

				if err1 == nil && err2 == nil {
					msg := WsMessage{
						Command: "mem",
						Used:    used,
						Total:   total,
					}

					lastMemStats = &msg

					for _, conn := range activeConnections {
						wsjson.Write(conn.context, conn.connection, msg)
					}
				} else {
					fmt.Println("Mem error:", err1, err2)
				}
			}
		}
		time.Sleep(1 * time.Second)
	}
}

func CPUStats(interval int) {
	for {
		cmd := exec.Command("mpstat", fmt.Sprintf("%v", interval), "1")
		stdout, _ := cmd.StdoutPipe()
		cmd.Start()

		buf, err := io.ReadAll(stdout)

		if err != nil {
			fmt.Println("Error:", err)
		} else if err == nil {
			pattern := regexp.MustCompile(`(?mi)(^Average.*\s)(\d?\d{1}.\d+$)`)
			match := pattern.FindStringSubmatch(string(buf))

			if len(match) > 2 {
				idleStr := match[2]
				idle, err := strconv.ParseFloat(idleStr, 64)
				if err == nil {
					msg := WsMessage{
						Command: "cpu" + strconv.Itoa(interval),
						Value:   100.0 - idle,
					}

					lastCpuStats[interval] = &msg

					for _, conn := range activeConnections {
						wsjson.Write(conn.context, conn.connection, msg)
					}
				}
			}
		}
	}

}

func ReadOutput() {
	id := activeProcess.id
	for {
		buf := make([]byte, 204800)
		new := ""

		if activeProcess == nil || activeProcess.id != id {
			break
		}

		if activeProcess != nil {
			n, err := (*activeProcess.bufioStdout).Read(buf)
			if err == nil && n > 0 {
				new = string(buf[:n])
				outputHistoryMutex.Lock()
				newOutputHistory += new
				outputHistory += new
				outputHistoryMutex.Unlock()
			} else if err != nil {
				break
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func OutputLoop() {
	for {
		outputHistoryMutex.Lock()
		for _, conn := range activeConnections {
			if conn != nil && conn.new {
				conn.new = false
				msg := WsMessage{
					Command: "output",
					Output:  outputHistory,
				}

				wsjson.Write(conn.context, conn.connection, msg)

				for _, value := range lastCpuStats {
					wsjson.Write(conn.context, conn.connection, value)
				}

				wsjson.Write(conn.context, conn.connection, lastMemStats)
				wsjson.Write(conn.context, conn.connection, lastDiskStats)
			}

			if newOutputHistory != "" {
				msg := WsMessage{
					Command: "output",
					Output:  newOutputHistory,
				}

				wsjson.Write(conn.context, conn.connection, msg)
			}
		}

		newOutputHistory = ""
		outputHistoryMutex.Unlock()

		time.Sleep(100 * time.Millisecond)
	}
}

func ExecuteCommand(command map[string]string) {
	fmt.Println("Executing command:", command)
	if command["command"] == "input" {
		if activeProcess != nil {
			(*activeProcess.stdin).Write([]byte(command["value"]))
		}
	} else if command["command"] == "start" || command["command"] == "restart" {
		StopProcess()
		StartProcess(command["target"], command["directory"])
	} else if command["command"] == "stop" {
		StopProcess()
	}
}

func CreateHttpCloudServer() *HttpCloudServer {
	hcs := &HttpCloudServer{}
	return hcs
}

func (hcs HttpCloudServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("I am accepting your connection")
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"localhost:5173"},
	})

	if err != nil {
		fmt.Println("Error:", err)
		conn.CloseNow()
	}

	ctx := context.WithoutCancel(r.Context())

	activeConnections = append(activeConnections, &ConnectionWrapper{
		connection: conn,
		context:    ctx,
		new:        true,
	})

	for {
		var v map[string]string
		err = wsjson.Read(ctx, conn, &v)
		if err != nil {
			if strings.Contains(err.Error(), "EOF") {
				fmt.Println("Client closed connection")
				break
			}

			fmt.Println("Error:", err)
			return
		} else {
			fmt.Println("Received:", v)
			ExecuteCommand(v)
		}
	}
}
