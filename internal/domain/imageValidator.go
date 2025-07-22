package domain

type ImageValidator struct {
	allowedTypes map[string]struct{}
}

type ImageValidator interface{}
