package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	client       = &http.Client{}
	mode         string
	token        string
	sleepAmount  int
	amountCycled = 0
)

type statusStruct struct {
	CustomStatus struct {
		Text string `json:"text"`
	} `json:"custom_status"`
}

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

func main() {
	fmt.Print("Enter your token: ")
	fmt.Scan(&token)
	fmt.Print("Enter the duration: ")
	fmt.Scan(&sleepAmount)
	fmt.Print("Finally, enter your mode (cycle/random/progression): ")
	fmt.Scan(&mode)
	if !checkToken() {
		panic("Not a proper token.")
	}
	for {
		stats := getStatus()
		switch mode {
		case "random":
			changeStatus(stats[rand.Intn(len(stats))])
		case "cycle":
			if amountCycled == len(stats) {
				amountCycled = 0
			}
			changeStatus(stats[amountCycled])
			amountCycled++
		case "progression":
			if amountCycled == 5 {
				amountCycled = 1
			}
			changeStatus(strings.Repeat(stats[1], amountCycled))
			amountCycled++
		default:
			fmt.Println("Not a mode")
			changeStatus(stats[0])
			time.Sleep(time.Duration(5) * time.Second)
			os.Exit(0)
		}
		time.Sleep(time.Duration(sleepAmount) * time.Second)
	}
}

func checkToken() bool {
	req, _ := http.NewRequest("GET", "https://discordapp.com/api/v8/users/@me", nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.130 Safari/537.36")
	req.Header.Set("Authorization", token)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	decodedJSON := decode(body)
	if len(decodedJSON["username"]) > 0 {
		return true
	}
	return false
}

func changeStatus(status string) {
	statusstruct := statusStruct{}
	statusstruct.CustomStatus.Text = status
	bytes := new(bytes.Buffer)
	json.NewEncoder(bytes).Encode(statusstruct)
	req, _ := http.NewRequest("PATCH", "https://discordapp.com/api/v8/users/@me/settings", bytes)
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", token)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.130 Safari/537.36")
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}

func getStatus() map[int]string {
	statuses := make(map[int]string)
	statuses[0] = "Made by https://github.com/not-cyrus"
	file, err := os.Open("Status.txt")
	if err != nil {
		ioutil.WriteFile("Status.txt", []byte(""), 0644)
		fmt.Println("Put your status(es) in the Status.txt file that has just been made and re-launch")
		time.Sleep(time.Duration(5) * time.Second)
		os.Exit(0)
	}
	defer file.Close()
	r := bufio.NewScanner(file)
	for r.Scan() {
		if len(statuses) == 0 {
			statuses[1] = r.Text()
		}
		statuses[len(statuses)] = r.Text()
	}
	return statuses
}

func decode(toDecode []byte) map[string]string {
	var a map[string]string
	json.Unmarshal([]byte(toDecode), &a)
	return a
}
