package errors

import (
	"fmt"
	"strings"
)

type MissingFilesError []string

func (e *MissingFilesError) Error() string {
	s := strings.Join(*e, ", ")
	return fmt.Sprintf("The following files could not be found: %s\n", s)
}
