package socket

import (
	"encoding/json"
	"fmt"
)

const create_array = "CreateArray"

func handleCreateArray(in Inbound) (Outbound, error) {
	var gn CreateArrayRequest
	if err := json.Unmarshal(in.Data, &gn); err != nil {
		return Outbound{
			ID:    in.ID,
			Error: &CreateArrayInvalidDataErr,
		}, loggingError(create_array, in.ID, fmt.Errorf("failed to unmarshal inbound data: %v", err))
	}
	if gn.Number < 0 {
		return Outbound{
			ID:    in.ID,
			Error: &CreateArrayNegativeNumberErr,
		}, loggingError(create_array, in.ID, fmt.Errorf("negative number provided"))
	}
	var arr []int
	for i := 0; i < gn.Number; i++ {
		arr = append(arr, i)
	}

	data, err := json.Marshal(&CreateArrayResponse{Numbers: arr})
	if err != nil {
		return Outbound{
			ID:    in.ID,
			Error: &InternalErrorMessageErr,
		}, loggingError(create_array, in.ID, fmt.Errorf("failed to marshal outbound response: %v", err))
	}

	return Outbound{
		ID:   in.ID,
		Data: string(data),
	}, nil
}

type CreateArrayRequest struct {
	Number int `json:"number"`
}

type CreateArrayResponse struct {
	Numbers []int `json:"numbers"`
}

var (
	CreateArrayInvalidDataErr    = Error{Code: "invalid_create_numbers_data", Message: "Invalid create array data."}
	CreateArrayEmptyNumberErr    = Error{Code: "create_numbers_no_number_provided", Message: "Number must be provided."}
	CreateArrayNegativeNumberErr = Error{Code: "create_numbers_negative_number", Message: "Negative number provided."}
)
