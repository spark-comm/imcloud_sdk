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
	APIADDR = "http://47.108.68.161:9099"
	WSADDR  = "ws://47.108.68.161:10001"
	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	UserID = "133374647734272"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxMzMzNzQ2NDc3MzQyNzJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxNDcxMzM4OTkzODkwNTA4OFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInRlbmFudElkXCI6XCI1MjI0NTAyMzIxNTIwNjRcIixcInNlcnZlcl9jb2RlXCI6XCJcIixcInJvbGVcIjpcIlVzZXJcIixcInNjb3BlXCI6XCJcIixcIm5vZGVJZFwiOlwiNTIyNDUwMjMyMTUyMDY0XCIsXCJvcHRpb25zXCI6bnVsbH0iLCJleHAiOjE3MTcxNTQ1MDMsIm5iZiI6MTcxNjc5NDUwMywiaWF0IjoxNzE2Nzk0NTAzfQ.4RqYyraSd8EJV4B2d-3GhI0YL9nxWJqnUWbuOJj-EQnHuQ_eX3S0ItkCU6xmxm1xfpeegqDIyQoeqmWa4SZ0hQ"
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
