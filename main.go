package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"github.com/a-h/templ"

	"remote-shell/templates"
)

func setupServer() {
	mux := http.NewServeMux()
	handlerManager(mux)

	if err := http.ListenAndServe(":9099", mux); err != nil {
		fmt.Println("failed to start server")
	}
}

func handlerManager(mux *http.ServeMux) {

	homePage := templates.HomePage()
	mux.Handle("/", templ.Handler(homePage))
	mux.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Write([]byte("Not Found"))
		}
		ctx := r.Context()

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
		}
		defer r.Body.Close()

		formattedBody := strings.Split(string(body), "input=")
		terminalCmd, err := url.QueryUnescape(formattedBody[1])
		if err != nil {
			log.Println(err)
		}

		vals := strings.Split(terminalCmd, " ")

		output := execCmd(vals[0], vals[1:])

		outputTempl := templates.Output(output)
		outputTempl.Render(ctx, w)
	})

	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func execCmd(name string, args []string) string {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	if err != nil {
		log.Println(err)
		fmt.Println("failed to run command")
	}

	return string(output)
}

func main() {
	fmt.Println("Remote shell executor")

	setupServer()
}
