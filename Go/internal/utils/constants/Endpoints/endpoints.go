// Package Endpoints provides constant endpoint paths used in the application.
package Endpoints

import (
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Paths"
)

const (
	ApiKey           = Paths.Util + constants.APIVersion + "/key"
	UserRegistration = Paths.Util + constants.APIVersion + "/user/register"
	UserDeletion     = Paths.Util + constants.APIVersion + "/user/delete"
	UserDeletionId   = Paths.Util + constants.APIVersion + "/user/delete/{ID}"
	RegistrationsID  = Paths.Dashboards + constants.APIVersion + "/registrations/{ID}"
	Registrations    = Paths.Dashboards + constants.APIVersion + "/registrations"
	DashboardsID     = Paths.Dashboards + constants.APIVersion + "/dashboard/{ID}"
	Dashboards       = Paths.Dashboards + constants.APIVersion + "/dashboard"
	NotificationsID  = Paths.Dashboards + constants.APIVersion + "/notifications/{ID}"
	Notifications    = Paths.Dashboards + constants.APIVersion + "/notifications"
	Status           = Paths.Dashboards + constants.APIVersion + "/status"
)
