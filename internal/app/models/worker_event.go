package models

type WorkerEvent struct {
	CorrelationID string
	Amount        float32
	ResultCh      chan *WorkerEventResult
}

type WorkerEventResult struct{}
