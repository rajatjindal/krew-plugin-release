package actions

import (
	"fmt"
	"os"
	"strings"
)

//GetInputForAction gets input to action
func GetInputForAction(key string) string {
	return os.Getenv(fmt.Sprintf("INPUT_%s", strings.ToUpper(key)))
}
