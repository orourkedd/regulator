package regulator

import "sync/atomic"

type Regulator struct {
	concurrency int
	sem         chan bool
	jobIndex    int32
	err         error
}

type RegulatorError struct {
	JobIndex int32
	Message  string
}

func (l RegulatorError) Error() string {
	return l.Message
}

func NewRegulator(concurrency int) *Regulator {
	return &Regulator{
		concurrency: concurrency,
		sem:         make(chan bool, concurrency),
	}
}

func (regulator *Regulator) Execute(job func() error) {
	index := atomic.AddInt32(&regulator.jobIndex, 1)
	regulator.sem <- true
	go func(jobIndex int32) {
		defer func() { <-regulator.sem }()
		if regulator.err != nil {
			return
		}
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
	return regulator.err
}