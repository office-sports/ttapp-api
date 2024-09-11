package data

import (
	"github.com/office-sports/ttapp-api/data/queries"
	"github.com/office-sports/ttapp-api/models"
	"github.com/slack-go/slack"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var digits = map[int]string{
	0: "zero",
	1: "one",
	2: "two",
	3: "three",
	4: "four",
	5: "five",
	6: "six",
	7: "seven",
	8: "eight",
	9: "nine",
}

var digitsOrder = map[int]string{
	1:  "first",
	2:  "second",
	3:  "third",
	4:  "fourth",
	5:  "fifth",
	6:  "sixth",
	7:  "seventh",
	8:  "eight",
	9:  "ninth",
	10: "tenth",
}

// AnnounceSchedule sends message with list of games to be played
func AnnounceSchedule() error {
	config, err := models.GetConfig("")
	if err != nil {
		panic(err)
	}
	md, err := GetOfficeByChannelId(config.MessageConfig.ChannelId)
	if err != nil {
		return err
	}

	rows, err := models.DB.Query(queries.GetTournamentGroupScheduleQuery(), md.ID)
	if err != nil {
		return err
	}

	defer rows.Close()

	// Initialize a map to store grouped game schedules
	groupedSchedules := make(map[string][]models.GameSchedule)

	// Loop through the database rows
	for rows.Next() {
		var gn string // Group name
		gs := new(models.GameSchedule)

		// Scan the row data into variables
		err := rows.Scan(&gn, &gs.HomePlayerName, &gs.AwayPlayerName, &gs.HomePlayerSlackName,
			&gs.AwayPlayerSlackName, &gs.GameWeek, &gs.DateOfMatch, &gs.IsFinished)
		if err != nil {
			return err
		}

		// Append the game schedule to the appropriate group
		groupedSchedules[gn] = append(groupedSchedules[gn], *gs)
	}

	// Convert the groupedSchedules map into a slice of TournamentGroupSchedule
	var tournamentGroups []models.TournamentGroupSchedule
	for groupName, schedules := range groupedSchedules {
		tournamentGroups = append(tournamentGroups, models.TournamentGroupSchedule{
			Name:         groupName,
			GameSchedule: schedules,
		})
	}

	pretext := ":table_tennis_paddle_and_ball: *Pending official matches*"
	text := ""

	for _, group := range tournamentGroups {
		text += "*" + group.Name + "*\n"
		for _, schedule := range group.GameSchedule {
			text += "Game Week " + strconv.Itoa(schedule.GameWeek) + ": "
			text += "<" + schedule.HomePlayerSlackName + "> vs <" + schedule.AwayPlayerSlackName + ">"
			text += "\n"
		}
		text += "\n"
	}

	SendSlackMessage(*config, pretext, text, "")

	return nil
}

func getStartMessageText(result *models.GameResult, config *models.Config) (pretext string, text string) {
	// playoffs match has a different message
	// each match has a name in playoffs, check it
	if result.Name != "" {
		pretext = ":table_tennis_paddle_and_ball: *" + result.GroupName + " Playoffs*, " + result.Name + " started!"
	} else {
		pretext = "*" + result.GroupName + " Group* match started!"
	}
	text = "*" + result.HomePlayerName + "* vs *" + result.AwayPlayerName + "*\n" +
		":eye: <https://" + config.Frontend.Url + "/game/" + strconv.Itoa(result.MatchId) + "/spectate|spectate>"

	return pretext, text
}

func getEndSetMessageText(result *models.GameResult, config *models.Config) (pretext string, text string) {
	setScores := " ("
	for _, s := range result.SetScores {
		setScores += strconv.Itoa(s.Home) + ":" + strconv.Itoa(s.Away) + ", "
	}
	setScores = strings.TrimSuffix(setScores, ", ") + ")"

	var lastSetIndex int
	currentSet := result.CurrentSet
	if currentSet == 1 {
		lastSetIndex = len(result.SetScores) - 1
	} else {
		lastSetIndex = currentSet - 2 // if current set is "2", last index (starting with 0) is current - 2
	}

	var msg string
	isHomeWinner := result.SetScores[lastSetIndex].Home > result.SetScores[lastSetIndex].Away
	if isHomeWinner {
		msg = result.HomePlayerName + " wins set `#" + strconv.Itoa(result.SetScores[lastSetIndex].Set) +
			"` with `" + strconv.Itoa(result.SetScores[lastSetIndex].Home) + ":" +
			strconv.Itoa(result.SetScores[lastSetIndex].Away) + "` score"
	} else {
		msg = result.AwayPlayerName + " wins set `#" + strconv.Itoa(result.SetScores[lastSetIndex].Set) +
			"` with `" + strconv.Itoa(result.SetScores[lastSetIndex].Away) + ":" +
			strconv.Itoa(result.SetScores[lastSetIndex].Home) + "` score"
	}

	pretext = msg
	text = ""

	return pretext, text
}

func getFinalMessageText(result *models.GameResult, config *models.Config) (pretext string, text string) {
	setScores := " ("
	for _, s := range result.SetScores {
		setScores += strconv.Itoa(s.Home) + ":" + strconv.Itoa(s.Away) + ", "
	}
	setScores = strings.TrimSuffix(setScores, ", ") + ")"

	if result.Name != "" {
		result.Level = strings.Replace(result.Level, "|LEAGUE|", result.GroupName, -1)
		if result.WinnerId == result.HomePlayerId {
			result.Level = strings.Replace(result.Level, "|WINNER|", result.HomePlayerName, -1)
		}
		if result.WinnerId == result.AwayPlayerId {
			result.Level = strings.Replace(result.Level, "|WINNER|", result.AwayPlayerName, -1)
		}
		pretext = ":table_tennis_paddle_and_ball: *" + result.GroupName + " Playoffs*, " + result.Name + " finished!\n" +
			result.Level

		text = "*" + result.HomePlayerName + "* " +
			strconv.Itoa(result.HomeScoreTotal) + ":" + strconv.Itoa(result.AwayScoreTotal) + " *" + result.AwayPlayerName +
			"* " + setScores + "\n" +
			"<https://" + config.Frontend.Url + "/game/" + strconv.Itoa(result.MatchId) + "/result|result> | " +
			"<https://" + config.Frontend.Url + "/tournament/" + strconv.Itoa(result.TournamentId) + "/ladders|ladders>"
	} else {
		pretext = "*" + result.GroupName + " Group* match finished"

		text = "*" + result.HomePlayerName + "* " +
			strconv.Itoa(result.HomeScoreTotal) + ":" + strconv.Itoa(result.AwayScoreTotal) + " *" + result.AwayPlayerName +
			"* " + setScores + "\n" +
			"<https://" + config.Frontend.Url + "/game/" + strconv.Itoa(result.MatchId) + "/result|result> | " +
			"<https://" + config.Frontend.Url + "/tournament/" + strconv.Itoa(result.TournamentId) + "/standings|standings>"
	}

	return pretext, text
}

func SendStartMessage(result *models.GameResult) {
	config, err := models.GetConfig("")
	md, err := GetOfficeById(result.OfficeId)
	if err != nil {
		panic(err)
	}

	// Important - overwrite channel id with the value from the db
	config.MessageConfig.ChannelId = *md.ChannelId

	// fetch starting message texts
	pretext, text := getStartMessageText(result, config)

	// break if we do not have config data
	if config.MessageConfig.Hook == "" {
		return
	}

	// send slack message and get the thread ts
	ts := SendSlackMessage(*config, pretext, text, "")

	// update the game to be announced with ts present
	SetTs(result.MatchId, ts)
}

// SendEndSetMessage either sends final score or update
func SendEndSetMessage(result *models.GameResult) {
	config, err := models.GetConfig("")
	md, err := GetOfficeById(result.OfficeId)
	if err != nil {
		panic(err)
	}

	// Important - overwrite channel id with the value from the db
	config.MessageConfig.ChannelId = *md.ChannelId

	if err != nil {
		panic(err)
	}
	if config.MessageConfig.Hook == "" || config.MessageConfig.ChannelId == "" || config.MessageConfig.Token == "" {
		return
	}

	// We need to check if the game is finished
	ann, err := IsAnnounced(result.MatchId)
	if err != nil {
		log.Println("Error fetching announcement: ", err)
	}

	if ann.IsAnnounced == 1 && ann.Ts != "0" {
		// we should have a thread id to post the message to (Ts)
		pretext, text := getEndSetMessageText(result, config)
		SendSlackMessage(*config, pretext, text, ann.Ts)
	}

	if result.IsFinished == 1 {
		pretext, text := getFinalMessageText(result, config)

		if ann.Ts == "0" {
			// There was no manual scoring and there is no thread id
			// which means we use scores form
			SendSlackMessage(*config, pretext, text, ann.Ts)
		} else {
			// The manual scoring was started and the thread id is present
			// so the original message needs to be update
			UpdateSlackMessage(*config, pretext, text, ann.Ts)
		}

		// Set announcement fields to final state
		SetAnnounced(result.MatchId, 1, "0")
	}

	if err != nil {
		panic(err)
	}
}

func SendSlackMessage(config models.Config, pretext string, text string, thread string) string {
	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(config.MessageConfig.Token, slack.OptionDebug(false))
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    text,
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
	}

	var ts string

	if thread != "" {
		_, _, _ = client.PostMessage( // resp, ts, err
			config.MessageConfig.ChannelId,
			slack.MsgOptionTS(thread),
			slack.MsgOptionAttachments(attachment),
		)
	} else {
		_, ts, _ = client.PostMessage( // resp, ts, err
			config.MessageConfig.ChannelId,
			slack.MsgOptionAttachments(attachment),
		)
	}

	return ts
}

func UpdateSlackMessage(config models.Config, pretext string, text string, thread string) {
	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(config.MessageConfig.Token, slack.OptionDebug(false))
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    text,
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
	}
	// PostMessage will send the message away.
	// First parameter is just the channelID, makes no sense to accept it
	_, _, _, err := client.UpdateMessage( // resp, ts, err
		config.MessageConfig.ChannelId,
		thread,
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
}

func CreateRecapMessage(recap []*models.GroupInfo) {
	totalGames := 0
	totalGamesRemaining := 0
	positionsUp := 0
	positionsStay := 0
	positionsString := ""

	for _, group := range recap {
		totalGames += group.GamesPlayed
		totalGamesRemaining += group.GamesRemaining
		positionsUp += group.PositionsUp
		positionsStay += group.PositionsStay
		positionsString += strconv.Itoa(group.GamesPlayed) + " in " + group.Name + ", "
	}

	for _, group := range recap {
		msg := ""

		if group.PositionsUp == 0 {
			msg += "There were no position changes over the course of the week. "
		} else {
			msg += getPositionsUpMessage(group)
			msg += getPositionsStayMessage(group)
			msg += getPositionsTopDropMessage(group)
			msg += getPositionsTopClimbMessage(group)
		}

		group.StatsMessage = msg

		msg = getSpotsMessage(group)
		group.CandidatesMessage = msg
	}
}

func getPositionsUpMessage(group *models.GroupInfo) string {
	sClimb := []string{
		"During the last week, players changed their table position a total of |NUM| times. ",
		"Over the course of the week, the players changed their table position |NUM| times. ",
		"Throughout last week's games, the players changed their table position |NUM| times. ",
		"Last week's matches were dynamic, with the players changing their table position |NUM| times throughout the week. "}

	msg := getRandomMessage(sClimb)
	if group.PositionsUp < 10 {
		msg = strings.Replace(msg, "|NUM|", digits[group.PositionsUp], -1)
	} else {
		msg = strings.Replace(msg, "|NUM|", strconv.Itoa(group.PositionsUp), -1)
	}

	return msg
}

func getPositionsStayMessage(group *models.GroupInfo) string {
	sStay := []string{
		"In case of |NUM| players, position did not change. ",
		"Despite some close games, |NUM| players managed to maintain their table position. ",
		"We saw |NUM| competitors who remained in the same position. "}

	msg := getRandomMessage(sStay)
	msg = strings.Replace(msg, "|NUM|", digits[group.PositionsStay], -1)

	return msg
}

func getPositionsTopDropMessage(group *models.GroupInfo) string {
	sBiggestDrop := []string{
		"It was a tough week for |PLAYERS|, who dropped the most by |NUM| places. ",
		"The most unfortunate movement in the rankings was by |PLAYERS|, dropping |NUM| places. ",
		"It was a disappointing week for |PLAYERS| in particular, who dropped the most in table position (|NUM| places) among all the participants. "}

	msg := getRandomMessage(sBiggestDrop)
	msg = strings.Replace(msg, "|PLAYERS|", "|="+group.TopDropPlayerName+"=|", -1)
	msg = strings.Replace(msg, "|NUM|", digits[group.TopDrop], -1)
	if group.TopDrop == 1 {
		msg = strings.Replace(msg, "places", "place", -1)
	}

	return msg
}

func getPositionsTopClimbMessage(group *models.GroupInfo) string {
	sBiggestClimb := []string{
		"Out of all the players, |PLAYERS| made the most progress and climbed |NUM| up in table over the course of the week. ",
		"Throughout last week's games, |PLAYERS| showed significant improvement and advanced |NUM| positions. ",
		"We saw some incredible performances, and |PLAYERS| in particular climbed |NUM| in table position. "}

	msg := getRandomMessage(sBiggestClimb)

	msg = strings.Replace(msg, "|PLAYERS|", "|="+group.TopClimbPlayerName+"=|", -1)
	msg = strings.Replace(msg, "|NUM|", digits[group.TopClimb], -1)
	if group.TopDrop == 1 {
		msg = strings.Replace(msg, "places", "place", -1)
		msg = strings.Replace(msg, "positions", "position", -1)
	}

	return msg
}

func getSpotsMessage(group *models.GroupInfo) string {
	noSpots := []string{
		"With just a few games left in the season, no competitor has secured a place in the playoffs yet. ",
		"Despite their strong performances so far, none of the competitors have a secured spot in the playoffs at this point. ",
		"The competition is intense, and there are no guaranteed places for the playoffs in the current table. ",
		"With several teams still in contention for the playoffs, no one has secured a place in the top positions of the table. ",
		"It's anyone's game at this point, as there are no secured places for the playoffs in the current standings. "}

	// count secured spots
	securedSpots := 0
	securedSpotsPlayoffs := 0
	var securedSpotsPlayers []string
	var securedSpotsPlayoffsPlayers []string
	var securedSpotsPositions []string
	securedSpotsPlayersNames := ""
	securedSpotsPlayoffsPlayersNames := ""
	for position, p := range group.PositionCandidates {
		if p.Secured == 1 {
			securedSpots++
			if len(p.PlayerNames) == 1 {
				securedSpotsPlayers = append(securedSpotsPlayers, p.PlayerNames[0])
			}

			// check if this is a playoffs spot
			if position <= group.GroupPromotions {
				securedSpotsPlayoffs++
				if len(p.PlayerNames) == 1 {
					securedSpotsPlayoffsPlayers = append(securedSpotsPlayoffsPlayers, p.PlayerNames[0])
					securedSpotsPositions = append(securedSpotsPositions, digitsOrder[position])
				}
			}
		}
	}

	// there might not be any secured spots (fixed) but there can be players
	// with sure advance to playoffs, those are with minimum position of advance group size
	if securedSpotsPlayoffs == 0 {
		for _, p := range group.PlayerInfo {
			if p.PositionMin <= group.GroupPromotions {
				securedSpotsPlayoffs++
				securedSpotsPlayoffsPlayers = append(securedSpotsPlayoffsPlayers, p.Name)
			}
		}
	}

	securedSpotsPlayersNames = strings.Join(securedSpotsPlayers, ", ")
	securedSpotsPlayoffsPlayersNames = strings.Join(securedSpotsPlayoffsPlayers, ", ")

	if group.GamesPlayed < len(group.PlayerInfo)/2 {
		return ""
	}

	// |NUM| of the most skilled
	msg := "We'll have total of " + digits[group.GroupPromotions] + " players advancing to playoffs from |GROUP|. "
	msg = strings.Replace(msg, "|GROUP|", group.Name, -1)

	if securedSpots == 0 {
		if securedSpotsPlayoffs == 0 {
			msg += getRandomMessage(noSpots)
		} else {
			msg += strings.Title(strings.ToLower(digits[securedSpotsPlayoffs])) + " already secured spot in playoffs ladder: |=" +
				securedSpotsPlayoffsPlayersNames + "=|, however final position in standings depend on results or remaining games. "
		}
	} else {
		if securedSpots == 1 {
			msg += "There is only " + digits[securedSpots] + " final position in the table so far: |=" +
				securedSpotsPlayersNames + "=|. "
			if securedSpotsPlayoffs > 0 {
				msg += "Player will advance to playoffs from " + strings.Join(securedSpotsPositions, ",") + " position in the table. "
			}
		} else {
			msg += "There are " + digits[securedSpots] + " final positions in the table for |=" +
				securedSpotsPlayersNames + "=|. "

			if securedSpotsPlayoffs > 0 {
				if securedSpots == securedSpotsPlayoffs {
					msg += "All of those competitors "
				} else {
					msg += strings.Title(strings.ToLower(digits[securedSpotsPlayoffs]))
					msg += " competitors: "
					msg += "|=" + securedSpotsPlayoffsPlayersNames + "=| "
				}
				if len(securedSpotsPositions) != 0 {
					msg += "advance to playoffs from " +
						strings.Join(securedSpotsPositions, ", ") +
						" positions in table, respectively. "
				} else {
					msg += "advance to playoffs, however their positions are still a subject to change. "
				}
				msg += "Congratulations! "
			}
		}
	}

	noPromo := []string{
		"Despite their best efforts, |PLAYERS| will not be advancing to the playoffs and will be fighting for |POSITION| position in the table. ",
		"It's been a tough season for |PLAYERS| and unfortunately |PLURAL_COMPETITORS| will not be moving on to the playoffs, but instead will be battling for |POSITION| position in the standings. ",
		"Although falling short of making the playoffs, |PLAYERS| |PLURAL_BE| determined to fight for |POSITION| position in the table. ",
		"It's a disappointing outcome for |PLAYERS|, who will not be advancing to the playoffs, and instead will be fighting for |POSITION| position in the standings. ",
		"While not having made it to the playoffs this season, |PLAYERS| |PLURAL_BE| not giving up and will be competing fiercely to take |POSITION| position in the table. "}

	promo := []string{
		"|PLAYERS| |PLURAL_BE| in a tough battle for |POSITION| position in the table, and |PLURAL_BE| determined to secure their spot in the playoffs. ",
		"With just a few games left in the season, |PLAYERS| |PLURAL_BE| fighting for a promotion to the playoffs from |POSITION| position in the table. ",
		"|PLAYERS| |PLURAL_BE| still very much in the playoff race and |PLURAL_BE| fighting hard to move up the table to secure |POSITION| position. ",
		"It's a close race for |POSITION| position, but |PLAYERS| |PLURAL_BE| not backing down and |PLURAL_BE| doing everything possible to secure the spot in the postseason. ",
		"The competition is fierce, but |PLAYERS| |PLURAL_BE| up for the challenge and |PLURAL_BE| focused on fighting for promotion to the playoffs from |POSITION| position in the standings. ",
		"Best case scenario for |PLAYERS| is |POSITION| position in standings. "}

	maxPos := []string{
		"With just a few games left in the season, this player surely will finish at least in |POSITION_MIN| position. ",
		"This player might end up on |POSITION_MIN| position in the final standings as well. ",
		"The bottom for this player is |POSITION_MIN|. ",
		"The worst scenario for this player would be |POSITION_MIN| position. "}

	for position, p := range group.PositionCandidates {
		if p.Secured != 0 {
			continue
		}

		if position <= group.GroupPromotions {
			msg += getRandomMessage(promo)

			if len(p.PlayerNames) == 1 {
				msg = strings.Replace(msg, "|PLURAL_BE|", "is", -1)
				msg = strings.Replace(msg, "|PLURAL_COMPETITORS|", "this competitor", -1)

				// get minimum position player can end up taking
				psi := getPlayerSituationInfo(strings.Join(p.PlayerNames, ""), group.PlayerInfo)
				msgMin := getRandomMessage(maxPos)
				msg += strings.Replace(msgMin, "|POSITION_MIN|", digitsOrder[psi.PositionMin], -1)
			} else {
				msg = strings.Replace(msg, "|PLURAL_BE|", "are", -1)
				msg = strings.Replace(msg, "|PLURAL_COMPETITORS|", "these competitors", -1)
			}

			msg = strings.Replace(msg, "|PLAYERS|", "|="+strings.Join(p.PlayerNames, ", ")+"=|", -1)
			msg = strings.Replace(msg, "|POSITION|", digitsOrder[position], -1)
		} else {
			msg += getRandomMessage(noPromo)
			msg = strings.Replace(msg, "|PLAYERS|", "|="+strings.Join(p.PlayerNames, ", ")+"=|", -1)
			msg = strings.Replace(msg, "|POSITION|", digitsOrder[position], -1)
			if len(p.PlayerNames) == 1 {
				msg = strings.Replace(msg, "|PLURAL_BE|", "is", -1)
				msg = strings.Replace(msg, "|PLURAL_COMPETITORS|", "this competitor", -1)
			} else {
				msg = strings.Replace(msg, "|PLURAL_BE|", "are", -1)
				msg = strings.Replace(msg, "|PLURAL_COMPETITORS|", "these competitors", -1)
			}
		}
	}

	return msg
}

func getPlayerSituationInfo(name string, playerInfo []models.PlayerInfo) *models.PlayerInfo {
	for _, p := range playerInfo {
		if p.Name == name {
			return &p
		}
	}

	return nil
}

func getRandomMessage(msgs []string) string {
	rand.Seed(time.Now().UnixNano())
	msg := msgs[rand.Intn(len(msgs))]

	return msg
}
