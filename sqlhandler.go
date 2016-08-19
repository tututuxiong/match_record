package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type TeamWithMemberInfo_SQL struct {
	memberId   string
	teamId     string
	teamNumber int
	teamName   string
	leader     string
	year       string
}

type PartnerWithMarchInfo_SQL struct {
	PairPersonId string
	PairGender   string
	MatchDate    string
	MajorRound   int
	SmallRound   int
	Score        int
	MatchId      string
}

type EnemyInfo_SQL struct {
	enemyName     string
	enemyGender   string
	enemyTeamName string
}

type teamMember_SQL struct {
	Year             string
	TeamNumber       string
	TeamName         string
	TeamId           string
	Leader           string
	TeamMember       string
	TeamMemberGender string
}

type badmintonTeam struct {
	teamId     string
	teamNumber int
	teamName   string
	leader     string
	year       string
}

type badmintonMatch struct {
	matchId     string
	majorRound  int
	smallRoaund int
	matchDate   string
}

type member struct {
	memberId string
	teamId   string
	personId string
}

type pair struct {
	pairId   string
	memberId string
}

type person struct {
	personId string
	gender   string
}

type dateAndMajorRound_SQL struct {
	date       string
	majorRound int
}

type matchSelfPartInfo_SQL struct {
	matchId    string
	pairId     string
	smallRound int
	score      int
}

type PairInfo_SQL struct {
	personId string
	gender   string
}

type EnemyTeamInfo_SQL struct {
	enemyTeamLeaderName string
	enemyTeamId         string
	enemyTeamNumber     int
	enemyTeamName       string
	year                string
}

func sqlQuery() (string, bool) {
	return "hello alan zha", true
}

var SQL_name = "./bm-all-utf8-0807-1815.s3db"

func loopMatchTableByMatchId(matchId string) []badmintonMatch {
	db, err := sql.Open("sqlite3", SQL_name)
	checkErr(err)

	rows, err := db.Query("SELECT * FROM Match where MatchId=?", matchId)
	checkErr(err)

	var matchs []badmintonMatch = make([]badmintonMatch, 0)

	for rows.Next() {
		var match badmintonMatch
		err = rows.Scan(&match.matchId, &match.matchDate, &match.majorRound, &match.smallRoaund)
		checkErr(err)
		matchs = append(matchs, match)
	}

	fmt.Println(matchs)
	db.Close()
	return matchs
}

func loopMatchTable() []badmintonMatch {
	db, err := sql.Open("sqlite3", SQL_name)
	checkErr(err)

	rows, err := db.Query("SELECT * FROM Match")
	checkErr(err)

	var matchs []badmintonMatch = make([]badmintonMatch, 0)

	for rows.Next() {
		var match badmintonMatch
		err = rows.Scan(&match.matchId, &match.matchDate, &match.majorRound, &match.smallRoaund)
		checkErr(err)
		matchs = append(matchs, match)
	}

	fmt.Println(matchs)
	db.Close()
	return matchs
}

func loopPersonTable() {
	db, err := sql.Open("sqlite3", SQL_name)
	checkErr(err)

	rows, err := db.Query("SELECT * FROM Person")
	checkErr(err)

	var persons []person = make([]person, 0)

	for rows.Next() {
		var p person
		err = rows.Scan(&p.personId, &p.gender)
		checkErr(err)

		persons = append(persons, p)
	}
	fmt.Println(persons)
	db.Close()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func getPersonGender(name string) string {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	var gender string

	person_rows, err := db.Query("SELECT Gender FROM Person WHERE PersonId=\"" + name + "\"")
	checkErr(err)
	for person_rows.Next() {
		err = person_rows.Scan(&gender)
		break
	}
	return gender
}

func getTeamWithMember(name string) []TeamWithMemberInfo_SQL {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	var t TeamWithMemberInfo_SQL
	var infos []TeamWithMemberInfo_SQL

	rows, err := db.Query("SELECT MemberId, Team.TeamId, TeamName, TeamNumber, LeaderPerson, Year FROM Member JOIN Team ON Member.TeamId = Team.TeamId WHERE PersonId = \"" + name + "\" ORDER BY Year DESC;")
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&t.memberId, &t.teamId, &t.teamName, &t.teamNumber, &t.leader, &t.year)
		checkErr(err)

		infos = append(infos, t)
	}

	return infos
}

func getPairId(memberId string) []string {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	var pairS []string
	var p string
	rows, err := db.Query("SELECT PairID from Pair WHERE MemberId = \"" + memberId + "\" ORDER BY PairID DESC;")
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&p)
		checkErr(err)

		pairS = append(pairS, p)
	}
	return pairS
}

func GetPartnerWithMatchInfo(pairId string, memberId string) PartnerWithMarchInfo_SQL {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	var p PartnerWithMarchInfo_SQL
	queryCmd := "SELECT Person.PersonId as PairPersonId, Person.Gender as PairGender, MatchDate, MajorRound, SmallRound, Score, Match.MatchId " +
		"FROM Score JOIN Pair ON Pair.PairId = Score.PairId " +
		"JOIN Member ON Pair.MemberId = Member.MemberID " +
		"JOIN Person ON Member.PersonId = Person.PersonId " +
		"JOIN Match ON Score.MatchId = Match.MatchId " +
		"WHERE Pair.MemberId <> \"" + memberId + "\" " +
		"and Pair.PairId = \"" + pairId + "\""

	fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&p.PairPersonId, &p.PairGender, &p.MatchDate, &p.MajorRound, &p.SmallRound, &p.Score, &p.MatchId)
		checkErr(err)
		break
	}
	fmt.Println(p)
	return p
}

func getEnemyScore(pairId string, matchId string) int {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT Score AS EnemyScore FROM Score WHERE MatchId = \"" +
		matchId + "\" AND PairId <> \"" + pairId + "\""

	fmt.Println(queryCmd)
	var s int
	rows, err := db.Query(queryCmd)

	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&s)
		checkErr(err)
		break
	}
	return s
}

func GetEnemyInfo(pairId string, matchId string) []EnemyInfo_SQL {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT Person.PersonId as EnemyName, Person.Gender AS EnemyGender, TeamName AS EnemyTeamName " +
		"FROM Team JOIN Member ON Member.TeamId = Team.TeamId " +
		"JOIN Person ON Member.PersonId = Person.PersonId " +
		"JOIN Pair ON Member.MemberId = Pair.MemberId " +
		"WHERE Pair.PairId = ( " +
		" SELECT PairId FROM SCORE WHERE MatchId =\"" + matchId + "\"" +
		" AND PairId <> \"" + pairId + "\" )" +
		"ORDER BY EnemyName"

	fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	checkErr(err)
	var infos []EnemyInfo_SQL
	var e EnemyInfo_SQL
	for rows.Next() {
		err = rows.Scan(&e.enemyName, &e.enemyGender, &e.enemyTeamName)
		checkErr(err)
		infos = append(infos, e)
	}
	return infos
}
func getPersonList() []personList_json {
	fmt.Println("Get Name list")
	var personList []personList_json

	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT * FROM Person ORDER BY PersonId"
	fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	checkErr(err)
	var p personList_json
	for rows.Next() {
		err = rows.Scan(&p.Name, &p.Gender)
		checkErr(err)
		personList = append(personList, p)
	}
	return personList
}

func getTeamNameById(teamId string) string {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT TeamName FROM Team WHERE TeamId = \"" + teamId + "\""
	fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	checkErr(err)
	var teamName string
	for rows.Next() {
		err = rows.Scan(&teamName)
		checkErr(err)
		break
	}
	return teamName
}

func getAllTeamMemberGenderInfo() []teamMember_SQL {

	var teamMemberList []teamMember_SQL
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT Year,TeamNumber, TeamName, Team.TeamId, LeaderPerson, Member.PersonId ,Person.Gender FROM Team JOIN Member on Member.TeamId == Team.TeamId JOIN Person on Member.PersonId == Person.PersonId ORDER BY Year DESC ,Team.TeamId"
	fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	checkErr(err)

	for rows.Next() {
		var t teamMember_SQL
		err = rows.Scan(&t.Year, &t.TeamNumber, &t.TeamName, &t.TeamId, &t.Leader, &t.TeamMember, &t.TeamMemberGender)
		checkErr(err)
		teamMemberList = append(teamMemberList, t)
	}

	return teamMemberList
}

func getDataAndMajorRound(teamId string) []dateAndMajorRound_SQL {
	var dataAndMajorRoundInfos []dateAndMajorRound_SQL
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT DISTINCT MatchDate, MajorRound FROM ScoreView WHERE TeamId = \"" + teamId + "\" ORDER BY MatchDate DESC, MajorRound ASC"
	fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		var p dateAndMajorRound_SQL
		err := rows.Scan(&p.date, &p.majorRound)
		checkErr(err)
		dataAndMajorRoundInfos = append(dataAndMajorRoundInfos, p)
	}
	return dataAndMajorRoundInfos
}

func getEachMatchInfo(teamId string, date string, majorRound int) []matchSelfPartInfo_SQL {
	var eachMatchInfos []matchSelfPartInfo_SQL
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT MatchId, PairId, SmallRound, Score FROM ScoreView WHERE TeamId = \"" + teamId +
		"\" AND MatchDate = \"" + date +
		"\" AND MajorRound = \"" + strconv.Itoa(majorRound) +
		"\" ORDER BY SmallRound"
	fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		var e matchSelfPartInfo_SQL
		err := rows.Scan(&e.matchId, &e.pairId, &e.smallRound, &e.score)
		checkErr(err)
		eachMatchInfos = append(eachMatchInfos, e)
	}
	return eachMatchInfos
}

func getPairInfos(pairId string) []PairInfo_SQL {
	var pairInfos []PairInfo_SQL
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT PersonId, Gender FROM PairView WHERE PairId = \"" + pairId + "\" ORDER By PersonId"
	fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		var p PairInfo_SQL
		err := rows.Scan(&p.personId, &p.gender)
		checkErr(err)
		pairInfos = append(pairInfos, p)
	}
	return pairInfos
}

func getEnemyTeamInfos(matchId string, teamId string) EnemyTeamInfo_SQL {
	var enemyTeamInfo EnemyTeamInfo_SQL
	//enemyTeamInfo.enemyTeamId = teamId

	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT Team.TeamId, TeamName, TeamNumber, LeaderPerson, Year FROM Team JOIN ScoreView ON Team.TeamId = ScoreView.TeamId WHERE ScoreView.TeamId <> \"" +
		teamId + "\" AND ScoreView.MatchId = \"" +
		matchId + "\""
	fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&enemyTeamInfo.enemyTeamId, &enemyTeamInfo.enemyTeamName, &enemyTeamInfo.enemyTeamNumber, &enemyTeamInfo.enemyTeamLeaderName, &enemyTeamInfo.year)
		checkErr(err)
		break
	}
	return enemyTeamInfo
}

func getEnemyTeamScore(matchId string, teamId string) int {
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT Score FROM ScoreView WHERE ScoreView.TeamId <>\"" + teamId + "\" AND ScoreView.MatchId = \"" + matchId + "\""
	fmt.Println(queryCmd)

	var score int
	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&score)
		checkErr(err)
		break
	}
	return score
}

func getEnemyPlayerData(matchId string, teamId string) []person {
	var enemyPlayers []person
	db, err := sql.Open("sqlite3", SQL_name)
	defer db.Close()
	checkErr(err)

	queryCmd := "SELECT PersonId, Gender FROM PairView JOIN ScoreView ON PairView.PairId = ScoreView.PairId " +
		"WHERE ScoreView.TeamId <> \"" + teamId +
		"\" AND ScoreView.MatchId = \"" + matchId +
		"\" Order BY PersonId"
	fmt.Println(queryCmd)

	var p person
	rows, err := db.Query(queryCmd)
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&p.personId, &p.gender)
		checkErr(err)
		enemyPlayers = append(enemyPlayers, p)

	}
	return enemyPlayers
}
