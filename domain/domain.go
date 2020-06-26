package domain

type Link struct {
	Key         string
	URL         string
	Visits      int64
	Creator     string
	CreatedTime int64
	ExpiredTime int64
}

type User struct {
	Email string
}
