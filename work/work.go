package work

import (
	"fmt"
	"time"

	"github.com/magiconair/properties"
)

// Work processes the migration.
func Work(i int) {
	time.Sleep(time.Second)
	fmt.Println(i)
	p := properties.MustLoadFile("connection.properties", properties.UTF8)
	x, _ := p.Get("x")
	fmt.Println(x)
}
