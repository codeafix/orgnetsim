package sim

import "fmt"

//NetworkOptions contains information about how the network is set up for the simulation
type NetworkOptions struct {
	LinkTeamPeers    bool     `json:"linkTeamPeers"`
	LinkedTeamList   []string `json:"linkedTeamList"`
	EvangelistList   []string `json:"evangelistList"`
	LoneEvangelist   []string `json:"loneEvangelist"`
	InitColors       []Color  `json:"initColors"`
	MaxColors        int      `json:"maxColors"`
	AgentsWithMemory bool     `json:"agentsWithMemory"`
}

//LinkTeamPeers links all Agents related to the same parent node to each other.
//Turns a strictly hierarchical network in to a more realistic communication network.
func LinkTeamPeers(rm RelationshipMgr) error {
	teams := map[string][]string{}

	for _, link := range rm.Links() {
		_, exists := teams[link.Agent1ID]
		if !exists {
			teams[link.Agent1ID] = []string{link.Agent2ID}
		} else {
			teams[link.Agent1ID] = append(teams[link.Agent1ID], link.Agent2ID)
		}
	}
	var err error
	for _, teamMembers := range teams {
		teamSize := len(teamMembers)
		for i := 0; i < teamSize; i++ {
			for j := i + 1; j < teamSize; j++ {
				err = addLink(rm, teamMembers[i], teamMembers[j])
				if err != nil {
					return fmt.Errorf("Error linking team peers: %s", err.Error())
				}
			}
		}
	}
	return nil
}

//AddEvangelists sets a list of individuals to Blue and increases their susceptibility
//so that they cannot be influenced by another Agent
func AddEvangelists(rm RelationshipMgr, o NetworkOptions) error {
	eTeamSize := len(o.EvangelistList)
	if eTeamSize > 0 {
		for i := 0; i < eTeamSize; i++ {
			id1 := o.EvangelistList[i]
			agent := rm.GetAgentByID(id1)
			if agent != nil {
				agent.State().Color = Blue
				agent.State().Susceptability = 5.0
			} else {
				return fmt.Errorf("Unrecognised entry in Evangelist List: %s", id1)
			}
		}
	}
	return nil
}

//LinkTeams creates links between a specified set of individuals from across teams in the network
func LinkTeams(rm RelationshipMgr, o NetworkOptions) error {
	lTeamSize := len(o.LinkedTeamList)
	var err error
	if lTeamSize > 0 {
		for i := 0; i < lTeamSize; i++ {
			id1 := o.LinkedTeamList[i]
			for j := i + 1; j < lTeamSize; j++ {
				id2 := o.LinkedTeamList[j]
				err = addLink(rm, id1, id2)
				if err != nil {
					return fmt.Errorf("Unrecognised entry in LinkTeams List: %s", err.Error())
				}
			}
		}
	}
	return nil
}

//AddLoneEvangelist links a single Agent to a list of other Agents across the Network.
//The first agent in the LoneEvangelist list is the Evangelist and all subsequent Agents
//are connected to her. If the Lone Evangelist Id does not exist in the network she is
//created.
func AddLoneEvangelist(rm RelationshipMgr, o NetworkOptions) error {
	leTeamSize := len(o.LoneEvangelist)
	var err error
	if leTeamSize > 0 {
		agent := rm.GetAgentByID(o.LoneEvangelist[0])
		if agent == nil {
			agent = GenerateRandomAgent(o.LoneEvangelist[0], o.InitColors, o.AgentsWithMemory)
			rm.AddAgent(agent)
			rm.(*Network).PopulateMaps()
		}
		agent.State().Susceptability = 5.0
		agent.State().Color = Blue
		for i := 1; i < leTeamSize; i++ {
			err = addLink(rm, agent.Identifier(), o.LoneEvangelist[i])
			if err != nil {
				return fmt.Errorf("Unrecognised entry in LoneEvangelist List: %s", err.Error())
			}
		}
	}
	return nil
}

//Convenience method to check agents exist before trying to add a link
func addLink(rm RelationshipMgr, id1 string, id2 string) error {
	a1 := rm.GetAgentByID(id1)
	if a1 == nil {
		return fmt.Errorf("Unrecognised Agent Id '%s'", id1)
	}
	a2 := rm.GetAgentByID(id2)
	if a2 == nil {
		return fmt.Errorf("Unrecognised Agent Id '%s'", id2)
	}
	rm.AddLink(a1, a2)
	return nil
}

//CloneModify clones the agents and links in the passed RelationshipMgr into a new
//RelationshipMgr changing the Agent type and initial colors of all Agents on the Network,
//then it modifies the links as specified in the passed Options struct.
func CloneModify(rm RelationshipMgr, o NetworkOptions) (RelationshipMgr, error) {
	ret, err := cloneNetwork(rm, o)
	if err != nil {
		return ret, err
	}
	err = ModifyNetwork(ret, o)
	if err != nil {
		return ret, err
	}
	err = ret.PopulateMaps()
	return ret, err
}

//ModifyNetwork takes a RelationshipMgr as input and adds links as specified in the passed Options
//struct. Note this method will ignore the InitColors and AgentsWithMemory options because
//a new set of Agents require to be generated in order to set these options. To do that use
//the CloneModify function instead.
func ModifyNetwork(rm RelationshipMgr, o NetworkOptions) error {
	rm.SetMaxColors(o.MaxColors)
	var err error
	if o.LinkTeamPeers {
		err = LinkTeamPeers(rm)
		if err != nil {
			return err
		}
	}
	err = AddEvangelists(rm, o)
	if err != nil {
		return err
	}
	err = LinkTeams(rm, o)
	if err != nil {
		return err
	}
	err = AddLoneEvangelist(rm, o)
	return err
}

//cloneNetwork creates a new network and creates copies of the nodes and links in it from the passed network
//The new Agents will be generated according to the settings in the passed Options struct
func cloneNetwork(rm RelationshipMgr, o NetworkOptions) (*Network, error) {
	ret := &Network{}
	for _, agent := range rm.Agents() {
		clone := GenerateRandomAgent(agent.Identifier(), o.InitColors, o.AgentsWithMemory)
		ret.AddAgent(clone)
	}
	ret.PopulateMaps()
	for _, link := range rm.Links() {
		err := addLink(ret, link.Agent1ID, link.Agent2ID)
		if err != nil {
			return nil, err
		}
	}
	ret.PopulateMaps()
	return ret, nil
}
