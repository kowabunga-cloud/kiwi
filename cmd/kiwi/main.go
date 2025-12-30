/*
 * Copyright (c) The Kowabunga Project
 * Apache License, Version 2.0 (see LICENSE or https://www.apache.org/licenses/LICENSE-2.0.txt)
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"os"

	"github.com/kowabunga-cloud/kiwi/internal/kiwi"
	"github.com/kowabunga-cloud/kowabunga/kowabunga/common/klog"
)

func main() {
	err := kiwi.Daemonize()
	if err != nil {
		klog.Error(err)
		os.Exit(1)
	}
}
