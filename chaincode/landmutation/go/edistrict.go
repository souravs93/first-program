package main

import (
	"encoding/json"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

/* Define the Land Mutation Application structure, with 6 properties.
   Structure tags are used by encoding/json library. Key consist of
   prefix + ApplicationID
*/
type LandMutationApplication struct {
	ApplicationID     string `json:"application_id"`
	AadharID          string `json:"aadhar_id"`
	UserName          string `json:"user_name"`
	PlotNumber        string `json:"plot_number"`
	DateOfApplication string `json:"date_of_application"`
	AssignTo          string `json:"assign_to"`
	Status            string `json:"status"`
}

/* Define the Citizen structure, with 5 properties.
   Structure tags are used by encoding/json library. Key consist of
   prefix +
*/
type Citizen struct {
	AadharID string `json:"aadhar_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	LastName string `json:"last_name"`
	Address  string `json:"address"`
}

func createLMA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid argument count.")
	}

	lma := LandMutationApplication{}
	err := json.Unmarshal([]byte(args[0]), &lma)
	if err != nil {
		return shim.Error(err.Error())
	}

	citizenKey, err := stub.CreateCompositeKey(prefixCitizen, []string{lma.AadharID})
	// Check if a user with the same username exists
	if err != nil {
		return shim.Error(err.Error())
	}
	citizenAsBytes, _ := stub.GetState(citizenKey)
	if citizenAsBytes == nil {
		return shim.Error("Citizen with this username does not exist.")
	}

	key, err := stub.CreateCompositeKey(prefixLMA, []string{lma.ApplicationID})
	if err != nil {
		return shim.Error(err.Error())
	}

	// Check if the user already exists
	lmaAsBytes, _ := stub.GetState(key)
	// User does not exist, attempting creation
	if len(lmaAsBytes) == 0 {
		lmaAsBytes, err = json.Marshal(lma)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(key, lmaAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return nil, if user is newly created
		return shim.Success(nil)
	}

	err = json.Unmarshal(lmaAsBytes, &lma)
	if err != nil {
		return shim.Error(err.Error())
	}

	lmaResponse := struct {
		ApplicationID string `json:"application_id"`
	}{
		ApplicationID: lma.ApplicationID,
	}

	lmaResponseAsBytes, err := json.Marshal(lmaResponse)
	if err != nil {
		return shim.Error(err.Error())
	}
	// Return the username and the password of the already existing user
	return shim.Success(lmaResponseAsBytes)
}

func listLMA(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	input := struct {
		ApplicationID string `json:"application_id"`
	}{}
	if len(args) == 1 {
		err := json.Unmarshal([]byte(args[0]), &input)
		if err != nil {
			return shim.Error(err.Error())
		}
	}
	filterByApplicationID := len(input.ApplicationID) > 0

	var resultsIterator shim.StateQueryIteratorInterface
	var err error
	// Filtering by username if required
	if filterByApplicationID {
		resultsIterator, err = stub.GetStateByPartialCompositeKey(prefixLMA, []string{input.ApplicationID})
	} else {
		resultsIterator, err = stub.GetStateByPartialCompositeKey(prefixLMA, []string{})
	}
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	results := []interface{}{}
	// Iterate over the results
	for resultsIterator.HasNext() {
		kvResult, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}

		// Construct response struct
		result := struct {
			*LandMutationApplication
		}{}

		err = json.Unmarshal(kvResult.Value, &result)
		if err != nil {
			return shim.Error(err.Error())
		}

		// Fetch key
		prefix, keyParts, err := stub.SplitCompositeKey(kvResult.Key)
		if len(keyParts) == 2 {
			result.ApplicationID = keyParts[1]
		} else {
			result.ApplicationID = prefix
		}
		results = append(results, result)
	}

	resultsAsBytes, err := json.Marshal(results)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(resultsAsBytes)
}

func createCitizen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid argument count.")
	}

	citizen := Citizen{}
	err := json.Unmarshal([]byte(args[0]), &citizen)
	if err != nil {
		return shim.Error(err.Error())
	}

	key, err := stub.CreateCompositeKey(prefixCitizen, []string{citizen.AadharID})
	if err != nil {
		return shim.Error(err.Error())
	}

	// Check if the user already exists
	citizenAsBytes, _ := stub.GetState(key)
	// User does not exist, attempting creation
	if len(citizenAsBytes) == 0 {
		citizenAsBytes, err = json.Marshal(citizen)
		if err != nil {
			return shim.Error(err.Error())
		}

		err = stub.PutState(key, citizenAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}

		// Return nil, if user is newly created
		return shim.Success(nil)
	}

	err = json.Unmarshal(citizenAsBytes, &citizen)
	if err != nil {
		return shim.Error(err.Error())
	}

	citizenResponse := struct {
		AadharID string `json:"aadhar_id"`
	}{
		AadharID: citizen.AadharID,
	}

	citizenResponseAsBytes, err := json.Marshal(citizenResponse)
	if err != nil {
		return shim.Error(err.Error())
	}
	// Return the username and the password of the already existing user
	return shim.Success(citizenResponseAsBytes)
}

func getCitizen(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid argument count.")
	}

	input := struct {
		AadharID string `json:"aadhar_id"`
	}{}

	err := json.Unmarshal([]byte(args[0]), &input)
	if err != nil {
		return shim.Error(err.Error())
	}

	userKey, err := stub.CreateCompositeKey(prefixCitizen, []string{input.AadharID})
	if err != nil {
		return shim.Error(err.Error())
	}
	userBytes, _ := stub.GetState(userKey)
	if len(userBytes) == 0 {
		return shim.Success(nil)
	}

	response := struct {
		AadharID string `json:"aadhar_id"`
		UserName string `json:"user_name"`
		Password string `json:"password"`
		LastName string `json:"last_name"`
		Address  string `json:"address"`
	}{}
	err = json.Unmarshal(userBytes, &response)
	if err != nil {
		return shim.Error(err.Error())
	}
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(responseBytes)
}

func citizenAcceptHearingDate(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("Invalid Arguments Count.")
	}

	input := struct {
		ApplicationID     string `json:"application_id"`
		AcceptHearingDate bool   `json:"accept_hearing_date"`
	}{}
	err := json.Unmarshal([]byte(args[0]), &input)

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

	if lma.AssignTo == "Citizen" {
		if input.AcceptHearingDate {
			lma.AssignTo = "EstateOfficer"
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
