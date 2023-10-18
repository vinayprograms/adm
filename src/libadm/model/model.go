package model

import (
	"errors"

	"github.com/cucumber/messages-go/v16"
)

type Model struct {
	Title string
	Tags  []string

	Actors      []*ModelActor
	Assumptions map[string]*Assumption
	Policies    map[string]*Policy
	Attacks     map[string]*Attack
	Defenses    map[string]*Defense
}

func (m *Model) Init(f *messages.Feature) error {

	if f == nil {
		return errors.New("cannot initialize a null model object")
	}

	if m.Assumptions == nil {
		m.Assumptions = make(map[string]*Assumption)
	}
	if m.Policies == nil {
		m.Policies = make(map[string]*Policy)
	}
	if m.Attacks == nil {
		m.Attacks = make(map[string]*Attack)
	}
	if m.Defenses == nil {
		m.Defenses = make(map[string]*Defense)
	}

	// Process gherkin content
	m.Title = f.Name
	m.ParseAndAddActors(f.Description)

	for _, child := range f.Children {
		if child.Rule != nil {
			_, err := m.initPolicy(child.Rule)
			if err != nil {
				return err
			}
		} else if child.Background != nil {
			_, err := m.initAssumption(child.Background)
			if err != nil {
				return err
			}
		} else if child.Scenario != nil {
			switch child.Scenario.Keyword {
			case "Attack":
				_, err := m.initAttack(child.Scenario)
				if err != nil {
					return err
				}
			case "Defense":
				_, err := m.initDefense(child.Scenario)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (m *Model) ParseAndAddActors(modelDescription string) {
	m.Actors = append(m.Actors, parseActors(modelDescription)...)
}

func (m *Model) AddActor(p *ModelActor) {
	m.Actors = append(m.Actors, p)
}

////////////////////////////////////////
// Internal functions used by 'Init' function

func (m *Model) initPolicy(r *messages.Rule) (*Policy, error) {
	p := Policy{}
	err := p.Init(r)
	if err != nil {
		// This code should be unreachable. But, if it does, contact author!
		return nil, err
	}
	if _, present := m.Policies[p.Title]; present {
		return nil, errors.New("policy - '" + p.Title + "' is already part of this model")
	}
	m.Policies[p.Title] = &p

	// There are 6 ways policy can be connected to other parts of the model
	//
	// Attack  ----> Policy Defense	[mitigation]             (1)
	// Attack  ----> Policy Defense	[response]               (2)
	// Attack  --T-> Policy         [mitigation by tag]      (3)
	// Attack  --T-> Policy Defense	[mitigation by tag]      (4)
	// Attack  <---- Policy Defense [thwart]                 (5)
	// Defense <---> Policy Defense [chain]                  (6)

	for _, a := range m.Attacks {
		ConnectAttackToDefenses(a, p.Defenses) // (1),(2),(5)
		p.ConnectAttackToPolicyByTags(a)       // (3)
	}
	MultiConnectAttacksAndDefensesByTag(m.Attacks, p.Defenses) // (4)
	for _, d := range p.Defenses {
		ChainDefenses(d, m.Defenses) // (5)
	}

	return &p, nil
}

func (m *Model) initAssumption(b *messages.Background) (*Assumption, error) {
	a := Assumption{}
	err := a.Init(b)
	if err != nil {
		return nil, err
	}
	if _, present := m.Assumptions[a.Title]; present {
		// This error will not be encountered under normal conditions because
		// gherkin parser would have failed when it sees 2 assumption sections.
		// This piece of code is added only for cases where gherkin structures
		// are filled in code before being passed to this model.
		return nil, errors.New("assumption - '" + a.Title + "' is already part of this model")
	}
	m.Assumptions[a.Title] = &a

	return &a, nil
}

func (m *Model) initAttack(s *messages.Scenario) (*Attack, error) {
	a := Attack{}
	err := a.Init(s)
	if err != nil {
		// This code should be unreachable. But, if it does, contact author!
		return nil, err
	}
	if _, present := m.Attacks[a.Title]; present {
		return nil, errors.New("attack - '" + a.Title + "' is already part of this model")
	}

	// There are 10 ways attack can be connected to other parts of the model
	//
	// Attack <---> Attack         [chain]                  (1)
	// Attack ----> Defense        [mitigation]             (2)
	// Attack ----> Policy Defense [mitigation]             (3)
	// Attack ----> Defense        [response]               (4)
	// Attack ----> Policy Defense [response]               (5)
	// Attack --T-> Defense        [mitigation by tag]      (6)
	// Attack --T-> Policy         [mitigation by tag]      (7)
	// Attack --T-> Policy Defense [mitigation by tag]      (8)
	// Attack <---- Defense        [thwart]                 (9)
	// Attack <---- Policy Defense [thwart]                 (10)

	ChainAttacks(&a, m.Attacks)                             // (1)
	ConnectAttackToDefenses(&a, m.Defenses)                 // (2),(4),(9)
	ConnectSingleAttackToMultiDefensesByTag(&a, m.Defenses) // (6)
	for _, p := range m.Policies {
		p.ConnectAttackToPolicyByTags(&a)                       // (7)
		ConnectAttackToDefenses(&a, p.Defenses)                 // (3),(5),(10)
		ConnectSingleAttackToMultiDefensesByTag(&a, p.Defenses) // (8)
	}

	m.Attacks[a.Title] = &a

	return &a, nil
}

func (m *Model) initDefense(s *messages.Scenario) (*Defense, error) {
	d := Defense{}
	err := d.Init(s)
	if err != nil {
		// This code should be unreachable. But, if it does, contact author!
		return nil, err
	}
	if _, present := m.Defenses[d.Title]; present {
		return nil, errors.New("defense - '" + d.Title + "' is already part of this model")
	}

	// There are 6 ways defense can be connected to other parts of the model
	//
	// Defense ----> Attack         [thwart]                 (1)
	// Defense <---> Defense        [chain]                  (2)
	// Defense <---> Policy Defense [chain]                  (3)
	// Defense <---- Attack         [mitigation]             (4)
	// Defense <---- Attack         [response]               (5)
	// Defense <-T-- Attack         [mitigation by tag]      (6)

	ConnectDefenseToAttacks(&d, m.Attacks) // (1),(4),(5)
	ChainDefenses(&d, m.Defenses)          // (2)
	for _, p := range m.Policies {
		ChainDefenses(&d, p.Defenses) // (3)
	}
	ConnectMultiAttacksToSingleDefensByTag(m.Attacks, &d) // (6)

	m.Defenses[d.Title] = &d

	return &d, nil
}
