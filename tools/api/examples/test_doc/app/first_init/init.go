package first_init

import (
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/app/constant"
	"github.com/haozzzzzzzz/go-rapid-development/tools/api/examples/test_doc/common/dependent"
)

func init() {
	dependent.ServiceName = constant.ServiceName
}
