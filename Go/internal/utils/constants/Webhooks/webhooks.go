package Webhooks

const (
	POSTTitle   = "New Country Data Registered to GlobeBoard"
	PUTTitle    = "Country Data changed on GlobeBoard"
	DELETETitle = "Country Data Deleted from GlobeBoard"
	GETTitle    = "Country Data Invoked from GlobeBoard"

	POSTColor   = 2664261  //Success Color
	PUTColor    = 16761095 //Update Color
	DELETEColor = 14431557 //Warning Color
	GETColor    = 1548984  //Info Color

	EventRegister = "REGISTER" // POST
	EventChange   = "CHANGE"   // PUT
	EventDelete   = "DELETE"   // DELETE
	EventInvoke   = "INVOKE"   // GET
)
