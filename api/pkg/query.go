package pkg

import (
	"strings"

	"streaming/internal/shared"
)

func QueryType(source string) shared.QueryType {
	if strings.HasPrefix(source, "rtsp") {
		return shared.TypeLink
	}
	return shared.TypeFile
}
