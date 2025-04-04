package filter

import (
	"quiz-builder/filter/auth"
	"quiz-builder/filter/cors"
	"quiz-builder/filter/database"

	"go.uber.org/fx"
)

// Module Filter exported
var Module = fx.Options(
	fx.Provide(cors.NewCorsFilter),
	fx.Provide(auth.NewAuthFilter),
	fx.Provide(database.NewDatabaseTrx),
	fx.Provide(NewFilters),
)

// IFilter filter interface
type Filter interface {
	Setup()
}

// Filters contains multiple filters
type Filters []Filter

// NewFilters creates new filters
// Register the filter that should be applied directly (globally)
func NewFilters(
	corsFilter cors.CorsFilter,
	dbTrxFilter database.DatabaseTrx,
) Filters {
	return Filters{
		corsFilter,
		dbTrxFilter,
	}
}

// Setup sets up filters
func (f Filters) Setup() {
	for _, filter := range f {
		filter.Setup()
	}
}
