package upcloud

/**
Represents a list of tags
*/
type Tags struct {
	Tags []Tag `xml:"tag"`
}

/**
Represents a server tag
*/
type Tag struct {
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Servers     []string `xml:"servers>server"`
}
