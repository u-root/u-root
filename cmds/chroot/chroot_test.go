package main

import (
  "fmt"
  "testing"
)

func TestUserSpecSet(t *testing.T) {
  var (
    user userSpec
    input string
    err error
  )
  input = "1000:1001"

  err = user.Set(input)
  if err != nil{
    t.Errorf("Unexpected error with input: %s", input)
  } else if user.uid != 1000 || user.gid != 1001 {
    t.Errorf("Expected uid 1000, gid 1001, got: uid %d, gid %d", user.uid, user.gid)
  }

  input = "test:1000"
  err = user.Set(input)
  if err == nil {
    t.Errorf("Expected error, input %s got: uid %d, gid %d", input, user.uid, user.gid)
  }

  input = "1000:1001:"
  err = user.Set(input)
  if err == nil {
    t.Errorf("Expected error, input %s got: uid %d, gid %d", input, user.uid, user.gid)
  }

  input = ":1000"
  err = user.Set(input)
  if err == nil {
    t.Errorf("Expected error, input %s got: uid %d, gid %d", input, user.uid, user.gid)
  }

  input = "1000:"
  err = user.Set(input)
  if err == nil {
    t.Errorf("Expected error, input %s got: uid %d, gid %d", input, user.uid, user.gid)
  }

}

func TestUserSpecString(t *testing.T) {
  var (
    user userSpec
    input string
    err error
  )
  input = "1000:1001"

  err = user.Set(input)
  if err != nil {
    t.Errorf("Unexpected error with input: %s", input)
  }

  str := user.String()
  if str != input {
    t.Errorf("Unexpected error with input: %s, String method returned: %s", input, str)
  }

}

func testGroupSet(input string, expected []uint32) error {
  var groups groupsSpec
  err := groups.Set(input)
  if err != nil {
    return fmt.Errorf("Unexpected error with input: %s, %s", input, err)
  } else if len(groups.groups) != len(expected) {
    return fmt.Errorf("Unexpected groups length with input %s, actual length %d", input, len(groups.groups))
  } else {
    for index, group := range groups.groups {
      if expected[index] != group {
        return fmt.Errorf("Unexpected error at index %d, was expecting %d , found %d", index, expected[index], group)
        }
    }
  }
  return nil
}

func TestGroupsSet(t *testing.T) {
  var (
    input string
    expected []uint32
    err error
  )

  input = "1000"
  expected = []uint32{1000}
  if err = testGroupSet(input, expected); err != nil {
    t.Errorf(err.Error())
  }

  input = "1000,1001"
  expected = []uint32{1000, 1001}
  if err = testGroupSet(input, expected); err != nil {
    t.Errorf(err.Error())
  }

  input = "1000,1001,"
  expected = []uint32{}
  if err = testGroupSet(input, expected); err == nil {
    t.Errorf("Expected error on input: %s, got: %v", input, groups.groups)
  }

  input = "test,1001"
  expected = []uint32{}
  if err = testGroupSet(input, expected); err == nil {
    t.Errorf("Expected error on input: %s, got: %v", input, groups.groups)
  }

  input = ",1000"
  expected = []uint32{}
  if err = testGroupSet(input, expected); err == nil {
    t.Errorf("Expected error on input: %s, got: %v", input, groups.groups)
  }

  input = "1000,"
  expected = []uint32{}
  if err = testGroupSet(input, expected); err == nil {
    t.Errorf("Expected error on input: %s, got: %v", input, groups.groups)
  }

}
