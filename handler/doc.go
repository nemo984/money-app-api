// Package classification Money Management API
//
// Documentation for Money Management API
//
// Schemes: http, ws
// BasePath: /api
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// SecurityDefinitions:
// bearerAuth:
//   type: apiKey
//   name: Authorization
//   in: header
//   description: Token obtained from loginUser endpoint
// cookieAuth:
//   type: apiKey
//   in: cookie
//   name: jwt-token
//
// swagger:meta
package handler

// swagger:response genericError
type genericError struct {
	// in: body
	Body struct {
		Error string `json:"error,omitempty"`
	}
}

// Username is already taken
// swagger:response usernameTakenError
type usernameTakenError struct {
	genericError
}

// Username not found or incorrect password
// swagger:response userLoginError
type userLoginError struct {
	genericError
}
