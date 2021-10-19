package cmd

import (
	util "Hybrid_Cluster/hybridctl/util"
	"encoding/json"
	"fmt"
	"log"
)

func checkErr(err error) {
	if err != nil {
		log.Println(err)
	}
}

//addon
func addonDisable(p util.AKSAddon) {
	httpPostUrl := "http://localhost:8080/addonDisable"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func addonEnable(p util.AKSAddon) {
	httpPostUrl := "http://localhost:8080/addonEnable"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func addonList(p util.AKSAddon) {
	httpPostUrl := "http://localhost:8080/addonList"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func addonListAvailable() {
	httpPostUrl := "http://localhost:8080/addonListAvailable"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, nil)
	checkErr(err)
	fmt.Println(string(bytes))
}

func addonShow(p util.AKSAddon) {
	httpPostUrl := "http://localhost:8080/addonShow"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func addonUpdate(p util.AKSAddon) {
	httpPostUrl := "http://localhost:8080/addonUpdate"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

// Pod-Identity
func podIdentityAdd(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityAdd"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityDelete(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityDelete"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityList(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityList"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityExceptionAdd(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityExceptionAdd"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, nil)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityExceptionDelete(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityExceptionDelete"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityExceptionList(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityExceptionList"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func podIdentityExceptionUpdate(p util.AKSPodIdentity) {
	httpPostUrl := "http://localhost:8080/podIdentityExceptionUpdate"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func aksStart(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksStart"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	if string(bytes) == "" {
		fmt.Println("Succeeded to start", p.ResourceName, "in", p.ResourceGroupName)
	} else {
		fmt.Println(string(bytes))
	}
}

func aksStop(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksStop"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	if string(bytes) == "" {
		fmt.Println("Succeeded to stop", p.ResourceName, "in", p.ResourceGroupName)
	} else {
		fmt.Println(string(bytes))
	}
}

func aksRotateCerts(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksRotateCerts"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	fmt.Println(string(bytes))
}

func aksGetOSoptions(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/aksGetOSoptions"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	CloudErr := util.CloudError{}
	errJson := json.Unmarshal(bytes, &CloudErr)
	fmt.Println(errJson)
	if errJson == nil {
		fmt.Println("Success")
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}

}

func maintenanceconfigurationCreateOrUpdate(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/maintenanceconfigurationCreateOrUpdate"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	CloudErr := util.CloudError{}
	errJson := json.Unmarshal(bytes, &CloudErr)
	fmt.Println(errJson)
	if errJson == nil {
		fmt.Println("Success")
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}
}

func maintenanceconfigurationList(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/maintenanceconfigurationList"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	CloudErr := util.CloudError{}
	errJson := json.Unmarshal(bytes, &CloudErr)
	fmt.Println(errJson)
	if errJson == nil {
		fmt.Println("Success")
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}
}

func maintenanceconfigurationDelete(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/maintenanceconfigurationDelete"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	CloudErr := util.CloudError{}
	errJson := json.Unmarshal(bytes, &CloudErr)
	fmt.Println(errJson)
	if errJson == nil {
		fmt.Println("Success")
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}
}

func maintenanceconfigurationShow(p util.EKSAPIParameter) {
	httpPostUrl := "http://localhost:8080/maintenanceconfigurationShow"
	bytes, err := util.GetResponseBody("POST", httpPostUrl, p)
	checkErr(err)
	CloudErr := util.CloudError{}
	errJson := json.Unmarshal(bytes, &CloudErr)
	fmt.Println(errJson)
	if errJson == nil {
		fmt.Println("Success")
		fmt.Println(string(bytes))
	} else {
		fmt.Println(err)
	}
}
