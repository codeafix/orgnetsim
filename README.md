# orgnetsim [![Build Status](https://travis-ci.org/codeafix/orgnetsim.svg?branch=master)](https://travis-ci.org/codeafix/orgnetsim) [![Coverage Status](http://codecov.io/github/codeafix/orgnetsim/coverage.svg?branch=master)](http://codecov.io/github/codeafix/orgnetsim?branch=master) [![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/github.com/codeafix/orgnetsim) [![MIT](https://img.shields.io/npm/l/express.svg)](https://github.com/codeafix/orgnetsim/blob/master/LICENSE)
A simulator for Organisational Networks

The simulator is created from a Network of Agents. The Network itself can be any arbitrary graph and contains a collection of Agents and a collection of links between those Agents. The simulator uses Colors to represent competing ideas on the Network. The default Color for an Agent is Grey. During a simulation Agents interact and decide whether or not to update their Color.

## Overview
The simulator runs over multiple iterations. Each iteration starts with every Agent in the Network attempting to contact another Agent in the list of Agents that are directly related to it (joined by a single link). The list of Agents related to it is obtained from the Network, and this list is always returned in a random order. Each Agent can only accept a single Mail in its Mail queue, so during this process each Agent will iterate over all other Agents directly related to it until it finds an Agent that has an empty Mail queue and can accept its Mail.

Once all Agents have completed the process of trying to send a Mail the next phase begins. Now each Agent reads the Mail it has in its Mail queue if any. The Mail contains the Identifier of the Agent that sent it, and the receiving Agent uses this to look up the sending Agent from the Network. Each Agent has three properties: Influence, Susceptibility, and Contrariness. The receiving Agent compares its properties to that of the sending Agent and uses a simple algorithm to decide how to update its Color. If the sending Agent has a Influence higher than the receiving Agent's Susceptibility then the receiving Agent will update its Color. If the receiving Agent's Contrariness is higher than the sending Agent's Influence then the receiving Agent will update to a random Color different from its previous Color, and from the Color of the sending Agent. If the receiving Agent's Contrariness is lower than the sending Agent's Influence then the receiving Agent updates its Color to the same as the sending Agent.

The RunSim function in orgnetsim.go is the entry point for a simulation. A RelationshipMgr (the interface to a Network) is passed into this function along with a number of iterations. As each iteration is performed, the Agents held within the RelationshipMgr are updated, and a log is taken of the number of Agents with each Color, and the number of "conversations" that happened in the iteration. At the end of the simulation, two slices are returned. The first slice is a two dimensional slice. The first dimension is Color, the second dimension is the number of iterations. Each element contains the count of the number of Agents with the given Color on the specified iteration. The second slice contains a count of the number of "conversations" that occured between all agents for each iteration.

After a simulation run has completed the Agents and Links can be accessed from the RelationshipMgr. Each Agent keeps a count of the number of times it updated its Color. Each Link keeps a count of the number of conversations that happen across that link. These can be accessed like this:
```
var n RelationshipMgr

for _, a := range n.Agents() {
	agc := a.State().ChangeCount
    ...
}

for _, l := range n.Links() {
	ls := l.Strength
    ...
}
```

## Getting started
An example run of a simulation is provided as a test in orgnetsim_test.go
```
func TestRunSim(t *testing.T) {

	s := HierarchySpec{
		4,                  //Levels
		5,                  //TeamSize
		3,                  //TeamLinkLevel
		true,               //LinkTeamPeers
		true,               //LinkTeams
		[]Color{Grey, Red}, //InitColors
		false,              //EvangelistAgents
		false,              //LoneEvangelist
		false,              //AgentsWithMemory
	}

	n, err := GenerateHierarchy(s)
	AssertSuccess(t, err)

	colors, conversations := RunSim(n, 500)
	WriteOutput(t, s, n, colors, conversations)
}
```
This simulation uses the GenerateHierarchy function to generate a hierarchal network, and then passes that network into the RunSim function. The simulation is run for 500 iterations and then the slices returned from the simulation, the change count of each node, and the strength of each link are written into out.csv along with the list of parameters provided in the HierarchySpec struct used to generate the network. This is a convenient output for loading the results into a spreadsheet so that you can plot graphs of the number of Agents with each Color.

## RelationshipMgr aka the Network
In the examples networks with regular features have been generated automatically using the networkgenerator. However the Network struct (implements RelationshipMgr) has been designed so that it can be created from a JSON description of the network. Here is an example:
```
json := '{"nodes":[{"id":"id_1"},{"id":"id_2"},{"id":"id_3"}],"links":[{"agent1Id":"id_1","agent2Id":"id_2"},{"agent1Id":"id_1","agent2Id":"id_3"}]}'
n, err := NewNetwork(json)
	
serJSON := n.Serialise()
```
In JSON each Agent in the "nodes" array has a "type" property which specifies which type of Agent to use (Agent or AgentWithMemory). If this property is not specified then it will default to Agent. If there is more than one Agent with the same "id" value then the last Agent encountered in the list of "nodes" will be the only one present in the network when unmarshalling is complete. Similarly if there are duplicate links (more than one link between the same pair of nodes) only the last link will be present in the network.

The network struct has been designed so that any organisational network (such as your team/business unit org chart etc) can be used as the basis of a simulation. Teams need not be regularly sized, any graph can be specified as individual nodes and links between them. After the simulation is complete it will be possible to marshal the Network back to json. When the network is returned as json the Agent states including the ChangeCount and Link strengths are recorded. This is done to allow simulations to be run for many iterations or even stepped through in smaller increments.

##Options to control a simulation
The number of possible competing ideas in a simulation is controlled by the MaxColors constant in color.go.

A networkgenerator can be used to generate hierarchical networks with different structures. The generator is controlled by the fields on a HierarchySpec.
```
s := HierarchySpec{
    4,                  //Levels
    5,                  //TeamSize
    3,                  //TeamLinkLevel
    true,               //LinkTeamPeers
    true,               //LinkTeams
    []Color{Grey, Red}, //InitColors
    false,              //EvangelistAgents
    false,              //LoneEvangelist
    false,              //AgentsWithMemory
}
```
The Network generated will always be Hierarchical with a single parent Agent. The number of Levels are the number of layers in the Hierarchy including the parent Agent. The TeamSize controls how many Agents are in each team. So in this example the parent Agent is the first layer, five Agents will be linked to the parent in the second layer, there will be five Agents linked to each agent in the second layer in the third layer, and in the fourth and final layer, five Agents will be linked to each Agent in the third layer. In total there will be 1 + 5 + 5*5 + 5*5*5 = 156 Agents.

The LinkTeamPeers option allows you to control how Agents within a team are connected on the Network. If this option is false, Agents will only be connected to their direct ancestors and children in the Hierarchy. With LinkTeamPeers set to true each team that is created will have each Agent within the team connected to every other Agent within that team. So if there are five Agents, within a team each of them will have an additional four links to connect them to their peers within the team.

The LinkTeams option allows the user to create additional connectivity between teams within a deep hierarchy. This is used together with the TeamLinkLevel option. When LinkTeams is set to true, a single member is selected from each team and additional links are created to connect that Agent with her equivalent on every other team within the same level. The TeamLinkLevel is used to control at which layer in the hierarchy these additional links are created between teams. These options can be used to model the effect of intentionally introducing a cross-organisational initiative to encourage communication between teams at lower levels of the hierarchy.

The InitColors option allows the user to set the Colors that the Agents on the network will be initialised to. If this is set to nil all Agents will be initialised with Grey. If Colors are specified within this slice then the Agents in the network will be randomly assigned a Color from this slice. The Colors will be randomly distributed over the Network.

EvangelistAgents is used to specify a set of Agents selected from teams at a specified level who are set to the Color Blue, and will never choose to change their Color. A single Agent is selected from each team in the level specified by the TeamLinkLevel option. The Agents are modified so that their Susceptibility score can never be beaten, and their Color is set to Blue. This option is used to model the effect of selecting a small determined set of individuals across the organisation to introduce a specific idea or cultural change.

LoneEvangelist is similar to EvangelistAgents except there is only one agent who is an evangelist for a particular idea but she is connected to a single individual from each team at the level specified by TeamLinkLevel. This is modeling a similar effect, but it is a particularly determined individual who is well connected across the organisation trying to introduce a new idea or cultural change.

The final option is AgentsWithMemory. When set to true this uses a different Agent model that also contains memory. An Agent remembers all the previous Colors it has updated itself to. When deciding to update to a new Color it will never choose a Color that it has already been set to in the past. Without Agent memory the simulation is useful for modeling the uptake of an idea or change that is less likely to be permanent, like the preference for wearing a particular colour, or perhaps political affiliations. Whereas using Agents with memory is more useful to model the introduction of ideas that are likely to involve a permanent change such as competing technologies where adopting the technology will result in a certain amount of lock-in.

## TODOs
- [ ] Integrate into a web service
- [ ] Create a network visualiser using D3