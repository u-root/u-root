// Copyright 2015 Google Inc. All Rights Reserved.
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

package timeutil

import (
	"errors"
	"fmt"
	"time"

	"github.com/jacobsa/oglematchers"
)

func timeEq(expected time.Time, c interface{}) error {
	actual, ok := c.(time.Time)
	if !ok {
		return errors.New("which is not a time")
	}

	// Make sure the times are the same instant.
	if diff := actual.Sub(expected); diff != 0 {
		return fmt.Errorf("which is off by %v", diff)
	}

	// Compare using == to capture other equality semantics; in particular
	// location.
	if actual != expected {
		return errors.New("")
	}

	return nil
}

// Return a matcher for times that are exactly equal to the given input time
// according to the == operator, which compares on location, instant, and
// monotonic clock reading.
//
// If you want to ignore location, canonicalize using Time.UTC. If you want to
// ignore ignore monotonic clock reading, strip it using Time.AddDate(0, 0, 0)
// (cf. https://goo.gl/rYU5UI).
func TimeEq(t time.Time) oglematchers.Matcher {
	return oglematchers.NewMatcher(
		func(c interface{}) error { return timeEq(t, c) },
		t.String())
}
