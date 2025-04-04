package app

import (
	"strings"

	"github.com/gofrs/uuid/v5"
)

func (a App) Uuid(name string) (string, error) {
	namespace, err := uuid.FromString(a.env.MY_NAMESPACE)
	if err != nil {
		return "", nil
	}
	u1 := uuid.NewV5(namespace, name)
	return strings.ReplaceAll(u1.String(), "-", ""), nil
}
