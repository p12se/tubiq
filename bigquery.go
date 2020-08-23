package main

import (
	"cloud.google.com/go/bigquery"
	"context"
	"google.golang.org/api/iterator"
	"strings"
)

type bqMeta interface {
	getName() string
	isDataset() bool
	isTable() bool
}

type bq struct {
	client  *bigquery.Client
	ctx     context.Context
	dataset []dataset
}

type table struct {
	name      string
	dataset   dataset
	projectID string
}

func (t table) isDataset() bool {
	return false
}

func (t table) isTable() bool {
	return true
}

func (t table) getName() string {
	return t.name
}

type dataset struct {
	name   string
	tables []table
}

func (d dataset) isTable() bool {
	return false
}

func (d dataset) getName() string {
	return d.name
}

func (d dataset) isDataset() bool {
	return true
}
func newBq(ctx context.Context, projectID string) *bq {
	client, err := bigquery.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}
	return &bq{client: client}
}

func (b *bq) datasets() ([]dataset, error) {
	datasets := []dataset{}
	it := b.client.Datasets(b.ctx)

	for {
		dt, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []dataset{}, err
		}
		datasets = append(datasets, dataset{name: dt.DatasetID})
	}
	return datasets, nil
}

//func (b bq)list(path string) []meta {
//	t := []table{}
//	if strings.EqualFold(path, ""){
//		return bq.datasets()
//	}
//
//	return
//}

func (b *bq) tables(dtset string) ([]table, error) {

	tables := []table{}
	dt := b.client.Dataset(dtset)

	it := dt.Tables(b.ctx)
	for {
		tb, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []table{}, err
		}
		tables = append(tables, table{name: tb.TableID, dataset: dataset{name: tb.DatasetID}, projectID: dt.ProjectID})
	}

	return tables, nil
}

func (b bq) list(path string) ([]bqMeta, error) {
	if strings.EqualFold(path, "") {
		dts := []bqMeta{}
		listDataset, err := b.datasets()
		if err != nil {
			return nil, err
		}
		for _, d := range listDataset {
			dts = append(dts, d)
		}
		return dts, nil
	}

	tbs := []bqMeta{}

	listTables, err := b.tables(path)

	if err != nil {
		return nil, err
	}

	for _, t := range listTables {
		tbs = append(tbs, t)
	}

	return tbs, nil
}
