// Copyright © 2017 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"os"

	"github.com/ciao-project/ciao/ciao-sdk"
	"github.com/ciao-project/ciao/ciao/cmd"

	"github.com/golang/glog"

)

var (
    tenantID       = new(string)
    controllerURL  = new(string)
    ciaoPort       = new(int)
    caCertFile     = new(string)
    clientCertFile = new(string)
)

const (
	ciaoControllerEnv     = "CIAO_CONTROLLER"
	ciaoCACertFileEnv     = "CIAO_CA_CERT_FILE"
	ciaoClientCertFileEnv = "CIAO_CLIENT_CERT_FILE"
)

func infof(format string, args ...interface{}) {
	if glog.V(1) {
		glog.InfoDepth(1, fmt.Sprintf("ciao-cli INFO: "+format, args...))
	}
}

func errorf(format string, args ...interface{}) {
	glog.ErrorDepth(1, fmt.Sprintf("ciao-cli ERROR: "+format, args...))
}

func fatalf(format string, args ...interface{}) {
	glog.FatalDepth(1, fmt.Sprintf("ciao-cli FATAL: "+format, args...))
}

func getCiaoEnvVariables() {
	controller := os.Getenv(ciaoControllerEnv)
	ca := os.Getenv(ciaoCACertFileEnv)
	clientCert := os.Getenv(ciaoClientCertFileEnv)

	infof("Ciao environment variables:\n")
	infof("\t%s:%s\n", ciaoControllerEnv, controller)
	infof("\t%s:%s\n", ciaoCACertFileEnv, ca)
	infof("\t%s:%s\n", ciaoClientCertFileEnv, clientCert)

	sdk.C.ControllerURL = controller
	sdk.C.CACertFile = ca
	sdk.C.ClientCertFile = clientCert

	if *controllerURL != "" {
		sdk.C.ControllerURL = *controllerURL
	}

	if *caCertFile != "" {
		sdk.C.CACertFile = *caCertFile
	}

	if *clientCertFile != "" {
		sdk.C.ClientCertFile = *clientCertFile
	}

	if *tenantID != "" {
		sdk.C.TenantID = *tenantID
	}
}

func main() {
	getCiaoEnvVariables()
	sdk.C.Init()
	cmd.Execute()
}
