package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type TokenContract struct {
	contractapi.Contract
}

type Token struct {
	Symbol      string  `json:"symbol"`
	Name        string  `json:"name"`
	TotalSupply float64 `json:"totalSupply"`
	Admin       string  `json:"admin"`
}

type User struct {
	Name       string             `json:"name"`
	Balance    float64            `json:"balance"`
	Allowances map[string]float64 `json:"allownaces"`
}

const tokenKey = "TOKEN"

func (s *TokenContract) Createuser(ctx contractapi.TransactionContextInterface, name string) error {
	existingUser, err := ctx.GetStub().GetState(name)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingUser != nil {
		return errors.New("user already exists")
	}
	user := User{
		Name:       name,
		Balance:    0,
		Allowances: make(map[string]float64),
	}
	tokenJSON, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(name, tokenJSON)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	return nil
}

func (s *TokenContract) InitLedger(ctx contractapi.TransactionContextInterface, symbol string, name string, initialSupply float64, admin string) error {
	existingToken, err := ctx.GetStub().GetState(tokenKey)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existingToken != nil {
		return errors.New("token already initialized")
	}
	token := Token{
		Symbol:      symbol,
		Name:        name,
		TotalSupply: initialSupply,
		Admin:       admin,
	}
	tokenJSON, err := json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token: %v", err)
	}

	err = ctx.GetStub().PutState(tokenKey, tokenJSON)
	if err != nil {
		return fmt.Errorf("failed to store token: %v", err)
	}

	adminBalance := User{Name: admin, Balance: initialSupply, Allowances: make(map[string]float64)}
	adminBalanceJSON, err := json.Marshal(adminBalance)
	if err != nil {
		return fmt.Errorf("failed to marshal admin balance: %v", err)
	}

	err = ctx.GetStub().PutState(admin, adminBalanceJSON)
	if err != nil {
		return fmt.Errorf("failed to store admin balance: %v", err)
	}

	return nil
}

func (s *TokenContract) MintTokens(ctx contractapi.TransactionContextInterface, admin string, amount float64) error {
	tokenData, err := ctx.GetStub().GetState(tokenKey)
	if err != nil || tokenData == nil {
		return errors.New("token not found")
	}

	var token Token
	err = json.Unmarshal(tokenData, &token)
	if err != nil {
		return fmt.Errorf("failed to unmarshal token data: %v", err)
	}

	if token.Admin != admin {
		return errors.New("only admin can mint tokens")
	}

	token.TotalSupply += amount
	tokenData, err = json.Marshal(token)
	if err != nil {
		return fmt.Errorf("failed to marshal token data: %v", err)
	}

	err = ctx.GetStub().PutState(tokenKey, tokenData)
	if err != nil {
		return fmt.Errorf("failed to update token supply: %v", err)
	}

	adminData, err := ctx.GetStub().GetState(admin)
	var adminAccount User
	if err == nil && adminData != nil {
		err = json.Unmarshal(adminData, &adminAccount)
		if err != nil {
			return fmt.Errorf("failed to unmarshal admin account data: %v", err)
		}
	} else {
		adminAccount = User{Balance: 0, Allowances: make(map[string]float64)}
	}

	adminAccount.Balance += amount

	adminData, err = json.Marshal(adminAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal admin account data: %v", err)
	}

	err = ctx.GetStub().PutState(admin, adminData)
	if err != nil {
		return fmt.Errorf("failed to update admin balance: %v", err)
	}

	return nil
}

func (s *TokenContract) ApproveSpender(ctx contractapi.TransactionContextInterface, owner string, spender string, amount float64) error {
	ownerData, err := ctx.GetStub().GetState(owner)
	if err != nil || ownerData == nil {
		return errors.New("owner not found")
	}

	var ownerAccount User
	err = json.Unmarshal(ownerData, &ownerAccount)
	if err != nil {
		return fmt.Errorf("failed to unmarshal owner data: %v", err)
	}
	if ownerAccount.Balance < amount {
		return errors.New("insufficient balance with owner")
	}

	spenderData, err := ctx.GetStub().GetState(spender)
	if err != nil || ownerData == nil {
		return errors.New("owner not found")
	}

	var spenderInfo User
	err = json.Unmarshal(spenderData, &spenderInfo)
	if err != nil {
		return fmt.Errorf("failed to unmarshal owner data: %v", err)
	}
	if ownerAccount.Balance < amount {
		return errors.New("insufficient balance with owner")
	}
	ownerAccount.Allowances[spender] += amount

	spenderInfo.Balance += amount
	ownerData, err = json.Marshal(ownerAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal owner data: %v", err)
	}

	err = ctx.GetStub().PutState(owner, ownerData)
	if err != nil {
		return fmt.Errorf("failed to update owner allowances: %v", err)
	}

	spenderInfoBytes, err := json.Marshal(spenderInfo)
	if err != nil {
		return fmt.Errorf("failed to marshal owner data: %v", err)
	}

	err = ctx.GetStub().PutState(spender, spenderInfoBytes)
	if err != nil {
		return fmt.Errorf("failed to update owner allowances: %v", err)
	}

	return nil
}

func (s *TokenContract) TransferTokens(ctx contractapi.TransactionContextInterface, from string, to string, amount float64) error {
	fromData, err := ctx.GetStub().GetState(from)
	if err != nil || fromData == nil {
		return errors.New("sender not found")
	}

	var sender User
	err = json.Unmarshal(fromData, &sender)
	if err != nil {
		return fmt.Errorf("failed to unmarshal sender data: %v", err)
	}

	if sender.Balance < amount {
		return errors.New("insufficient balance")
	}

	receiver := User{}
	toData, err := ctx.GetStub().GetState(to)
	if err == nil && toData != nil {
		err = json.Unmarshal(toData, &receiver)
		if err != nil {
			return fmt.Errorf("failed to unmarshal receiver data: %v", err)
		}
	}

	sender.Balance -= amount
	receiver.Balance += amount

	fromData, err = json.Marshal(sender)
	if err != nil {
		return fmt.Errorf("failed to marshal sender data: %v", err)
	}

	toData, err = json.Marshal(receiver)
	if err != nil {
		return fmt.Errorf("failed to marshal receiver data: %v", err)
	}

	err = ctx.GetStub().PutState(from, fromData)
	if err != nil {
		return fmt.Errorf("failed to update sender balance: %v", err)
	}

	err = ctx.GetStub().PutState(to, toData)
	if err != nil {
		return fmt.Errorf("failed to update receiver balance: %v", err)
	}

	return nil
}

func (s *TokenContract) GetBalance(ctx contractapi.TransactionContextInterface, user string) (float64, error) {
	userData, err := ctx.GetStub().GetState(user)
	if err != nil || userData == nil {
		return 0, errors.New("user not found")
	}

	var userAccount User
	err = json.Unmarshal(userData, &userAccount)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	return userAccount.Balance, nil
}

func (s *TokenContract) GetUser(ctx contractapi.TransactionContextInterface, user string) (User, error) {
	var userAccount User
	userData, err := ctx.GetStub().GetState(user)
	if err != nil || userData == nil {
		return userAccount, errors.New("user not found")
	}

	err = json.Unmarshal(userData, &userAccount)
	if err != nil {
		return userAccount, fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	return userAccount, nil
}

func (s *TokenContract) BurnTokens(ctx contractapi.TransactionContextInterface, user string, amount float64) error {
	userData, err := ctx.GetStub().GetState(user)
	if err != nil || userData == nil {
		return errors.New("user not found")
	}

	var userAccount User
	err = json.Unmarshal(userData, &userAccount)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user data: %v", err)
	}

	if userAccount.Balance < amount {
		return errors.New("insufficient balance to burn")
	}

	userAccount.Balance -= amount

	userData, err = json.Marshal(userAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal user data: %v", err)
	}

	err = ctx.GetStub().PutState(user, userData)
	if err != nil {
		return fmt.Errorf("failed to update user balance: %v", err)
	}

	return nil
}

func (s *TokenContract) TransferFromApprovedSpenders(ctx contractapi.TransactionContextInterface, owner, spender string, recipient string, amount float64) error {
	ownerData, err := ctx.GetStub().GetState(owner)
	if err != nil || ownerData == nil {
		return errors.New("owner not found")
	}

	var ownerAccount User
	err = json.Unmarshal(ownerData, &ownerAccount)
	if err != nil {
		return fmt.Errorf("failed to unmarshal owner data: %v", err)
	}

	spenderData, err := ctx.GetStub().GetState(spender)
	if err != nil || ownerData == nil {
		return errors.New("owner not found")
	}

	var spenderUser User
	err = json.Unmarshal(spenderData, &spenderUser)
	if err != nil {
		return fmt.Errorf("failed to unmarshal owner data: %v", err)
	}

	allowance, ok := ownerAccount.Allowances[spender]
	if !ok || allowance < amount {
		return errors.New("spender allowance exceeded or not approved")
	}

	if ownerAccount.Balance < amount {
		return errors.New("insufficient balance")
	}

	ownerAccount.Allowances[spender] -= amount
	ownerAccount.Allowances[recipient] += amount

	recipientData, err := ctx.GetStub().GetState(recipient)
	var recipientAccount User
	if err == nil && recipientData != nil {
		err = json.Unmarshal(recipientData, &recipientAccount)
		if err != nil {
			return fmt.Errorf("failed to unmarshal recipient data: %v", err)
		}
	}
	spenderUser.Balance -= amount
	recipientAccount.Balance += amount

	ownerData, err = json.Marshal(ownerAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal owner data: %v", err)
	}
	recipientData, err = json.Marshal(recipientAccount)
	if err != nil {
		return fmt.Errorf("failed to marshal recipient data: %v", err)
	}
	spenderData, err = json.Marshal(spenderUser)
	if err != nil {
		return fmt.Errorf("failed to marshal recipient data: %v", err)
	}

	err = ctx.GetStub().PutState(owner, ownerData)
	if err != nil {
		return fmt.Errorf("failed to update owner balance: %v", err)
	}
	err = ctx.GetStub().PutState(recipient, recipientData)
	if err != nil {
		return fmt.Errorf("failed to update recipient balance: %v", err)
	}
	err = ctx.GetStub().PutState(spender, spenderData)
	if err != nil {
		return fmt.Errorf("failed to update spender balance: %v", err)
	}

	return nil
}

func main() {
	// Create a new chaincode and start it
	chaincode, err := contractapi.NewChaincode(&TokenContract{})
	if err != nil {
		fmt.Printf("Error creating chaincode: %v", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %v", err)
	}
}