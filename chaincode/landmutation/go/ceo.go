package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// CEO
type CEO struct {
	CEOID          string `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DepartmentName string `json:"department_name"`
	Address        string `json:"address"`
}

func processLMACEO(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid Arguments Count.")
	}

	input := struct {
		ApplicationID     string `json:"application_id"`
		SupervisorComment string `json:"comment"`
	}{}
	err := json.Unmarshal([]byte(args[0]), &input)

	// CEO Comment
	ceoComment := []string{input.SupervisorComment}
	ceoComment = append(ceoComment, "true")

	lmaKey, err := stub.CreateCompositeKey(prefixLMA, []string{input.ApplicationID})
	if err != nil {
		return shim.Error(err.Error())
	}

	lmaBytes, _ := stub.GetState(lmaKey)
	if len(lmaBytes) == 0 {
		return shim.Error("Land Mutation Application ID does not exist")
	}

	lma := LandMutationApplication{}
	err = json.Unmarshal(lmaBytes, &lma)
	if err != nil {
		return shim.Error(err.Error())
	}

	if lma.AssignTo == "CEO" {
		lma.AssignTo = "FinanceOfficer"
	}

	lmaBytes, err = json.Marshal(lma)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = stub.PutState(lmaKey, lmaBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}
