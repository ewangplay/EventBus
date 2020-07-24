/*
 * Copyright arxanfintech.com. 2017 All Rights Reserved.
 *
 * Author:  ZhaoFang Han(frank@arxanfintech.com)
 *
 * Purpose: chain-mgmt service metadata info, like version
 *
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.txt', which is part of this source code package.
 */

package metadata

import (
	"fmt"
	"runtime"
)

const release = 1
const fixpack = 0
const hotfix = 0

// Version ...
var Version = fmt.Sprintf("%d.%d.%d", release, fixpack, hotfix)

func GetVersionInfo(progName string) string {
	return fmt.Sprintf("%s:\n Version: %s\n BuildNumber: %s\n Go version: %s\n OS/Arch: %s",
		progName, Version, BuildNumber, runtime.Version(),
		fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH))
}
