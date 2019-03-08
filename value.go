// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package lookslike

// ValueResult represents the result of checking a leaf value.
type ValueResult struct {
	Valid   bool
	Message string // Reason this is invalid
}

// A ValueValidator is used to validate a value in a Map.
type ValueValidator func(path Path, v interface{}) *Results

// An IsDef defines the type of Check to do.
// Generally only Name and Checker are set. Optional and CheckKeyMissing are
// needed for weird checks like key presence.
type IsDef struct {
	Name            string
	Checker         ValueValidator
	Optional        bool
	CheckKeyMissing bool
}

// Check runs the IsDef at the given value at the given path
func (id IsDef) Check(path Path, v interface{}, keyExists bool) *Results {
	if id.CheckKeyMissing {
		if !keyExists {
			return ValidResult(path)
		}

		return SimpleResult(path, false, "this key should not exist")
	}

	if !id.Optional && !keyExists {
		return KeyMissingResult(path)
	}

	if id.Checker != nil {
		return id.Checker(path, v)
	}

	return ValidResult(path)
}

// ValidResult is a convenience value for Valid results.
func ValidResult(path Path) *Results {
	return SimpleResult(path, true, "is valid")
}

// ValidVR is a convenience value for Valid results.
var ValidVR = ValueResult{true, "is valid"}

// KeyMissingResult is emitted when a key was expected, but was not present.
func KeyMissingResult(path Path) *Results {
	return SingleResult(path, KeyMissingVR)
}

// KeyMissingVR is emitted when a key was expected, but was not present.
var KeyMissingVR = ValueResult{
	false,
	"expected this key to be present",
}

// StrictFailureResult is emitted when Strict() is used, and an unexpected field is found.
func StrictFailureResult(path Path) *Results {
	return SingleResult(path, StrictFailureVR)
}

// StrictFailureVR is emitted when Strict() is used, and an unexpected field is found.
var StrictFailureVR = ValueResult{
	false,
	"unexpected field encountered during strict validation",
}
