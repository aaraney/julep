package build

type Builder interface {
	Build(job Job) error
	Tag(job Job) error
}

type BuilderPusher interface {
	Builder
	Push() error
}

// TODO: implement DefaultBuilder and DefaultBuilderPusher
