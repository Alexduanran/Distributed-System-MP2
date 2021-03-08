package tcp

import (
	"fmt"
	"os"
)

// checkError checks for error in err, exits
// and prints given error message errMSG to the console if err is not nil
func checkError(err error, errMsg string) {
	if err != nil {
		if errMsg == "" {
			errMsg = "Fatal error"
		}
		fmt.Println(errMsg, ":", err.Error())
		os.Exit(1)
	}
}