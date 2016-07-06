package upcloud

/**
Represents a /zone response
*/
type Zones struct {
	Zones []Zone `xml:"zone"`
}

/**
Represents a zone
*/
type Zone struct {
	Id          string `xml:"id"`
	Description string `xml:"description"`
}
