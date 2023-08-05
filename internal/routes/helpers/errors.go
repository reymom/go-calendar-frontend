package helpers

type constError string

func (c constError) Error() string {
	return string(c)
}

const (
	ErrTemplateEmpty     constError = "template is empty"
	ErrUnknownFilterMode constError = "unknown time filter mode"
	ErrIncorrectIsoWeek  constError = "incorrect ISO week"
	ErrIncorrectIsoYear  constError = "incorrect ISO year"
)
