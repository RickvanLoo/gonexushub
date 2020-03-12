package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	callClear()
	// Make a Regex to say we only want letters and numbers
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("-> ")
		input, _ := reader.ReadString('\n')
		// convert CRLF to LF
		input = strings.Replace(input, "\n", "", -1)
		if strings.HasPrefix(input, ":q") {
			os.Exit(0)
		}
		callClear()

		input = reg.ReplaceAllString(input, "-")
		
		//Windows Bullshit
		input = strings.TrimSuffix(input, "-")
		
		input = strings.ToLower(input)
		
		fmt.Println("Searching for: " + input)
		processInput(input)
	}
}

func processInput(input string) {

	EntryPoint := "https://api.nexushub.co/wow-classic/v1/items/mograine-alliance/"

	//Get Item

	url := EntryPoint + input
	nexusClient := http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "gonexushub")

	res, getErr := nexusClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	item := Item{}
	jsonErr := json.Unmarshal(body, &item)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	//
	url = url + "/prices"

	req, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	res, getErr = nexusClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	body, readErr = ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	prices := Prices{}
	jsonErr = json.Unmarshal(body, &prices)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	item.print()
	prices.print()
}

func intToWoWString(i int) string {
	m := i
	copper := m % 100
	m = (m - copper) / 100
	silver := m % 100
	gold := (m - silver) / 100

	copperStr := strconv.Itoa(copper)
	silverStr := strconv.Itoa(silver)
	goldStr := strconv.Itoa(gold)

	return goldStr + "g " + silverStr + "s " + copperStr + "c"
}

var clear map[string]func() //create a map for storing clear funcs

func init() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func callClear() {
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
