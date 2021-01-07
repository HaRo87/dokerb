// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "HaRo87",
            "url": "https://github.com/HaRo87"
        },
        "license": {
            "name": "MIT",
            "url": "https://github.com/HaRo87/dokerb/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/docs": {
            "get": {
                "description": "Get a list of helpful documentation resources",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "documentation"
                ],
                "summary": "Get the documentation info",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.DocResponse"
                        }
                    }
                }
            }
        },
        "/sessions": {
            "post": {
                "description": "Creates a new Doker session and responds with the corresponding token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Create a new Doker session",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}": {
            "delete": {
                "description": "Deletes a existing Doker session based on the provided token",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "session"
                ],
                "summary": "Delete a existing Doker session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/estimates": {
            "get": {
                "description": "Gets all estimates of all existing users of all existing work packages inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "Get the estimates of all users for all work packages",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.PerUserEstimateResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a estimate of a existing user of a existing work package inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "Add the estimate of a user for a work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New Estimate",
                        "name": "estimate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apiserver.PerUserEstimate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/estimates/{id}": {
            "get": {
                "description": "Gets the average estimate of all existing users of a existing work package inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "Get the average estimate of all users for a specific work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Work Package ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.CalcEstimate"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/estimates/{id}/users/distance": {
            "get": {
                "description": "Gets the users with max distance in their estimates of a existing work package inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "Get the users with max distance between their estimates for a specific work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Work Package ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.UsersResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/estimates/{user}/{id}": {
            "delete": {
                "description": "Removes a estimate of a existing user of a existing work package inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "estimate"
                ],
                "summary": "Remove the estimate of a user for a work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "User Name",
                        "name": "user",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Work Package ID",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/users": {
            "get": {
                "description": "Gets all users of an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Get the users of an existing session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.UsersResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new (non-existing) user to an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Add a new user to a existing session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New User",
                        "name": "user",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apiserver.User"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/users/{name}": {
            "delete": {
                "description": "Removes a existing user from an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "user"
                ],
                "summary": "Remove a user from a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "Name of the user",
                        "name": "name",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/workpackages": {
            "get": {
                "description": "Gets all work packages of an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workpackage"
                ],
                "summary": "Get the work packages of a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.WorkPackagesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Adds a new (non-existing) work package to an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workpackage"
                ],
                "summary": "Add a new work package to a existing session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New Work Package",
                        "name": "workpackage",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apiserver.WorkPackage"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/workpackages/{id}": {
            "put": {
                "description": "Updates a estimate of a existing work package inside a existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workpackage"
                ],
                "summary": "Update the estimate of a work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the work package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "New Estimate",
                        "name": "estimate",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/apiserver.Estimate"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Removes a existing work package from an existing session",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workpackage"
                ],
                "summary": "Remove a work package from a session",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the work package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/sessions/{token}/workpackages/{id}/estimate": {
            "delete": {
                "description": "Removes the estimate from an existing work package",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "workpackage"
                ],
                "summary": "Delete the estimate from a work package",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Session Token",
                        "name": "token",
                        "in": "path",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "ID of the work package",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/apiserver.GeneralResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/apiserver.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "apiserver.CalcEstimate": {
            "type": "object",
            "properties": {
                "estimate": {
                    "format": "Estimate",
                    "$ref": "#/definitions/apiserver.Estimate"
                },
                "hint": {
                    "type": "string",
                    "format": "string",
                    "example": "not all users provided estimates"
                },
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "warning"
                },
                "users": {
                    "type": "array",
                    "format": "[]string",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Tigger"
                    ]
                }
            }
        },
        "apiserver.DocEntry": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "format": "string",
                    "example": "GitHub"
                },
                "url": {
                    "type": "string",
                    "format": "string",
                    "example": "https://github.com/HaRo87"
                }
            }
        },
        "apiserver.DocResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "ok"
                },
                "results": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/apiserver.DocEntry"
                    }
                }
            }
        },
        "apiserver.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "error"
                },
                "reason": {
                    "type": "string",
                    "format": "string",
                    "example": "oops, something went wrong"
                }
            }
        },
        "apiserver.Estimate": {
            "type": "object",
            "properties": {
                "effort": {
                    "type": "number",
                    "format": "float64",
                    "example": 1.5
                },
                "standarddeviation": {
                    "type": "number",
                    "format": "float64",
                    "example": 0.2
                }
            }
        },
        "apiserver.GeneralResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "ok"
                },
                "route": {
                    "type": "string",
                    "format": "string",
                    "example": "/sessions/token"
                }
            }
        },
        "apiserver.PerUserEstimate": {
            "type": "object",
            "properties": {
                "b": {
                    "type": "number",
                    "format": "float64",
                    "example": 1.5
                },
                "id": {
                    "type": "string",
                    "format": "string",
                    "example": "TEST01"
                },
                "m": {
                    "type": "number",
                    "format": "float64",
                    "example": 2
                },
                "user": {
                    "type": "string",
                    "format": "string",
                    "example": "Tigger"
                },
                "w": {
                    "type": "number",
                    "format": "float64",
                    "example": 3.6
                }
            }
        },
        "apiserver.PerUserEstimateResponse": {
            "type": "object",
            "properties": {
                "estimates": {
                    "type": "array",
                    "format": "[]datastore.Estimate",
                    "items": {
                        "$ref": "#/definitions/datastore.Estimate"
                    }
                },
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "ok"
                }
            }
        },
        "apiserver.User": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string",
                    "format": "string",
                    "example": "Tigger"
                }
            }
        },
        "apiserver.UsersResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "ok"
                },
                "users": {
                    "type": "array",
                    "format": "[]string",
                    "items": {
                        "type": "string"
                    },
                    "example": [
                        "Tigger",
                        "Rabbit"
                    ]
                }
            }
        },
        "apiserver.WorkPackage": {
            "type": "object",
            "properties": {
                "id": {
                    "type": "string",
                    "format": "string",
                    "example": "TEST01"
                },
                "summary": {
                    "type": "string",
                    "format": "string",
                    "example": "a sample task"
                }
            }
        },
        "apiserver.WorkPackagesResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string",
                    "format": "string",
                    "example": "ok"
                },
                "workpackages": {
                    "type": "array",
                    "format": "[]datastore.WorkPackage",
                    "items": {
                        "$ref": "#/definitions/datastore.WorkPackage"
                    }
                }
            }
        },
        "datastore.Estimate": {
            "type": "object",
            "properties": {
                "bestCase": {
                    "type": "number"
                },
                "mostLikelyCase": {
                    "type": "number"
                },
                "userName": {
                    "type": "string"
                },
                "workPackageID": {
                    "type": "string"
                },
                "worstCase": {
                    "type": "number"
                }
            }
        },
        "datastore.WorkPackage": {
            "type": "object",
            "properties": {
                "effort": {
                    "type": "number"
                },
                "id": {
                    "type": "string"
                },
                "standardDeviation": {
                    "type": "number"
                },
                "summary": {
                    "type": "string"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "0.1.0",
	Host:        "localhost:5000",
	BasePath:    "/api",
	Schemes:     []string{},
	Title:       "Doker Backend API",
	Description: "A backend for playing Planning Poker with Delphi estimate method.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register(swag.Name, &s{})
}
