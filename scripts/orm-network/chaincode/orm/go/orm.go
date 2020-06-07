/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-chaincode-go/shim/entitydefinition"
	pb "github.com/hyperledger/fabric-protos-go/peer"
	"strconv"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type User struct {
	ID       string
	Name     string
	Email    string    `gorm:"type:varchar(100);unique_index"`
	Accounts []Account `ormdb:"entity"`
}

type Account struct {
	ID     string
	Number string
	Amount sql.NullFloat64 `ormdb:"datatype"`
	UserId string
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
		return t.invoke(stub, args)
	} else if function == "delete" {
		return t.delete(stub, args)
	} else if function == "query" {
		return t.query(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

func (t *SimpleChaincode) invoke(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	user := &User{}
	user.ID = args[3]
	user.Name = args[0]
	user.Email = args[1]

	accounts := make([]Account, 2)
	for i, _ := range accounts {
		account := &Account{}
		account.ID = args[4] + strconv.Itoa(i)
		account.Number = args[2] + strconv.Itoa(i)
		account.UserId = user.ID
		account.Amount = sql.NullFloat64{Float64: 100.00, Valid: true}
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

func (t *SimpleChaincode) delete(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]

	accounts := make([]Account, 0)
	search := &entitydefinition.Search{}
	err := search.Where("user_id = ?", userId)
	search.Offset(0).Limit(1)

	if err != nil {
		return shim.Error("create search condition failed")
	}
	err = stub.ConditionQuery(accounts, search)
	if err != nil {
		return shim.Error("condition query failed")
	}
	accountsByte, err := json.Marshal(accounts)
	if err != nil {
		return shim.Error("marshal accounts failed")
	}

	for _, account := range accounts {
		err = stub.Delete(account)
		if err != nil {
			return shim.Error("delete account failed")
		}
	}

	return shim.Success(accountsByte)
}

func (t *SimpleChaincode) query(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	userId := args[0]

	user := &User{ID: userId}
	err := stub.Get(user)
	if err != nil {
		return shim.Error("query user failed")
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		return shim.Error("marshal user failed")
	}
	return shim.Success(userBytes)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
