package domain

type ImageValidator interface {
	Validate(image []byte) error
	AllowedTypes() []string
}
