// Copyright © 2023 OpenIM SDK. All rights reserved.
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

package testv2

import "open_im_sdk/sdk_struct"

const (
	//APIADDR = "http://43.154.157.177:10002"
	//WSADDR  = "ws://43.154.157.177:10001"
	//UserID  = "kernaltestuid2"

	//APIADDR      = "http://59.36.173.89:10002"
	APIADDR = "http://8.137.13.1:9099"
	//WSADDR       = "ws://59.36.173.89:10001"
	WSADDR = "ws://8.137.13.1:10001"
	//UserID       = "2688118337"
	//UserID       = "7204255074"
	UserID       = "45778745637736448"
	friendUserID = "3281432310"
	// APIADDR = "http://192.168.44.128:10002"
	// WSADDR  = "ws://192.168.44.128:10001"
	// UserID  = "100"

	//APIADDR = "http://59.36.173.89:10002"
	//WSADDR  = "ws://59.36.173.89:10001"
	//UserID  = "kernaltestuid9"
	token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcl9pZFwiOlwiNDU3Nzg3NDU2Mzc3MzY0NDhcIixcInBsYXRmb3JtXCI6XCJXaW5kb3dzXCIsXCJyb2xlXCI6XCJcIn0iLCJleHAiOjE2OTA4MTE3ODYsIm5iZiI6MTY5MDQ1MTc4NiwiaWF0IjoxNjkwNDUxNzg2fQ.DgoKXjrQuJZFkWtnkItBED-nfiC3OxVSH9kVrbR09s4kOql4FGsPZuNgYY8oKSpuZpCTgwJrw87--qwmm3u7dw"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "../"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = 1
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
