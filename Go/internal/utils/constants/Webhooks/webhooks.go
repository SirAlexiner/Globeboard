package Webhooks

const (
	POSTTitle   = "Registered New Country Data to GlobeBoard"
	PUTTitle    = "Changed Country Data on GlobeBoard"
	DELETETitle = "Deleted Country Data from GlobeBoard"
	GETTitle    = "Invoked Country Data from GlobeBoard"

	POSTColor   = 2664261  //Success Color
	PUTColor    = 16761095 //Update Color
	DELETEColor = 14431557 //Warning Color
	GETColor    = 1548984  //Info Color

	EventRegister = "REGISTER" // POST
	EventChange   = "CHANGE"   // PUT
	EventDelete   = "DELETE"   // DELETE
	EventInvoke   = "INVOKE"   // GET
)
