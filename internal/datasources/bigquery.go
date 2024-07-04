package datasources

import (
	"context"
	"os"
	"sync"

	"cloud.google.com/go/bigquery"
	"github.com/benmizrahi/bigquery-metadata-builder/internal/models"
	milvus "github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	log "github.com/sirupsen/logrus"
)

type Bigquery struct {
	ctx     context.Context
	client  *bigquery.Client
	project string
	mapping []*models.Mapping
	vector  milvus.Client
}

// private contractor
func newBigquery() IDatasource {
	ctx := context.Background()
	project := os.Getenv("PROJECT_ID")

	client, err := bigquery.NewClient(ctx, project)
	if err != nil {
		log.Fatal("unable to start BQ client,", err)
	}

	cmilvus, err := milvus.NewClient(ctx, milvus.Config{
		Address: os.Getenv("MILVUS_PATH"),
	})

	if err != nil {
		log.Fatal("unable to start BQ client,", err)
	}

	return Bigquery{
		ctx:     context.Background(),
		client:  client,
		project: project,
		mapping: make([]*models.Mapping, 0),
		vector:  cmilvus,
	}
}

// Explore implements IDatasource.
func (b Bigquery) Explore() IDatasource {

	log.Info("Analyze Datasets in project:", b.project)

	var wg sync.WaitGroup

	b.exploreDatasets(wg)
	b.exploreTables(wg)

	wg.Wait()

	return b
}

func (b Bigquery) exploreDatasets(_ sync.WaitGroup) {

	datasetsIterator := b.client.Datasets(b.ctx)
	dataset, err := datasetsIterator.Next()

	if err != nil {
		log.Error("unable to get datasets, ", err)
	}
	for {
		if dataset == nil {
			break
		}
		b.mapping = append(b.mapping, &models.Mapping{
			LinkDS: dataset,
			Tables: make([]*bigquery.Table, 0),
		})

		dataset, err = datasetsIterator.Next()
		if err != nil {
			log.Error("unable to get datasets, ", err)
		}
	}

	log.Info("Explore datasets done found: ", len(b.mapping))

}

func (b Bigquery) exploreTables(wg sync.WaitGroup) {
	for _, ds := range b.mapping {
		wg.Add(1)
		go b.mapTables(ds, wg)
	}

}

func (b Bigquery) mapTables(ds *models.Mapping, wg sync.WaitGroup) {
	defer wg.Done()
	tableIterable := ds.LinkDS.Tables(b.ctx)
	table, err := tableIterable.Next()
	if err != nil {
		log.Error("unable to get datasets, ", err)
	}
	for {

		if table == nil {
			break
		}

		ds.Tables = append(ds.Tables, table)

		table, err = tableIterable.Next()
		if err != nil {
			log.Fatal("unable to get datasets, ", err)
		}
	}

	log.Info("found ", len(ds.Tables), " tables in ds ", ds.LinkDS.DatasetID)

}

// BuildRAG implements IDatasource.
func (b Bigquery) BuildVector() IDatasource {
	exists, err := b.vector.HasCollection(context.Background(), os.Getenv("DATABASE_COLLECTION"))
	if err != nil {
		log.Fatal("unable to check is collection exists, ", err)
	}
	if !exists {
		schema := &entity.Schema{
			CollectionName: "",
			Description:    "Test book search",
			Fields: []*entity.Field{
				{
					Name:       "table_id",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: true,
					AutoID:     false,
				},
				{
					Name:       "table_name",
					DataType:   entity.FieldTypeInt64,
					PrimaryKey: false,
					AutoID:     false,
				},
				{
					Name:     "table_type",
					DataType: entity.FieldTypeFloatVector,
					TypeParams: map[string]string{
						"dim": "2",
					},
				},
			},
		}

		err = b.vector.CreateCollection(
			context.Background(),
			schema,
			2,
		)
		if err != nil {
			log.Fatal("failed to create collection:", err.Error())
		}
	}

	return nil
}

// Persist implements IDatasource.
func (b Bigquery) Persist() IDatasource {
	panic("unimplemented")
}
