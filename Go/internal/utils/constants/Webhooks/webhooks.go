package Webhooks

const (
	POSTTitle    = "New Country Data Registered to GlobeBoard"
	GETRegTitle  = "Registered Country Data Requested from GlobeBoard"
	GETDashTitle = "Populated Country Data Requested from GlobeBoard"
	PUTTitle     = "Country Data Updated on GlobeBoard"
	DELETETitle  = "Country Data Deleted from GlobeBoard"

	POSTColor   = 2664261  //Success Color
	GETColor    = 1548984  //Info Color
	PUTColor    = 16761095 //Update Color
	DELETEColor = 14431557 //Warning Color

	EventRegister = "REGISTER"
	EventChange   = "CHANGE"
	EventDelete   = "DELETE"
	EventInvoke   = "INVOKE"
)
