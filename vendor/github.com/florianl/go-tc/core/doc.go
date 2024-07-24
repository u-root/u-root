/*
Package core contains some generic helper functions for the package github.com/florianl/go-tc.
Depending on the actual hardware in use, parameters to filters and qdiscs can vary to achieve the
same behaviour.

If PROC_ROOT is set it will be used to lookup packet scheduler information from net/psched.
*/
package core
