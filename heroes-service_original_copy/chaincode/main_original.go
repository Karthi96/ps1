package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type student struct {
	RollNo  string `json:"rollno"` //rollno is used to distinguish the various types of objects in  database i.e primary key
	Name    string `json:"name"`
	Course  string `json:"course"`
	Grade   string `json:"grade"`
	Details string `json:"details"`
	Date    string `json:"date"`
}

// ===================================================================================
// Main Function
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

	// Handle different functions
	if function == "createStudent" { //create a new student record
		return t.createStudent(stub, args)
	} else if function == "readStudent" { //read a particular student record
		return t.readStudent(stub, args)
	} else if function == "getHistoryForStudent" { //get history of values for a marble
		return t.getHistoryForStudent(stub, args)
	} else if function == "updateStudent" { // Update a student Based on Rollno
		return t.updateStudent(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error
	return shim.Error("Received unknown function invocation")
}

// ============================================================
// createStudent - create a new student, store into chaincode state
// ============================================================
func (t *SimpleChaincode) createStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6")
	}

	// ==== Input sanitation ====
	for i := 0; i < 6; i++ {
		if len(args[i]) <= 0 {
			return shim.Error("argument must be a non-empty string")
		}
	}

	rollno := strings.ToLower(args[0])
	name := args[1]
	course := args[2]
	//date, err := time.Parse(time.RFC3339, args[4])
	//date1 := date.String()
	date1 := args[4]
	details := args[3]
	grade := args[5]
	if err != nil {
		return shim.Error("4rd argument must be a date string" + err.Error())
	}

	// ==== Check if rollno already exists ====
	studentAsBytes, err := stub.GetState(rollno + "_" + course)
	if err != nil {
		return shim.Error("Failed to get student: " + err.Error())
	} else if studentAsBytes != nil {
		fmt.Println("This rollno already exists with this course " + rollno + course)
		return shim.Error("This rollno already exists: with this course" + rollno)
	}

	// ==== Create student object and marshal to JSON ====
	student := &student{rollno, name, course, grade, details, date1}
	studentJSONasBytes, err := json.Marshal(student)
	if err != nil {
		return shim.Error(err.Error())
	}

	// === Save marble to state ===
	err = stub.PutState(rollno+"_"+course, studentJSONasBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err = stub.SetEvent("eventputstudents", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	//  ==== Index the marble to enable color-based range queries, e.g. return all blue marbles ====
	//  An 'index' is a normal key/value entry in state.
	//  The key is a composite key, with the elements that you want to range query on listed first.
	//  In our case, the composite key is based on indexName~color~name.
	//  This will enable very efficient state range queries based on composite keys matching indexName~color~*
	//	indexName := "color~name"
	//	colorNameIndexKey, err := stub.CreateCompositeKey(indexName, []string{marble.Color, marble.Name})
	//	if err != nil {
	//		return shim.Error(err.Error())
	//	}
	//  Save index entry to state. Only the key name is needed, no need to store a duplicate copy of the marble.
	//  Note - passing a 'nil' value will effectively delete the key from state, therefore we pass null character as value
	//	value := []byte{0x00}
	//	stub.PutState(colorNameIndexKey, value)

	// ==== Marble saved and indexed. Return success ====
	fmt.Println("- end createstudent")
	return shim.Success(nil)

}

// ============================================================
// updateStudent - Update a student record based on rollno
// ============================================================
func (t *SimpleChaincode) updateStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error

	if len(args) != 6 {
		return shim.Error("Incorrect number of arguments. Expecting 6 ")
	}

	// ==== Input sanitation ====
	for i := 0; i < 6; i++ {
		if len(args[i]) <= 0 {
			return shim.Error("argument must be a non-empty string")
		}
	}

	rollno := strings.ToLower(args[0])
	name := args[1]
	course := args[2]
	date1 := args[4]
	details := args[3]
	grade := args[5]
	if err != nil {
		return shim.Error("4rd argument must be a date string" + err.Error())
	}

	// ==== Check if rollno already exists to update====
	studentAsBytes, err := stub.GetState(rollno + "_" + course)
	if err != nil {
		return shim.Error("Failed to get student: " + err.Error())
	} else if studentAsBytes == nil {
		return shim.Error("Unable to find the rollno " + rollno + "for this course=" + course)
	}

	recordToUpdate := student{}
	err = json.Unmarshal(studentAsBytes, &recordToUpdate) //unmarshal it aka JSON.parse()
	if err != nil {
		return shim.Error(err.Error())
	}

	recordToUpdate.Name = name
	recordToUpdate.Course = course
	recordToUpdate.Date = date1
	recordToUpdate.Details = details
	recordToUpdate.Grade = grade

	studentUpdateAsBytes, _ := json.Marshal(recordToUpdate)
	err = stub.PutState(rollno+"_"+course, studentUpdateAsBytes) //rewrite the student
	if err != nil {
		return shim.Error(err.Error())
	}

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err = stub.SetEvent("eventupdatestudent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	fmt.Println("- end updatestudent")
	return shim.Success(nil)

}

//Read particular student value
func (t *SimpleChaincode) readStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) < 2 {
		return shim.Error("readStudent operation must include two argument, a rollno and course")
	}
	Rollno := args[0]
	Course := args[1]
	var errorCount int = 0
	var keyStr []string
	var keysOutput []string
	if Course == "All" {

		keyStr = append(keyStr, Rollno+"_Blockchain")
		keyStr = append(keyStr, Rollno+"_Machine Learning")
		keyStr = append(keyStr, Rollno+"_Data Science")
		keyStr = append(keyStr, Rollno+"_Artificial Intelligence")

	} else {
		keyStr = append(keyStr, Rollno+"_"+Course)
	}

	for index, eachKey := range keyStr {
		fmt.Printf("Index=%d", index)
		value, err := stub.GetState(eachKey)
		if err != nil {
			return shim.Error(fmt.Sprintf("get operation failed. Error accessing state: %s", err))
		}
		if value == nil {
			errorCount = errorCount + 1
			continue
		}
		keysOutput = append(keysOutput, string(value))
	}
	if len(keyStr) == errorCount {
		return shim.Error("Unable to find the ROllno= " + Rollno + "with this Course=" + Course)
	}

	jsonVal, err := json.Marshal(strings.Join(keysOutput, ","))

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err = stub.SetEvent("eventgetstudents", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsonVal)
}

//Get the history of transaction on particular key
func (t *SimpleChaincode) getHistoryForStudent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	fmt.Println("Get history")
	Rollno := args[0]
	Course := args[1]
	var errorCount int = 0
	var keyStr []string
	var keysOutput []string
	if Course == "All" {
		keyStr = append(keyStr, Rollno+"_Blockchain")
		keyStr = append(keyStr, Rollno+"_Machine Learning")
		keyStr = append(keyStr, Rollno+"_Data Science")
		keyStr = append(keyStr, Rollno+"_Artificial Intelligence")

	} else {
		keyStr = append(keyStr, Rollno+"_"+Course)
	}
	for index, eachKey := range keyStr {
		fmt.Println("Index=%d", index)
		value, err := stub.GetState(eachKey)
		// if err != nil {
		// 	return shim.Error(fmt.Sprintf("get operation failed. Error accessing state: %s", err))
		// }

		if value == nil {
			errorCount = errorCount + 1
			continue
		}
		keysIter, err := stub.GetHistoryForKey(eachKey)

		defer keysIter.Close()
		for keysIter.HasNext() {
			response, iterErr := keysIter.Next()
			if iterErr != nil {
				return shim.Error(fmt.Sprintf("getHistoryForStudent history operation failed. Error accessing state: %s", err))
			}
			timestamp := response.GetTimestamp()
			time := time.Unix(timestamp.Seconds, 0)
			var m map[string]interface{}
			err := json.Unmarshal(response.GetValue(), &m)
			if err != nil {
				return shim.Error(err.Error())
			}
			m["timestamp"] = time.String()
			m["txtid"] = response.GetTxId()

			newData, err1 := json.Marshal(m)
			if err1 != nil {
				return shim.Error(err1.Error())
			}

			keysOutput = append(keysOutput, string(newData))
		}
	}
	if len(keyStr) == errorCount {
		return shim.Error("Unable to find the ROllno= " + Rollno + "with this Course=" + Course)
	}

	jsonKeys, err := json.Marshal(strings.Join(keysOutput, ","))
	if err != nil {
		return shim.Error(fmt.Sprintf("getHistoryForStudent history operation failed. Error marshaling JSON: %s", err))
	}

	// Notify listeners that an event "eventInvoke" have been executed (check line 19 in the file invoke.go)
	err = stub.SetEvent("eventhistorystudent", []byte{})
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(jsonKeys)
}
