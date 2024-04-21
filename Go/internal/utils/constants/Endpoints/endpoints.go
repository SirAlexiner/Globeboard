// Package Endpoints provides constants for fully constructed endpoint paths
// using base paths and versions defined in other packages.
package Endpoints

import (
	"globeboard/internal/utils/constants"       // Imports constants for the API version.
	"globeboard/internal/utils/constants/Paths" // Imports base paths for endpoints.
)

const (
	// ApiKey endpoint for API key operations.
	ApiKey = Paths.Util + constants.APIVersion + "/key"
	// UserRegistration endpoint for user registration operations.
	UserRegistration = Paths.Util + constants.APIVersion + "/user/register"
	// UserDeletion endpoint URL for user registration operations without the ID wildcard.
	UserDeletion = Paths.Util + constants.APIVersion + "/user/delete"
	// UserDeletionID endpoint for user deletion operations by ID.
	UserDeletionID = Paths.Util + constants.APIVersion + "/user/delete/{ID}"
	// RegistrationsID endpoint for accessing specific registration by ID.
	RegistrationsID = Paths.Dashboards + constants.APIVersion + "/registrations/{ID}"
	// Registrations endpoint for accessing registrations without ID.
	Registrations = Paths.Dashboards + constants.APIVersion + "/registrations"
	// DashboardsID endpoint for accessing specific dashboard by ID.
	DashboardsID = Paths.Dashboards + constants.APIVersion + "/dashboard/{ID}"
	// Dashboards endpoint URL for dashboard operations without the ID wildcard.
	Dashboards = Paths.Dashboards + constants.APIVersion + "/dashboard"
	// NotificationsID endpoint for accessing specific notification by ID.
	NotificationsID = Paths.Dashboards + constants.APIVersion + "/notifications/{ID}"
	// Notifications endpoint for accessing notifications operations without ID.
	Notifications = Paths.Dashboards + constants.APIVersion + "/notifications"
	// Status endpoint for checking the status of the dashboard services.
	Status = Paths.Dashboards + constants.APIVersion + "/status"
)
