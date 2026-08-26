package shared

import "context"

// GeoResolverInterface provides countries and timezones. The host
// implements this against whatever geo data source it uses
// (dracory/geostore, a static list, an external API, etc.).
type GeoResolverInterface interface {
	// Countries returns all countries.
	Countries(ctx context.Context) ([]Country, error)

	// Timezones returns timezones for the given country code. The
	// country code is optional — pass no argument or an empty string
	// when no country is selected; implementations should return an
	// empty list in that case.
	Timezones(ctx context.Context, countryCode ...string) ([]Timezone, error)
}

// Country is a single country entry returned by GeoResolverInterface.
type Country struct {
	IsoCode2 string
	Name     string
}

// Timezone is a single timezone entry returned by GeoResolverInterface.
type Timezone struct {
	Code string
}
