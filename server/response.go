//
//  response.go
//  mercury
//
//  Copyright (c) 2017 Miguel Ángel Ortuño. All rights reserved.
//

package server

type ResponseSuccess struct {
	Success bool		`json:"success"`
	Data	interface{}	`json:"data,omitempty"`
}

type ResponseError struct {
	Success bool	`json:"success"`
	Error	string	`json:"error"`
}

func NewSuccessResponse(data interface{}) *ResponseSuccess {
	return &ResponseSuccess{Success: true, Data: data}
}

func NewErrorResponse(error string) *ResponseError {
	return &ResponseError{Success: false, Error: error}
}
