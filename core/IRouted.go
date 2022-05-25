package core

type IRouted interface {
	Next() any
}

type Routed struct {
	NextRouted any
}

func (r *Routed) Next() any {
	return r.NextRouted
}
