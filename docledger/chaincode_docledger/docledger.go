/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Writing Your First Blockchain Application
 */

package main

/* Imports
 * 4 utility libraries for formatting, handling bytes, reading and writing JSON, and string manipulation
 * 2 specific Hyperledger Fabric specific libraries for Smart Contracts
 */
import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the StudentDoc structure, with 3 properties.  Structure tags are used by encoding/json library
type StudentDoc struct {
	DocStatus  string `json:"docStatus"`
	Owner string `json:"owner"`
}

/*
 * The Init method is called when the Smart Contract "docledger" is instantiated by the blockchain network
 * Best practice is to have any Ledger initialization in separate function -- see initLedger()
 */
func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

/*
 * The Invoke method is called as a result of an application request to run the Smart Contract "docledger"
 * The calling application program has also specified the particular smart contract function to be called, with arguments
 */
func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "queryStudentDoc" {
		return s.queryStudentDoc(APIstub, args)
	} else if function == "initLedger" {
		return s.initLedger(APIstub)
	} else if function == "createStudentDoc" {
		return s.createStudentDoc(APIstub, args)
	} else if function == "queryAllStudentDocs" {
		return s.queryAllStudentDocs(APIstub)
	} else if function == "changeStudentDocOwner" {
		return s.changeStudentDocOwner(APIstub, args)
	} else if function == "changeStudentDocStatus" {
		return s.changeStudentDocStatus(APIstub, args)
	} else if function == "getHistoryForStudentDoc" {
		return s.getHistoryForStudentDoc(APIstub, args)
	}

	return shim.Error(function + ": Invalid Smart Contract function name.")
}

func (s *SmartContract) queryStudentDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	studentDocAsBytes, _ := APIstub.GetState(args[0])
	return shim.Success(studentDocAsBytes)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	studentDocs := []StudentDoc{
		StudentDoc{DocStatus: "scanned", Owner: "Tomoko"},
		StudentDoc{DocStatus: "transmitted responses", Owner: "Jack"},
		StudentDoc{DocStatus: "received responses", Owner: "John"},
		StudentDoc{DocStatus: "machine scored", Owner: "Mark"},
		StudentDoc{DocStatus: "human scored", Owner: "Tim"},
		StudentDoc{DocStatus: "scores exported", Owner: "Jane"},
		StudentDoc{DocStatus: "transmitted scores", Owner: "Peter"},
		StudentDoc{DocStatus: "received scores", Owner: "Sid"},
		StudentDoc{DocStatus: "scores reported", Owner: "Mesut"},
	}

	i := 0
	for i < len(studentDocs) {
		fmt.Println("i is ", i)
		studentDocAsBytes, _ := json.Marshal(studentDocs[i])
		APIstub.PutState("StudentDoc"+strconv.Itoa(i), studentDocAsBytes)
		fmt.Println("Added", studentDocs[i])
		i = i + 1
	}

	return shim.Success(nil)
}

func (s *SmartContract) createStudentDoc(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	var studentDoc = StudentDoc{DocStatus: args[1], Owner: args[2]}

	studentDocAsBytes, _ := json.Marshal(studentDoc)
	APIstub.PutState(args[0], studentDocAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) queryAllStudentDocs(APIstub shim.ChaincodeStubInterface) sc.Response {

	startKey := "StudentDoc0"
	endKey := "StudentDoc999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllStudentDocs:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) changeStudentDocOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2")
	}

	studentDocAsBytes, _ := APIstub.GetState(args[0])
	studentDoc := StudentDoc{}

	json.Unmarshal(studentDocAsBytes, &studentDoc)
	studentDoc.Owner = args[1]

	studentDocAsBytes, _ = json.Marshal(studentDoc)
	APIstub.PutState(args[0], studentDocAsBytes)

	return shim.Success(nil)
}

func (s *SmartContract) changeStudentDocStatus(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) != 3 {
		return shim.Error("Incorrect number of arguments. Expecting 3")
	}

	studentDocAsBytes, _ := APIstub.GetState(args[0])
	studentDoc := StudentDoc{}

	json.Unmarshal(studentDocAsBytes, &studentDoc)
	studentDoc.Owner = args[1]
	studentDoc.DocStatus = args[2]

	studentDocAsBytes, _ = json.Marshal(studentDoc)
	APIstub.PutState(args[0], studentDocAsBytes)

	return shim.Success(nil)
}

func (t *SmartContract) getHistoryForStudentDoc(stub shim.ChaincodeStubInterface, args []string) sc.Response {

	if len(args) < 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docId := args[0]

	fmt.Printf("- start getHistoryForStudentDoc: %s\n", docId)

	resultsIterator, err := stub.GetHistoryForKey(docId)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getHistoryForStudentDoc returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
