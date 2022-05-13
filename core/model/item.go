package model

// ContentItem is a workaround due to problem with data json & bson encode and decode with abstract type
type ContentItem = map[string]interface{}
