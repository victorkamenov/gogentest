//go:generate go run gen.go
//go:generate go run main.go

package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var lastOut []byte

func main() {
	replace()
	lastOut, _ = ioutil.ReadFile("out")
	f, err := os.Create("main.go")
	die(err)
	defer f.Close()

	program := getPre() + getBody() + getPost()
	_, werror := f.WriteString(program)
	die(werror)

}

func getBody() string {

	return `

	fmt.Println("my LAST OUTPUT WAS:  ` + fmt.Sprintf("%s", lastOut) + `")
	` + getOut()
}

func die(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func getPost() string {
	return `
	executeCmd("go", "generate")
}

func executeCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	}

	defer stdOut.Close()

	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StderrPipe for Cmd", err)
	}

	defer stdErr.Close()

	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScanner.Scan() {

			txt := stdErrScanner.Text()

			if !strings.Contains(txt, "no buildable Go source files in") {
				fmt.Printf("%s\n", txt)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	// go generate command will fail when no generate command find.
	if err != nil {
		if err.Error() != "exit status 1" {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
}
`
}

func getPre() string {
	return `
//go:generate go run gen.go
//go:generate go run main.go ` + fmt.Sprintf("%s", lastOut) + " " + time.Now().Format(time.ANSIC) + `
package main
			
import ("fmt"
"log"
"os"
	
"strings"
"os/exec"
"bufio"
)

func main() {
`
}

func getOut() string {
	return `
	f, _ := os.Create("out")
	defer f.Close()
	f.WriteString(fmt.Sprintf("%s", os.Args[1:]))
	f.Sync()
	`
}

func replace() {
	if _, err := os.Stat("gen.go.old"); !os.IsNotExist(err) {
		return
	}
	input, err := ioutil.ReadFile("gen.go")
	if err != nil {
		log.Fatalln(err)
	}
	err = ioutil.WriteFile("gen.go.old", input, 0644)

	lines := strings.Split(string(input), "\n")

	lines[1] = ""
	lines[0] = "// +build ignore"
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile("gen.go", []byte(output), 0644)

	if err != nil {
		log.Fatalln(err)
	}
}

func executeCmd(command string, args ...string) {
	cmd := exec.Command(command, args...)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StdoutPipe for Cmd", err)
	}

	defer stdOut.Close()

	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s\n", scanner.Text())
		}
	}()

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(os.Stderr, "Error creating StderrPipe for Cmd", err)
	}

	defer stdErr.Close()

	stdErrScanner := bufio.NewScanner(stdErr)
	go func() {
		for stdErrScanner.Scan() {

			txt := stdErrScanner.Text()

			if !strings.Contains(txt, "no buildable Go source files in") {
				fmt.Printf("%s\n", txt)
			}
		}
	}()

	err = cmd.Start()
	if err != nil {
		log.Fatal(os.Stderr, "Error starting Cmd", err)
	}

	err = cmd.Wait()
	// go generate command will fail when no generate command find.
	if err != nil {
		if err.Error() != "exit status 1" {
			fmt.Println(err)
			log.Fatal(err)
		}
	}
}
