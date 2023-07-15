package model

import (
	"errors"
	"strings"

	"github.com/cucumber/messages-go/v16"
)

type Attack struct {
	Title 					string
	Tags 						[]string
	PreConditions 	map[string]*ModelLink[[]*Attack]
	Actions 				map[string]*ModelLink[[]*Defense]
	Results 				map[string]*Step
}

////////////////////////////////////////
// Attack specific functions

func (a *Attack) Init(s *messages.Scenario) error {

	a.initializeMembers()

	if (s.Keyword != "Attack") {
		return errors.New("Expected 'Attack', got '" + strings.TrimSpace(s.Keyword) + "'")
	}
	
	a.Title = s.Name
	a.AddTags(s.Tags)

	var currentStepKey string
	for _, step := range s.Steps {
		item := Step{}
		item.Init(step)

		switch(step.Keyword) {
		case "Given":
			currentStepKey = step.Keyword
			if _, present := a.PreConditions[item.Statement]; present {
				return errors.New("precondition - '" + item.Statement + "' is already part of this attack")
			}
			a.PreConditions[item.Statement] = &ModelLink[[]*Attack]{Step:&item}
		
		case "When":
			currentStepKey = step.Keyword
			if _, present := a.Actions[item.Statement]; present {
				return errors.New("action - '" + item.Statement + "' is already part of this attack")
			}
			a.Actions[item.Statement] = &ModelLink[[]*Defense]{Step:&item}
		
		case "Then":
			currentStepKey = step.Keyword
			if _, present := a.Results[item.Statement]; present {
				return errors.New("result - '" + item.Statement + "' is already part of this attack")
			}
			a.Results[item.Statement] = &item
		
		case "And":
			switch currentStepKey {
			case "Given":
				if _, present := a.PreConditions[item.Statement]; present {
					return errors.New("precondition - '" + item.Statement + "' is already part of this attack")
				}
				a.PreConditions[item.Statement] = &ModelLink[[]*Attack]{Step:&item}
			case "When":
				if _, present := a.Actions[item.Statement]; present {
					return errors.New("action - '" + item.Statement + "' is already part of this attack")
				}
				a.Actions[item.Statement] = &ModelLink[[]*Defense]{Step:&item}
			case "Then":
				if _, present := a.Results[item.Statement]; present {
					return errors.New("result - '" + item.Statement + "' is already part of this attack")
				}
				a.Results[item.Statement] = &item
			}
		default:
			return errors.New("unsupported keyword - '" + step.Keyword + "'")
		}
	}

	return nil
}

////////////////////////////////////////
// Helper functions

func (a *Attack) initializeMembers() {
	if a.PreConditions == nil	{ a.PreConditions = make(map[string]*ModelLink[[]*Attack]) }
	if a.Actions == nil 			{ a.Actions = make(map[string]*ModelLink[[]*Defense]) }
	if a.Results == nil 			{ a.Results = make(map[string]*Step) }
}

func (a *Attack) AddTags(t []*messages.Tag) {
	for _, tag := range t {
		a.Tags = append(a.Tags, tag.Name)
	}
}