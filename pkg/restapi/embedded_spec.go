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
    "/config": {
      "get": {
        "produces": [
          "application/json",
          "application/xml"
        ],
        "tags": [
          "xapp"
        ],
        "summary": "Returns the configuration of all xapps",
        "operationId": "getXappConfigList",
        "responses": {
          "200": {
            "description": "successful query of xApp config",
            "schema": {
              "$ref": "#/definitions/XappConfigList"
            }
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Returns list of subscriptions",
        "operationId": "getAllSubscriptions",
        "responses": {
          "200": {
            "description": "successful query of subscriptions",
            "schema": {
              "$ref": "#/definitions/SubscriptionList"
            }
          },
          "500": {
            "description": "Internal error"
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Subscribe a list of X2AP event triggers to receive messages sent by RAN",
        "operationId": "Subscribe",
        "parameters": [
          {
            "description": "Subscription parameters",
            "name": "SubscriptionParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SubscriptionParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResponse"
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
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Unsubscribe X2AP events from Subscription Manager",
        "operationId": "Unsubscribe",
        "parameters": [
          {
            "type": "string",
            "description": "The subscriptionId received in the Subscription Response",
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
            "description": "Invalid subscriptionId supplied"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    }
  },
  "definitions": {
    "ActionDefinition": {
      "description": "E2SM Octet string. ActionDefinition is an OPTIONAL IE",
      "type": "object",
      "properties": {
        "OctetString": {
          "type": "string"
        }
      }
    },
    "ActionToBeSetup": {
      "type": "object",
      "required": [
        "ActionID",
        "ActionType"
      ],
      "properties": {
        "ActionDefinition": {
          "$ref": "#/definitions/ActionDefinition"
        },
        "ActionID": {
          "type": "integer",
          "maximum": 255
        },
        "ActionType": {
          "type": "string",
          "enum": [
            "insert",
            "policy",
            "report"
          ]
        },
        "SubsequentAction": {
          "$ref": "#/definitions/SubsequentAction"
        }
      }
    },
    "ActionsToBeSetup": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/ActionToBeSetup"
      }
    },
    "ConfigMetadata": {
      "type": "object",
      "required": [
        "xappName",
        "configType"
      ],
      "properties": {
        "configType": {
          "description": "The type of the content",
          "type": "string",
          "enum": [
            "json",
            "xml",
            "other"
          ]
        },
        "xappName": {
          "description": "Name of the xApp",
          "type": "string"
        }
      }
    },
    "EventTriggerDefinition": {
      "description": "E2SM Octet string",
      "type": "object",
      "properties": {
        "OctetString": {
          "type": "string"
        }
      }
    },
    "SubscriptionData": {
      "type": "object",
      "properties": {
        "Endpoint": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "Meid": {
          "type": "string"
        },
        "SubscriptionId": {
          "type": "integer"
        }
      }
    },
    "SubscriptionDetails": {
      "type": "object",
      "required": [
        "EventTriggerList",
        "ActionToBeSetupList"
      ],
      "properties": {
        "ActionToBeSetupList": {
          "$ref": "#/definitions/ActionsToBeSetup"
        },
        "EventTriggerList": {
          "$ref": "#/definitions/EventTriggerDefinition"
        }
      }
    },
    "SubscriptionDetailsList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/SubscriptionDetails"
      }
    },
    "SubscriptionInstance": {
      "type": "object",
      "required": [
        "RequestorId",
        "InstanceId",
        "ErrorCause"
      ],
      "properties": {
        "ErrorCause": {
          "description": "Empty string when no error.",
          "type": "string"
        },
        "InstanceId": {
          "type": "integer",
          "maximum": 65535
        },
        "RequestorId": {
          "type": "integer",
          "maximum": 65535
        }
      }
    },
    "SubscriptionList": {
      "description": "A list of subscriptions",
      "type": "array",
      "items": {
        "$ref": "#/definitions/SubscriptionData"
      }
    },
    "SubscriptionParams": {
      "type": "object",
      "required": [
        "ClientEndpoint",
        "Meid",
        "RequestorId",
        "InstanceId",
        "RANFunctionID",
        "SubscriptionDetails"
      ],
      "properties": {
        "ClientEndpoint": {
          "description": "xApp service address and port",
          "type": "object",
          "properties": {
            "Port": {
              "description": "xApp service address port",
              "type": "integer",
              "maximum": 65535
            },
            "ServiceName": {
              "description": "xApp service address name like 'service-ricxapp-xappname-http.ricxapp'",
              "type": "string"
            }
          }
        },
        "InstanceId": {
          "type": "integer",
          "maximum": 65535
        },
        "Meid": {
          "type": "string"
        },
        "RANFunctionID": {
          "type": "integer",
          "maximum": 4095
        },
        "RequestorId": {
          "type": "integer",
          "maximum": 65535
        },
        "SubscriptionDetails": {
          "$ref": "#/definitions/SubscriptionDetailsList"
        }
      }
    },
    "SubscriptionResponse": {
      "type": "object",
      "required": [
        "SubscriptionId",
        "SubscriptionInstances"
      ],
      "properties": {
        "SubscriptionId": {
          "type": "string"
        },
        "SubscriptionInstances": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/SubscriptionInstance"
          }
        }
      }
    },
    "SubsequentAction": {
      "description": "SubsequentAction is an OPTIONAL IE",
      "type": "object",
      "required": [
        "SubsequentActionType",
        "TimeToWait"
      ],
      "properties": {
        "SubsequentActionType": {
          "type": "string",
          "enum": [
            "continue",
            "wait"
          ]
        },
        "TimeToWait": {
          "type": "string",
          "enum": [
            "zero",
            "w1ms",
            "w2ms",
            "w5ms",
            "w10ms",
            "w20ms",
            "w30ms",
            "w40ms",
            "w50ms",
            "w100ms",
            "w200ms",
            "w500ms",
            "w1s",
            "w2s",
            "w5s",
            "w10s",
            "w20s",
            "w60s"
          ]
        }
      }
    },
    "XAppConfig": {
      "type": "object",
      "required": [
        "metadata",
        "config"
      ],
      "properties": {
        "config": {
          "description": "Configuration in JSON format",
          "type": "object"
        },
        "metadata": {
          "$ref": "#/definitions/ConfigMetadata"
        }
      }
    },
    "XappConfigList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/XAppConfig"
      }
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
    "/config": {
      "get": {
        "produces": [
          "application/json",
          "application/xml"
        ],
        "tags": [
          "xapp"
        ],
        "summary": "Returns the configuration of all xapps",
        "operationId": "getXappConfigList",
        "responses": {
          "200": {
            "description": "successful query of xApp config",
            "schema": {
              "$ref": "#/definitions/XappConfigList"
            }
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    },
    "/subscriptions": {
      "get": {
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Returns list of subscriptions",
        "operationId": "getAllSubscriptions",
        "responses": {
          "200": {
            "description": "successful query of subscriptions",
            "schema": {
              "$ref": "#/definitions/SubscriptionList"
            }
          },
          "500": {
            "description": "Internal error"
          }
        }
      },
      "post": {
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Subscribe a list of X2AP event triggers to receive messages sent by RAN",
        "operationId": "Subscribe",
        "parameters": [
          {
            "description": "Subscription parameters",
            "name": "SubscriptionParams",
            "in": "body",
            "schema": {
              "$ref": "#/definitions/SubscriptionParams"
            }
          }
        ],
        "responses": {
          "201": {
            "description": "Subscription successfully created",
            "schema": {
              "$ref": "#/definitions/SubscriptionResponse"
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
        "consumes": [
          "application/json"
        ],
        "produces": [
          "application/json"
        ],
        "tags": [
          "common"
        ],
        "summary": "Unsubscribe X2AP events from Subscription Manager",
        "operationId": "Unsubscribe",
        "parameters": [
          {
            "type": "string",
            "description": "The subscriptionId received in the Subscription Response",
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
            "description": "Invalid subscriptionId supplied"
          },
          "500": {
            "description": "Internal error"
          }
        }
      }
    }
  },
  "definitions": {
    "ActionDefinition": {
      "description": "E2SM Octet string. ActionDefinition is an OPTIONAL IE",
      "type": "object",
      "properties": {
        "OctetString": {
          "type": "string"
        }
      }
    },
    "ActionToBeSetup": {
      "type": "object",
      "required": [
        "ActionID",
        "ActionType"
      ],
      "properties": {
        "ActionDefinition": {
          "$ref": "#/definitions/ActionDefinition"
        },
        "ActionID": {
          "type": "integer",
          "maximum": 255,
          "minimum": 0
        },
        "ActionType": {
          "type": "string",
          "enum": [
            "insert",
            "policy",
            "report"
          ]
        },
        "SubsequentAction": {
          "$ref": "#/definitions/SubsequentAction"
        }
      }
    },
    "ActionsToBeSetup": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/ActionToBeSetup"
      }
    },
    "ConfigMetadata": {
      "type": "object",
      "required": [
        "xappName",
        "configType"
      ],
      "properties": {
        "configType": {
          "description": "The type of the content",
          "type": "string",
          "enum": [
            "json",
            "xml",
            "other"
          ]
        },
        "xappName": {
          "description": "Name of the xApp",
          "type": "string"
        }
      }
    },
    "EventTriggerDefinition": {
      "description": "E2SM Octet string",
      "type": "object",
      "properties": {
        "OctetString": {
          "type": "string"
        }
      }
    },
    "SubscriptionData": {
      "type": "object",
      "properties": {
        "Endpoint": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "Meid": {
          "type": "string"
        },
        "SubscriptionId": {
          "type": "integer"
        }
      }
    },
    "SubscriptionDetails": {
      "type": "object",
      "required": [
        "EventTriggerList",
        "ActionToBeSetupList"
      ],
      "properties": {
        "ActionToBeSetupList": {
          "$ref": "#/definitions/ActionsToBeSetup"
        },
        "EventTriggerList": {
          "$ref": "#/definitions/EventTriggerDefinition"
        }
      }
    },
    "SubscriptionDetailsList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/SubscriptionDetails"
      }
    },
    "SubscriptionInstance": {
      "type": "object",
      "required": [
        "RequestorId",
        "InstanceId",
        "ErrorCause"
      ],
      "properties": {
        "ErrorCause": {
          "description": "Empty string when no error.",
          "type": "string"
        },
        "InstanceId": {
          "type": "integer",
          "maximum": 65535,
          "minimum": 0
        },
        "RequestorId": {
          "type": "integer",
          "maximum": 65535,
          "minimum": 0
        }
      }
    },
    "SubscriptionList": {
      "description": "A list of subscriptions",
      "type": "array",
      "items": {
        "$ref": "#/definitions/SubscriptionData"
      }
    },
    "SubscriptionParams": {
      "type": "object",
      "required": [
        "ClientEndpoint",
        "Meid",
        "RequestorId",
        "InstanceId",
        "RANFunctionID",
        "SubscriptionDetails"
      ],
      "properties": {
        "ClientEndpoint": {
          "description": "xApp service address and port",
          "type": "object",
          "properties": {
            "Port": {
              "description": "xApp service address port",
              "type": "integer",
              "maximum": 65535,
              "minimum": 0
            },
            "ServiceName": {
              "description": "xApp service address name like 'service-ricxapp-xappname-http.ricxapp'",
              "type": "string"
            }
          }
        },
        "InstanceId": {
          "type": "integer",
          "maximum": 65535,
          "minimum": 0
        },
        "Meid": {
          "type": "string"
        },
        "RANFunctionID": {
          "type": "integer",
          "maximum": 4095,
          "minimum": 0
        },
        "RequestorId": {
          "type": "integer",
          "maximum": 65535,
          "minimum": 0
        },
        "SubscriptionDetails": {
          "$ref": "#/definitions/SubscriptionDetailsList"
        }
      }
    },
    "SubscriptionResponse": {
      "type": "object",
      "required": [
        "SubscriptionId",
        "SubscriptionInstances"
      ],
      "properties": {
        "SubscriptionId": {
          "type": "string"
        },
        "SubscriptionInstances": {
          "type": "array",
          "items": {
            "$ref": "#/definitions/SubscriptionInstance"
          }
        }
      }
    },
    "SubsequentAction": {
      "description": "SubsequentAction is an OPTIONAL IE",
      "type": "object",
      "required": [
        "SubsequentActionType",
        "TimeToWait"
      ],
      "properties": {
        "SubsequentActionType": {
          "type": "string",
          "enum": [
            "continue",
            "wait"
          ]
        },
        "TimeToWait": {
          "type": "string",
          "enum": [
            "zero",
            "w1ms",
            "w2ms",
            "w5ms",
            "w10ms",
            "w20ms",
            "w30ms",
            "w40ms",
            "w50ms",
            "w100ms",
            "w200ms",
            "w500ms",
            "w1s",
            "w2s",
            "w5s",
            "w10s",
            "w20s",
            "w60s"
          ]
        }
      }
    },
    "XAppConfig": {
      "type": "object",
      "required": [
        "metadata",
        "config"
      ],
      "properties": {
        "config": {
          "description": "Configuration in JSON format",
          "type": "object"
        },
        "metadata": {
          "$ref": "#/definitions/ConfigMetadata"
        }
      }
    },
    "XappConfigList": {
      "type": "array",
      "items": {
        "$ref": "#/definitions/XAppConfig"
      }
    }
  }
}`))
}
