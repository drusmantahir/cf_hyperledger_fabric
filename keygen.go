package main

import (
	"encoding/json"
	"errors" // Add this import
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// KeyData represents the structure stored in world state for each key
type KeyData struct {
	Consumed bool   `json:"consumed"`
	TxID     string `json:"txID"`     // Transaction ID associated with the key
	ImageURL string `json:"imageURL"` // URL or path to the associated image
}

// KeyContract provides functions for managing keys
type KeyContract struct {
	contractapi.Contract
}

// CreateKey stores a key with consumed=false and associates it with the current transaction ID.
func (kc *KeyContract) CreateKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	// Check if the key already exists
	existingBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if len(existingBytes) != 0 {
		return "", fmt.Errorf("key %s already exists", key)
	}

	// Get the current transaction ID
	txID := ctx.GetStub().GetTxID()

	// Create KeyData with consumed=false and the current transaction ID
	data := KeyData{
		Consumed: false,
		TxID:     txID,
	}

	dataBytes, _ := json.Marshal(data)

	// Store the key in the world state
	if err := ctx.GetStub().PutState(key, dataBytes); err != nil {
		return "", fmt.Errorf("failed to put key %s in world state: %v", key, err)
	}

	// Build the response
	response := map[string]interface{}{
		"message":      "Key created successfully",
		"key":          key,
		"transactionID": txID,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// CreateBulkKeys generates multiple keys with consumed=false and associates them with the current transaction ID.
func (kc *KeyContract) CreateBulkKeys(ctx contractapi.TransactionContextInterface, keys []string) (string, error) {
	createdKeys := []string{}
	failedKeys := map[string]string{}

	for _, key := range keys {
		// Check if the key already exists
		existingBytes, err := ctx.GetStub().GetState(key)
		if err != nil {
			failedKeys[key] = fmt.Sprintf("Failed to read from world state: %v", err)
			continue
		}
		if len(existingBytes) != 0 {
			failedKeys[key] = "Key already exists"
			continue
		}

		// Get the current transaction ID
		txID := ctx.GetStub().GetTxID()

		// Create KeyData with consumed=false and the current transaction ID
		data := KeyData{
			Consumed: false,
			TxID:     txID,
		}

		dataBytes, _ := json.Marshal(data)

		// Store the key in the world state
		if err := ctx.GetStub().PutState(key, dataBytes); err != nil {
			failedKeys[key] = fmt.Sprintf("Failed to put key in world state: %v", err)
			continue
		}

		createdKeys = append(createdKeys, key)
	}

	// Build the response
	response := map[string]interface{}{
		"message":      "Bulk key creation completed",
		"createdKeys":  createdKeys,
		"failedKeys":   failedKeys,
		"transactionID": ctx.GetStub().GetTxID(),
	}

	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// CreateKeyWithImage stores a key with consumed=false, image provided by query and associates it with the current transaction ID.
func (kc *KeyContract) CreateKeyWithImage(ctx contractapi.TransactionContextInterface, key string, imageURL string) (string, error) {
	// Check if key already exists
	existingBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if len(existingBytes) != 0 {
		return "", fmt.Errorf("key %s already exists", key)
	}

	// Get the current transaction ID
	txID := ctx.GetStub().GetTxID()

	// Create KeyData with consumed=false, transaction ID, and the image URL
	data := KeyData{
		Consumed: false,
		TxID:     txID,
		ImageURL: imageURL,
	}

	dataBytes, _ := json.Marshal(data)

	// PutState
	if err := ctx.GetStub().PutState(key, dataBytes); err != nil {
		return "", fmt.Errorf("failed to put key %s in world state: %v", key, err)
	}

	response := map[string]interface{}{
		"message":      "Key with image created successfully",
		"key":          key,
		"transactionID": txID,
		"imageURL":     imageURL,
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ConsumeKey marks an existing key as consumed=true
func (kc *KeyContract) ConsumeKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	// Read existing data
	existingBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to get state for %s: %v", key, err)
	}
	if len(existingBytes) == 0 {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	var data KeyData
	if err := json.Unmarshal(existingBytes, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal data for key %s: %v", key, err)
	}

	if data.Consumed {
		return "", errors.New("the key is already consumed")
	}

	// Update consumed to true
	data.Consumed = true
	dataBytes, _ := json.Marshal(data)

	if err := ctx.GetStub().PutState(key, dataBytes); err != nil {
		return "", fmt.Errorf("failed to update key %s: %v", key, err)
	}

	response := map[string]interface{}{
		"message":      "Key consumed successfully",
		"key":          key,
		"transactionID": ctx.GetStub().GetTxID(),
	}
	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// ReadKey retrieves the consumption status, transaction ID, and image URL (if exists) of the key
func (kc *KeyContract) ReadKey(ctx contractapi.TransactionContextInterface, key string) (string, error) {
	dataBytes, err := ctx.GetStub().GetState(key)
	if err != nil {
		return "", fmt.Errorf("failed to read from world state: %v", err)
	}
	if len(dataBytes) == 0 {
		return "", fmt.Errorf("key %s does not exist", key)
	}

	var data KeyData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return "", err
	}

	response := map[string]interface{}{
		"key":          key,
		"consumed":     data.Consumed,
		"transactionID": data.TxID,
	}

	if data.ImageURL != "" {
		response["imageURL"] = data.ImageURL
	}

	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

// QueryAllKeys retrieves all keys from the world state along with their transaction IDs
func (kc *KeyContract) QueryAllKeys(ctx contractapi.TransactionContextInterface) (string, error) {
	iterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return "", fmt.Errorf("failed to get keys from world state: %v", err)
	}
	defer iterator.Close()

	var keys []map[string]interface{}

	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return "", fmt.Errorf("failed to iterate through keys: %v", err)
		}

		var data KeyData
		if err := json.Unmarshal(queryResponse.Value, &data); err != nil {
			return "", fmt.Errorf("failed to unmarshal key data: %v", err)
		}

		keys = append(keys, map[string]interface{}{
			"key":          queryResponse.Key,
			"transactionID": data.TxID,
			"consumed":     data.Consumed,
		})
	}

	response := map[string]interface{}{
		"message": "All keys retrieved successfully",
		"keys":    keys,
	}

	responseJSON, _ := json.Marshal(response)
	return string(responseJSON), nil
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(KeyContract))
	if err != nil {
		fmt.Printf("Error creating keychaincode: %v\n", err)
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting keychaincode: %v\n", err)
	}
}