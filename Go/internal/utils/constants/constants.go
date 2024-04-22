// Package constants provide constant values that are used within the application.
package constants

const (
	APIVersion   = "v1" // APIVersion specifies the version of the API being used.
	ApiKeyLength = 20   // ApiKeyLength specifies the length of API keys generated.
	DocIdLength  = 24   // DocIdLength specifies the length of document identifiers.
	IdLength     = 20   // IdLength specifies the length of general purpose identifiers.

	// ClientConnectUnsupported formats an error message for when a client tries to connect using an unsupported method.
	ClientConnectUnsupported = "Client attempted to connect to %s with unsupported method: %s\n"
	// ClientConnectNoToken formats an error message for connection attempts where no token is provided.
	ClientConnectNoToken = "Failed %s attempt to %s: No Token.\n"
	// ClientConnectNoID formats an error message for connection attempts where no ID is provided.
	ClientConnectNoID = "Failed %s attempt to %s: No ID.\n"
	// ClientConnectUnauthorized formats an error message for unauthorized connection attempts.
	ClientConnectUnauthorized = "Unauthorized %s attempted to %s.\n"
	// ClientConnectEmptyBody formats an error message for connection attempts with no body content.
	ClientConnectEmptyBody = "Failed %s attempt to %s: No Body.\n"
)
