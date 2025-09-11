package reflect

type TagType string

const (
	TagTypePath   TagType = "path"
	TagTypeHeader TagType = "header"
	TagTypeQuery  TagType = "query"
)
