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
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/admin/image": {
            "post": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Uploads an image to AWS S3",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminUploadImage",
                "parameters": [
                    {
                        "description": "path - path within the S3 bucket",
                        "name": "path",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "width - width of the image to resize. If width and height are missing - then the new image will use the original size",
                        "name": "width",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "height - height of the image to resize. If width and height are missing - then the new image will use the original size",
                        "name": "height",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "quality - quality of the image. Default: 90",
                        "name": "quality",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "fileName - the uploaded file name",
                        "name": "fileName",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/admin/student_guides": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves  all items",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminGetStudentGuides",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Coma separated IDs of the desired records",
                        "name": "ids",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "post": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves  all items",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminCreateStudentGuide",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/admin/student_guides/{id}": {
            "get": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Retrieves  all items",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminGetStudentGuide",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "put": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Updates a student guide with the specified id",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminUpdateStudentGuide",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            },
            "delete": {
                "security": [
                    {
                        "AdminUserAuth": []
                    }
                ],
                "description": "Deletes a student guide with the specified id",
                "tags": [
                    "Admin"
                ],
                "operationId": "AdminDeleteStudentGuide",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/image": {
            "post": {
                "security": [
                    {
                        "RokwireAuth": []
                    }
                ],
                "description": "Uploads an image to AWS S3",
                "consumes": [
                    "multipart/form-data"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "AdminUpdateStudentGuide",
                "parameters": [
                    {
                        "description": "path - path within the S3 bucket",
                        "name": "path",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "width - width of the image to resize. If width and height are missing - then the new image will use the original size",
                        "name": "width",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "height - height of the image to resize. If width and height are missing - then the new image will use the original size",
                        "name": "height",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "quality - quality of the image. Default: 90",
                        "name": "quality",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    },
                    {
                        "description": "fileName - the uploaded file name",
                        "name": "fileName",
                        "in": "body",
                        "schema": {
                            "type": "string"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/student_guides": {
            "get": {
                "security": [
                    {
                        "RokwireAuth": []
                    }
                ],
                "description": "Retrieves  all items",
                "consumes": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "GetStudentGuides",
                "parameters": [
                    {
                        "type": "string",
                        "description": "Coma separated IDs of the desired records",
                        "name": "ids",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/student_guides/{id}": {
            "get": {
                "security": [
                    {
                        "RokwireAuth": []
                    }
                ],
                "description": "Retrieves  all items",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "GetStudentGuide",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/twitter/posts": {
            "get": {
                "security": [
                    {
                        "RokwireAuth": []
                    }
                ],
                "description": "Retrieves top most Twitter posts",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "GetTweeterPosts",
                "parameters": [
                    {
                        "type": "string",
                        "description": "count - the number of the tweets that will be retrieved. Default: 5",
                        "name": "count",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "force - Forced refresh. Default: false",
                        "name": "force",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        },
        "/version": {
            "get": {
                "description": "Gives the service version.",
                "produces": [
                    "text/plain"
                ],
                "tags": [
                    "Client"
                ],
                "operationId": "Version",
                "responses": {
                    "200": {
                        "description": ""
                    }
                }
            }
        }
    },
    "securityDefinitions": {
        "AdminGroupAuth": {
            "type": "apiKey",
            "name": "GROUP",
            "in": "header"
        },
        "AdminUserAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header (add Bearer prefix to the Authorization value)"
        },
        "RokwireAuth": {
            "type": "apiKey",
            "name": "ROKWIRE-API-KEY",
            "in": "header"
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
	Version:     "1.0.7",
	Host:        "localhost",
	BasePath:    "/content",
	Schemes:     []string{"https"},
	Title:       "Rokwire Content Building Block API",
	Description: "Rokwire Content Building Block API Documentation.",
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
