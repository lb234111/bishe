package common

import "fmt"

type Worker struct {
	HostName string
}

func NewWorker(IP string) *Worker {
	WorkerHostName := fmt.Sprintf("http://%s:11110", IP)
	return &Worker{HostName: WorkerHostName}
}
