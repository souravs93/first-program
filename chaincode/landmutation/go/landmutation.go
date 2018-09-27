package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

const prefixLMA = "lma"
const prefixCitizen = "citizen"

var logger = shim.NewLogger("main")

type SmartContract struct {
}

var bcFunctions = map[string]func(shim.ChaincodeStubInterface, []string) pb.Response{
	// eDistrict
	"lma_create":     createLMA,
	"query_lma":      listLMA,
	"citizen_create": createCitizen,
	"query_citizen":  getCitizen,
	"accept_citizen": citizenAcceptHearingDate,

	// CEO
	"poa_ceo": processLMACEO,

	// eState Manager
	"poa_estate_manager":     processLMAEstateManager,
	"estate_manager_hearing": estateManagerHearing,

	// Supervisor
	"poa_supervisor": processLMASupervisor,

	// Finance Officer
	"poa_finance_officer": processLMAFinanceOfficer,
}

// Init callback representing the invocation of a chaincode
func (t *SmartContract) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke Function accept blockchain code invocations.
func (t *SmartContract) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println(function)
	fmt.Println(args)
	if function == "init" {
		return t.Init(stub)
	}
	if function == "init_ledger" {
		return t.initLedger(stub)
	}

	bcFunc := bcFunctions[function]
	if bcFunc == nil {
		return shim.Error("Invalid invoke function.")
	}
	return bcFunc(stub, args)
}

func (t *SmartContract) initLedger(stub shim.ChaincodeStubInterface) pb.Response {
	user := []Citizen{
		Citizen{AadharID: "a100", UserName: "Sumit", Password: "pass0", LastName: "Maity", Address: "pwc0"},
		Citizen{AadharID: "a101", UserName: "Sourav", Password: "pass1", LastName: "Singh", Address: "pwc1"},
		Citizen{AadharID: "a102", UserName: "Saurabh", Password: "pass2", LastName: "Dargude", Address: "pwc2"},
		Citizen{AadharID: "a103", UserName: "Tirtha", Password: "pass3", LastName: "Ghosh", Address: "pwc3"},
		Citizen{AadharID: "a104", UserName: "Neelabjo", Password: "pass4", LastName: "Mukherjee", Address: "pwc4"},
	}

	i := 0
	for i < len(user) {
		fmt.Println("i is ", i)
		userAsBytes, _ := json.Marshal(user[i])
		key, err := stub.CreateCompositeKey(prefixCitizen, []string{user[i].AadharID})
		if err != nil {
			return shim.Error(err.Error())
		}
		stub.PutState(key, userAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Added:", user[i])
		i = i + 1
	}

	lma := []LandMutationApplication{
		LandMutationApplication{ApplicationID: "0000", AadharID: "a100", UserName: "Sumit", PlotNumber: "p0", DateOfApplication: "01/01/18", AssignTo: "Not_Assigned", Status: "In_Progress"},
		LandMutationApplication{ApplicationID: "0001", AadharID: "a101", UserName: "Sourav", PlotNumber: "p1", DateOfApplication: "02/01/18", AssignTo: "Not_Assigned", Status: "In_Progress"},
		LandMutationApplication{ApplicationID: "0002", AadharID: "a102", UserName: "Saurabh", PlotNumber: "p2", DateOfApplication: "03/01/18", AssignTo: "Not_Assigned", Status: "In_Progress"},
		LandMutationApplication{ApplicationID: "0003", AadharID: "a103", UserName: "Tirtha", PlotNumber: "p3", DateOfApplication: "04/01/18", AssignTo: "Not_Assigned", Status: "In_Progress"},
		LandMutationApplication{ApplicationID: "0004", AadharID: "a104", UserName: "Neelabjo", PlotNumber: "p4", DateOfApplication: "05/01/18", AssignTo: "Not_Assigned", Status: "In_Progress"},
	}

	i = 0
	for i < len(lma) {
		fmt.Println("i is ", i)
		lmaAsBytes, _ := json.Marshal(lma[i])
		key, err := stub.CreateCompositeKey(prefixLMA, []string{lma[i].ApplicationID})
		if err != nil {
			return shim.Error(err.Error())
		}
		stub.PutState(key, lmaAsBytes)
		if err != nil {
			return shim.Error(err.Error())
		}
		fmt.Println("Application Added:", lma[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
