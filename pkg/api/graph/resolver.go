package graph

import (
	"github.com/phyrwork/benevolent-dictator/pkg/api/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB *database.DB
}
