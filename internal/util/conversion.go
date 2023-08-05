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

package util

func Batch[T any, V any](fn func(T) V, ts []T) []V {
	if ts == nil {
		return nil
	}
	res := make([]V, 0, len(ts))
	for i := range ts {
		res = append(res, fn(ts[i]))
	}
	return res
}
func BatchConversion[T any, V any](fn func(T) (V, error), ts []T) ([]V, error) {
	if ts == nil {
		return make([]V, 0), nil
	}
	res := make([]V, 0, len(ts))
	for i := range ts {
		v, err := fn(ts[i])
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
