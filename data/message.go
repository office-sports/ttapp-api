package data

import (
	"github.com/office-sports/ttapp-api/models"
	"github.com/slack-go/slack"
	"strconv"
	"strings"
)

func SendMessage(result *models.GameResult) {
	config, err := models.GetConfig("")
	if err != nil {
		panic(err)
	}

	setScores := " ("
	for _, s := range result.SetScores {
		setScores += strconv.Itoa(s.Home) + ":" + strconv.Itoa(s.Away) + ", "
	}
	setScores = strings.TrimSuffix(setScores, ", ") + ")"

	pretext := "*" + result.GroupName + " Group* match finished"
	txt := "*" + result.HomePlayerName + "* " +
		strconv.Itoa(result.HomeScoreTotal) + ":" + strconv.Itoa(result.AwayScoreTotal) + " *" + result.AwayPlayerName +
		"* " + setScores + "\n" +
		"<https://" + config.Frontend.Url + "/game/" + strconv.Itoa(result.MatchId) + "/result|result> | " +
		"<https://" + config.Frontend.Url + "/tournament/" + strconv.Itoa(result.TournamentId) + "/standings|standings>"

	SetAnnounced(result.MatchId)

	if config.MessageConfig.Hook == "" {
		return
	}

	// Create a new client to slack by giving token
	// Set debug to true while developing
	client := slack.New(config.MessageConfig.Token, slack.OptionDebug(true))
	attachment := slack.Attachment{
		Pretext: pretext,
		Text:    txt,
		// Color Styles the Text, making it possible to have like Warnings etc.
		Color: "#36a64f",
	}
	// PostMessage will send the message away.
	// First parameter is just the channelID, makes no sense to accept it
	_, _, err = client.PostMessage( // resp, ts, err
		config.MessageConfig.ChannelId,
		// uncomment the item below to add an extra Header to the message, try it out :)
		//slack.MsgOptionText("New message from bot", false),
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
}
