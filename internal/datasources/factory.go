package datasources

type IDatasource interface {
	Explore() IDatasource
	BuildVector() IDatasource
	Persist() IDatasource
}

func Resolver(dstype *string) IDatasource {
	return newBigquery()
}
