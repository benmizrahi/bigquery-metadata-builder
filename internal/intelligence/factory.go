package intelligence

import "github.com/benmizrahi/bigquery-metadata-builder/internal/datasources"

type IIntelligence interface {
	SuggestMetadata(datasources.IDatasource) IIntelligence
}

func Resolver(aitype *string) IIntelligence {
	return newGimini()
}
