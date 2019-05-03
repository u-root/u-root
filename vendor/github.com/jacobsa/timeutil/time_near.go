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
	"math"
	"time"

	"github.com/jacobsa/oglematchers"
)

func timeNear(t time.Time, d time.Duration, c interface{}) error {
	actual, ok := c.(time.Time)
	if !ok {
		return errors.New("which is not a time")
	}

	absDiff := time.Duration(math.Abs(float64(actual.Sub(t))))
	if absDiff >= d {
		return fmt.Errorf("which differs by %v", absDiff)
	}

	return nil
}

// Return a matcher for times whose absolute distance from t is less than d.
func TimeNear(t time.Time, d time.Duration) oglematchers.Matcher {
	return oglematchers.NewMatcher(
		func(c interface{}) error { return timeNear(t, d, c) },
		fmt.Sprintf("within %v of %v", d, t))
}
