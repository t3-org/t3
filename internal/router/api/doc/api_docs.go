// Package doc.
//
// Space API docs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: space.app
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
