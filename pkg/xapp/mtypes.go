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
