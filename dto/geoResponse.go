package dto

type GeoResponse struct {
	CountryOffset uintptr `maxminddb:"country"`

	Location struct {
		Latitude float64 `maxminddb:"latitude"`
		// Longitude is directly nested within the parent map.
		LongitudeOffset uintptr `maxminddb:"longitude"`
		// TimeZone is indirected via a pointer.
		TimeZoneOffset uintptr `maxminddb:"time_zone"`
	} `maxminddb:"location"`
}
