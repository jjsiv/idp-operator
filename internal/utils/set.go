package utils

type Set map[string]struct{}

func (s Set) Insert(v string) {
	s[v] = struct{}{}
}

func (s Set) Has(v string) bool {
	_, ok := s[v]
	return ok
}
