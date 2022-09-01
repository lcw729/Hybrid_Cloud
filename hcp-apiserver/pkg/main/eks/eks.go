package eks

import (
	"Hybrid_Cloud/hcp-apiserver/pkg/handler"
	"Hybrid_Cloud/hcp-apiserver/pkg/util"
	cobrautil "Hybrid_Cloud/hybridctl/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aws/aws-sdk-go/service/eks"
)

func CreateCluster(w http.ResponseWriter, req *http.Request) {

	var hcpCreateClusterInput cobrautil.HCPCreateClusterInput

	util.Parser(req, &hcpCreateClusterInput)
	out, err := handler.EKSCreateCluster(hcpCreateClusterInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DeleteCluster(w http.ResponseWriter, req *http.Request) {

	var hcpDeleteClusterInput cobrautil.HCPDeleteClusterInput

	util.Parser(req, &hcpDeleteClusterInput)
	fmt.Println(hcpDeleteClusterInput)
	out, err := handler.EKSDeleteCluster(hcpDeleteClusterInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func CreateAddon(w http.ResponseWriter, req *http.Request) {

	var createAddonInput eks.CreateAddonInput

	util.Parser(req, &createAddonInput)
	out, err := handler.EKSCreateAddon(createAddonInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DeleteAddon(w http.ResponseWriter, req *http.Request) {

	var deleteAddonInput eks.DeleteAddonInput

	util.Parser(req, &deleteAddonInput)
	out, err := handler.EKSDeleteAddon(deleteAddonInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DescribeAddon(w http.ResponseWriter, req *http.Request) {

	var describeAddonInput eks.DescribeAddonInput

	util.Parser(req, &describeAddonInput)
	out, err := handler.EKSDescribeAddon(describeAddonInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DescribeAddonVersions(w http.ResponseWriter, req *http.Request) {

	var describeAddonVersionsInput eks.DescribeAddonVersionsInput

	util.Parser(req, &describeAddonVersionsInput)
	out, err := handler.EKSDescribeAddonVersions(describeAddonVersionsInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))

}

func ListAddon(w http.ResponseWriter, req *http.Request) {

	var listAddonInput eks.ListAddonsInput
	util.Parser(req, &listAddonInput)
	out, err := handler.EKSListAddon(listAddonInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func UpdateAddon(w http.ResponseWriter, req *http.Request) {

	var updateAddonInput eks.UpdateAddonInput

	util.Parser(req, &updateAddonInput)
	out, err := handler.EKSUpdateAddon(updateAddonInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func ListUpdate(w http.ResponseWriter, req *http.Request) {
	var listUpdateInput eks.ListUpdatesInput

	util.Parser(req, &listUpdateInput)
	out, err := handler.EKSListUpdate(listUpdateInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DescribeUpdate(w http.ResponseWriter, req *http.Request) {
	var describeUpdateInput eks.DescribeUpdateInput

	util.Parser(req, &describeUpdateInput)
	out, err := handler.EKSDescribeUpdate(describeUpdateInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func ListTagsForResource(w http.ResponseWriter, req *http.Request) {
	var listTagsForResourceInput eks.ListTagsForResourceInput

	util.Parser(req, &listTagsForResourceInput)
	out, err := handler.EKSListTagsForResource(listTagsForResourceInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func AssociateIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var associateIdentityProviderConfigInput eks.AssociateIdentityProviderConfigInput

	util.Parser(req, &associateIdentityProviderConfigInput)
	out, err := handler.EKSAssociateIdentityProviderConfig(associateIdentityProviderConfigInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func AssociateEncryptionConfig(w http.ResponseWriter, req *http.Request) {
	var associateEncryptionConfigInput eks.AssociateEncryptionConfigInput

	util.Parser(req, &associateEncryptionConfigInput)
	out, err := handler.EKSAssociateEncryptionConfig(associateEncryptionConfigInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DisassociateIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var disassociateIdentityProviderConfigInput eks.DisassociateIdentityProviderConfigInput

	util.Parser(req, &disassociateIdentityProviderConfigInput)
	out, err := handler.EKSDisassociateIdentityProviderConfig(disassociateIdentityProviderConfigInput)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func DescribeIdentityProviderConfig(w http.ResponseWriter, req *http.Request) {
	var input eks.DescribeIdentityProviderConfigInput

	util.Parser(req, &input)
	out, err := handler.EKSDescribeIdentityProviderConfig(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func ListIdentityProviderConfigs(w http.ResponseWriter, req *http.Request) {
	var input eks.ListIdentityProviderConfigsInput

	util.Parser(req, &input)
	out, err := handler.EKSListIdentityProviderConfigs(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func TagResource(w http.ResponseWriter, req *http.Request) {
	var input eks.TagResourceInput

	util.Parser(req, &input)
	out, err := handler.EKSTagResource(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func UntagResource(w http.ResponseWriter, req *http.Request) {
	var input eks.UntagResourceInput

	util.Parser(req, &input)
	out, err := handler.EKSUntagResource(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func UpdateClusterConfig(w http.ResponseWriter, req *http.Request) {
	var input eks.UpdateClusterConfigInput

	util.Parser(req, &input)
	out, err := handler.EKSUpdateClusterConfig(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}

func UpdateNodegroupConfig(w http.ResponseWriter, req *http.Request) {
	var input eks.UpdateNodegroupConfigInput

	util.Parser(req, &input)
	out, err := handler.EKSUpdateNodeGroupConfig(input)
	var jsonData []byte
	if err != nil {
		jsonData, _ = json.Marshal(&err)
	} else {
		jsonData, _ = json.Marshal(&out)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(jsonData))
}
