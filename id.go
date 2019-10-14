package timy

import (
	"github.com/lithammer/shortuuid"
)

func generateNewID() string {
	return shortuuid.New()
}
