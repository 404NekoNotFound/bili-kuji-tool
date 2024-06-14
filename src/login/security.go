package login

import (
	"fmt"
	"github.com/google/uuid"
	"strings"
)

func buvid() string {
	u := uuid.NewString()
	u = strings.ReplaceAll(u, "-", "")

	return fmt.Sprintf("XY%v%v%v%v", string(u[2]), string(u[12]), string(u[22]), u)
}
