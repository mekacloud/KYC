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
	"crypto/sha256"
	"encoding/hex"
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

var customerIndexStr = "_customerindex" //name for the key/value that will store a list of all known customers
var brokerIndexStr = "_brokerindex"     //name for the key/value that will store a list of all known customers

// key of customer
var customerKey = "cus_"

// BrokerKey key of broker
var BrokerKey = "bro_"

// GauranteeIDKey key of gaurantee id
var GauranteeIDKey = "gau_"

// CusGauIDKey key of something
var CusGauIDKey = "cgi_"

//Customer is Customer
type Customer struct {
	Name       string `json:"name"` //the fieldtags are needed to keep case from bouncing around
	CardID     string `json:"cardid"`
	TelNo      string `json:"telno"`
	Age        int    `json:"age"`
	Occupation string `json:"occupation"`
	Address    string `json:"address"`
	Creator    string `json:"creator"`
}

//GauranteeID generate from Customer
type GauranteeID struct {
	GauranteeID  string `json:"gauranteeid"`
	CustomerID   string `json:"customerid"`
	AllowBroke   []int  `json:"allowbroke"`
	PendingBroke []int  `json:"pendingbroke"`
}

// type CusGauID struct {
// 	CardID      string      `json:"cardid"`
// 	Customer    Customer    `json:"customer"`
// 	GauranteeID GauranteeID `json:"gauranteeid"`
// }

// Broker Contain Name and Number and AllowCustomer
type Broker struct {
	Name            string   `json:"name"`
	BrokerNo        int      `json:"brokerno"`
	AllowCustomer   []string `json:"allowcustomer"`
	PendingCustomer []string `json:"pendingcustomer"`
}

// Main
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init - reset all the things
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

	err = stub.PutState(brokerIndexStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Run - Our entry point for Invocations - [LEGACY] obc-peer 4/25/2016
// ============================================================================================================================
func (t *SimpleChaincode) Run(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("run is running " + function)
	return t.Invoke(stub, function, args)
}

// Invoke - Our entry point for Invocations
// ============================================================================================================================
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" { //writes a value to the chaincode state
		return t.Write(stub, args)
	} else if function == "newcustomer" { //create a new customer
		return t.newcustomer(stub, args)
	} else if function == "newbroke" {
		return t.newbroke(stub, args)
	} else if function == "requestPermission" {
		return t.requestPermission(stub, args)
	} else if function == "customerallow" {
		return t.customerallow(stub, args)
	}
	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query - Our entry point for Queries
// ============================================================================================================================
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	if function == "readcustomer" { //read a variable
		return t.readcustomer(stub, args)
	}
	if function == "readcustomergid" {
		return t.readcustomergid(stub, args)
	}
	if function == "readbroker" { //read a variable
		return t.readbroker(stub, args)
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
// Read - read a variable from chaincode state by cardid
// ============================================================================================================================
func (t *SimpleChaincode) readCustomerByCus(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var cardid, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	cardid = args[0]
	valAsbytes, err := stub.GetState(customerKey + cardid) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cardid + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Read - read a variable from chaincode state by cardid
// ============================================================================================================================
func (t *SimpleChaincode) readcustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var cardid, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	cardid = args[0]
	valAsbytes, err := stub.GetState(customerKey + cardid) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + cardid + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Read - read a variable from chaincode state by gauranteeid
// ============================================================================================================================
func (t *SimpleChaincode) readcustomergid(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var gid, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	gid = args[0]
	valAsbytes, err := stub.GetState(GauranteeIDKey + gid) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + gid + "\"}"
		return nil, errors.New(jsonResp)
	}

	gau := GauranteeID{}
	json.Unmarshal(valAsbytes, &gau)

	return valAsbytes, nil //send it onward
}

// ============================================================================================================================
// Read - read a variable from chaincode state by brokeno
// ============================================================================================================================
func (t *SimpleChaincode) readbroker(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var brokeno, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	brokeno = args[0]
	valAsbytes, err := stub.GetState(BrokerKey + brokeno) //get the var from chaincode state
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + brokeno + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil //send it onward
}

func (t *SimpleChaincode) requestPermission(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	// var brokerno, gid string
	// var err error

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting gaurantee id and brokeno")
	}

	gid := args[0]
	//brokeno := args[1]
	brokeNoAsString := args[1]
	brokeNo, err := strconv.Atoi(brokeNoAsString)
	brokerAsBytes, err := stub.GetState(BrokerKey + brokeNoAsString)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for broker " + brokeNoAsString + "\"}"
		return nil, errors.New(jsonResp)
	}
	broker := Broker{}
	json.Unmarshal(brokerAsBytes, &broker)

	gidAsbytes, err := stub.GetState(GauranteeIDKey + gid)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + gid + "\"}"
		return nil, errors.New(jsonResp)
	}
	gau := GauranteeID{}
	json.Unmarshal(gidAsbytes, &gau)

	already := false
	fmt.Printf("GID %s\n", gidAsbytes)
	fmt.Printf("brokeno %d\n", brokeNo)
	for _, s := range gau.AllowBroke {
		//set[s] = struct{}{}
		if s == brokeNo {
			already = true
		}
	}

	fmt.Println("already" + strconv.FormatBool(already))

	if already {
		jsonResp := "{\"Error\":\"Already Allowed " + brokeNoAsString + "\"}"
		return nil, errors.New(jsonResp)
	}

	gau.PendingBroke = append(gau.PendingBroke, brokeNo)
	broker.PendingCustomer = append(broker.PendingCustomer, gau.GauranteeID)

	fmt.Println("gau.AllowBroke")
	fmt.Println(gau.AllowBroke)

	jsonAsBytes, err := json.Marshal(gau)
	fmt.Printf("GID %s\n", jsonAsBytes)
	err = stub.PutState(GauranteeIDKey+gid, jsonAsBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	jsonBrokerAsBytes, err := json.Marshal(broker)
	fmt.Printf("Broke %s\n", jsonBrokerAsBytes)
	err = stub.PutState(BrokerKey+brokeNoAsString, jsonBrokerAsBytes) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// ============================================================================================================================
// Set User Permission on Customer
// ============================================================================================================================
func (t *SimpleChaincode) customerallow(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error
	//gid, brokeno

	if len(args) < 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2")
	}

	gid := args[0]
	brokeNoAsString := args[1]
	brokeNo, err := strconv.Atoi(brokeNoAsString)

	gidAsBytes, err := stub.GetState(GauranteeIDKey + gid)
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	gau := GauranteeID{}
	json.Unmarshal(gidAsBytes, &gau)

	brokeAsBytes, err := stub.GetState(BrokerKey + brokeNoAsString)
	if err != nil {
		return nil, errors.New("Failed to get thing")
	}
	broke := Broker{}
	json.Unmarshal(brokeAsBytes, &broke)

	for i := len(gau.PendingBroke) - 1; i >= 0; i-- {
		if gau.PendingBroke[i] == brokeNo {
			gau.PendingBroke = append(gau.PendingBroke[:i], gau.PendingBroke[i+1:]...)
			break
		}
	}
	for i := len(broke.PendingCustomer) - 1; i >= 0; i-- {
		if broke.PendingCustomer[i] == gid {
			broke.PendingCustomer = append(broke.PendingCustomer[:i], broke.PendingCustomer[i+1:]...)
			break
		}
	}
	already := false
	for _, s := range gau.AllowBroke {
		if s == brokeNo {
			already = true
		}
	}
	if !already {
		gau.AllowBroke = append(gau.AllowBroke, brokeNo)
		broke.AllowCustomer = append(broke.AllowCustomer, gid)
	}

	//write state

	jsonAsBytes, _ := json.Marshal(broke)
	err = stub.PutState(BrokerKey+brokeNoAsString, jsonAsBytes) //rewrite the customer with id as key
	if err != nil {
		return nil, err
	}

	jsonGauAsBytes, _ := json.Marshal(gau)
	err = stub.PutState(GauranteeIDKey+gid, jsonGauAsBytes) //rewrite the customer with id as key
	if err != nil {
		return nil, err
	}

	fmt.Println("- end set allow permission")
	return nil, nil
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
// Init Customer - create a new customer, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) newcustomer(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0        1       2        3           4          5
	// "name", "telno", "age", "occupation", "cardid", "creator"
	if len(args) != 6 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init customer")
	if len(args[0]) <= 0 {
		fmt.Println("- not pass0 init customer")
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		fmt.Println("- not pass1 init customer")
		return nil, errors.New("2nd argument must be a non-empty string")
	}
	if len(args[2]) <= 0 {
		fmt.Println("- not pass2 init customer")
		return nil, errors.New("3rd argument must be a non-empty string")
	}
	if len(args[3]) <= 0 {
		fmt.Println("- not pass3 init customer")
		return nil, errors.New("4th argument must be a non-empty string")
	}
	if len(args[4]) <= 0 {
		fmt.Println("- not pass4 init customer")
		return nil, errors.New("5th argument must be a non-empty string")
	}
	if len(args[5]) <= 0 {
		fmt.Println("- not pass5 init customer")
		return nil, errors.New("6th argument must be a non-empty string")
	}

	fmt.Println("- pass1 init customer")
	name := args[0]
	telno := strings.ToLower(args[1])
	age, err := strconv.Atoi(args[2])
	if err != nil {
		return nil, errors.New("3rd argument must be a numeric string")
	}
	occupation := strings.ToLower(args[3])
	cardid := args[4]
	creator := args[5]

	fmt.Println("- pass2 init customer")
	//check if customer already exists
	customerAsBytes, err := stub.GetState(customerKey + name)
	if err != nil {
		return nil, errors.New("Failed to get customer name")
	}
	res := Customer{}
	json.Unmarshal(customerAsBytes, &res)
	if res.Name == name {
		fmt.Println("This customer arleady exists: " + name)
		fmt.Println(res)
		return nil, errors.New("This customer arleady exists") //all stop a customer by this name exists
	}

	res.Name = name
	res.TelNo = telno
	res.Age = age
	res.Occupation = occupation
	res.CardID = cardid
	res.Creator = creator
	//build the customer json string manually
	//str := `{"name": "` + name + `", "telno": "` + telno + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	//err = stub.PutState(name, []byte(str)) //store customer with id as key
	str, err := json.Marshal(res)
	err = stub.PutState(customerKey+name, str)
	if err != nil {
		return nil, err
	}
	err = stub.PutState(customerKey+cardid, str)
	if err != nil {
		return nil, err
	}

	sha256AsByte := sha256.Sum256(str)

	gauranteeID := GauranteeID{}
	var emptyIntArray []int
	gauranteeID.GauranteeID = strings.ToUpper(hex.EncodeToString(sha256AsByte[:]))
	gauranteeID.CustomerID = res.CardID
	gauranteeID.AllowBroke = emptyIntArray
	gauranteeID.PendingBroke = emptyIntArray

	str, err = json.Marshal(gauranteeID)
	err = stub.PutState(GauranteeIDKey+gauranteeID.GauranteeID, str)
	if err != nil {
		return nil, err
	}

	//gid := GauranteeID{}
	// h := sha1.New()
	// h.Write([]byte(res))
	// gid.GauranteeID = h.Sum()

	//get the customer index-
	customersAsBytes, err := stub.GetState(customerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get customer index")
	}
	var customerIndex []string
	json.Unmarshal(customersAsBytes, &customerIndex) //un stringify it aka JSON.parse()

	//append
	customerIndex = append(customerIndex, cardid) //add customer name to index list
	fmt.Println("! customer index: ", customerIndex)
	jsonAsBytes, _ := json.Marshal(customerIndex)
	err = stub.PutState(customerIndexStr, jsonAsBytes) //store name of customer

	fmt.Println("- end init customer")
	return nil, nil
}

// ============================================================================================================================
// Init Broke - create a new broke, store into chaincode state
// ============================================================================================================================
func (t *SimpleChaincode) newbroke(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var err error

	//   0        1
	// "name", "brokeID"
	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4")
	}

	//input sanitation
	fmt.Println("- start init customer")
	if len(args[0]) <= 0 {
		return nil, errors.New("1st argument must be a non-empty string")
	}
	if len(args[1]) <= 0 {
		return nil, errors.New("2nd argument must be a non-empty string")
	}

	name := args[0]
	brokeNoAsString := args[1]
	brokeNo, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, errors.New("2nd argument must be a numeric string")
	}

	//check if broker already exists
	brokerAsBytes, err := stub.GetState(brokeNoAsString)
	if err != nil {
		return nil, errors.New("Failed to get broker name")
	}
	res := Broker{}
	json.Unmarshal(brokerAsBytes, &res)
	if res.BrokerNo == brokeNo {
		fmt.Println("This broker arleady exists: " + name)
		fmt.Println(res)
		return nil, errors.New("This broker arleady exists") //all stop a broker by this name exists
	}

	res.Name = name
	res.BrokerNo = brokeNo
	//build the broker json string manually
	//str := `{"name": "` + name + `", "telno": "` + telno + `", "size": ` + strconv.Itoa(size) + `, "user": "` + user + `"}`
	//err = stub.PutState(name, []byte(str)) //store broker with id as key
	str, err := json.Marshal(res)
	err = stub.PutState(BrokerKey+brokeNoAsString, str)
	if err != nil {
		return nil, err
	}

	//get the broker index-
	brokersAsBytes, err := stub.GetState(brokerIndexStr)
	if err != nil {
		return nil, errors.New("Failed to get broker index")
	}
	var brokerIndex []string
	json.Unmarshal(brokersAsBytes, &brokerIndex) //un stringify it aka JSON.parse()

	//append
	brokerIndex = append(brokerIndex, brokeNoAsString) //add broker name to index list
	fmt.Println("! broker index: ", brokerIndex)
	jsonAsBytes, _ := json.Marshal(brokerIndex)
	err = stub.PutState(brokerIndexStr, jsonAsBytes) //store name of broker

	fmt.Println("- end init broker")
	return nil, nil
}
