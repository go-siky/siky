package util

import (
	"testing"

	"github.com/bmizerany/assert"
	"github.com/qiniu/log"
)

func init() {
	log.SetOutputLevel(0)

}
func TestJoinURL(t *testing.T) {
	url := JoinURL("http://localhost", "xxx", "yyyy")
	assert.Equal(t, "http://localhost/xxx/yyyy", url)
}
