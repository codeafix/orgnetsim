package sim

import "testing"

func TestParseDelimAgents(t *testing.T) {
	data := []string{
		"Header row is always skipped |check_this_is_not_an_Id||",
		"Strips ws around Id leaves ws in Id|   id _ 1  |   |  ",
		"Should be ignored|||",
		"Strips ws around Id| my_id ||",
		"Should also be ignored|   ||",
		"Pulls entire Id without ws|1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+={}[]()!£$%^&*@~#<>/?\\||",
		"Should also be ignored|",
		"Regular Id|Some_id7|",
	}
	po := ParseOptions{
		Delimiter:  "|",
		Identifier: 1,
		Parent:     -1,
	}
	rm, err := po.ParseDelim(data)
	AssertSuccess(t, err)

	agents := rm.Agents()
	AreEqual(t, 4, len(agents), "Wrong number of agents parsed from source data")
	AreEqual(t, "id _ 1", agents[0].Identifier(), "Wrong first agent Id")
	AreEqual(t, "my_id", agents[1].Identifier(), "Wrong first agent Id")
	AreEqual(t, "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+={}[]()!£$%^&*@~#<>/?\\", agents[2].Identifier(), "Wrong first agent Id")
	AreEqual(t, "Some_id7", agents[3].Identifier(), "Wrong first agent Id")
}

func TestParseDelimLinks(t *testing.T) {
	data := []string{
		"Header row is always skipped |check_this_is_not_an_Id||",
		"Strips ws around Id leaves ws in Id|   id _ 1  |   |  ",
		"Should be ignored|||",
		"Strips ws around Id| my_id ||id _ 1",
		"Blank lines are ignored",
		"Should also be ignored|   ||",
		"Pulls entire Id without ws|1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+={}[]()!£$%^&*@~#<>/?\\",
		"Should also be ignored|",
		"Regular Id|Some_id7||my_id|",
		"Repeated Agent Id is ignored|Some_id7||id _ 1|",
		"Another Agent child of my_id|NewId18 || my_id ",
		"Ignores a link to itself|Id69 || Id69 ",
		"Child of really long Id|NewId35 || 1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+={}[]()!£$%^&*@~#<>/?\\",
	}
	po := ParseOptions{
		Delimiter:  "|",
		Identifier: 1,
		Parent:     3,
	}
	rm, err := po.ParseDelim(data)
	AssertSuccess(t, err)

	links := rm.Links()
	AreEqual(t, 4, len(links), "Wrong number of links parsed from source data")
	for _, link := range links {
		switch link.Agent2ID {
		case "my_id":
			AreEqual(t, "id _ 1", link.Agent1ID, "Wrong parent for Agent 'my_id'")
		case "Some_id7":
			AreEqual(t, "my_id", link.Agent1ID, "Wrong parent for Agent 'Some_id7'")
		case "NewId18":
			AreEqual(t, "my_id", link.Agent1ID, "Wrong parent for Agent 'NewId18'")
		case "NewId35":
			AreEqual(t, "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-+={}[]()!£$%^&*@~#<>/?\\", link.Agent1ID, "Wrong parent for Agent 'NewId35'")
		default:
			t.Errorf("Unexpected link found: %s -> %s", link.Agent1ID, link.Agent2ID)
		}
	}
}

func TestParseDelimLinksThrowsErrorWhenLinkTargetDoesntExist(t *testing.T) {
	data := []string{
		"Header row is always skipped |check_this_is_not_an_Id||",
		"Strips ws around Id leaves ws in Id|   id _ 1  |   |  ",
		"Should be ignored|||",
		"Strips ws around Id| my_id ||id _ 1",
		"Blank lines are ignored",
		"Parent doesn't exist!|Some_id7||Id615|",
		"Another Agent child of my_id|NewId18 || my_id ",
		"Ignores a link to itself|Id69 || Id69 ",
	}
	po := ParseOptions{
		Delimiter:  "|",
		Identifier: 1,
		Parent:     3,
	}
	_, err := po.ParseDelim(data)
	IsFalse(t, err == nil, "Expecting an error to be thrown")
}

func TestParseDelimCsv(t *testing.T) {
	data := []string{
		"Header row is always skipped ,check_this_is_not_an_Id,,",
		"Should be ignored|||",
		"",
		"Strips ws around Id, my_id",
		"Blank lines are ignored",
	}
	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 1,
		Parent:     3,
	}
	rm, err := po.ParseDelim(data)
	AssertSuccess(t, err)

	agents := rm.Agents()
	AreEqual(t, 1, len(agents), "Wrong number of agents parsed from source data")
	AreEqual(t, "my_id", agents[0].Identifier(), "Wrong first agent Id")
}

func TestParseDelimTsv(t *testing.T) {
	data := []string{
		"Header row is always skipped	check_this_is_not_an_Id		",
		"Should be ignored|||",
		"",
		"Strips ws around Id	  my_id ",
		"Blank lines are ignored",
	}
	po := ParseOptions{
		Delimiter:  "\t",
		Identifier: 1,
		Parent:     3,
	}
	rm, err := po.ParseDelim(data)
	AssertSuccess(t, err)

	agents := rm.Agents()
	AreEqual(t, 1, len(agents), "Wrong number of agents parsed from source data")
	AreEqual(t, "my_id", agents[0].Identifier(), "Wrong first agent Id")
}

func TestParseDelimOnlyIncludesRowsWithOtherMatchingColumn(t *testing.T) {
	data := []string{
		"Header row is always skipped |check_this_is_not_an_Id||Include",
		"Strips ws around Id leaves ws in Id|   id _ 1  |   | Include  ",
		"Should be ignored| id_2 ||",
		"Should be ignored| id_3",
		"Blank lines are ignored",
		"Should be included|id_4|   |Include|",
	}
	po := ParseOptions{
		Delimiter:  "|",
		Identifier: 1,
		Regex: map[string]string{
			"3": `\S+`,
		},
	}
	rm, err := po.ParseDelim(data)
	AssertSuccess(t, err)

	agents := rm.Agents()
	AreEqual(t, 2, len(agents), "Wrong number of agents parsed from source data")
	AreEqual(t, "id _ 1", agents[0].Identifier(), "Wrong first agent Id")
	AreEqual(t, "id_4", agents[1].Identifier(), "Wrong second agent Id")
}

func TestParseEdgesCsv(t *testing.T) {
	data := []string{
		"Header row is always skipped ,check_this_is_not_an_Id,,",
		"Should be ignored|||",
		"",
		"Blank lines are ignored",
		",Lines with no Id should be ignored",
		"Lines with no parent should be ignored,",
		"Child,Parent",
		"  my_id,  my_parent  ",
	}
	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 0,
		Parent:     1,
	}
	n := Network{}
	n.AddAgent(GenerateRandomAgent("my_id", "An agent", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("my_parent", "Parent agent", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("Child", "Another agent", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("Parent", "Another parent agent", []Color{}, false))
	n.PopulateMaps()

	rm, err := po.ParseEdges(data, &n)
	AssertSuccess(t, err)

	links := rm.Links()
	AreEqual(t, 2, len(links), "Wrong number of links parsed from source data")
	for _, link := range links {
		switch link.Agent2ID {
		case "my_id":
			AreEqual(t, "my_parent", link.Agent1ID, "Wrong parent for Agent 'my_id'")
		case "Child":
			AreEqual(t, "Parent", link.Agent1ID, "Wrong parent for Agent 'Some_id7'")
		default:
			t.Errorf("Unexpected link found: %s -> %s", link.Agent1ID, link.Agent2ID)
		}
	}
}

func TestParseEdgesThrowsErrorWhenAgentDoesntExist(t *testing.T) {
	data := []string{
		"Header row is always skipped ,check_this_is_not_an_Id,,",
		"my_id, my_parent ",
	}
	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 0,
		Parent:     1,
	}
	n := Network{}
	n.AddAgent(GenerateRandomAgent("my_id", "An agent", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("Child", "Another agent", []Color{}, false))
	n.PopulateMaps()

	_, err := po.ParseEdges(data, &n)
	IsFalse(t, err == nil, "Expecting an error to be thrown")
}

func TestParseEdgesAddsMultipleParentLinks(t *testing.T) {
	data := []string{
		"u1,u2,<4,4-7,8+,total",
		"1,1349,1,6,163,170",
		"1,35,0,0,5,5",
		"16,35,0,0,5,5",
	}

	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 0,
		Parent:     1,
	}

	n := Network{}
	n.AddAgent(GenerateRandomAgent("1", "agent 1", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("1349", "agent 2", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("35", "agent 4", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("16", "agent 5", []Color{}, false))
	n.PopulateMaps()

	rm, err := po.ParseEdges(data, &n)
	AssertSuccess(t, err)

	links := rm.Links()
	AreEqual(t, 3, len(links), "Wrong number of links parsed from source data")
	for _, link := range links {
		switch link.Agent2ID {
		case "1":
			IsTrue(t, link.Agent1ID == "1349" || link.Agent1ID == "35", "Wrong parent for Agent '1'")
		case "16":
			AreEqual(t, "35", link.Agent1ID, "Wrong parent for Agent '16'")
		default:
			t.Errorf("Unexpected link found: %s -> %s", link.Agent1ID, link.Agent2ID)
		}
	}
}

func TestParseEdgesOnlyIncludesLinksWithOtherMatchingColumn(t *testing.T) {
	data := []string{
		"u1,u2,<4,4-7,8+,total",
		"1,1349,1,6,163,170",
		"1,248,0,1,28,29",
		"1,35,0,0,5,5",
	}

	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 0,
		Parent:     1,
		Regex: map[string]string{
			"2": `[1-9]\d*`,
		},
	}

	n := Network{}
	n.AddAgent(GenerateRandomAgent("1", "agent 1", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("1349", "agent 2", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("248", "agent 3", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("35", "agent 4", []Color{}, false))
	n.PopulateMaps()

	rm, err := po.ParseEdges(data, &n)
	AssertSuccess(t, err)

	links := rm.Links()
	AreEqual(t, 1, len(links), "Wrong number of links parsed from source data")
	AreEqual(t, "1349", links[0].Agent1ID, "Wrong parent for first link")
	AreEqual(t, "1", links[0].Agent2ID, "Wrong child for first link")
}

func TestParseEdgesOnlyIncludesLinksWithOtherMatchingColumn2(t *testing.T) {
	data := []string{
		"u1,u2,<4,4-7,8+,total",
		"1,1349,1,6,163,170",
		"1,248,0,1,28,29",
		"1,35,0,0,5,5",
	}

	po := ParseOptions{
		Delimiter:  ",",
		Identifier: 0,
		Parent:     1,
		Regex: map[string]string{
			"3": `[1-9]\d*`,
		},
	}

	n := Network{}
	n.AddAgent(GenerateRandomAgent("1", "agent 1", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("1349", "agent 2", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("248", "agent 3", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("35", "agent 4", []Color{}, false))
	n.AddAgent(GenerateRandomAgent("351", "agent 5", []Color{}, false))
	n.PopulateMaps()

	rm, err := po.ParseEdges(data, &n)
	AssertSuccess(t, err)

	links := rm.Links()
	AreEqual(t, 2, len(links), "Wrong number of links parsed from source data")
	for _, link := range links {
		switch link.Agent1ID {
		case "1349":
			AreEqual(t, "1", link.Agent2ID, "Wrong child '1349' for Agent '1'")
		case "248":
			AreEqual(t, "1", link.Agent2ID, "Wrong child '248' for Agent '1'")
		default:
			t.Errorf("Unexpected link found: %s -> %s", link.Agent1ID, link.Agent2ID)
		}
	}
}
