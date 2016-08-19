package main

type personData_json struct {
	Date           string
	Partner        string
	Opponenter     [2]string
	OpponenterTeam string
	Score          string
	Combin_type    string
	Team           string
	TeamId         string
}

type personList_json struct {
	Name   string
	Gender string
}

type teamInfo_json struct {
	Year                string
	TeamNumber          string
	TeamName            string
	TeamId              string
	Score               string
	TeamLeader          string
	MaleTeamMemberNum   int
	MaleTeamMember      []string
	FemaleTeamMemberNum int
	FeMaleTeamMember    []string
}

type PersonHtmlTemplate struct {
	Name string
}

type MatchResult_json struct {
	OurPlayerNum      int
	OurPlayer         []string
	OurPlayerGender   []string
	EnemyPlayerNum    int
	EnemyPlayer       []string
	EnemyPlayerGender []string
	Score             string
}

type EnemyTeamInfo_json struct {
	EnemyTeamName   string
	EnemyTeamId     string
	EnemyTeamNumber int
	Date            string
	MatchResultNum  int
	MatchResults    []MatchResult_json
}

type teamResultInfo_json struct {
	TeamName          string
	TotalScore        int
	EnemyTeamInfoNum  int
	EnemyTeamInfoList []EnemyTeamInfo_json
}
