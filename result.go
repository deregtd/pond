package pond

import (
	"github.com/alitto/pond/v2/internal/future"
)

// ResultPool is a pool that can be used to submit tasks that return a result.
type ResultPool[R any] interface {
	basePool

	// Submits a task to the pool and returns a future that can be used to wait for the task to complete and get the result.
	Submit(task func() R) Result[R]

	// Submits a task to the pool and returns a future that can be used to wait for the task to complete and get the result.
	SubmitErr(task func() (R, error)) Result[R]

	// Creates a new subpool with the specified maximum concurrency.
	NewSubpool(maxConcurrency int) ResultPool[R]

	// Creates a new task group.
	NewGroup() ResultTaskGroup[R]
}

type resultPool[R any] struct {
	*pool
}

func (p *resultPool[R]) NewGroup() ResultTaskGroup[R] {
	return newResultTaskGroup[R](p.pool)
}

func (p *resultPool[R]) Submit(task func() R) Result[R] {
	return p.submit(task)
}

func (p *resultPool[R]) SubmitErr(task func() (R, error)) Result[R] {
	return p.submit(task)
}

func (p *resultPool[R]) submit(task any) Result[R] {
	future, resolve := future.NewValueFuture[R](p.Context())

	wrapped := wrapTask[R, func(R, error)](task, resolve)

	p.dispatcher.Write(wrapped)

	return future
}

func (p *resultPool[R]) NewSubpool(maxConcurrency int) ResultPool[R] {
	return newResultSubpool[R](maxConcurrency, p.Context(), p.pool)
}

func NewResultPool[R any](maxConcurrency int, options ...Option) ResultPool[R] {
	return &resultPool[R]{
		pool: newPool(maxConcurrency, options...),
	}
}