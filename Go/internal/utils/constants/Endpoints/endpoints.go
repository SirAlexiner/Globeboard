// Package Endpoints provides constant endpoint paths used in the application.
package Endpoints

import (
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Paths"
)

// Registrations represent the path to the registrations' endpoint.
// Dashboards represent the path to the dashboards' endpoint.
// Notifications represent the path to the notifications' endpoint.
// Status represents the path to the status endpoint.
const (
	ApiKey        = Paths.Util + constants.APIVersion + "/key/"
	Registrations = Paths.Dashboard + constants.APIVersion + "/registrations/"
	Dashboards    = Paths.Dashboard + constants.APIVersion + "/dashboards/"
	Notifications = Paths.Dashboard + constants.APIVersion + "/notifications/"
	Status        = Paths.Dashboard + constants.APIVersion + "/status/"
)
