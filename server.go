package main

import (
	"fmt"
	"encoding/json"
	"net"
	"bufio"
	"strings"
)

type task struct {
	User	string
	Importance	int
	Task	string
	Duedate	string
	Duetime	string
	Completiontime	int
}

func doTheThing(connection net.Conn){
	recept, err := bufio.NewReader(connection).ReadString('\n')
	if (err != nil){
		connection.Close()
		return
	}
	strs := strings.Split(recept,"{") //Splits the string along the { to divide up the json
	jason := make([]task,0) //stores the objects that are parsed in
	if (recept[0] == 49 || recept[0] == 50){ //Ascii value for 1 (change after testing)
		for i := 1; i < len(strs);i++{
			//This loop goes through the split string and creates json objects
			strs[i] = "{" + strs[i]
			tmp := task{}
			json.Unmarshal([]byte(strs[i]),&tmp)
			jason = append(jason,tmp)
		}
	} else {
		//Send the list of tasks
	}
	fmt.Println(jason)
	connection.Close()
}

func main(){
	ln, _ := net.Listen("tcp",":6666")
	for {
		conn, err := ln.Accept()
		if(err != nil){
			continue
		}
		go doTheThing(conn)
	}
}
