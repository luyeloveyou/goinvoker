package core

type IRouted interface {
	Next() Object
}

type Routed struct {
	NextRouted Object
}

func (r *Routed) Next() Object {
	return r.NextRouted
}
