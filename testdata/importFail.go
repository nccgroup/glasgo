package main

import(
	"doesNotExist"
	"crypto/md5"
)

func ImportFail() {
	doesNotExist.Stuff();
	/* this check might seem redundant because
	 * there is already crypto check
	 * so the point is to make sure tests are still run
	 * even after the import fails.
	*/
	h := md5.New();
        if h != nil {
                return;
        }
	return
}

