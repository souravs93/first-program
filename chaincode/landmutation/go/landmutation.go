package main

import (
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

	if function == "init" {
		return t.Init(stub)
	}
	bcFunc := bcFunctions[function]
	if bcFunc == nil {
		return shim.Error("Invalid invoke function.")
	}
	return bcFunc(stub, args)
}

func main() {
	logger.SetLevel(shim.LogInfo)

	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
