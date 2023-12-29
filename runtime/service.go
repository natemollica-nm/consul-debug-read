// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: BUSL-1.1

package service_os

var chanGraceExit = make(chan int)

func ShutdownChannel() <-chan int {
	return chanGraceExit
}
