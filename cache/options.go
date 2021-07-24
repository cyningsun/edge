package cache

const (
	maxSegments = 1 << 16
	maxCapacity = 1 << 30
)

type options struct {
	concurrency int
	capacity    int
}

type Opt func(*options)

func WithConcurrency(c int) Opt {
	return func(o *options) {
		o.concurrency = c
	}
}

func WithCapacity(c int) Opt {
	return func(o *options) {
		o.capacity = c
	}
}
