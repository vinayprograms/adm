package model

import (
	"errors"

	"github.com/cucumber/messages-go/v16"
)

type Policy struct {
	Title string
	Tags map[string][]*Attack
	Assumptions map[string]*Assumption
	Defenses map[string]*Defense
}

////////////////////////////////////////
// Policy specific functions

func (p *Policy) Init(r *messages.Rule) error {

	if p.Assumptions == nil		{ p.Assumptions = make(map[string]*Assumption) }
	if p.Defenses == nil 			{ p.Defenses = make(map[string]*Defense) }
	if p.Tags == nil					{ p.Tags = make(map[string][]*Attack) }

	p.Title = r.Name

	for _, c := range r.Children {
		if c.Background != nil 		{ 
			err := p.initAssumption(c.Background) 
			if err != nil { return err }
		} else if c.Scenario != nil { 
			err := p.initDefense(c.Scenario) 
			if err != nil { return err }
		}
	}

	p.AddTags(r.Tags)

	return nil
}

func (p *Policy) ConnectAttackToPolicyByTags(a *Attack) {
	for _, tag := range a.Tags {
		if _, present := p.Tags[tag]; present {
			p.Tags[tag] = append(p.Tags[tag], a)
		}
	}
}

////////////////////////////////////////
// Internal functions used by 'Init' function

func (p *Policy) initAssumption(b *messages.Background) error {
	a := Assumption{}
	err := a.Init(b)
	if err != nil {
		// This code should be unreachable. But, if it does, contact author!
		return err
	}
	if _, present := p.Assumptions[a.Title]; present {
		// This error will not be encountered under normal conditions because 
		// gherkin parser would have failed when it sees 2 assumption sections.
		// This piece of code is added only for cases where gherkin structures 
		// are filled in code before being passed to this structure.
		return errors.New("assumption - '" + a.Title + "' is already part of this model")
	}
	p.Assumptions[a.Title] = &a

	return nil
}

func (p *Policy) initDefense(s *messages.Scenario) error {
	d := Defense{}
	err := d.Init(s)
	if err != nil {
		// This code should be unreachable. But, if it does, contact author!
		return err
	}
	if _, present := p.Defenses[d.Title]; present {
		return errors.New("defense - '" + d.Title + "' is already part of this model")
	}

	ChainDefenses(&d, p.Defenses)

	p.Defenses[d.Title] = &d

	return nil
}

func (p *Policy) AddTags(t []*messages.Tag) {
	for _, tag := range t {
		if _, present := p.Tags[tag.Name]; !present { // if this is a new tag
			p.Tags[tag.Name] = nil
		}
	}
}