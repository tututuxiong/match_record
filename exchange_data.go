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
	Score               int
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
	OurScore          int
	EnemyScore        int
	SmallRound        int
}

type EnemyTeamInfo_json struct {
	EnemyTeamName   string
	EnemyTeamId     string
	EnemyTeamNumber int
	Date            string
	MatchResultNum  int
	MajorRound      int
	MatchResults    []MatchResult_json
}

type teamResultInfo_json struct {
	TeamName          string
	TeamId            string
	TotalScore        int
	EnemyTeamInfoNum  int
	EnemyTeamInfoList []EnemyTeamInfo_json
}

type updateScoreRedirction_json struct {
	Address string
}

type UpdateScorePageTemlate struct {
	Date          string
	TypeId        string
	MajorRound    string
	SmallRound    string
	OurTeamId     string
	OurTeamName   string
	EnemyTeamId   string
	EnemyTeamName string
	OurScore      string
	EnemyScore    string
	OurPlayer0    string
	OurPlayer1    string
	EnemyPlayer0  string
	EnemyPlayer1  string
}

type latestTwoMatchInfo_json struct {
	Date   string
	Round1 []latestMatchInfo_json
	Round2 []latestMatchInfo_json
}

type latestMatchInfo_json struct {
	TeamName string
	TeamId   string
	Scores   []int
}
