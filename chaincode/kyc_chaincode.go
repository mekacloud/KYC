/*
Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

var customerIndexStr = "_marbleindex" //name for the key/value that will store a list of all known marbles
var openTradesStr = "_opentrades"     //name for the key/value that will store all open trades

type Customer struct {
	Name        string   `json:"name"` //the fieldtags are needed to keep case from bouncing around
	TelNo       string   `json:"telno"`
	Age         int      `json:"age"`
	Ocupation   string   `json:"ocupation"`
	AllowBroke  []Broker `json:"allowbroke"`
	GauranteeID string   `json:"gauranteeid"`
}

type Broker struct {
	Name     string `json:"name"`
	BrokerNo int    `json:"brokerno"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// ============================================================================================================================
// Init - reset all the things
// ============================================================================================================================
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var Aval int
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// Initialize the chaincode
	Aval, err = strconv.Atoi(args[0])
	if err != nil {
		return nil, errors.New("Expecting integer value for asset holding")
	}

	// Write the state to the ledger
	err = stub.PutState("kyc", []byte(strconv.Itoa(Aval))) //making a test var "kyc", I find it handy to read/write to it right away to test the network
	if err != nil {
		return nil, err
	}

	var empty []string
	jsonAsBytes, _ := json.Marshal(empty) //marshal an emtpy array of strings to clear the index
	err = stub.PutState(customerIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	/*var trades AllTrades
	jsonAsBytes, _ = json.Marshal(trades) //clear the open trade struct
	err = stub.PutState(openTradesStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	*/
	return nil, nil
}

// ============================================================================================================================
// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// ============================================================================================================================
// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
		/*} else if function == "delete" { //deletes an entity from its state
		res, err := t.Delete(stub, args)
		cleanTrades(stub) //lets make sure all open trades are still valid
		return res, err*/
	} else if function == "write" { //writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "new_customer" { //create a new marble
		return t.new_customer(stub, args)
	} else if function == "set_user" { //change owner of a marble
		res, err := t.set_user(stub, args)
		//cleanTrades(stub) //lets make sure all open trades are still valid
		return res, err
	} else if function == "open_trade" { //create a new trade order
		//return t.open_trade(stub, args)
	} else if function == "perform_trade" { //forfill an open trade order
		// res, err := t.perform_trade(stub, args)
		// cleanTrades(stub) //lets clean just in case
		// return res, err
	} else if function == "remove_trade" { //cancel an open trade order
		// return t.remove_trade(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// ============================================================================================================================
// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// ============================================================================================================================
// Read - read a variable from chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Write - write variable into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) Write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var name, value string // Entities
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0] //rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// ============================================================================================================================
// Init Customer - create a new marble, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) new_customer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1       2     3
	// "asdf", "blue", "35", "bob"
	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init marble")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		return nil, errors.New("4th argument must be a non-empty string")
	}
	name := args[0]
	telno := strings.ToLower(args[1])
	size, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	ocupation := strings.ToLower(args[3])

	//check if marble already exists
	marbleAsBytes, err := stub.GetState(name)
	if err != nil {
		return nil, errors.New("Failed to get marble name")
	}
	res := Customer{}
	json.Unmarshal(marbleAsBytes, &res)
	if res.Name == name {
		fmt.Println("This marble arleady exists: " + name)
		fmt.Println(res)
		return nil, errors.New("This marble arleady exists") //all stop a marble by this name exists
	}

	res.Name = name
	res.TelNo = telno
	res.Age = size
	res.Ocupation = ocupation
	//build the marble json string manually
	//str := `{"name": "` + name + `", "telno": "` + telno + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	//err = stub.PutState(name, []byte(str)) //store marble with id as key
	str, err := json.Marshal(res)
	err = stub.PutState(name, str)
	if err != nil {
		return nil, err
	}

	//get the marble index
	marblesAsBytes, err := stub.GetState(customerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get marble index")
	}
	var marbleIndex []string
	json.Unmarshal(marblesAsBytes, &marbleIndex) //un stringify it aka JSON.parse()

	//append
	marbleIndex = append(marbleIndex, name) //add marble name to index list
	fmt.Println("! marble index: ", marbleIndex)
	jsonAsBytes, _ := json.Marshal(marbleIndex)
	err = stub.PutState(customerIndexStr, jsonAsBytes) //store name of marble

	fmt.Println("- end init marble")
	return nil, nil
}

// ============================================================================================================================
// Set User Permission on Customer
// ============================================================================================================================
func (t *SimpleChaincode) set_user(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0       1
	// "name", "bob"
	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	fmt.Println("- start set user")
	fmt.Println(args[0] + " - " + args[1])
	marbleAsBytes, err := stub.GetState(args[0])
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	res := Customer{}
	json.Unmarshal(marbleAsBytes, &res) //un stringify it aka JSON.parse()
	//res.User = args[1]                  //change the user

	jsonAsBytes, _ := json.Marshal(res)
	err = stub.PutState(args[0], jsonAsBytes) //rewrite the marble with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end set user")
	return nil, nil
}
