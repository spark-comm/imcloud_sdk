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
	APIADDR = "http://8.137.13.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	UserID  = "55122332682817536"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjMzMjY4MjgxNzUzNlwiLFwicGxhdGZvcm1cIjpcIkFuZHJvaWRcIixcInRlbmFudElkXCI6XCJcIixcInNlcnZlcl9jb2RlXCI6XCJcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE3MDM1MTYzMzEsIm5iZiI6MTcwMzE1NjMzMSwiaWF0IjoxNzAzMTU2MzMxfQ.gMz4O2ejMJ6oBxXUyHb6tB4v-oxL4N5yhK3AWgpujsbK1jP_y-8kSH2E_0HLoo2_xY7pJp2C67I6LixOYG8mPw"
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
