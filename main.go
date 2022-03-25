package main

import (
	"context"
	"fmt"
)

func main() {
	o := Onion{
		NewStep(
			func(ctx context.Context) error {
				fmt.Println("forward 1")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("backward 1")
				return nil
			},
		),
		NewStep(
			func(ctx context.Context) error {
				fmt.Println("forward 2")
				return nil
			},
			func(ctx context.Context) error {
				fmt.Println("backward 2")
				return nil
			},
		),
		NewStep(
			func(ctx context.Context) error {
				return fmt.Errorf("ack")
				// fmt.Println("forward 3")
				// return nil
			},
			func(ctx context.Context) error {
				fmt.Println("backward 3")
				return nil
			},
		),
	}

	_ = o.Run(context.Background())
}

type Onion []Step

func (o Onion) Run(ctx context.Context) error {
	var err error
	for i, step := range o {
		err = step.Forward(ctx)
		if err != nil {
			err = o.Unwind(ctx, i)
			if err != nil {
				panic(err)
			}
		}

	}
	return nil
}

func (o Onion) Unwind(ctx context.Context, start int) error {
	var err error
	for i := start; i >= 0; i-- {
		err = o[i].Backward(ctx)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

type Step interface {
	Forward(context.Context) error
	Backward(context.Context) error
}

type step struct {
	forward  func(context.Context) error
	backward func(context.Context) error
}

func (s *step) Forward(ctx context.Context) error {
	return s.forward(ctx)
}

func (s *step) Backward(ctx context.Context) error {
	return s.backward(ctx)
}

func NewStep(forward func(context.Context) error, backward func(context.Context) error) Step {
	s := &step{
		forward:  forward,
		backward: backward,
	}
	return s
}
