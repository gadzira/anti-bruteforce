package helpers

import (
	"strings"

	"github.com/gadzira/anti-bruteforce/internal/db"
)

func MakeEntry(s, t string) *db.Entry {
	ipMaskList := strings.Split(s, ":")
	return &db.Entry{
		IP:   ipMaskList[0],
		Mask: ipMaskList[1],
		List: t,
	}
}
