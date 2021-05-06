package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	valid := true
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)
	s.Scan()
	add := strings.Split(s.Text(), "|")
	addresses := make(map[string]int)
	for _, elem := range add {
		addresses[strings.Split(elem, ":")[1]] = -1
	}

	for s.Scan() {
		line := s.Text()
		if strings.Contains(line, "Message from") {
			buff := strings.Split(line, "Msg ")

			value := strings.Split(buff[1], " ")[0]
			ip := strings.Split(buff[1], " ")[1]
			// ip = strings.Split(ip, " ")[0]

			intValue, _ := strconv.Atoi(value)
			if intValue != addresses[ip]+1 {
				println(line)
				println("Ip:", ip)
				println("Expected ", addresses[ip]+1)
				println("Got ", intValue)
				valid = false
				break
			} else {
				addresses[ip] += 1
			}
		}
	}
	if !valid {
		fmt.Println("Validation Error!")
	} else {
		fmt.Println("Success!")
	}
	err = s.Err()
	if err != nil {
		log.Fatal(err)
	}
}
