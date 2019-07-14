/*
==================================================================================
  Copyright (c) 2019 AT&T Intellectual Property.
  Copyright (c) 2019 Nokia

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
==================================================================================
*/

package xapp

/*
#include <rmr/RIC_message_types.h>
*/
import "C"

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
var RICMessageTypes = map[string]int{
	"RIC_SUB_REQ":                  C.RIC_SUB_REQ,
	"RIC_SUB_RESP":                 C.RIC_SUB_RESP,
	"RIC_SUB_FAILURE":              C.RIC_SUB_FAILURE,
	"RIC_SUB_DEL_REQ":              C.RIC_SUB_DEL_REQ,
	"RIC_SUB_DEL_RESP":             C.RIC_SUB_DEL_RESP,
	"RIC_SUB_DEL_FAILURE":          C.RIC_SUB_DEL_FAILURE,
	"RIC_SERVICE_UPDATE":           C.RIC_SERVICE_UPDATE,
	"RIC_SERVICE_UPDATE_ACK":       C.RIC_SERVICE_UPDATE_ACK,
	"RIC_SERVICE_UPDATE_FAILURE":   C.RIC_SERVICE_UPDATE_FAILURE,
	"RIC_CONTROL_REQ":              C.RIC_CONTROL_REQ,
	"RIC_CONTROL_ACK":              C.RIC_CONTROL_ACK,
	"RIC_CONTROL_FAILURE":          C.RIC_CONTROL_FAILURE,
	"RIC_INDICATION":               C.RIC_INDICATION,
	"RIC_SERVICE_QUERY":            C.RIC_SERVICE_QUERY,
	"RIC_X2_SETUP_REQ":             C.RIC_X2_SETUP_REQ,
	"RIC_X2_SETUP_RESP":            C.RIC_X2_SETUP_RESP,
	"RIC_X2_SETUP_FAILURE":         C.RIC_X2_SETUP_FAILURE,
	"RIC_X2_RESET":                 C.RIC_X2_RESET,
	"RIC_X2_RESET_RESP":            C.RIC_X2_RESET_RESP,
	"RIC_ENDC_X2_SETUP_REQ":        C.RIC_ENDC_X2_SETUP_REQ,
	"RIC_ENDC_X2_SETUP_RESP":       C.RIC_ENDC_X2_SETUP_RESP,
	"RIC_ENDC_X2_SETUP_FAILURE":    C.RIC_ENDC_X2_SETUP_FAILURE,
	"RIC_ENDC_CONF_UPDATE":         C.RIC_ENDC_CONF_UPDATE,
	"RIC_ENDC_CONF_UPDATE_ACK":     C.RIC_ENDC_CONF_UPDATE_ACK,
	"RIC_ENDC_CONF_UPDATE_FAILURE": C.RIC_ENDC_CONF_UPDATE_FAILURE,
	"RIC_RES_STATUS_REQ":           C.RIC_RES_STATUS_REQ,
	"RIC_RES_STATUS_RESP":          C.RIC_RES_STATUS_RESP,
	"RIC_RES_STATUS_FAILURE":       C.RIC_RES_STATUS_FAILURE,
	"RIC_ENB_CONF_UPDATE":          C.RIC_ENB_CONF_UPDATE,
	"RIC_ENB_CONF_UPDATE_ACK":      C.RIC_ENB_CONF_UPDATE_ACK,
	"RIC_ENB_CONF_UPDATE_FAILURE":  C.RIC_ENB_CONF_UPDATE_FAILURE,
	"RIC_ENB_LOAD_INFORMATION":     C.RIC_ENB_LOAD_INFORMATION,
	"RIC_GNB_STATUS_INDICATION":    C.RIC_GNB_STATUS_INDICATION,
	"RIC_RESOURCE_STATUS_UPDATE":   C.RIC_RESOURCE_STATUS_UPDATE,
	"RIC_ERROR_INDICATION":         C.RIC_ERROR_INDICATION,
	"DC_ADM_INT_CONTROL":           C.DC_ADM_INT_CONTROL,
	"DC_ADM_INT_CONTROL_ACK":       C.DC_ADM_INT_CONTROL_ACK,
}

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
const (
	RIC_SUB_REQ                  = C.RIC_SUB_REQ
	RIC_SUB_RESP                 = C.RIC_SUB_RESP
	RIC_SUB_FAILURE              = C.RIC_SUB_FAILURE
	RIC_SUB_DEL_REQ              = C.RIC_SUB_DEL_REQ
	RIC_SUB_DEL_RESP             = C.RIC_SUB_DEL_RESP
	RIC_SUB_DEL_FAILURE          = C.RIC_SUB_DEL_FAILURE
	RIC_SERVICE_UPDATE           = C.RIC_SERVICE_UPDATE
	RIC_SERVICE_UPDATE_ACK       = C.RIC_SERVICE_UPDATE_ACK
	RIC_SERVICE_UPDATE_FAILURE   = C.RIC_SERVICE_UPDATE_FAILURE
	RIC_CONTROL_REQ              = C.RIC_CONTROL_REQ
	RIC_CONTROL_ACK              = C.RIC_CONTROL_ACK
	RIC_CONTROL_FAILURE          = C.RIC_CONTROL_FAILURE
	RIC_INDICATION               = C.RIC_INDICATION
	RIC_SERVICE_QUERY            = C.RIC_SERVICE_QUERY
	RIC_X2_SETUP_REQ             = C.RIC_X2_SETUP_REQ
	RIC_X2_SETUP_RESP            = C.RIC_X2_SETUP_RESP
	RIC_X2_SETUP_FAILURE         = C.RIC_X2_SETUP_FAILURE
	RIC_X2_RESET                 = C.RIC_X2_RESET
	RIC_X2_RESET_RESP            = C.RIC_X2_RESET_RESP
	RIC_ENDC_X2_SETUP_REQ        = C.RIC_ENDC_X2_SETUP_REQ
	RIC_ENDC_X2_SETUP_RESP       = C.RIC_ENDC_X2_SETUP_RESP
	RIC_ENDC_X2_SETUP_FAILURE    = C.RIC_ENDC_X2_SETUP_FAILURE
	RIC_ENDC_CONF_UPDATE         = C.RIC_ENDC_CONF_UPDATE
	RIC_ENDC_CONF_UPDATE_ACK     = C.RIC_ENDC_CONF_UPDATE_ACK
	RIC_ENDC_CONF_UPDATE_FAILURE = C.RIC_ENDC_CONF_UPDATE_FAILURE
	RIC_RES_STATUS_REQ           = C.RIC_RES_STATUS_REQ
	RIC_RES_STATUS_RESP          = C.RIC_RES_STATUS_RESP
	RIC_RES_STATUS_FAILURE       = C.RIC_RES_STATUS_FAILURE
	RIC_ENB_CONF_UPDATE          = C.RIC_ENB_CONF_UPDATE
	RIC_ENB_CONF_UPDATE_ACK      = C.RIC_ENB_CONF_UPDATE_ACK
	RIC_ENB_CONF_UPDATE_FAILURE  = C.RIC_ENB_CONF_UPDATE_FAILURE
	RIC_ENB_LOAD_INFORMATION     = C.RIC_ENB_LOAD_INFORMATION
	RIC_GNB_STATUS_INDICATION    = C.RIC_GNB_STATUS_INDICATION
	RIC_RESOURCE_STATUS_UPDATE   = C.RIC_RESOURCE_STATUS_UPDATE
	RIC_ERROR_INDICATION         = C.RIC_ERROR_INDICATION
	DC_ADM_INT_CONTROL           = C.DC_ADM_INT_CONTROL
	DC_ADM_INT_CONTROL_ACK       = C.DC_ADM_INT_CONTROL_ACK
)

//-----------------------------------------------------------------------------
//
//-----------------------------------------------------------------------------
var RicMessageTypeToName = map[int]string{
	RIC_SUB_REQ:                  "RIC SUBSCRIPTION REQUEST",
	RIC_SUB_RESP:                 "RIC SUBSCRIPTION RESPONSE",
	RIC_SUB_FAILURE:              "RIC SUBSCRIPTION FAILURE",
	RIC_SUB_DEL_REQ:              "RIC SUBSCRIPTION DELETE REQUEST",
	RIC_SUB_DEL_RESP:             "RIC SUBSCRIPTION DELETE RESPONSE",
	RIC_SUB_DEL_FAILURE:          "RIC SUBSCRIPTION DELETE FAILURE",
	RIC_SERVICE_UPDATE:           "RIC SERVICE UPDATE",
	RIC_SERVICE_UPDATE_ACK:       "RIC SERVICE UPDATE ACKNOWLEDGE",
	RIC_SERVICE_UPDATE_FAILURE:   "RIC SERVICE UPDATE FAILURE",
	RIC_CONTROL_REQ:              "RIC CONTROL REQUEST",
	RIC_CONTROL_ACK:              "RIC CONTROL ACKNOWLEDGE",
	RIC_CONTROL_FAILURE:          "RIC CONTROL FAILURE",
	RIC_INDICATION:               "RIC INDICATION",
	RIC_SERVICE_QUERY:            "RIC SERVICE QUERY",
	RIC_X2_SETUP_REQ:             "RIC X2 SETUP REQUEST",
	RIC_X2_SETUP_RESP:            "RIC X2 SETUP RESPONSE",
	RIC_X2_SETUP_FAILURE:         "RIC X2 SETUP FAILURE",
	RIC_X2_RESET:                 "RIC X2 RESET REQUEST",
	RIC_X2_RESET_RESP:            "RIC X2 RESET RESPONSE",
	RIC_ENDC_X2_SETUP_REQ:        "RIC EN-DC X2 SETUP REQUEST",
	RIC_ENDC_X2_SETUP_RESP:       "RIC EN-DC X2 SETUP RESPONSE",
	RIC_ENDC_X2_SETUP_FAILURE:    "RIC EN-DC X2 SETUP FAILURE",
	RIC_ENDC_CONF_UPDATE:         "RIC EN-DC CONFIGURATION UPDATE",
	RIC_ENDC_CONF_UPDATE_ACK:     "RIC EN-DC CONFIGURATION UPDATE ACKNOWLEDGE",
	RIC_ENDC_CONF_UPDATE_FAILURE: "RIC EN-DC CONFIGURATION UPDATE FAILURE",
	RIC_RES_STATUS_REQ:           "RIC RESOURCE STATUS REQUEST",
	RIC_RES_STATUS_RESP:          "RIC RESOURCE STATUS RESPONSE",
	RIC_RES_STATUS_FAILURE:       "RIC RESOURCE STATUS FAILURE",
	RIC_ENB_CONF_UPDATE:          "RIC ENB CONFIGURATION UPDATE",
	RIC_ENB_CONF_UPDATE_ACK:      "RIC ENB CONFIGURATION UPDATE ACKNOWLEDGE",
	RIC_ENB_CONF_UPDATE_FAILURE:  "RIC ENB CONFIGURATION UPDATE FAILURE",
	RIC_ENB_LOAD_INFORMATION:     "RIC ENB LOAD INFORMATION",
	RIC_GNB_STATUS_INDICATION:    "RIC GNB STATUS INDICATION",
	RIC_RESOURCE_STATUS_UPDATE:   "RIC RESOURCE STATUS UPDATE",
	RIC_ERROR_INDICATION:         "RIC ERROR INDICATION",
	DC_ADM_INT_CONTROL:           "DC ADMISSION INTERVAL CONTROL",
	DC_ADM_INT_CONTROL_ACK:       "DC ADMISSION INTERVAL CONTROL ACK",
}