// Package Endpoints provides constant endpoint paths used in the application.
package Endpoints

import (
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Paths"
)

const (
	ApiKey           = Paths.Util + constants.APIVersion + "/key"
	UserRegistration = Paths.Util + constants.APIVersion + "/register"
	RegistrationsID  = Paths.Dashboard + constants.APIVersion + "/registrations/{id}"
	Registrations    = Paths.Dashboard + constants.APIVersion + "/registrations"
	Dashboards       = Paths.Dashboard + constants.APIVersion + "/dashboards/{id}"
	NotificationsID  = Paths.Dashboard + constants.APIVersion + "/notifications/{id}"
	Notifications    = Paths.Dashboard + constants.APIVersion + "/notifications"
	Status           = Paths.Dashboard + constants.APIVersion + "/status"
)
