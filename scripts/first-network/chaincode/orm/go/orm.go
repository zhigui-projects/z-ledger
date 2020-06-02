/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"database/sql"
	"fmt"
	"github.com/hyperledger/fabric/common/util"
	"math/rand"
	"strings"
	"time"

	"github.com/hyperledger/fabric-chaincode-go/shim"
	pb "github.com/hyperledger/fabric-protos-go/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	ID        string
	Name      string
	Email     string    `gorm:"type:varchar(100);unique_index"`
	CreatedAt time.Time `ormdb:"datatype"`
	UpdatedAt time.Time `ormdb:"datatype"`
	Accounts  []Account `ormdb:"entity"`
}

type Account struct {
	ID        string
	Number    string
	Amount    sql.NullFloat64 `ormdb:"datatype"`
	UserId    string
	CreatedAt time.Time `ormdb:"datatype"`
	UpdatedAt time.Time `ormdb:"datatype"`
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Init")
	err := stub.CreateTable(&Account{}, 1)
	if err != nil {
		return shim.Error("Error create account table")
	}
	err = stub.CreateTable(&User{}, 2)
	if err != nil {
		return shim.Error("Error create user table")
	}
	return shim.Success(nil)
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "invoke" {
		// Make payment of X units from A to B
		return t.invoke(stub, args)
	} else if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "query" {
		// the old "Query" is now implemtned in invoke
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Save user and two account
func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	user := &User{}
	user.ID = strings.ReplaceAll(util.GenerateUUID(), "-", "")
	user.Name = args[0]
	user.Email = args[1]
	now := time.Now().UTC()
	user.CreatedAt = now
	user.UpdatedAt = now

	accounts := make([]Account, 2)
	for i, _ := range accounts {
		account := &Account{}
		account.ID = strings.ReplaceAll(util.GenerateUUID(), "-", "")
		ra, _ := randStr()
		account.Number = args[2] + ra
		account.UserId = user.ID
		account.Amount = sql.NullFloat64{Float64: 100.00}
		accounts[i] = *account
		err := stub.Save(account)
		if err != nil {
			fmt.Println(err)
			return shim.Error("Failed to save account")
		}
	}

	user.Accounts = accounts

	err := stub.Save(user)
	if err != nil {
		fmt.Println(err)
		return shim.Error("Failed to save user")
	}

	return shim.Success([]byte(user.ID))
}

// Deletes an entity from state
func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]

	stub.Delete()
	// Delete the key from the state in ledger
	err := stub.Delete(userId)
	if err != nil {
		return shim.Error("Failed to delete user")
	}

	return shim.Success(nil)
}

// query callback representing the query of a chaincode
func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var A string // Entities
	var err error

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting name of the person to query")
	}

	A = args[0]

	// Get the state from the ledger
	Avalbytes, err := stub.GetState(A)
	if err != nil {
		jsonResp := "{\"Error\":\"Failed to get state for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	if Avalbytes == nil {
		jsonResp := "{\"Error\":\"Nil amount for " + A + "\"}"
		return shim.Error(jsonResp)
	}

	jsonResp := "{\"Name\":\"" + A + "\",\"Amount\":\"" + string(Avalbytes) + "\"}"
	fmt.Printf("Query Response:%s\n", jsonResp)
	return shim.Success(Avalbytes)
}

func randStr() (string, error) {
	c := 10
	b := make([]byte, c)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("error:", err)
		return "", err
	}
	return string(b), nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
