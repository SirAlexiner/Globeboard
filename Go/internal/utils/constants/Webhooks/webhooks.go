// Package Webhooks defines constants for webhook titles, colors, and event types used for Webhook notifications.
package Webhooks

const (
	POSTTitle   = "Registered New Country Data to GlobeBoard" // POSTTitle defines the title for POST webhook events.
	PATCHTitle  = "Changed Country Data on GlobeBoard"        // PATCHTitle defines the title for PATCH webhook events.
	DELETETitle = "Deleted Country Data from GlobeBoard"      // DELETETitle defines the title for DELETE webhook events.
	GETTitle    = "Invoked Country Data from GlobeBoard"      // GETTitle defines the title for GET webhook events.

	POSTColor   = 2664261  // Success Color - light green
	PUTColor    = 16761095 // Update Color - bright orange
	DELETEColor = 14431557 // Warning Color - pale red
	GETColor    = 1548984  // Info Color - light blue

	EventRegister = "REGISTER" // EventRegister defines the event type for POST operations.
	EventChange   = "CHANGE"   // EventChange defines the event type for PATCH operations.
	EventDelete   = "DELETE"   // EventDelete defines the event type for DELETE operations.
	EventInvoke   = "INVOKE"   // EventInvoke defines the event type for GET operations.
)
