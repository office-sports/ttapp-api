package data

import (
	"github.com/office-sports/ttapp-api/models"
	"github.com/slack-go/slack"
	"log"
	"strconv"
	"strings"
)

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
	if err != nil {
		panic(err)
	}

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
