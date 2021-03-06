package regulator

import "sync/atomic"

type RegulatorError struct {
	JobIndex int32
	Message  string
}

func (err RegulatorError) Error() string {
	return err.Message
}

type Regulator struct {
	concurrency int
	sem         chan bool
	jobIndex    int32
	err         error
}

func (regulator *Regulator) Execute(job func() error) {
	index := atomic.AddInt32(&regulator.jobIndex, 1)
	regulator.sem <- true
	// if error occurred while waiting for channel to unblock, return
	if regulator.err != nil {
		<-regulator.sem
		return
	}
	go func(jobIndex int32) {
		defer func() { <-regulator.sem }()
		err := job()
		if err != nil {
			regulator.err = RegulatorError{JobIndex: jobIndex, Message: err.Error()}
		}
	}(index)
}

func (regulator *Regulator) Wait() error {
	for i := 0; i < cap(regulator.sem); i++ {
		regulator.sem <- true
	}
	close(regulator.sem)
	return regulator.err
}

func NewRegulator(concurrency int) *Regulator {
	return &Regulator{
		concurrency: concurrency,
		sem:         make(chan bool, concurrency),
	}
}
