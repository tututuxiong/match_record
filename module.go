package main

import (
	"fmt"
	"strconv"
)

func getTypeByGender(gender1 string, gender2 string) string {
	if gender1 == "M" && gender2 == "M" {
		return "男双"
	} else if gender1 == "F" && gender2 == "F" {
		return "女双"
	} else {
		return "混双"
	}

}

func getPersonRecordDatas(name string) []personData_json {
	fmt.Println("Query record for" + name)
	var personRecords []personData_json

	Gender := getPersonGender(name)
	//fmt.Println(Gender)

	teamWithMemberS := getTeamWithMember(name)
	//fmt.Println(teamWithMemberS)

	for _, t_m := range teamWithMemberS {
		var pairS []string
		pairS = getPairId(t_m.memberId)
		//fmt.Println(pairS)

		for _, p := range pairS {
			PartnerWithMatch := GetPartnerWithMatchInfo(p, t_m.memberId)
			//fmt.Println(PartnerWithMatch)

			enemyScore := getEnemyScore(p, PartnerWithMatch.MatchId)
			//fmt.Println(enemyScore)

			enemyPlayerInfo := GetEnemyInfo(p, PartnerWithMatch.MatchId)
			//fmt.Println(enemyPlayerInfo)

			var record personData_json
			record.Date = PartnerWithMatch.MatchDate
			record.Partner = PartnerWithMatch.PairPersonId
			record.Score = strconv.Itoa(PartnerWithMatch.Score) + " : " + strconv.Itoa(enemyScore)
			record.Team = t_m.teamName
			record.Combin_type = getTypeByGender(Gender, PartnerWithMatch.PairGender)
			record.TeamId = t_m.teamId

			for i, e := range enemyPlayerInfo {
				if i < 2 {
					record.Opponenter[i] = e.enemyName
					record.OpponenterTeam = e.enemyTeamName
				} else {
					//fmt.Println("")
				}

			}
			personRecords = append(personRecords, record)
		}
	}
	//fmt.Println(personRecords)
	return personRecords
}

func getTeamList(teamMemberList []teamMember_SQL) []teamInfo_json {
	var teamInfos []teamInfo_json

	for _, teamMemberList_range := range teamMemberList {
		var haveTeamAlready = false
		for _, teamInfos_range := range teamInfos {
			if teamInfos_range.TeamId == teamMemberList_range.TeamId {
				haveTeamAlready = true
				break
			}
		}

		if haveTeamAlready == false {
			var t teamInfo_json
			t.Year = teamMemberList_range.Year
			t.TeamLeader = teamMemberList_range.Leader
			t.TeamId = teamMemberList_range.TeamId
			t.TeamName = teamMemberList_range.TeamName
			t.TeamNumber = teamMemberList_range.TeamNumber
			teamInfos = append(teamInfos, t)
		}
	}

	return teamInfos
}
func getTeamInfoList() []teamInfo_json {
	teamMemberList := getAllTeamMemberGenderInfo()

	var return_teamInfos []teamInfo_json
	teamInfos := getTeamList(teamMemberList)

	teamScores := getTeamScore()

	for _, teamInfos_range := range teamInfos {
		for _, TeamScore := range teamScores {
			if teamInfos_range.TeamId == TeamScore.TeamId {
				teamInfos_range.Score = TeamScore.Score
			}
		}

		for _, teamMemberList_range := range teamMemberList {
			if teamInfos_range.TeamId == teamMemberList_range.TeamId {
				if teamMemberList_range.TeamMemberGender == "M" {
					teamInfos_range.MaleTeamMemberNum += 1
					teamInfos_range.MaleTeamMember = append(teamInfos_range.MaleTeamMember, teamMemberList_range.TeamMember)
				} else {
					teamInfos_range.FemaleTeamMemberNum += 1
					teamInfos_range.FeMaleTeamMember = append(teamInfos_range.FeMaleTeamMember, teamMemberList_range.TeamMember)
				}
			}
		}
		return_teamInfos = append(return_teamInfos, teamInfos_range)
	}
	//fmt.Println(return_teamInfos)
	return return_teamInfos
}

func getTeamRrcordInfos(teamId string) teamResultInfo_json {
	//fmt.Println("Query team record info")

	var teamResultInfo teamResultInfo_json
	teamResultInfo.TeamName = getTeamNameById(teamId)
	teamResultInfo.TotalScore = 0

	dateAndMajorRoundInfos := getDataAndMajorRound(teamId)
	//fmt.Println(dateAndMajorRoundInfos)

	for _, dateAndMajorRound := range dateAndMajorRoundInfos {
		teamResultInfo.EnemyTeamInfoNum++

		var e EnemyTeamInfo_json

		matchs := getEachMatchInfo(teamId, dateAndMajorRound.date, dateAndMajorRound.majorRound)
		//fmt.Println(matchs)

		var winNum = 0
		var loseNum = 0

		for _, oneMatch := range matchs {
			e.MatchResultNum++

			var m MatchResult_json

			e.Date = dateAndMajorRound.date
			e.MajorRound = dateAndMajorRound.majorRound

			pairInfos := getPairInfos(oneMatch.pairId)
			//fmt.Println(pairInfos)

			enemyTeamInfo := getEnemyTeamInfos(oneMatch.matchId, teamId)
			//fmt.Println(enemyTeamInfo)

			e.EnemyTeamName = enemyTeamInfo.enemyTeamName
			e.EnemyTeamId = enemyTeamInfo.enemyTeamId
			e.EnemyTeamNumber = enemyTeamInfo.enemyTeamNumber

			Enemyscore := getEnemyTeamScore(oneMatch.matchId, teamId)
			//fmt.Println(Enemyscore)

			m.SmallRound = oneMatch.smallRound
			m.OurScore = oneMatch.score
			m.EnemyScore = Enemyscore

			if oneMatch.score != 0 || Enemyscore != 0 {
				if oneMatch.score > Enemyscore {
					winNum++
				} else {
					loseNum++
				}
			}

			enemyPlays := getEnemyPlayerData(oneMatch.matchId, teamId)
			//fmt.Println(enemyPlays)

			for _, player := range pairInfos {
				m.OurPlayerNum++
				m.OurPlayer = append(m.OurPlayer, player.personId)
				m.OurPlayerGender = append(m.OurPlayerGender, player.gender)
			}

			for _, EnemyPlayer := range enemyPlays {
				m.EnemyPlayerNum++
				m.EnemyPlayer = append(m.EnemyPlayer, EnemyPlayer.personId)
				m.EnemyPlayerGender = append(m.EnemyPlayerGender, EnemyPlayer.gender)
			}

			e.MatchResults = append(e.MatchResults, m)
		}
		teamResultInfo.EnemyTeamInfoList = append(teamResultInfo.EnemyTeamInfoList, e)

		if winNum+loseNum == 5 {
			if winNum > loseNum {
				teamResultInfo.TotalScore += 3

			} else {
				teamResultInfo.TotalScore += 1
			}
		}
	}
	teamResultInfo.TeamId = teamId
	//fmt.Println(teamResultInfo)

	return teamResultInfo
}

func handleUpdateScoreSql(TeamId string, Date string, MajorRound string, SmallRound string, Score string) {

	matchId, pairId := getMatchIdAndPairId(TeamId, Date, MajorRound, SmallRound)
	updateScore(matchId, pairId, Score)
}

func getLatestRoundInfo() latestTwoMatchInfo_json {
	matchDate := getLatestMatchDate()
	var infos latestTwoMatchInfo_json
	infos.Date = matchDate
	infos.Round1 = getLatestTeamScoreInfo(matchDate, 1)
	infos.Round2 = getLatestTeamScoreInfo(matchDate, 2)

	return infos
}

func getLatestTeamScoreInfo(matchDate string, round int) []latestMatchInfo_json {
	var LMIs []latestMatchInfo_json

	TeamScores := getTeamScoreInfo(matchDate, round)
	TeamIdList := getLatestTeamList(matchDate, round)

	for _, oneTeam := range TeamIdList {
		var oneteamScore latestMatchInfo_json
		oneteamScore.TeamName = getTeamNameById(oneTeam)
		oneteamScore.TeamId = oneTeam

		for _, oneMatchInfo := range TeamScores {
			if oneMatchInfo.TeamId == oneTeam {
				oneteamScore.Scores = append(oneteamScore.Scores, oneMatchInfo.Score)
			}
		}
		LMIs = append(LMIs, oneteamScore)
	}

	//caculate big score

	for i := 0; i < len(LMIs); {
		var team1Win = 0
		var team2Win = 0
		team1 := LMIs[i]
		team2 := LMIs[i+1]

		for j := 0; j < len(team1.Scores); j++ {
			if team1.Scores[j] == 0 && team2.Scores[j] == 0 {
			} else if team1.Scores[j] > team2.Scores[j] {
				team1Win++
			} else {
				team2Win++
			}
		}
		LMIs[i].Scores = append(LMIs[i].Scores, team1Win)
		LMIs[i+1].Scores = append(LMIs[i+1].Scores, team2Win)
		i = i + 2
	}

	return LMIs
}
