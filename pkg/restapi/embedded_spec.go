// Code generated by go-swagger; DO NOT EDIT.

package restapi

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"
)

var (
	// SwaggerJSON embedded version of the swagger document used at generation time
	SwaggerJSON json.RawMessage
	// FlatSwaggerJSON embedded flattened version of the swagger document used at generation time
	FlatSwaggerJSON json.RawMessage
)

func init() {
	SwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the initial REST API for RIC subscription",
    "title": "RIC subscription",
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "0.0.1"
  },
  "host": "hostname",
  "basePath": "/ric/v1",
  "paths": {
    "/subscriptions/control": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "control"
        ],
        "summary": "Subscribe and send \"CONTROL\" message to RAN to initiate or resume call processing in RAN",
        "operationId": "subscribeControl",
        "parameters": [
          {
            "description": "Subscription control parameters",
            "name": "ControlParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/ControlParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/policy": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "policy"
        ],
        "summary": "Subscribe and send \"POLICY\" message to RAN to execute a specific POLICY during call processing in RAN after each occurrence of a defined SUBSCRIPTION",
        "operationId": "subscribePolicy",
        "parameters": [
          {
            "description": "Subscription policy parameters",
            "name": "PolicyParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PolicyParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/report": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "report"
        ],
        "summary": "Subscribe a list of X2AP event triggers to receive \"REPORT\" messages sent by RAN",
        "operationId": "subscribeReport",
        "parameters": [
          {
            "description": "Subscription report parameters",
            "name": "ReportParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/ReportParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/{subscriptionId}": {
      "delete": {
        "tags": [
          "common"
        ],
        "summary": "Unsubscribe X2AP events from Subscription Manager",
        "operationId": "Unsubscribe",
        "parameters": [
          {
            "type": "integer",
            "description": "The subscriptionId to be unsubscribed",
            "name": "subscriptionId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Operation done successfully"
          },
          "400": {
            "description": "Invalid requestorId supplied"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    }
  },
  "definitions": {
    "ControlParams": {
      "type": "object",
      "properties": {
        "RequestorId": {
          "type": "integer"
        },
        "TBD": {
          "type": "string"
        }
      }
    },
    "EventTrigger": {
      "type": "object",
      "required": [
        "InterfaceDirection",
        "ProcedureCode",
        "TypeOfMessage"
      ],
      "properties": {
        "ENBId": {
          "type": "integer"
        },
        "InterfaceDirection": {
          "type": "integer"
        },
        "PlmnId": {
          "type": "string"
        },
        "ProcedureCode": {
          "type": "integer"
        },
        "TypeOfMessage": {
          "type": "integer"
        }
      }
    },
    "EventTriggerList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/EventTrigger"
      }
    },
    "PolicyParams": {
      "type": "object",
      "properties": {
        "RequestorId": {
          "type": "integer"
        },
        "TBD": {
          "type": "string"
        }
      }
    },
    "ReportParams": {
      "type": "object",
      "required": [
        "RequestorId",
        "EventTriggers"
      ],
      "properties": {
        "EventTriggers": {
          "$ref": "#/definitions/EventTriggerList"
        },
        "RequestorId": {
          "type": "integer"
        }
      }
    },
    "SubscriptionResult": {
      "description": "A list of unique IDs",
      "type": "array",
      "items": {
        "type": "integer"
      }
    },
    "SubscriptionType": {
      "type": "string",
      "enum": [
        "control",
        "insert",
        "policy",
        "report"
      ]
    }
  }
}`))
	FlatSwaggerJSON = json.RawMessage([]byte(`{
  "schemes": [
    "http"
  ],
  "swagger": "2.0",
  "info": {
    "description": "This is the initial REST API for RIC subscription",
    "title": "RIC subscription",
    "license": {
      "name": "Apache 2.0",
      "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
    },
    "version": "0.0.1"
  },
  "host": "hostname",
  "basePath": "/ric/v1",
  "paths": {
    "/subscriptions/control": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "control"
        ],
        "summary": "Subscribe and send \"CONTROL\" message to RAN to initiate or resume call processing in RAN",
        "operationId": "subscribeControl",
        "parameters": [
          {
            "description": "Subscription control parameters",
            "name": "ControlParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/ControlParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/policy": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "policy"
        ],
        "summary": "Subscribe and send \"POLICY\" message to RAN to execute a specific POLICY during call processing in RAN after each occurrence of a defined SUBSCRIPTION",
        "operationId": "subscribePolicy",
        "parameters": [
          {
            "description": "Subscription policy parameters",
            "name": "PolicyParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/PolicyParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/report": {
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "report"
        ],
        "summary": "Subscribe a list of X2AP event triggers to receive \"REPORT\" messages sent by RAN",
        "operationId": "subscribeReport",
        "parameters": [
          {
            "description": "Subscription report parameters",
            "name": "ReportParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/ReportParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResult"
            }
          },
          "400": {
            "description": "Invalid input"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions/{subscriptionId}": {
      "delete": {
        "tags": [
          "common"
        ],
        "summary": "Unsubscribe X2AP events from Subscription Manager",
        "operationId": "Unsubscribe",
        "parameters": [
          {
            "type": "integer",
            "description": "The subscriptionId to be unsubscribed",
            "name": "subscriptionId",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "204": {
            "description": "Operation done successfully"
          },
          "400": {
            "description": "Invalid requestorId supplied"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    }
  },
  "definitions": {
    "ControlParams": {
      "type": "object",
      "properties": {
        "RequestorId": {
          "type": "integer"
        },
        "TBD": {
          "type": "string"
        }
      }
    },
    "EventTrigger": {
      "type": "object",
      "required": [
        "InterfaceDirection",
        "ProcedureCode",
        "TypeOfMessage"
      ],
      "properties": {
        "ENBId": {
          "type": "integer"
        },
        "InterfaceDirection": {
          "type": "integer"
        },
        "PlmnId": {
          "type": "string"
        },
        "ProcedureCode": {
          "type": "integer"
        },
        "TypeOfMessage": {
          "type": "integer"
        }
      }
    },
    "EventTriggerList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/EventTrigger"
      }
    },
    "PolicyParams": {
      "type": "object",
      "properties": {
        "RequestorId": {
          "type": "integer"
        },
        "TBD": {
          "type": "string"
        }
      }
    },
    "ReportParams": {
      "type": "object",
      "required": [
        "RequestorId",
        "EventTriggers"
      ],
      "properties": {
        "EventTriggers": {
          "$ref": "#/definitions/EventTriggerList"
        },
        "RequestorId": {
          "type": "integer"
        }
      }
    },
    "SubscriptionResult": {
      "description": "A list of unique IDs",
      "type": "array",
      "items": {
        "type": "integer"
      }
    },
    "SubscriptionType": {
      "type": "string",
      "enum": [
        "control",
        "insert",
        "policy",
        "report"
      ]
    }
  }
}`))
}
