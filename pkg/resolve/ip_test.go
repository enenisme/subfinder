package resolve

import (
	"fmt"
	"testing"
)

func TestQueryIP(t *testing.T) {
	domain := "info.sungrow.cn"
	fmt.Println(QueryIP(domain))
}
