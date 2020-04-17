package broker

const (
	exchangeName       = "micropic"
	RequestQueueName   = "previewer.requestQueue"
	RequestRoutingKey  = "previewer.requestQueue.tasks"
	ResponseQueueName  = "previewer.responseQueue"
	ResponseRoutingKey = "previewer.responseQueue.previews"
)
