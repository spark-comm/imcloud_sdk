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
	UserID = "14743920172863488"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjM2NTY2NjgyNDE5MlwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE2OTk0NDI4MzksIm5iZiI6MTY5ODcyMjgzOSwiaWF0IjoxNjk4NzIyODM5fQ.iGmBGdYtMI1E4Tq6wKjZTczhVYqpxQOLaaVT2XbyEnUrs_6rRfan3lURXKaXBOkww4gE4Sk6QyFf19DEr99cTw"
	token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxNDc0MzkyMDE3Mjg2MzQ4OFwiLFwiY2VudGVyX3VzZXJfaWRcIjpcIjE0NzEyOTIyMTU2NTY4NTc2XCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiNTIyNDUwMjMyMTUyMDY0XCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjUyMjQ1MDIzMjE1MjA2NFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzE1ODY3NDM3LCJuYmYiOjE3MTU1MDc0MzcsImlhdCI6MTcxNTUwNzQzN30.Q7Mjg38bq-3R7shqP_nCm5w4K7Kf0d69DqqCrjz6v1cmxqEY5xoqWNDIglLJfmR34dkFXTo36iAeIJptvDvAeg"
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
