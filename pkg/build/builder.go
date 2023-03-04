package build

import (
	"bufio"
	"context"
	"encoding/json"
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

type BuilderMessage struct {
	ID      string `json:"ID"`
	Payload []byte `json:"payload"`
}

// TODO: implement DefaultBuilder and DefaultBuilderPusher
type DefaultBuilder struct {
	output_stream chan<- BuilderMessage
}

func NewDefaultBuilder(stream chan<- BuilderMessage) DefaultBuilder {
	return DefaultBuilder{output_stream: stream}
}

type errorPayload struct {
	Error string `json:"error"`
}

type statusPayload struct {
	Status string `json:"status"`
}

func (b DefaultBuilder) Build(job Job) error {
	cli, _ := client.NewClientWithOpts()

	// TODO: refactor
	dockerfile := filepath.Base(job.Path())
	ctx := filepath.Dir(job.Path())

	img_opts := types.ImageBuildOptions{
		Tags:       []string{job.FullName()},
		Dockerfile: dockerfile,
		BuildID:    job.FullName(),
	}

	// TODO: not sure what the guarantees are for the context of this request
	res, err := cli.ImageBuild(context.Background(), getContext(ctx), img_opts)

	if err != nil {
		payload, _ := json.Marshal(errorPayload{Error: err.Error()})
		b.output_stream <- BuilderMessage{ID: job.FullName(), Payload: payload}
		return err
	}
	defer res.Body.Close()

	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		b.output_stream <- BuilderMessage{ID: job.FullName(), Payload: scanner.Bytes()}
		// fmt.Printf("%s\n", scanner.Text())
		// fmt.Println(scanner.Text()) // Println will add back the final '\n'
	}

	// body, _ := io.ReadAll(res.Body)
	// fmt.Printf("%s\n", body)
	done_msg, _ := json.Marshal(statusPayload{Status: "complete"})
	b.output_stream <- BuilderMessage{ID: job.FullName(), Payload: done_msg}

	return nil
}

// TODO: implement
func (DefaultBuilder) Tag(job Job) error {
	return nil
}

func (b DefaultBuilder) Do(job Job) error {
	return b.Build(job)
}

func (DefaultBuilder) Cancel(job Job) {
	cli, _ := client.NewClientWithOpts()
	err := cli.BuildCancel(context.Background(), job.FullName())
	fmt.Println(err.Error())
}
