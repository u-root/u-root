// Copyright 2018 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

// Your service.go file will be heavily tailored for the problem you are trying to solve.
// a type representing your service and a constructor for that service are all that is
// specifically required. If you are depending on something outside of the program, such
// as text files or system settings, add an update function which loads that data into your
// service struct.

// the one field in our example service is a string, set by the user.
type ExampleService struct {
	Example string
}

// here is where you'd define an update function. Since this demo does not depend on any
// system settings, I have omitted it.

// define any service functionality here.
// ExampleServiceFunc1: Set service variable to input
func (es *ExampleService) ExampleServiceFunc1(input ExampleJsonMsg) error {
	es.Example = input.Example
	return nil
}

// ExampleServiceFunc2: Reset service variable to empty string
func (es *ExampleService) ExampleServiceFunc2() error {
	es.Example = ""
	return nil
}

// NewExampleService: construct a new service
// If your service controls system settings, make a getCurrentSetting() function to
// initialize your service.
func NewExampleService() (*ExampleService, error) {
	return &ExampleService{
		Example: "",
	}, nil
}
