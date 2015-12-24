package src

//from yadisk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ErrorResponse represents erroneous API response.
// Implements go's built in `error`.
type ErrorResponse struct {
	ErrorName   string `json:"error"`
	Description string `json:"description"`
	Message     string `json:"message"`

	StatusCode int `json:""`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("[%d - %s] %s (%s)", e.StatusCode, e.ErrorName, e.Description, e.Message)
}

// ProccessErrorResponse tries to represent data passed as
// an ErrorResponse object.
func ProccessErrorResponse(data io.Reader) (*ErrorResponse, error) {
	dec := json.NewDecoder(data)
	var errorResponse ErrorResponse

	if err := dec.Decode(&errorResponse); err == io.EOF {
		// ok
	} else if err != nil {
		return nil, err
	}

	// TODO: check if there is any trash data after JSON and crash if there is.

	return &errorResponse, nil
}

// CheckAPIError is a convenient function to turn erroneous
// API response into go error.
func CheckAPIError(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	errorResponse, err := ProccessErrorResponse(resp.Body)
	if err != nil {
		return err
	}
	errorResponse.StatusCode = resp.StatusCode

	//defer resp.Body.Close()
	defer func() {
		err = resp.Body.Close()
	}()

	return errorResponse
}

// ProccessErrorString tries to represent data passed as
// an ErrorResponse object.
func ProccessErrorString(data string) (*ErrorResponse, error) {
	var errorResponse ErrorResponse
	if err := json.Unmarshal([]byte(data), &errorResponse); err == nil {
		// ok
	} else if err != nil {
		return nil, err
	}

	// TODO: check if there is any trash data after JSON and crash if there is.

	return &errorResponse, nil
}

//Parse API error
func (c *Client) ParseAPIError(jsonErr string) (error, string) { //ErrorName
	errorResponse, err := ProccessErrorString(jsonErr)
	if err != nil {
		return err, err.Error()
	}

	return nil, errorResponse.ErrorName
}
