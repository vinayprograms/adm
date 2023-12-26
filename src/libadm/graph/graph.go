package graph

import (
	"errors"
	"libadm/model"
)

type Graph struct {
	Models             map[string]*model.Model
	RealitySuccessors  map[string]interface{}
	UnmitigatedAttacks map[string]interface{}
}

// A graph should atleast have one model
func (g *Graph) Init() {
	g.Models = make(map[string]*model.Model)
	g.RealitySuccessors = make(map[string]interface{})
	g.UnmitigatedAttacks = make(map[string]interface{})
}

func (g *Graph) AddModel(m *model.Model) error {
	if m == nil {
		return errors.New("cannot add a null model to graph")
	}

	// Apply ADM rules first to make inter-model connections
	g.applyADMRules(m)

	// Add model to graph
	g.ConnectToReality(m)
	g.CollectUnmitigatedAttacks(m)

	g.Models[m.Title] = m

	return nil
}

func (g *Graph) ContainsModel(title string) bool {
	if g.Models[title] != nil {
		return true
	} else {
		return false
	}
}

func (g *Graph) GetModel(title string) *model.Model {
	return g.Models[title]
}

// TODO: Replacement and deletion requires logic to disconnect model contents from the graph
/*
func (g *Graph) ReplaceModel(title string, m *model.Model) error {
	if g.models[title] != nil && title == m.Title {
		g.models[title] = m
	} else {
		return errors.New("graph doesn't contain a model with title '" + title + "'")
	}

	return nil
}

func (g *Graph) DeleteModel(title string) error {
	if g.models[title] != nil {
		delete(g.models, title)
	} else {
		return errors.New("graph doesn't contain a model with title '" + title + "'")
	}

	return nil
}*/

// entities that don't have a predecessor are connected to Reality.
func (g *Graph) ConnectToReality(m *model.Model) error {
	for _, a := range m.Assumptions {
		g.connectAssumptionToReality(a)
	}
	for _, p := range m.Policies {
		g.connectPolicyToReality(p)
	}
	for _, a := range m.Attacks {
		g.connectAttackToReality(a)
	}
	for _, d := range m.Defenses {
		g.connectDefenseToReality(d)
	}

	return nil
}

func (g *Graph) CollectUnmitigatedAttacks(m *model.Model) {
	for title, attack := range m.Attacks {
		if len(attack.PreConditions) == 0 && len(attack.Actions) == 0 && len(attack.Results) == 0 { // Empty specs cannot be considered
			continue
		} else if _, present := g.UnmitigatedAttacks[title]; present { // if already listed
			continue
		}

		isUnmitigated := false

		// Attacks with "@success" tag are considered unmitigated
		if contains("@success", attack.Tags) {
			isUnmitigated = true
		} else if !isMitigated(attack, m.Defenses, m.Policies) { // if unmitigated within the model
			if len(g.Models) > 0 {
				for _, other := range g.Models {
					if !isMitigated(attack, other.Defenses, other.Policies) { // if unmitigated across models
						isUnmitigated = true
					}
				}
			} else {
				isUnmitigated = true
			}
		}

		if isUnmitigated {
			g.UnmitigatedAttacks[title] = attack
		}
	}
}

////////////////////////////////////////
// Internal functions

// Apply ADM rules across models
func (g *Graph) applyADMRules(m *model.Model) {
	for _, existing := range g.Models {
		//////////
		// There are 17 ways attacks and defenses can be connected.
		//  Here 'T' represents tag based relationship (always attack -> defense).
		//      NEW MODEL     |    EXISTING MODEL
		// -------------------|------------------------
		//  Attack					<--->  Attack             (1)
		//  Attack					---->  Defense            (2)
		//  Attack					--T->  Defense            (3)
		//  Attack					---->  Policy Defense     (4)
		//  Attack					--T->  Policy             (5)
		//  Attack					--T->  Policy Defense     (6)
		//  Attack					<----  Defense            (7)
		//  Attack					<----  Policy Defense     (8)
		//
		//  Defense					---->  Attack             (9)
		//  Defense					<--->  Defense            (10)
		//  Defense					<--->  Policy Defense     (11)
		//  Defense					<----  Attack             (12)
		//  Defense					<-T--  Attack             (13)
		//
		//  Policy Defense	---->  Attack             (14)
		//  Policy Defense	<--->  Defense            (15)
		//  Policy Defense	<--->  Policy Defense     (16)
		//  Policy Defense	<----  Attack             (17)
		//  Policy Defense	<-T--  Attack             (18)
		//  Policy          <-T--  Attack             (19)

		// new attacks
		for _, attack := range m.Attacks {
			model.ChainAttacks(attack, existing.Attacks)             // (1)
			model.ConnectAttackToDefenses(attack, existing.Defenses) // (2),(7)
			for _, p := range existing.Policies {
				model.ConnectAttackToDefenses(attack, p.Defenses) // (4),(8)
				p.ConnectAttackToPolicyByTags(attack)             // (5)
			}
		}
		model.MultiConnectAttacksAndDefensesByTag(m.Attacks, existing.Defenses) // (3)
		for _, p := range existing.Policies {
			model.MultiConnectAttacksAndDefensesByTag(m.Attacks, p.Defenses) // (6)
		}

		// new defenses
		for _, defense := range m.Defenses {
			model.ConnectDefenseToAttacks(defense, existing.Attacks) // (9),(12)
			model.ChainDefenses(defense, existing.Defenses)          // (10)
			for _, p := range existing.Policies {
				model.ChainDefenses(defense, p.Defenses) // (11)
			}
		}
		model.MultiConnectAttacksAndDefensesByTag(existing.Attacks, m.Defenses) // (13)

		// new policies
		for _, p := range m.Policies {
			for _, defense := range p.Defenses {
				model.ConnectDefenseToAttacks(defense, existing.Attacks) // (14), (15)
				model.ChainDefenses(defense, existing.Defenses)          // (17)
				for _, ep := range existing.Policies {
					model.ChainDefenses(defense, ep.Defenses) // (16)
				}
			}
			for _, attack := range existing.Attacks {
				p.ConnectAttackToPolicyByTags(attack) // (19)
			}
			model.MultiConnectAttacksAndDefensesByTag(existing.Attacks, p.Defenses) // (18)
		}
	}
}

// All assumptions & its preconditions are connected to reality
func (g *Graph) connectAssumptionToReality(a *model.Assumption) {
	g.RealitySuccessors[a.Title] = a
}

// Connect all preconditions of a policy to 'reality' along with defense that don't have any predecessors
func (g *Graph) connectPolicyToReality(p *model.Policy) {
	// Connect defenses that don't have predecessor to reality
	for _, d := range p.Defenses {
		g.connectDefenseToReality(d)
	}
	// All assumptions & its preconditions are connected to reality
	for _, a := range p.Assumptions {
		g.connectAssumptionToReality(a)
	}
}

// Connect attack to 'reality' if it has no predecessors or if it doesn't thwart anything.
func (g *Graph) connectAttackToReality(a *model.Attack) {
	hasNoPredecessor := true
	for _, pre := range a.PreConditions {
		if pre.Item != nil {
			hasNoPredecessor = false
			break
		}
	}

	if hasNoPredecessor {
		for _, a := range a.Actions {
			if a.Item != nil {
				hasNoPredecessor = false
				break
			}
		}
	}

	if hasNoPredecessor {
		g.RealitySuccessors[a.Title] = a
	}
}

// Connect defense to 'reality' if it has no predecessors or if it doesn't mitigate any attack.
func (g *Graph) connectDefenseToReality(d *model.Defense) {
	hasNoPredecessor := true
	for _, pre := range d.PreConditions {
		if pre.Item != nil {
			hasNoPredecessor = false
			break
		}
	}

	if hasNoPredecessor {
		for _, a := range d.Actions {
			if a.Item != nil {
				hasNoPredecessor = false
				break
			}
		}
	}

	if hasNoPredecessor {
		g.RealitySuccessors[d.Title] = d
	}
}

// Check if one of the defenses mitigate the given attack
func isMitigated(attack *model.Attack, defenses map[string]*model.Defense, policies map[string]*model.Policy) bool {
	// check with model defenses
	for _, defense := range defenses {
		if len(defense.PreConditions) == 0 && len(defense.Actions) == 0 && len(defense.Results) == 0 { // empty defense spec. cannot be considered
			continue
		}
		for _, defenseAction := range defense.Actions {
			if contains(attack, defenseAction.Item) {
				return true
			}
		}
		for _, taggedAttacks := range defense.Tags {
			if contains(attack, taggedAttacks) {
				return true
			}
		}
	}

	// check with policy defenses
	for _, policy := range policies {
		for _, defense := range policy.Defenses {
			if len(defense.PreConditions) == 0 && len(defense.Actions) == 0 && len(defense.Results) == 0 { // empty defense spec. cannot be considered
				continue
			}
			for _, defenseAction := range defense.Actions {
				if contains(attack, defenseAction.Item) {
					return true
				}
			}
			for _, taggedAttacks := range defense.Tags {
				if contains(attack, taggedAttacks) {
					return true
				}
			}
		}
	}

	return false
}

////////////////////////////////////////
// Helper functions

// Check array membership
func contains[T comparable](item T, array []T) bool {
	for _, x := range array {
		if item == x {
			return true
		}
	}

	return false
}
