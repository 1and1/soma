package somaproto

type Grant struct {
	RecipientType string        `json:"recipientType"` //user,team,tool
	RecipientId   string        `json:"recipientId"`
	Permission    string        `json:"permission"`
	RepositoryId  string        `json:"repositoryId"`
	ObjectType    string        `json:"objectType"`
	ObjectId      string        `json:"objectId"`
	Details       *GrantDetails `json:"details, omitempty"`
}

type GrantDetails struct {
	DetailsCreation
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
