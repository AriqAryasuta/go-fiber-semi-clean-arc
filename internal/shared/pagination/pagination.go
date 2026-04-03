package pagination

type Query struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (q Query) Normalize() Query {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 {
		q.Limit = 10
	}
	if q.Limit > 100 {
		q.Limit = 100
	}
	return q
}

func (q Query) Offset() int {
	return (q.Page - 1) * q.Limit
}
