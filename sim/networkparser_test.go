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
