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
	"flag"
	"fmt"
	"github.com/imCloud/im/pkg/common/mcontext"
	"open_im_sdk/internal/file"
	"open_im_sdk/open_im_sdk"
	"path/filepath"
	"testing"
)

type FilePutCallback struct {
}

func (c *FilePutCallback) Open(size int64) {
	//TODO implement me
	fmt.Println("Open")
}

func (c *FilePutCallback) PartSize(partSize int64, num int) {
	//TODO implement me
	fmt.Println("PartSize")
}

func (c *FilePutCallback) HashPartProgress(index int, size int64, partHash string) {
	//TODO implement me
	fmt.Println("HashPartProgress")
}

func (c *FilePutCallback) HashPartComplete(partsHash string, fileHash string) {
	//TODO implement me
	fmt.Println("HashPartComplete")
}

func (c *FilePutCallback) UploadID(uploadID string) {
	//TODO implement me
	fmt.Println("UploadID")
}

func (c *FilePutCallback) UploadPartComplete(index int, partSize int64, partHash string) {
	//TODO implement me
	fmt.Println("UploadPartComplete")
}

func (c *FilePutCallback) UploadComplete(fileSize int64, streamSize int64, storageSize int64) {
	//TODO implement me
	fmt.Println("UploadComplete")
}

func (c *FilePutCallback) Complete(size int64, url string, typ int) {
	//TODO implement me
	fmt.Println("Complete")
}

func TestPut(t *testing.T) {
	ctx := mcontext.NewCtx("123456")
	req := &file.UploadFileReq{
		Name:     "icon.png",
		Filepath: "/Users/tang/workspace/icon.png",
	}
	req.Name = filepath.Base(req.Filepath)
	callback := FilePutCallback{}
	str, err := open_im_sdk.UserForSDK.File().UploadFileFullPath(ctx, req, &callback)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("url", str)
}

func Test_Fmt(t *testing.T) {
	i := flag.Int("sn", 2, "sender num")
	fmt.Println(i)
}
