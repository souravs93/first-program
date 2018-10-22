package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// EstateManager
type EstateManager struct {
	EstateManagerID string `json:"id"`

	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	DepartmentName string `json:"department_name"`
	Address        string `json:"address"`
}

func processLMAEstateManager(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid Arguments Count.")
	}

	input := struct {
		ApplicationID       string `json:"application_id"`
		EstateMangerComment string `json:"comment"`
		EstateMangerAction  string `json:"action"`
		DateOfHearing       string `json:"date_of_hearing"`
	}{}
	err := json.Unmarshal([]byte(args[0]), &input)

	// Estate  Manager Comment
	estateManagerComment := []string{input.EstateMangerComment}
	estateManagerComment = append(estateManagerComment, "true")

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

	if lma.AssignTo == "EstateManager" {
		if input.EstateMangerAction == "SetHearingDate" || input.EstateMangerAction == "ApplicationSentForCorrection" {
			lma.AssignTo = "Citizen"
			lma.Status = "Inprogress"
		} else if input.EstateMangerAction == "ApplicationRejected" {
			lma.AssignTo = ""
			lma.Status = "Rejected"
		}
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

func estateManagerHearing(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid Arguments Count.")
	}

	input := struct {
		ApplicationID        string `json:"application_id"`
		EstateManagerComment string `json:"comment"`
	}{}
	err := json.Unmarshal([]byte(args[0]), &input)

	// Estate Manager Hearing Comment
	estateManagerHearingComment := []string{input.EstateManagerComment}
	estateManagerHearingComment = append(estateManagerHearingComment, "true")

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

	if lma.AssignTo == "EstateManager" {
		lma.AssignTo = "CEO"
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
