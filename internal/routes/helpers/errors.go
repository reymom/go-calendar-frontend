package helpers

type constError string

func (c constError) Error() string {
	return string(c)
}

const (
	ErrTemplateEmpty constError = "template is empty"
)
