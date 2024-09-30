package handlers

type Storage interface {
	Update(name string, typeMetric string, val string) error
	String() string
}
