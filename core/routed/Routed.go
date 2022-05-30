package routed

type Routed struct {
	NextRouted any
}

func (r *Routed) Next() any {
	return r.NextRouted
}
