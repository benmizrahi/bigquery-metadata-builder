package models

import "cloud.google.com/go/bigquery"

type Mapping struct {
	LinkDS *bigquery.Dataset
	Tables []*bigquery.Table
}
