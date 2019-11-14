package integrit

import "github.com/google/uuid"

type validation interface {
	match(appid uuid.UUID, apiUpload string, agentUpload string) (bool error)
}
