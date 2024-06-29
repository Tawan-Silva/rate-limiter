// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "termsOfService": "http://swagger.io/terms/",
        "contact": {
            "name": "Tawan Silva",
            "url": "https://www.linkedin.com/in/tawan-silva-684b581b7/",
            "email": "tawan.tls43@gmail.com"
        },
        "license": {
            "name": "Rate Limiter License",
            "url": "http://www.ratelimiter.com.br"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/get-all-rate-limiter": {
            "get": {
                "description": "get all rate limiter settings",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate limiter"
                ],
                "summary": "Get all rate limiter settings",
                "responses": {
                    "200": {
                        "description": "Successfully retrieved all rate limiter settings",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/middleware.LimitData"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/home": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get index",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "home"
                ],
                "summary": "Welcome to the rate limited index page!",
                "responses": {
                    "200": {
                        "description": "Welcome to the rate limited index page!",
                        "schema": {
                            "$ref": "#/definitions/server.IndexResponse"
                        }
                    }
                }
            }
        },
        "/token": {
            "get": {
                "description": "get authToken",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "token"
                ],
                "summary": "Generates a new auth token",
                "responses": {
                    "200": {
                        "description": "token detail",
                        "schema": {
                            "$ref": "#/definitions/server.AuthTokenResponse"
                        }
                    }
                }
            }
        },
        "/update-rate-limiter/": {
            "put": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "update rate limiter settings for a specific key (ip or token)",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate limiter"
                ],
                "summary": "Update rate limiter settings",
                "parameters": [
                    {
                        "description": "Update rate limiter settings",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/middleware.LimitDataInput"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Successfully updated rate limiter settings",
                        "schema": {
                            "$ref": "#/definitions/middleware.LimitDataInput"
                        }
                    },
                    "400": {
                        "description": "Invalid request body",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal server error",
                        "schema": {
                            "$ref": "#/definitions/middleware.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "middleware.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "middleware.LimitData": {
            "description": "Struct to store rate limiter data",
            "type": "object",
            "properties": {
                "block_duration": {
                    "type": "integer"
                },
                "id": {
                    "type": "string"
                },
                "key": {
                    "type": "string"
                },
                "max_requests": {
                    "type": "integer"
                },
                "seconds": {
                    "type": "integer"
                }
            }
        },
        "middleware.LimitDataInput": {
            "description": "Struct to store rate limiter data for Swagger documentation",
            "type": "object",
            "properties": {
                "block_duration": {
                    "type": "integer"
                },
                "max_requests": {
                    "type": "integer"
                },
                "seconds": {
                    "type": "integer"
                }
            }
        },
        "server.AuthTokenResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "server.IndexResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "API_KEY",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{},
	Title:            "Rate Limiter API Example",
	Description:      "Rate Limiter API with Redis",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
