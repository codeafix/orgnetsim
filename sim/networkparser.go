package sim

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//ParseOptions contains details about how to parse data from the incoming file
//Identifier provides the index of the column to use as an Agnent Identifier
//Parent provides the index of the column to use as the Identifier of a Parent Agent
//Regex provides the regular expressions to use to extract data from the columns.
//It is in the form of a map, the key being the index of the column, the value being
//the regex to apply.
//When a Regex is supplied for the parent or identifier columns it will be used to
//extract the value to use as the Identifier. If it is applied to another column then
//the row will be skipped where the regex does not match the contents in that column.
//Any regex supplied for a column that doesn't exist within a row will result in that
//row being skipped.
//If no regex is supplied for the parent and identifier columns then a default is applied
//which will strip leading and trailing whitespace.
//Delimiter is the delimiter to use when slicing rows into columns
type ParseOptions struct {
	Identifier int               `json:"identifier"`
	Parent     int               `json:"parent"`
	Regex      map[string]string `json:"regex"`
	Delimiter  string            `json:"delimiter"`
}

//IdentifierRegex returns the Regexp that must be applied to the Identifier column
func (po *ParseOptions) IdentifierRegex() *regexp.Regexp {
	regex, _ := po.GetColRegex(po.Identifier)
	return regex
}

//ParentRegex returns the Regexp that must be applied to the Parent column
func (po *ParseOptions) ParentRegex() *regexp.Regexp {
	regex, _ := po.GetColRegex(po.Parent)
	return regex
}

//GetColRegex returns the Regexp that must be applied to the indicated column
func (po *ParseOptions) GetColRegex(col int) (*regexp.Regexp, bool) {
	regex, exists := po.Regex[strconv.Itoa(col)]
	if !exists || whitespace().MatchString(regex) {
		regex = `^\s*(\S.*\S)\s*$`
	}
	return regexp.MustCompile(regex), exists
}

//GetOtherRegex returns the compiled regular expressions for columns other
//than the parent and identifier columns
func (po *ParseOptions) GetOtherRegex() map[int]*regexp.Regexp {
	ret := make(map[int]*regexp.Regexp, len(po.Regex))
	for index, regex := range po.Regex {
		i, err := strconv.Atoi(index)
		if err != nil || i == po.Identifier || i == po.Parent {
			continue
		}
		ret[i] = regexp.MustCompile(regex)
	}
	return ret
}

//Returns a Regex that matches a string that only contains whitespace or is empty
func whitespace() *regexp.Regexp {
	return regexp.MustCompile(`^\s*$`)
}

//ParseDelim takes a hierarchy expressed in a comma separated file and generates a Network out of it
//po are the ParseOptions that control how the parse will operate.
//returns a RelationshipMgr containing all the Agents with links to parent Agents as described in the input
//data. If the same Id is listed in multiple rows as specified after any regular expressions is applied the
//first row is used and subsequent rows are ignored.
func (po *ParseOptions) ParseDelim(data []string) (RelationshipMgr, error) {
	ws := whitespace()
	idre := po.IdentifierRegex()
	pre := po.ParentRegex()

	n := Network{}

	//Using a map of maps here because the links can be duplicated in the source data
	//and this will allow me to ignore the duplicates
	links := map[string]map[string]struct{}{}
	agents := map[string]Agent{}

	for i := 1; i < len(data); i++ {
		cols := strings.Split(data[i], po.Delimiter)
		if po.Identifier < 0 || po.Identifier >= len(cols) {
			continue
		}
		skip := false
		for index, re := range po.GetOtherRegex() {
			if index < 0 || index >= len(cols) || !re.MatchString(cols[index]) {
				skip = true
				break
			}
		}
		if skip {
			continue
		}

		id := idre.ReplaceAllString(cols[po.Identifier], "$1")
		idParent := ""
		if po.Parent > 0 && po.Parent < len(cols) {
			idParent = pre.ReplaceAllString(cols[po.Parent], "$1")
		}

		if ws.MatchString(id) {
			continue
		}
		_, exists := agents[id]
		if exists {
			continue
		}
		a := GenerateRandomAgent(id, []Color{}, false)
		n.AddAgent(a)
		agents[id] = a
		if !ws.MatchString(idParent) {
			if id != idParent {
				id1entry, exists := links[idParent]
				if exists {
					_, exists := id1entry[id]
					if !exists {
						id1entry[id] = struct{}{}
					}
				} else {
					links[idParent] = map[string]struct{}{id: struct{}{}}
				}
			}
		}
	}

	for id1, id1Entries := range links {
		for id2 := range id1Entries {
			agent1, agent2, err := getAgents(agents, id1, id2)
			if err != nil {
				return nil, err
			}
			n.AddLink(agent1, agent2)
		}
	}

	err := n.PopulateMaps()
	return &n, err
}

//Convenience method to get Agents and check they exist
func getAgents(agents map[string]Agent, id1 string, id2 string) (Agent, Agent, error) {
	agent1, exists := agents[id1]
	if !exists {
		return nil, nil, fmt.Errorf("Agent id1 '%s' not found when adding link. Agent id2 is '%s'", id1, id2)
	}
	agent2, exists := agents[id2]
	if !exists {
		return agent1, nil, fmt.Errorf("Agent id2 '%s' not found when adding link. Agent id1 is '%s'", id2, id1)
	}
	return agent1, agent2, nil
}
