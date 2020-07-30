package work

import (
	"fmt"
	"time"
)

// Work processes the migration.
func Work(i int) {
	time.Sleep(time.Second)
	fmt.Println(i)
}
