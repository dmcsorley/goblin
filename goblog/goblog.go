package goblog

import (
	"fmt"
	"os"
	"time"
)

func RawLog(message string) {
	os.Stdout.WriteString(
		fmt.Sprintf("%s %s\n",
			time.Now().Format(time.RFC3339),
			message,
		),
	)
}

func Log(prefix string, message string) {
	os.Stdout.WriteString(
		fmt.Sprintf("%s %s %s\n",
			time.Now().Format(time.RFC3339),
			prefix,
			message,
		),
	)
}
