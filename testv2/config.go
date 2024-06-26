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
	APIADDR = "http://43.243.73.185:9099"
	WSADDR  = "ws://43.243.73.185:10001"
	//APIADDR = "http://127.0.0.1:9099"
	//WSADDR  = "ws://127.0.0.1:10001"
	//UserID = "14743920172863488"
	//token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxNDc0MzkyMDE3Mjg2MzQ4OFwiLFwiY2VudGVyX3VzZXJfaWRcIjpcIjE0NzEyOTIyMTU2NTY4NTc2XCIsXCJwbGF0Zm9ybVwiOlwiV2luZG93c1wiLFwidGVuYW50SWRcIjpcIjUyMjQ1MDIzMjE1MjA2NFwiLFwic2VydmVyX2NvZGVcIjpcIlwiLFwicm9sZVwiOlwiVXNlclwiLFwic2NvcGVcIjpcIlwiLFwibm9kZUlkXCI6XCI1MjI0NTAyMzIxNTIwNjRcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxNjg3NTAzNSwibmJmIjoxNzE2NTE1MDM1LCJpYXQiOjE3MTY1MTUwMzV9.RLQvUupF-7xRdzL73xG8MVwwrHCv4Ywo90HUsaoOdRdXDS9B2D6lhCap85I61pew9UtnVMbslc5-Xmrhhi7yWQ"
	UserID = "1873637812473856"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCIxODczNjM3ODEyNDczODU2XCIsXCJjZW50ZXJfdXNlcl9pZFwiOlwiMTg3MzYzNzgxMjQ3Mzg1NlwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInRlbmFudElkXCI6XCIxODAwNzE5ODE2NDYyMzM2XCIsXCJzZXJ2ZXJfY29kZVwiOlwiXCIsXCJyb2xlXCI6XCJVc2VyXCIsXCJzY29wZVwiOlwiXCIsXCJub2RlSWRcIjpcIjE4MDA3MTk4MTY0NjIzMzZcIixcIm9wdGlvbnNcIjpudWxsfSIsImV4cCI6MTcxOTM1Mjk5NiwibmJmIjoxNzE4OTkyOTk2LCJpYXQiOjE3MTg5OTI5OTZ9.4ZTpjHkvjt1TL1rHHu2dWNt2DD6LRZwNaoZqg7dz_0nG12Z1ryGg-fxNHoH3oCekfgOVFfqSldcfjA9-Nwmi_A"
)

func getConf(APIADDR, WSADDR string) sdk_struct.IMConfig {
	var cf sdk_struct.IMConfig
	cf.ApiAddr = APIADDR
	cf.WsAddr = WSADDR
	cf.DataDir = "./"
	cf.LogLevel = 3
	cf.IsExternalExtensions = true
	cf.PlatformID = 1
	cf.LogFilePath = ""
	cf.IsLogStandardOutput = true
	return cf
}
