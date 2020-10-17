// (c) AlenPaulVarghese
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const nekobinURL string = "https://nekobin.com/"
const dogbinURL string = "https://del.dog/"

/*
Available flags :-
		-f : read from a file
				usage : -f test.txt
		-n : use nekobin service
		-d : use dogbin service
WARNING: nekobin flag supersedes dogbin flag.
--------------------------------------------------------------------------
Read from file -->
dogbin: ./pastes -d -f filename.txt or ./pastes -f filename.txt
nekobin: ./pastes -n -f filename.txt
WARNING: make sure file flage `-f` should'nt be placed before other flags.
--------------------------------------------------------------------------
To get short links -->
dogbin: ./pastes https://example.com or ./pastes -d https://example.com
nekobin: ./pastes -n https://example.com
--------------------------------------------------------------------------
For multiline -->
dobin: ./pastes -d or ./pastes
nekobin: ./pastes -n
paste the content in nano editor and save the file without renaming.
--------------------------------------------------------------------------
*/

var responseJSON map[string]interface{}

func main() {
	file := flag.String("f", "", "file path to read from")
	dog := flag.Bool("d", false, "use dogbin")
	neko := flag.Bool("n", false, "use nekobin")
	flag.Parse()
	switch true {
	case *file != "" && *neko:
		nekobin(filereader(*file))
	case *file != "" && *dog, (*file != ""):
		dogbin(filereader(*file))
	default:
		var message string
		if len(os.Args) <= 2 {
			message = nano()
		} else {
			message = fmt.Sprint(os.Args[2:])
		}
		if *neko {
			nekobin(message)
		} else {
			dogbin(message)
		}
	}
}

func nekobin(message string) {
	payload := url.Values{"content": {message}}
	resp, err := http.PostForm(nekobinURL+"api/documents", payload)
	defer resp.Body.Close()
	if err != nil {
		log.Fatal("Failed to connect nekobin server!")
	}
	json.NewDecoder(resp.Body).Decode(&responseJSON)
	fmt.Println("Your Link --> " + nekobinURL + fmt.Sprint(
		responseJSON["result"].(map[string]interface{})["key"]))
}

func dogbin(message string) {
	status := strings.NewReader(message)
	if resp, err := http.Post(dogbinURL+"documents", "text/plain; charset=UTF-8", status); err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(&responseJSON)
		fmt.Println("Your Link --> " + dogbinURL + fmt.Sprint(responseJSON["key"]))
	}
}

func filereader(filename string) string {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(file)
}

func nano() string {
	filename := "PasteYouContentHere"
	cmd := exec.Command("nano", filename)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Start(); err != nil {
		log.Fatal("Please install nano !")
	}
	cmd.Wait()
	if _, err := os.Stat(filename); err != nil {
		log.Fatal(err)
	}
	defer os.Remove(filename)
	return filereader(filename)
}