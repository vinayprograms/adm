package model

import (
	"errors"
	"strings"

	"github.com/cucumber/messages-go/v16"
)

type Assumption struct {
	Title         string
	PreConditions map[string]*Step
}

func (a *Assumption) Init(b *messages.Background) error {
	if a.PreConditions == nil {
		a.PreConditions = make(map[string]*Step)
	}
	if b == nil {
		return errors.New("expected 'Assumption' spec. Got 'nil'")
	}

	a.Title = b.Name
	for _, s := range b.Steps {
		if s.Keyword != "Given" && s.Keyword != "And" {
			return errors.New("Unexpected keyword - '" + strings.TrimSpace(s.Keyword) + "' in Assumption specification")
		}
		step := Step{}
		err := step.Init(s)
		if err != nil {
			// This code should be unreachable. But, if it does, contact author!
			return err
		}

		if _, present := a.PreConditions[step.Statement]; present {
			return errors.New("precondition - '" + step.Statement + "' is already part of this assumption")
		}
		a.PreConditions[step.Statement] = &step
	}

	return nil
}
