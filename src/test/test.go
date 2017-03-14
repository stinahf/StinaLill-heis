package main

import(
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type dumb struct{
	int2 int
	string2 string
}


func DumpToHardware(d dumb) {
	file, _ := json.Marshal(d)
	err := ioutil.WriteFile("reg_redundancy", file, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func LoadFromHardware() dumb{
	var d dumb
	file, err := ioutil.ReadFile("reg_redundancy")
	if err != nil {
		return dumb{}
	}
	if err := json.Unmarshal(file, &d); err != nil {
		return dumb{}
	}
	return d
}

func main(){
	dumbest := dumb{1,"yo"}
	DumpToHardware(dumbest)
	fmt.Println(LoadFromHardware())
}