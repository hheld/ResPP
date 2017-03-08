package main

import (
	"fmt"
	"strings"
)

type missingFilesError []string

func (e *missingFilesError) Error() string {
	s := strings.Join(*e, ", ")
	return fmt.Sprintf("The following files could not be found: %s\n", s)
}
