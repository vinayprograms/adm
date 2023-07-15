package model

import (
	"errors"
	"strings"

	"github.com/cucumber/messages-go/v16"
)

type Defense struct {
	Title 					string
	Tags 						map[string][]*Attack
	PreConditions 	map[string]*ModelLink[[]*Defense]
	Actions 				map[string]*ModelLink[[]*Attack]
	Results 				map[string]*Step
}

////////////////////////////////////////
// Defense specific functions

func (d *Defense) Init(s *messages.Scenario) error {
	
	d.initializeMembers()

	if (s.Keyword != "Defense") {
		return errors.New("Expected 'Defense', got '" + strings.TrimSpace(s.Keyword) + "'")
	}
	
	d.Title = s.Name
	d.AddTags(s.Tags)

	var currentStepKey string
	for _, step := range s.Steps {
		item := Step{}
		item.Init(step)
		switch(step.Keyword) {
		case "Given":
			currentStepKey = step.Keyword
			if _, present := d.PreConditions[item.Statement]; present {
				return errors.New("precondition - '" + item.Statement + "' is already part of this defense")
			}
			d.PreConditions[item.Statement] = &ModelLink[[]*Defense]{Step:&item}
		case "When":
			currentStepKey = step.Keyword
			if _, present := d.Actions[item.Statement]; present {
				return errors.New("action - '" + item.Statement + "' is already part of this defense")
			}
			d.Actions[item.Statement] = &ModelLink[[]*Attack]{Step:&item}
		case "Then":
			currentStepKey = step.Keyword
			if _, present := d.Results[item.Statement]; present {
				return errors.New("result - '" + item.Statement + "' is already part of this defense")
			}
			d.Results[item.Statement] = &item
		case "And":
			switch currentStepKey {
			case "Given":
				if _, present := d.PreConditions[item.Statement]; present {
					return errors.New("precondition - '" + item.Statement + "' is already part of this defense")
				}
				d.PreConditions[item.Statement] = &ModelLink[[]*Defense]{Step:&item}
			case "When":
				if _, present := d.Actions[item.Statement]; present {
					return errors.New("action - '" + item.Statement + "' is already part of this defense")
				}
				d.Actions[item.Statement] = &ModelLink[[]*Attack]{Step:&item}
			case "Then":
				if _, present := d.Results[item.Statement]; present {
					return errors.New("result - '" + item.Statement + "' is already part of this defense")
				}
				d.Results[item.Statement] = &item
			}
		default:
			return errors.New("unsupported keyword - '" + step.Keyword + "'")
		}
	}

	return nil
}

////////////////////////////////////////
// Helper functions

func (d *Defense) initializeMembers() {
	if d.Tags == nil					{ d.Tags = make(map[string][]*Attack) }
	if d.PreConditions == nil	{ d.PreConditions = make(map[string]*ModelLink[[]*Defense]) }
	if d.Actions == nil 			{ d.Actions = make(map[string]*ModelLink[[]*Attack]) }
	if d.Results == nil 			{ d.Results = make(map[string]*Step) }
}

func (d *Defense) AddTags(t []*messages.Tag) {
	for _, tag := range t {
		d.Tags[tag.Name] = nil
	}
}