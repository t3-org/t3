package doc

type ignoreId struct {
	Id string `json:"-"`
}

type ignoreUserID struct {
	UserID string `json:"-"`
}

type Pagination struct {
	Page    int `json:"page"`
	PerPage int `json:"per_page"`
}
