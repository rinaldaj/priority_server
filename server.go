package main

import (
	"fmt"
	"encoding/json"
	"net"
	"bufio"
	"strings"
	"os"
	"sync"
)

type task struct {
	User	string
	Importance	int
	Task	string
	Duedate	string
	Duetime	string
	Completiontime	int
}
	var mutx = &sync.Mutex{} //makes sure that threads aren't fighting over jason
	var jason = make([]task,0) //stores the objects that are parsed in

func eq(x task, y task) bool{
	return strings.EqualFold(x.User, y.User) && x.Importance == y.Importance && strings.EqualFold(x.Task, y.Task) && strings.EqualFold(x.Duedate,y.Duedate) && strings.EqualFold(x.Duetime,y.Duetime) && x.Completiontime == y.Completiontime
}

func doTheThing(connection net.Conn){
	recept, err := bufio.NewReader(connection).ReadString('\n')
	if (err != nil){
		connection.Close()
		return
	}
	strs := strings.Split(recept,"{") //Splits the string along the { to divide up the json
	if (recept[0] == 49 || recept[0] == 50){ //Ascii value for 1 (change after testing)
		for i := 1; i < len(strs);i++{
			//This loop goes through the split string and creates json compatible strings
			strs[i] = "{" + strs[i]
		}
	} else {
		//Send the list of tasks
		fmt.Println("test")
	}

	if (recept[0] == 49){
		for i:=1;i<len(strs);i++{
			//parses strings and adds them to json
			tmp := task{}
			json.Unmarshal([]byte(strs[i]),&tmp)
			mutx.Lock()
			jason = append(jason,tmp)
			mutx.Unlock()
		}
	} else {
		mutx.Lock()
		tmper := make([]task,0)
		for i:=0;i<len(jason);i++ {
			for j:=1;j<len(strs);j++{
				tmp := task{}
				json.Unmarshal([]byte(strs[j]),&tmp)
				if(!eq(tmp,jason[i])){
					tmper = append(tmper,jason[i])
				}
			}
		}
		jason = tmper
		mutx.Unlock()
	}
	fmt.Println(jason)
	connection.Close()
	return
}

func main(){
	args := os.Args
	if(len(args) != 2){
		fmt.Println("Useage:",args[0]," <port>")
		return
	}
	ln, errs := net.Listen("tcp", ":" + args[1])
	if (errs != nil){
		fmt.Println("something went wrong when binding the socket")
		return
	}
	for {
		conn, err := ln.Accept()
		if(err != nil){
			continue
		}
		go doTheThing(conn)
	}
}
