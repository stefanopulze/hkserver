package characteristic

import "github.com/brutella/hc/characteristic"

const (
	TypePower      = "032B12CF-D4E8-4277-9021-188816FD00C6"
	TypeTotalPower = "032B12CA-D4E8-4277-9021-188816FD00C6"
)

type Power struct {
	*characteristic.Int
}

func NewPower(val int) *Power {
	p := Power{characteristic.NewInt("")}
	p.Value = val
	p.Format = characteristic.FormatUInt64
	p.Perms = characteristic.PermsRead()

	return &p
}
