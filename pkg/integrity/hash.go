package integrit

import "github.com/google/uuid"

type hashValidation struct {
}

func (hv hashValidation) match(appid uuid.UUID, hashFromAPI string, hashFromAgent string) (bool, error) {
	if hashFromAPI == hashFromAgent {
		return true, nil
	}
	return false, nil
}