package proto

type User struct {
	Id             string           `json:"id,omitempty"`
	UserName       string           `json:"userName,omitempty"`
	FirstName      string           `json:"firstName,omitempty"`
	LastName       string           `json:"lastName,omitempty"`
	EmployeeNumber string           `json:"employeeNumber,omitempty"`
	MailAddress    string           `json:"mailAddress,omitempty"`
	IsActive       bool             `json:"isActive,omitempty"`
	IsSystem       bool             `json:"isSystem,omitempty"`
	IsDeleted      bool             `json:"isDeleted,omitempty"`
	TeamId         string           `json:"teamId,omitempty"`
	Details        *UserDetails     `json:"details,omitempty"`
	Credentials    *UserCredentials `json:"credentials,omitempty"`
}

type UserCredentials struct {
	Reset          bool   `json:"reset,omitempty"`
	ForcedPassword string `json:"forcedPassword,omitempty"`
}

type UserDetails struct {
	DetailsCreation
}

type UserFilter struct {
	UserName  string `json:"userName,omitempty"`
	IsActive  bool   `json:"isActive,omitempty"`
	IsSystem  bool   `json:"isSystem,omitempty"`
	IsDeleted bool   `json:"isDeleted,omitempty"`
}

func NewUserRequest() Request {
	return Request{
		Flags: &Flags{},
		User:  &User{},
	}
}

func NewUserFilter() Request {
	return Request{
		Filter: &Filter{
			User: &UserFilter{},
		},
	}
}

func NewUserResult() Result {
	return Result{
		Errors: &[]string{},
		Users:  &[]User{},
	}
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
