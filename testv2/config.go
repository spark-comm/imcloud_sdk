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
	//APIADDR = "http://47.108.68.161:9099"
	//WSADDR  = "ws://47.108.68.161:10001"
	//UserID  = "kernaltestuid2"

	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	//预生产
	APIADDR = "http://8.137.13.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	UserID  = "405431151235072"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI0MDU0MzExNTEyMzUwNzJcIixcImNlbnRlcl91c2VyX2lkXCI6XCI1MzE2OTg4MTkyNzI3MDRcIixcInBsYXRmb3JtXCI6XCJJT1NcIixcInRlbmFudElkXCI6XCIxMjIyOTgxNzg3MzYxMjhcIixcInNlcnZlcl9jb2RlXCI6XCJcIixcInJvbGVcIjpcIlVzZXJcIixcInNjb3BlXCI6XCJcIixcIm5vZGVJZFwiOlwiMTIyMjk4MTc4NzM2MTI4XCIsXCJvcHRpb25zXCI6bnVsbH0iLCJleHAiOjE3MDYyNjU0NDksIm5iZiI6MTcwNTkwNTQ0OSwiaWF0IjoxNzA1OTA1NDQ5fQ.dsP-7vpr-lWO2lavU5yU6YxGfpG1pS1uRw8r0LgKbsbFQF3IvpQCG5F5h-Pi4bgLWHesqi9fXG3q_dHGFBry1Q"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = 1
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
