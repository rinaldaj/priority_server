package main

import (
	"fmt"
	"encoding/json"
	"net"
	"bufio"
	"strings"
	"os"
	"sync"
	"io/ioutil"
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
	defer connection.Close()
	defer saveJson()
	recept, err := bufio.NewReader(connection).ReadString('\n')
	if (err != nil){
		return
	}
	strs := strings.Split(recept,"{") //Splits the string along the { to divide up the json
		for i := 1; i < len(strs);i++{
			//This loop goes through the split string and creates json compatible strings
			strs[i] = "{" + strs[i]
		}
	 if (recept[0] == 51 ) {
		//Send the list of tasks
		kero := task{}
		json.Unmarshal([]byte(strs[1]),&kero)
		mutx.Lock()
		sendSting := ""
		for i := 0;i<len(jason);i++ {
			if(kero.User == jason[i].User){
				jss,_ := json.Marshal(jason[i])
				sendSting = sendSting + string(jss)
			}
		}
		mutx.Unlock()
		fmt.Fprint(connection,sendSting)
		return
	}

	if (recept[0] == 49){
		for i:=1;i<len(strs);i++{
			//parses strings and adds them to json
			tmp := task{}
			json.Unmarshal([]byte(strs[i]),&tmp)
			mutx.Lock()
			dup := false //Check for duplicates
			for j := 0; j < len(jason); j++{
				if(eq(tmp,jason[j])){
					dup = true
					break
				}
			}
			if (!dup){
				jason = append(jason,tmp)
			}
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
	return
}

func saveJson(){
	file,_ := os.Create(".jsondump")
	mutx.Lock()
	jazz:=jason
	mutx.Unlock()
	defer file.Close()
	for i := 0; i<len(jazz);i++ {
		jss,_ := json.Marshal(jazz[i])
		if(!eq(task{},jazz[i])){
			fmt.Fprintln(file,string(jss))
		}
	}
	file.Close()
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

	//Read in json from file
	datas, _:=ioutil.ReadFile(".jsondump")
	strs := strings.Split(string(datas),"{") //Splits the string along the { to divide up the json
	for i:=0;i<len(strs);i++{
		strs[i] = "{" + strs[i]
		tmp := task{}
		json.Unmarshal([]byte(strs[i]),&tmp)
		jason=append(jason,tmp)
	}
	for {
		conn, err := ln.Accept()
		if(err != nil){
			continue
		}
		go doTheThing(conn)
	}
}
