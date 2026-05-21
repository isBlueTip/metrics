package sender

import (
	"github.com/sethgrid/pester"
)

type Sender struct {
	HC   *pester.Client
	Addr string
	Key  string
}
