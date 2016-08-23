package main

import (
	"database/sql"
	"fmt"
	//"log"
	"strconv"
	"time"

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

type teamScorer_SQL struct {
	TeamId string
	Score  int
}

type badmintonTeam struct {
	teamId     string
	teamNumber int
	teamName   string
	leader     string
	year       string
	Score      int
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

type TeamScoreView_SQL struct {
	TeamId string
	Score  int
}

func sqlQuery() (string, bool) {
	return "hello alan zha", true
}

var SQL_name = "./bm-all-utf8-0807-1815.s3db"

func loopMatchTableByMatchId(matchId string) []badmintonMatch {

	rows, err := db.Query("SELECT * FROM Match where MatchId=?", matchId)
	defer rows.Close()
	checkErr(err)

	var matchs []badmintonMatch = make([]badmintonMatch, 0)

	for rows.Next() {
		var match badmintonMatch
		err = rows.Scan(&match.matchId, &match.matchDate, &match.majorRound, &match.smallRoaund)
		checkErr(err)
		matchs = append(matchs, match)
	}

	return matchs
}

func loopMatchTable() []badmintonMatch {
	rows, err := db.Query("SELECT * FROM Match")
	defer rows.Close()
	checkErr(err)

	var matchs []badmintonMatch = make([]badmintonMatch, 0)

	for rows.Next() {
		var match badmintonMatch
		err = rows.Scan(&match.matchId, &match.matchDate, &match.majorRound, &match.smallRoaund)
		checkErr(err)
		matchs = append(matchs, match)
	}

	return matchs
}

func loopPersonTable() {

	rows, err := db.Query("SELECT * FROM Person")
	defer rows.Close()
	checkErr(err)

	var persons []person = make([]person, 0)

	for rows.Next() {
		var p person
		err = rows.Scan(&p.personId, &p.gender)
		checkErr(err)

		persons = append(persons, p)
	}
	//fmt.Println(persons)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
		//log.Fatal(err)
	}
}

func getPersonGender(name string) string {

	var gender string

	person_rows, err := db.Query("SELECT Gender FROM Person WHERE PersonId=\"" + name + "\"")
	defer person_rows.Close()
	checkErr(err)
	for person_rows.Next() {
		err = person_rows.Scan(&gender)
		break
	}
	return gender
}

func getTeamWithMember(name string) []TeamWithMemberInfo_SQL {
	var t TeamWithMemberInfo_SQL
	var infos []TeamWithMemberInfo_SQL

	rows, err := db.Query("SELECT MemberId, Team.TeamId, TeamName, TeamNumber, LeaderPerson, Year FROM Member JOIN Team ON Member.TeamId = Team.TeamId WHERE PersonId = \"" + name + "\" ORDER BY Year DESC;")
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&t.memberId, &t.teamId, &t.teamName, &t.teamNumber, &t.leader, &t.year)
		checkErr(err)

		infos = append(infos, t)
	}

	return infos
}

func getPairId(memberId string) []string {

	var pairS []string
	var p string
	rows, err := db.Query("SELECT PairID from Pair WHERE MemberId = \"" + memberId + "\" ORDER BY PairID DESC;")
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&p)
		checkErr(err)

		pairS = append(pairS, p)
	}
	return pairS
}

func GetPartnerWithMatchInfo(pairId string, memberId string) PartnerWithMarchInfo_SQL {
	var p PartnerWithMarchInfo_SQL
	queryCmd := "SELECT Person.PersonId as PairPersonId, Person.Gender as PairGender, MatchDate, MajorRound, SmallRound, Score, Match.MatchId " +
		"FROM Score JOIN Pair ON Pair.PairId = Score.PairId " +
		"JOIN Member ON Pair.MemberId = Member.MemberID " +
		"JOIN Person ON Member.PersonId = Person.PersonId " +
		"JOIN Match ON Score.MatchId = Match.MatchId " +
		"WHERE Pair.MemberId <> \"" + memberId + "\" " +
		"and Pair.PairId = \"" + pairId + "\""

	//fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&p.PairPersonId, &p.PairGender, &p.MatchDate, &p.MajorRound, &p.SmallRound, &p.Score, &p.MatchId)
		checkErr(err)
		break
	}
	//fmt.Println(p)
	return p
}

func getEnemyScore(pairId string, matchId string) int {

	queryCmd := "SELECT Score AS EnemyScore FROM Score WHERE MatchId = \"" +
		matchId + "\" AND PairId <> \"" + pairId + "\""

	//fmt.Println(queryCmd)
	var s int
	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err = rows.Scan(&s)
		checkErr(err)
		break
	}
	return s
}

func GetEnemyInfo(pairId string, matchId string) []EnemyInfo_SQL {

	queryCmd := "SELECT Person.PersonId as EnemyName, Person.Gender AS EnemyGender, TeamName AS EnemyTeamName " +
		"FROM Team JOIN Member ON Member.TeamId = Team.TeamId " +
		"JOIN Person ON Member.PersonId = Person.PersonId " +
		"JOIN Pair ON Member.MemberId = Pair.MemberId " +
		"WHERE Pair.PairId = ( " +
		" SELECT PairId FROM SCORE WHERE MatchId =\"" + matchId + "\"" +
		" AND PairId <> \"" + pairId + "\" )" +
		"ORDER BY EnemyName"

	//fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	defer rows.Close()
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
	//fmt.Println("Get Name list")
	var personList []personList_json
	queryCmd := "SELECT * FROM Person ORDER BY PersonId"
	//fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	defer rows.Close()
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
	queryCmd := "SELECT TeamName FROM Team WHERE TeamId = \"" + teamId + "\""
	//fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	defer rows.Close()
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

	queryCmd := "SELECT Year,TeamNumber, TeamName, Team.TeamId, LeaderPerson, Member.PersonId ,Person.Gender FROM Team JOIN Member on Member.TeamId == Team.TeamId JOIN Person on Member.PersonId == Person.PersonId ORDER BY Year DESC ,Team.TeamId"
	//fmt.Println(queryCmd)
	rows, err := db.Query(queryCmd)
	defer rows.Close()
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

	queryCmd := "SELECT DISTINCT MatchDate, MajorRound FROM ScoreView WHERE TeamId = \"" + teamId + "\" ORDER BY MatchDate DESC, MajorRound ASC"
	//fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		var p dateAndMajorRound_SQL
		err := rows.Scan(&p.date, &p.majorRound)
		checkErr(err)
		//fmt.Println(p)
		dataAndMajorRoundInfos = append(dataAndMajorRoundInfos, p)
	}
	return dataAndMajorRoundInfos
}

func getEachMatchInfo(teamId string, date string, majorRound int) []matchSelfPartInfo_SQL {
	var eachMatchInfos []matchSelfPartInfo_SQL

	queryCmd := "SELECT MatchId, PairId, SmallRound, Score FROM ScoreView WHERE TeamId = \"" + teamId +
		"\" AND MatchDate = \"" + date +
		"\" AND MajorRound = \"" + strconv.Itoa(majorRound) +
		"\" ORDER BY SmallRound"
	//fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
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

	queryCmd := "SELECT PersonId, Gender FROM PairView WHERE PairId = \"" + pairId + "\" ORDER By PersonId"
	//fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
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

	queryCmd := "SELECT Team.TeamId, TeamName, TeamNumber, LeaderPerson, Year FROM Team JOIN ScoreView ON Team.TeamId = ScoreView.TeamId WHERE ScoreView.TeamId <> \"" +
		teamId + "\" AND ScoreView.MatchId = \"" +
		matchId + "\""
	//fmt.Println(queryCmd)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&enemyTeamInfo.enemyTeamId, &enemyTeamInfo.enemyTeamName, &enemyTeamInfo.enemyTeamNumber, &enemyTeamInfo.enemyTeamLeaderName, &enemyTeamInfo.year)
		checkErr(err)
		break
	}
	return enemyTeamInfo
}

func getEnemyTeamScore(matchId string, teamId string) int {

	queryCmd := "SELECT Score FROM ScoreView WHERE ScoreView.TeamId <>\"" + teamId + "\" AND ScoreView.MatchId = \"" + matchId + "\""
	//fmt.Println(queryCmd)

	var score int
	rows, err := db.Query(queryCmd)
	defer rows.Close()
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

	queryCmd := "SELECT PersonId, Gender FROM PairView JOIN ScoreView ON PairView.PairId = ScoreView.PairId " +
		"WHERE ScoreView.TeamId <> \"" + teamId +
		"\" AND ScoreView.MatchId = \"" + matchId +
		"\" Order BY PersonId"
	//fmt.Println(queryCmd)

	var p person
	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&p.personId, &p.gender)
		checkErr(err)
		enemyPlayers = append(enemyPlayers, p)

	}
	return enemyPlayers
}

func getMatchIdAndPairId(teamId string, matchDate string, majorRound string, smallRound string) (string, string) {

	queryCmd := "select matchId, pairId from scoreView " +
		" where teamId = '" + teamId + "' and matchDate = '" + matchDate + "' and majorRound = " + majorRound + " and smallRound = " + smallRound
	//fmt.Println(queryCmd)

	var matchId = ""
	var pairId = ""
	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&matchId, &pairId)
		checkErr(err)
		//fmt.Println("round")
	}

	return matchId, pairId
}
func updateScore(matchId string, pairId string, score string) {
	score_int, err := strconv.Atoi(score)
	checkErr(err)

	fmt.Println("Update score:" + matchId + ", " + pairId + ", " + score)
	stmt, err := db.Prepare("UPDATE Score SET Score =? WHERE MatchId =? and PairId =?")
	defer stmt.Close()
	checkErr(err)

	_, err = stmt.Exec(score_int, matchId, pairId)
	checkErr(err)

	stmt, err = db.Prepare("INSERT INTO UpdateScoreRecord(MatchId, PairId, Score, UpdateTime) VALUES( ?, ?, ?, ?)")
	checkErr(err)

	timenow := time.Now().String()

	_, err = stmt.Exec(matchId, pairId, score_int, timenow)
	checkErr(err)
}

func getTeamScore() []teamScorer_SQL {
	var teamScores []teamScorer_SQL

	queryCmd := "SELECT TeamId, Sum(ScoreInMajorRound) AS TeamScore " +
		//"TeamNumber, TeamName, LeaderPerson, Year " +
		"FROM ( " +
		"SELECT TeamId AS TheTeamId, MatchDate, MajorRound, " +
		"Case " +
		"WHEN Count(TeamId) = 0 THEN 0 " +
		"WHEN Count(TeamId) < 3 THEN 1 " +
		"ELSE 3 " +
		"END As ScoreInMajorRound " +
		"FROM ScoreView JOIN (select MatchId AS TheMatchId, Max(Score) AS MaxScore from Score Group BY MatchId) AS MaxScoreTable " +
		"ON ScoreVIew.MatchId = TheMatchId AND ScoreView.Score = MaxScore AND MaxScore > 0 " +
		"GROUP BY TeamId, MatchDate, MajorRound " +
		") JOIN Team ON TheTeamId = TeamId " +
		"GROUP BY TeamId " +
		"Order By Year DESC, TeamNumber ASC"

	//fmt.Println(queryCmd)

	var t teamScorer_SQL
	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)
	for rows.Next() {
		err := rows.Scan(&t.TeamId, &t.Score)
		checkErr(err)
		teamScores = append(teamScores, t)
	}

	return teamScores
}

func getTeamScoreInfo(matchDate string, majorRound int) []TeamScoreView_SQL {
	var TeamScores []TeamScoreView_SQL
	queryCmd := "select TeamId, Score from Scoreview where MatchDate ='" + matchDate + "' and MajorRound=" + strconv.Itoa(majorRound)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)

	var t TeamScoreView_SQL

	for rows.Next() {
		err := rows.Scan(&t.TeamId, &t.Score)
		checkErr(err)
		TeamScores = append(TeamScores, t)
		//fmt.Println("round")
	}

	return TeamScores
}

func getLatestTeamList(matchDate string, majorRound int) []string {
	var TeamIdList []string

	queryCmd := "select DISTINCT TeamId from Scoreview where MatchDate ='" + matchDate + "' and MajorRound=" + strconv.Itoa(majorRound)

	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)

	var t string

	for rows.Next() {
		err := rows.Scan(&t)
		checkErr(err)
		TeamIdList = append(TeamIdList, t)
	}

	return TeamIdList
}

func getLatestMatchDate() string {
	queryCmd := "select DISTINCT MatchDate from Match ORDER BY MatchDate DESC"

	rows, err := db.Query(queryCmd)
	defer rows.Close()
	checkErr(err)

	var t string

	for rows.Next() {
		err := rows.Scan(&t)
		checkErr(err)
		break
	}

	return t
}

var db *sql.DB

func Sql_open() {
	var err error
	db, _ = sql.Open("sqlite3", SQL_name)
	checkErr(err)
}
