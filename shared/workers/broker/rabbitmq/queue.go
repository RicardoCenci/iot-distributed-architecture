package rabbitmq

type Queue struct {
	queueName string
	options   *QueueOptions
}

type QueueOptions struct {
	durable          bool
	deleteWhenUnused bool
	exclusive        bool
	noWait           bool
	arguments        map[string]interface{}
}

type Option func(*QueueOptions)

func WithDurable(durable bool) Option {
	return func(options *QueueOptions) {
		options.durable = durable
	}
}

func WithDeleteWhenUnused(deleteWhenUnused bool) Option {
	return func(options *QueueOptions) {
		options.deleteWhenUnused = deleteWhenUnused
	}
}

func WithExclusive(exclusive bool) Option {
	return func(options *QueueOptions) {
		options.exclusive = exclusive
	}
}

func WithNoWait(noWait bool) Option {
	return func(options *QueueOptions) {
		options.noWait = noWait
	}
}

func WithArguments(arguments map[string]interface{}) Option {
	return func(options *QueueOptions) {
		options.arguments = arguments
	}
}

func NewQueue(queueName string, options ...Option) *Queue {
	defaultOptions := &QueueOptions{
		durable:          true,
		deleteWhenUnused: false,
		exclusive:        false,
		noWait:           false,
		arguments:        nil,
	}

	for _, option := range options {
		option(defaultOptions)
	}

	return &Queue{
		queueName: queueName,
		options:   defaultOptions,
	}
}

func (q *Queue) GetName() string {
	return q.queueName
}
