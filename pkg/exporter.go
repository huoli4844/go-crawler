package pkg

import "context"

type Exporter interface {
	Export(context.Context, Item) error
	PipelineNames() map[uint8]string
	Pipelines() []Pipeline
	SetPipeline(Pipeline, uint8)
	DelPipeline(int)
	CleanPipelines()
	WithDumpPipeline()
	WithFilePipeline()
	WithImagePipeline()
	WithFilterPipeline()
	WithCsvPipeline()
	WithJsonLinesPipeline()
	WithMongoPipeline()
	WithMysqlPipeline()
	WithKafkaPipeline()
	WithCustomPipeline(Pipeline)
}
