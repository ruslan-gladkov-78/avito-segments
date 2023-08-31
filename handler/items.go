package handler

type Segment struct {
	Slug string `json:"slug"`
}

type User struct {
	ID int `json:"id"`
}

type Segments struct {
	Slugs []string `json:"segments"`
}

type ChangeRequest struct {
	User     User     `json:"user"`
	ToAdd    Segments `json:"to_add"`
	ToDelete Segments `json:"to_delete"`
	TTL      int
}
