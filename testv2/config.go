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
	UserID  = "5324390178754560"
	token   = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1MzI0MzkwMTc4NzU0NTYwXCIsXCJjZW50ZXJfdXNlcl9pZFwiOlwiNTMyNDM5MDE3ODc1NDU2MFwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwidGVuYW50SWRcIjpcIjQxNjYzNjc4ODAxOTIwMFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI0MTY2MzY3ODgwMTkyMDBcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxMzY5NjgxMCwibmJmIjoxNzEzMzM2ODEwLCJpYXQiOjE3MTMzMzY4MTB9.Pag6jjPQMCEM-DD2qb47dZhEuNhXKUe8SRknlyDmVZq9RBNoAOjFgFLO93-YKFJZtp-Fk7mjHr2uRmlqFfDdEA"
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
