package helpers

import (
	"strings"

	"github.com/gadzira/anti-bruteforce/internal/database"
)

func MakeEntry(s, t string) *database.Entry {
	ipMaskList := strings.Split(s, ":")
	return &database.Entry{
		IP:   ipMaskList[0],
		Mask: ipMaskList[1],
		List: t,
	}
}
