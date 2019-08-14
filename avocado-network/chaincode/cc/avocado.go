/*
 SPDX-License-Identifier: Apache-2.0
*/

// ====CHAINCODE EXECUTION SAMPLES (CLI) ==================
package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type TransportData struct{
	ObjectType string `json:"docType"`
	Make string `json:"make"`
	Model string `json:"model"`
	Type string `json:"type"`
	Owner string `json:"owner"`
	RefrigeratedTemp string `json:"refrigeratedTemp"`
} 

type box struct {
	ObjectType string `json:"docType"` //docType is used to distinguish the various types of objects in state database
	BoxID      string `json:"boxID"`    //the fieldtags are needed to keep case from bouncing around
	ProducerName string `json:"producerName"`
	DateTimeCreated	string `json:"dateTimeCreated"`
	DateTimeCooled string `json:"dateTimeCooled"`
	TargetTemp string `json:"targetTemp"`
	TotalAvocados int `json:"totalAvocados"`
	TransportDataAsLot TransportData `json:"TransportDataAsLot"`
}

// ===================================================================================
// Main
// ===================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init initializes chaincode
// ===========================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

// Invoke - Our entry point for Invocations
// ========================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "sortToBox" { //create a new box
		return t.sortToBox(stub, args)
	} else if function == "getBox" { //change owner of a specific box
		return t.getBox(stub, args)
	} else if function == "preCoolBox" { //transfer all boxs of a certain color
		return t.preCoolBox(stub, args)
	} else if function == "LoadForRefigeratedTransport" { //transfer all boxs of a certain color
		return t.LoadForRefigeratedTransport(stub, args)
	} 

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// initMarble - create a new box, store into chaincode state
// ============================================================
func (t *SimpleChaincode) sortToBox(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	//   0       		1       		2     			 3
	// "boxID", "producerName", "dateTimeCreated", "totalAvocados"
	if len(args) != 4 {
		return shim.Error("Incorrect number of arguments. Expecting 4")
	}

	// ==== Input sanitation ====
	fmt.Println("- start sorting to box")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty numeric string")
	}

	boxID := args[0]
	producerName := args[1]
	dateTimeCreated := args[2]
	totalAvocados, err := strconv.Atoi(args[3])
	if err != nil {
		return shim.Error("4th argument must be a numeric string")
	}

	// ==== Check if box already exists ====
	boxAsBytes, err := stub.GetState(boxID)
	if err != nil {
		return shim.Error("Failed to get box: " + err.Error())
	} else if boxAsBytes != nil {
		fmt.Println("This box already exists: " + boxID)
		return shim.Error("This box already exists: " + boxID)
	}

	// ==== Create box object and marshal to JSON ====
	objectType := "box"
	dateTimeCooled := ""
	targetTemp := ""
	transportDataAsLot := TransportData{"transportData","","","","",""}
	box := &box{objectType, boxID, producerName, dateTimeCreated, dateTimeCooled, targetTemp, totalAvocados, transportDataAsLot} // creates the json struct for the box
	boxJSONasBytes, err := json.Marshal(box)
	if err != nil {
		return shim.Error(err.Error())
	}
	//Alternatively, build the box json string manually if you don't want to use struct marshalling
	//boxJSONasString := `{"docType":"Marble",  "name": "` + boxName + `", "color": "` + color + `", "size": ` + strconv.Itoa(size) + `, "owner": "` + owner + `"}`
	//boxJSONasBytes := []byte(str)

	// === Save box to state ===
	err = stub.PutState(boxID, boxJSONasBytes) //saves it by boxID
	if err != nil {
		return shim.Error(err.Error())
	}

	// ==== Box saved Return success ====
	fmt.Println("- end init box")
	return shim.Success(nil)
}

// ===============================================
// getBox - read a box from chaincode state
// ===============================================
func (t *SimpleChaincode) getBox(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var boxID, jsonResp string
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting ID of the box to query")
	}

	boxID = args[0]
	valAsbytes, err := stub.GetState(boxID) //get the box from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + boxID + "\"}"
		return shim.Error(jsonResp)
	} else if valAsbytes == nil {
		jsonResp = "{\"Error\":\"Marble does not exist: " + boxID + "\"}"
		return shim.Error(jsonResp)
	}

	return shim.Success(valAsbytes)
}


// ===========================================================
// Pre cool a box by setting its previously unset parameters
// ===========================================================
func (t *SimpleChaincode) preCoolBox(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	//    0       		1				2
	// "boxID", "dateTimeCooled", "targetTemp"
	if len(args) < 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	boxID := args[0]
	dateTimeCooled := args[1]
	targetTemp := args[2]
	fmt.Println("- start preCoolBox ", boxID)

	boxAsBytes, err := stub.GetState(boxID)
	if err != nil {
		return shim.Error("Failed to get box:" + err.Error())
	} else if boxAsBytes == nil {
		return shim.Error("Box does not exist")
	}

	boxToCool := box{}
	err = json.Unmarshal(boxAsBytes, &boxToCool) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	boxToCool.DateTimeCooled = dateTimeCooled //add the date time cooled
	boxToCool.TargetTemp = targetTemp

	boxJSONasBytes, _ := json.Marshal(boxToCool)
	err = stub.PutState(boxID, boxJSONasBytes) //rewrite the box
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end preCoolBox (success)")
	return shim.Success(nil)
}
// ===========================================================
// Load for refrigerated transport as lot
// ===========================================================
func (t *SimpleChaincode) LoadForRefigeratedTransport(stub shim.ChaincodeStubInterface, args []string) pb.Response{
	// Make string `json:"make"`
	// Model string `json:"model"`
	// Type string `json:"type"`
	// Owner string `json:"owner"`
	// RefrigeratedTemp

	//   0        1        2       3		4				5
	// "boxID", "Make", "Model", "Type", "Owner", "Refrigerated Temperature"

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	fmt.Println("- start sorting to box")
	if len(args[0]) <= 0 {
		return shim.Error("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return shim.Error("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return shim.Error("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return shim.Error("4th argument must be a non-empty numeric string")
	}
	if len(args[4]) <= 0 {
		return shim.Error("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		return shim.Error("6th argument must be a non-empty numeric string")
	}
	boxID := args[0]
	make := args[1]
	model := args[2]
	typeOfVehicle := args[3]
	owner:= args[4]
	refrigeratedTemp := args[5]

	boxAsBytes, err := stub.GetState(boxID)
	if err != nil {
		return shim.Error("Failed to get box:" + err.Error())
	} else if boxAsBytes == nil {
		return shim.Error("Box does not exist")
	}

	boxToTransport := box{}
	err = json.Unmarshal(boxAsBytes, &boxToTransport) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}
	transportData := TransportData{"transportData",make,model,typeOfVehicle,owner,refrigeratedTemp}
	boxToTransport.TransportDataAsLot = transportData 


	boxJSONasBytes, _ := json.Marshal(boxToTransport)
	err = stub.PutState(boxID, boxJSONasBytes) //rewrite the box
	if err != nil {
		return shim.Error(err.Error())
	}
	fmt.Println("- end Load for refrigerated transport as lot (success)")
	return shim.Success(nil)
}


