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

	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	//预生产
	APIADDR = "http://47.108.68.161:9099"
	WSADDR  = "ws://47.108.68.161:10001"
	UserID  = "100583720527859712"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxMDA1ODM3MjA1Mjc4NTk3MTJcIixcInBsYXRmb3JtXCI6XCJBbmRyb2lkXCIsXCJ0ZW5hbnRJZFwiOlwiMTIxNjExODE0NDQwOTYwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVU0VSXCJ9IiwiZXhwIjoxNzAzMzUyMzk3LCJuYmYiOjE3MDI5OTIzOTcsImlhdCI6MTcwMjk5MjM5N30.g4j8-e4LSqdgSU83JgwMo3c72vaYf_seDDoZ6_7diK4VFkHa5Ri64NmWLj1OGzYHLGB69vyXF8zgorDaXFMbvw"
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
