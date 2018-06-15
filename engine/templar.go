package engine

type Templar interface {
	GenerateTemplate() (string, error)
}
