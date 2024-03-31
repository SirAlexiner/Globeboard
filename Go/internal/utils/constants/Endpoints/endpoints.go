// Package Endpoints provides constant endpoint paths used in the application.
package Endpoints

import (
	"globeboard/internal/utils/constants"
	"globeboard/internal/utils/constants/Paths"
)

// Library represents the path to the library endpoint.
// SupportedLanguages represents the path to the supported_languages endpoint.
// BookCount represents the path to the bookcount endpoint.
// Readership represents the path to the readership endpoint.
// Status represents the path to the status endpoint.
const (
	Library            = Paths.Library + constants.APIVersion + "/"
	ApiKey             = Paths.Library + constants.APIVersion + "/key"
	SupportedLanguages = Paths.LibraryStats + constants.APIVersion + "/supported_languages/"
	BookCount          = Paths.LibraryStats + constants.APIVersion + "/bookcount/"
	Readership         = Paths.LibraryStats + constants.APIVersion + "/readership/"
	Status             = Paths.LibraryStats + constants.APIVersion + "/status/"
)
