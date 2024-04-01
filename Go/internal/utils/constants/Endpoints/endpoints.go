// Package Endpoints provides constant endpoint paths used in the application.
package Endpoints

import (
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Paths"
)

// Util represents the path to the util endpoint.
// Registrations represents the path to the supported_languages endpoint.
// Dashboards represents the path to the bookcount endpoint.
// Notifications represents the path to the readership endpoint.
// Status represents the path to the status endpoint.
const (
	Util          = Paths.Util + constants.APIVersion + "/"
	ApiKey        = Paths.Util + constants.APIVersion + "/key"
	Registrations = Paths.Dashboard + constants.APIVersion + "/registrations/"
	Dashboards    = Paths.Dashboard + constants.APIVersion + "/dashboards/"
	Notifications = Paths.Dashboard + constants.APIVersion + "/notifications/"
	Status        = Paths.Dashboard + constants.APIVersion + "/status/"
)
