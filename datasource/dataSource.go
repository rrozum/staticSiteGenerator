package datasource

type DataSource interface {
	Fetch(from, to string) ([] string, error)
}

func New() DataSource {
	return &GitDataSource{}
}