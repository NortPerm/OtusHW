package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

type Pool struct {
	stages   []Stage
	in, done In
	out      Bi
}

func NewPool(stages []Stage, in In, done In) *Pool {
	return &Pool{
		stages: stages,
		in:     in,
		done:   done,
		out:    make(Bi),
	}
}

func (p *Pool) execute() Out {
	go func() {
		defer close(p.out)
		for {
			select {
			case value, ok := <-p.in:
				if ok {
					p.out <- value
				} else {
					return
				}
			case <-p.done:
				return
			}
		}
	}()

	return p.out
}

func (p *Pool) composeStages() {
	for _, stage := range p.stages {
		p.in = stage(p.in)
	}
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	p := NewPool(stages, in, done)
	p.composeStages()
	return p.execute()
}
