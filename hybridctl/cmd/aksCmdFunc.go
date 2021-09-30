package cmd

import (
	util "Hybrid_Cluster/hybridctl/util"
	"fmt"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

func aksStart(p util.EksAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksStart"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	if string(bytes) == "" {
		fmt.Println("Succeeded to start", p.ResourceName, "in", p.ResourceGroupName)
	} else {
		fmt.Println(string(bytes))
	}
}

func aksStop(p util.EksAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksStop"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	if string(bytes) == "" {
		fmt.Println("Succeeded to stop", p.ResourceName, "in", p.ResourceGroupName)
	} else {
		fmt.Println(string(bytes))
	}
}
