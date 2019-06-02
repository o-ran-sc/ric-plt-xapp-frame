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
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -lmdclog
#
#include <mdclog/mdclog.h>
void xAppMgr_mdclog_write(mdclog_severity_t severity, const char *msg) {
     mdclog_write(severity, "%s", msg);
}
*/
import "C"

import (
	"fmt"
	"log"
	"time"
)

type Log struct {
}

const (
	LogLvlErr   = C.MDCLOG_ERR
	LogLvlWarn  = C.MDCLOG_WARN
	LogLvlInfo  = C.MDCLOG_INFO
	LogLvlDebug = C.MDCLOG_DEBUG
)

func WriteLog(lvl C.mdclog_severity_t, msg string) {
	t := time.Now().Format("2019-01-02 15:04:05")
	text := fmt.Sprintf("%s:: %s ", t, msg)

	C.xAppMgr_mdclog_write(lvl, C.CString(text))
}

func (Log) SetLevel(level int) {
	l := C.mdclog_severity_t(level)
	C.mdclog_level_set(l)
}

func (Log) SetMdc(key string, value string) {
	C.mdclog_mdc_add(C.CString(key), C.CString(value))
}

func (Log) Fatal(pattern string, args ...interface{}) {
	WriteLog(LogLvlErr, fmt.Sprintf(pattern, args...))
	log.Panic("Fatal error occured, exiting ...")
}

func (Log) Error(pattern string, args ...interface{}) {
	WriteLog(LogLvlErr, fmt.Sprintf(pattern, args...))
}

func (Log) Warn(pattern string, args ...interface{}) {
	WriteLog(LogLvlWarn, fmt.Sprintf(pattern, args...))
}

func (Log) Info(pattern string, args ...interface{}) {
	WriteLog(LogLvlInfo, fmt.Sprintf(pattern, args...))
}

func (Log) Debug(pattern string, args ...interface{}) {
	WriteLog(LogLvlDebug, fmt.Sprintf(pattern, args...))
}
