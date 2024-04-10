package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type HttpCloudServer struct{}

type ProcessWrapper struct {
	cmd    *exec.Cmd
	stdin  *io.WriteCloser
	stdout *io.ReadCloser
	stderr *io.ReadCloser
}

type Output struct {
	Type   string `json:"type"`
	Output string `json:"output"`
}

func StartProcess() (*ProcessWrapper, error) {
	//cmd := exec.Command("../mockprocess/mock.out", "../mockprocess/mock.log", "50", "100")
	cmd := exec.Command("../mockprocess/echo.out")
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	cmd.Start()

	pw := &ProcessWrapper{
		cmd:    cmd,
		stdin:  &stdin,
		stdout: &stdout,
		stderr: &stderr,
	}

	return pw, nil
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
	defer conn.CloseNow()

	ctx := context.WithoutCancel(r.Context())

	process, err := StartProcess()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	go func() {
		for {
			buf := make([]byte, 204800)
			n, err := (*process.stdout).Read(buf)
			if err != nil {
				fmt.Println("Error:", err)
				break
			}

			if n > 0 {
				fmt.Println("Output:", string(buf[:n]))
				output := Output{
					Type:   "output",
					Output: string(buf[:n]),
				}

				err = wsjson.Write(ctx, conn, output)
				if err != nil {
					fmt.Println("Error:", err)
					break
				}
			}
		}
	}()

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
			if v["type"] == "input" {
				fmt.Println("Input:", v["input"])
				_, err = (*process.stdin).Write([]byte(v["input"]))
				if err != nil {
					fmt.Println("Error:", err)
					break
				}
			}
		}
	}
}
