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

import (
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/spark-comm/imcloud_sdk/sdk_struct"
)

const (
	APIADDR = "http://8.137.13.1:9099"
	WSADDR  = "ws://8.137.13.1:10001"
	//UserID  = "931422227402752"
	//token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MzE0MjIyMjc0MDI3NTJcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM3ODEyNDczODU2XCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIwOTQ2MDkwLCJuYmYiOjE3MjA1ODYwOTAsImlhdCI6MTcyMDU4NjA5MH0.52O0J5dMlXqIu56BiR1dAK4CBipk2ddF9hUy6Y7U0XuJilz60EdYZaeabHYtgdQOg6Pl91Yguq_QHVRJGSkqzg"
	UserID = "922670631751680"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI5MjI2NzA2MzE3NTE2ODBcIixcImNlbnRlcl91c2VyX2lkXCI6XCIxODczNjM2NjEyOTAyOTEyXCIsXCJwbGF0Zm9ybVwiOlwiSU9TXCIsXCJ0ZW5hbnRJZFwiOlwiOTExMzU1NzYyNzA4NDgwXCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjkxMTM1NTc2MjcwODQ4MFwiLFwib3B0aW9uc1wiOm51bGx9IiwiZXhwIjoxNzIzNDcxMzU0LCJuYmYiOjE3MjA4NzkzNTQsImlhdCI6MTcyMDg3OTM1NH0.32IEp5vo0Dxpu4SksPsn44I8FM99k7ydOyrr1yWzTV7LoEkOPHi8f4RsiVNDeglsAe4vEdVTKq2jipJWRAvXTg"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "../"
	cf.LogLevel = 6
	cf.IsExternalExtensions = true
	cf.PlatformID = constant.LinuxPlatformID
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
