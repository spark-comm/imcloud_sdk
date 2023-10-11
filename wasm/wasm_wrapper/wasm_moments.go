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

//go:build js && wasm
// +build js,wasm

package wasm_wrapper

import (
	"open_im_sdk/open_im_sdk"
	"open_im_sdk/pkg/utils"
	"open_im_sdk/wasm/event_listener"
	"syscall/js"
)

// ------------------------------------group---------------------------
type WrapperMoments struct {
	*WrapperCommon
}

func NewWrapperMoments(wrapperCommon *WrapperCommon) *WrapperGroup {
	return &WrapperGroup{WrapperCommon: wrapperCommon}
}

func (w *WrapperGroup) GetMomentsList(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.GetMomentsList, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) PublishMoments(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.PublishMoments, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) DeleteMoments(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.DeleteMoments, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) CommentMoments(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.CommentMoments, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) LikeMoments(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.LikeMoments, callback, &args).AsyncCallWithCallback()
}

func (w *WrapperGroup) UnLikeMoments(_ js.Value, args []js.Value) interface{} {
	callback := event_listener.NewBaseCallback(utils.FirstLower(utils.GetSelfFuncName()), w.commonFunc)
	return event_listener.NewCaller(open_im_sdk.UnLikeMoments, callback, &args).AsyncCallWithCallback()
}
