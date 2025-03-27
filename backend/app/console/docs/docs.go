// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/auth/login/": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Login user",
                "parameters": [
                    {
                        "description": "Login request fields",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.LoginRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "403": {
                        "description": "Forbidden",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        },
        "/auth/register/": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Auth"
                ],
                "summary": "Register new user",
                "parameters": [
                    {
                        "description": "Register request fields",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.RegisterRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        },
        "/info/messages": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Info"
                ],
                "summary": "temp endpoint for testing, provides json object with say field containing 'Hello, World!'.",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/info.MessagesResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        },
        "/users/": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Get user from Authorization token",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token to authorize access",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.UserView"
                        }
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        },
        "/users/change-password/": {
            "patch": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Update user's password",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token to authorize access",
                        "name": "Authorization",
                        "in": "header",
                        "required": true
                    },
                    {
                        "description": "Update password fields",
                        "name": "request",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/users.UpdatePasswordRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK"
                    },
                    "401": {
                        "description": "Unauthorized",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        },
        "/users/{id}/": {
            "get": {
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Users"
                ],
                "summary": "Provides user public info by id",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Bearer token to authorize access",
                        "name": "Authorization",
                        "in": "header"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/users.UserPublicView"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/common.ErrResponseCode"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "common.ErrResponseCode": {
            "type": "object",
            "properties": {
                "code": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                },
                "reason": {
                    "type": "string"
                }
            }
        },
        "info.MessagesResponse": {
            "type": "object",
            "properties": {
                "say": {
                    "type": "string"
                }
            }
        },
        "users.AuthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                },
                "user": {
                    "$ref": "#/definitions/users.UserView"
                }
            }
        },
        "users.LoginRequest": {
            "type": "object",
            "properties": {
                "identifier": {
                    "description": "INFO: Email or phone number.",
                    "type": "string"
                },
                "password": {
                    "type": "string"
                }
            }
        },
        "users.RegisterRequest": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "password": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "post": {
                    "type": "string"
                },
                "postDepartment": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "users.UpdatePasswordRequest": {
            "type": "object",
            "properties": {
                "newPass": {
                    "type": "string"
                },
                "oldPass": {
                    "type": "string"
                }
            }
        },
        "users.UserPublicView": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "fileName": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        },
        "users.UserView": {
            "type": "object",
            "properties": {
                "city": {
                    "type": "string"
                },
                "email": {
                    "type": "string"
                },
                "fileName": {
                    "type": "string"
                },
                "firstName": {
                    "type": "string"
                },
                "id": {
                    "type": "string"
                },
                "lastName": {
                    "type": "string"
                },
                "phoneNumber": {
                    "type": "string"
                },
                "post": {
                    "type": "string"
                },
                "postDepartment": {
                    "type": "string"
                },
                "website": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "v0",
	Host:             "",
	BasePath:         "/api/v0",
	Schemes:          []string{},
	Title:            "Swagger API for One-Help app.",
	Description:      "This is Swagger's auto-generated documentation for One-Help app.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
