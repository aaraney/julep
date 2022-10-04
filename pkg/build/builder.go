package build

import (
	"context"
	"fmt"
	"io"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

// TODO: rethink this interface. need a way to send status updates and cancel a build. At face
// value, I am thinking a status channel and a context.
type Builder interface {
	Build(job Job) error
	Tag(job Job) error
}

type BuilderPusher interface {
	Builder
	Push() error
}

func getContext(filePath string) io.Reader {
	ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
	return ctx
}

// TODO: implement DefaultBuilder and DefaultBuilderPusher
type DefaultBuilder struct{}

func (DefaultBuilder) Build(job Job) error {
	cli, _ := client.NewClientWithOpts()

	// TODO: refactor
	dockerfile := filepath.Base(job.Path())
	ctx := filepath.Dir(job.Path())

	img_opts := types.ImageBuildOptions{
		Tags:       []string{job.FullName()},
		Dockerfile: dockerfile,
		BuildID:    job.FullName(),
	}

	res, err := cli.ImageBuild(context.Background(), getContext(ctx), img_opts)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	fmt.Printf("%s\n", body)

	return nil
}

// TODO: implement
func (DefaultBuilder) Tag(job Job) error {
	return nil
}
