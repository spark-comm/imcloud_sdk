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
	//UserID       = "2688118337"
	//UserID       = "7204255074"
	//UserID = "50122626445611008"
	//UserID = "55122331994951680"
	////friendUserID = "3281432310"
	//// APIADDR = "http://192.168.44.128:10002"
	//// WSADDR  = "ws://192.168.44.128:10001"
	//// UserID  = "100"
	//
	////APIADDR = "http://59.36.173.89:10002"
	////WSADDR  = "ws://59.36.173.89:10001"
	////UserID  = "kernaltestuid9"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjMzMTk5NDk1MTY4MFwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwicm9sZVwiOlwiVVNFUlwifSIsImV4cCI6MTY5NjA3NDM0NSwibmJmIjoxNjk1NzE0MzQ1LCJpYXQiOjE2OTU3MTQzNDV9.4khO81UwFgN4rOX11N3Iy5mi7VT90hBFXW0f2CTqtJd-qjCrLnHtkTlzeTOQaoYUuWUL9BdcjcgQdcoJ2_Gnqg"
	//UserID = "55122365234810880"
	//UserID = "55122365549383680"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjM2NTU0OTM4MzY4MFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE2OTkxNjQ4MjIsIm5iZiI6MTY5ODgwNDgyMiwiaWF0IjoxNjk4ODA0ODIyfQ.o46dTaGGurutExkXWwLN8ChfduXfNgCmZ77DGOGiLQHgQN3_TPylWh4IutVX2StTeI4NUmyDKh_qKhUtoNJdYg"
	//token = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjM2NTIzNDgxMDg4MFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE2OTkxNjM5NTYsIm5iZiI6MTY5ODgwMzk1NiwiaWF0IjoxNjk4ODAzOTU2fQ.1Nv9ph46oY8IIAa0GBjIKTe9QwmyDJXmyelbLp5Yq1YWWoxej5aRwmcE_S4kcka18e9qQz31eJdk7kGHgdzmzg"
	//UserID = "55122365549383680"
	//token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjM2NTU0OTM4MzY4MFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlVTRVJcIn0iLCJleHAiOjE2OTk1OTg1NzAsIm5iZiI6MTY5OTIzODU3MCwiaWF0IjoxNjk5MjM4NTcwfQ.nVJfjGnXfIhIR-I0lnhPXiyhELdaCjzXsXmF8MXPPQtO0v5raTTgHRvup1oSKEgMxkqBvn-PGRi4OJACDjMhUw"
	UserID = "55122331994951680"
	token  = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI1NTEyMjMzMTk5NDk1MTY4MFwiLFwicGxhdGZvcm1cIjpcIklPU1wiLFwidG9rZW5cIjpcIlwiLFwidGVuYW50SWRcIjpcIlwiLFwiZXhwaXJlX3RpbWVfc2Vjb25kc1wiOjAsXCJyb2xlXCI6XCJVU0VSXCJ9IiwiZXhwIjoxNzAwODEwMTM3LCJuYmYiOjE3MDA0NTAxMzcsImlhdCI6MTcwMDQ1MDEzN30.fU3GGWuyMpn0nvYf4AVJ9JQ-lO4nZ6-9-Woe5aVl2WYNGiY5auU5xjyTncm8B-YuAK_rKbjPVOXHfKZVpUqQgg"
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
