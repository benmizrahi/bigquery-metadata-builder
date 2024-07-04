package intelligence

import "github.com/benmizrahi/bigquery-metadata-builder/internal/datasources"

type GiminiIntelligence struct{}

// SuggestMetadata implements IIntelligence.
func (g *GiminiIntelligence) SuggestMetadata(datasources.IDatasource) IIntelligence {
	panic("unimplemented")
}

func newGimini() IIntelligence {
	return &GiminiIntelligence{}
}
