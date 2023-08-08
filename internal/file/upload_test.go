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

package file

import (
	"context"
	"fmt"
	"open_im_sdk/pkg/ccontext"
	"open_im_sdk/pkg/db"
	"open_im_sdk/sdk_struct"
	"path/filepath"
	"testing"
)

func TestName(t *testing.T) {
	userID := `49389272901357568`
	ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
		UserID: userID,
		Token:  `eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI0OTM4OTI3MjkwMTM1NzU2OFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlwifSIsImV4cCI6MTY5MTgyNzE0NCwibmJmIjoxNjkxNDY3MTQ0LCJpYXQiOjE2OTE0NjcxNDR9._EQqOgC0JblO-CE9oBETgE62SbEt55bQCx__bumoPAAV8U3iegTkaLziTD-czjHdGoraN5zujdVpxRkquDkqpA`,
		IMConfig: sdk_struct.IMConfig{
			ApiAddr: "http://8.137.13.1:9099",
		},
	})
	//userID := `49383675594280960`
	//ctx := ccontext.WithInfo(context.Background(), &ccontext.GlobalConfig{
	//	UserID: userID,
	//	Token:  `eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJpbUNsb3VkIiwic3ViIjoie1widXNlcklkXCI6XCI0OTM4MzY3NTU5NDI4MDk2MFwiLFwicGxhdGZvcm1cIjpcIldpbmRvd3NcIixcInJvbGVcIjpcIlwifSIsImV4cCI6MTY5MTc1MjYzNiwibmJmIjoxNjkxMzkyNjM2LCJpYXQiOjE2OTEzOTI2MzZ9.FyBCcmzThoD9RFZCu9fwr_W8pvQ5VhMa-lisKOo5fv_gNqhGAFtpK_DnmOBgWC47JVIRvu2d6mksytGS-LykFw`,
	//	IMConfig: sdk_struct.IMConfig{
	//		ApiAddr: "http://localhost:9099",
	//	},
	//})
	ctx = ccontext.WithOperationID(ctx, `test`)

	database, err := db.NewDataBase(ctx, userID, `/Users/tang/workspace/go/openim/openim-sdk-core`)
	if err != nil {
		panic(err)
	}
	f := NewFile(database, userID)

	path := `/Users/tang/workspace`
	path = filepath.Join(path, `icon.png`)
	resp, err := f.UploadFile(ctx, &UploadFileReq{
		Filepath: path,
		Name:     filepath.Base(path),
		Cause:    "test",
	}, nil)
	if err != nil {
		t.Logf("%+v\n", err)
		return
	}
	t.Logf("%+v\n", resp)
}

func TestName1(t *testing.T) {
	p := make([]byte, 10)

	a := []byte("12345")

	copy(p, a)

	fmt.Println(p)

}
