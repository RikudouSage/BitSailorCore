package main

/*
#include "bw_common.h"
*/
import "C"

const (
	// BitwardenSuccess indicates that a C API call completed successfully.
	BitwardenSuccess C.BitwardenResult = iota
	// BitwardenError indicates that a C API call failed and last error is available.
	BitwardenError
)
