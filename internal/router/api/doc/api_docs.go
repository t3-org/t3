// Package doc.
//
// Shield API docs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: shield.ir
//     BasePath:
//     Version: 0.1.0
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - bearerAuth:
//
//     SecurityDefinitions:
//     bearerAuth:
//          type: apiKey
//          name: Authorization
//          in: header
//
// swagger:meta
package doc

type replyCode struct {
	// response code
	Code string `json:"code"`
}
