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

package syncutil

import (
	"sync"

	"golang.org/x/net/context"
)

// A collection of concurrently-executing operations, each of which may fail.
//
// Operations are defined by functions that accept a context and return an
// error. If any operation returns a non-nil error, all concurrent and future
// operations will have their contexts cancelled . When Join() is called on a
// bundle with one or more operations that returned an error, it always returns
// the first error (i.e. that which led to the cancellation of others).
//
// Bundles can be used to set up pipelines of concurrent actors sending data to
// each other, conveniently cancelling the pipeline if anything fails. A
// typical use looks like the following:
//
//     // Run a pipeline that consists of one goroutine listing object names,
//     // while N goroutines concurrently delete the listed objects one by one.
//     // If any listing or deletion operation fails, cancel the whole pipeline
//     // and return the error.
//     func deleteAllObjects(ctx context.Context, N int) error {
//       bundle := syncutil.NewBundle(ctx)
//
//       // List objects into a channel. Assuming that listObjects responds to
//       // cancellation of its context, it will not get stuck blocking forever
//       // on a write into objectNames if the deleters return early in error
//       // before draining the channel.
//       objectNames := make(chan string)
//       bundle.Add(func(ctx context.Context) error {
//         defer close(objectNames)
//         return listObjects(ctx, objectNames)
//       })
//
//       // Run N deletion workers.
//       for i := 0; i < N; i++ {
//         bundle.Add(func(ctx context.Context) error {
//           for name := range objectNames {
//             if err := deleteObject(ctx, name); err != nil {
//               return err
//             }
//           }
//         })
//       }
//
//       // Wait for the whole pipeline to finish, and return its status.
//       return bundle.Join()
//    }
//
// Bundles must be created using NewBundle. Bundle methods must not be called
// concurrently.
type Bundle struct {
	context context.Context
	cancel  context.CancelFunc

	waitGroup sync.WaitGroup

	errorOnce  sync.Once
	firstError error
}

// Add a new operation to the bundle. The operation will be invoked with a
// context that will be cancelled if any other operation fails or has already
// failed.
func (b *Bundle) Add(f func(context.Context) error) {
	b.waitGroup.Add(1)

	// Run the function in the background.
	go func() {
		defer b.waitGroup.Done()

		err := f(b.context)
		if err == nil {
			return
		}

		// On first error, cancel the context and save the error.
		b.errorOnce.Do(func() {
			b.firstError = err
			b.cancel()
		})
	}()
}

// Wait for all previously-added operations to complete. Return nil if all
// operations succeeded. Otherwise return the first error.
//
// Add must not be called concurrently with or after Join.
func (b *Bundle) Join() error {
	b.waitGroup.Wait()

	// context.WithCancel requires that we arrange for this to be called
	// eventually in order to avoid leaking resources. Since everything is done,
	// to do so now is harmless.
	b.cancel()

	return b.firstError
}

// Create a bundle whose operations are fed a context inheriting from the given
// parent context, which must be non-nil. The bundle must eventually be joined
// with Join.
func NewBundle(parent context.Context) *Bundle {
	b := &Bundle{}
	b.context, b.cancel = context.WithCancel(parent)

	return b
}
