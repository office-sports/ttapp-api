package data

import (
	"bytes"
	"encoding/json"
	"github.com/office-sports/ttapp-api/models"
	"net/http"
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

	txt := "> *" + result.GroupName + " Group* match finished\n" +
		"> *" + result.HomePlayerName + "* " +
		strconv.Itoa(result.HomeScoreTotal) + ":" + strconv.Itoa(result.AwayScoreTotal) + " *" + result.AwayPlayerName +
		"* " + setScores + "\n"

	payload := map[string]string{
		"text":        txt,
		"channel":     config.MessageConfig.ChannelId,
		"method":      "post",
		"contentType": "application/json",
		"username":    "tabletennisbot",
		"icon_emoji":  ":table_tennis_paddle_and_ball:",
	}
	jsonValue, _ := json.Marshal(payload)

	SetAnnounced(result.MatchId)

	if config.MessageConfig.Hook == "" {
		return
	}

	_, err = http.Post(config.MessageConfig.Hook, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
}
