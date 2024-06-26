package bot

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	chg "github.com/gophers-latam/challenges/http"
	"github.com/gophers-latam/challenges/storage"
)

func GetChallenge(level, topic string) (*chg.Challenge, error) {
	var res []chg.Challenge

	err := storage.Get().Find(&res, "level = ? and challenge_type = ? and active = ?", level, topic, 1).Error
	if err != nil {
		return &chg.Challenge{}, err
	}

	l := len(res)
	if l == 0 {
		return &chg.Challenge{}, sql.ErrNoRows
	}

	i, err := intnCrypt(l)

	return &res[i], err
}

func GetFact() (*chg.Fact, error) {
	var res []chg.Fact

	err := storage.Get().Find(&res).Error
	if err != nil {
		return &chg.Fact{}, err
	}

	l := len(res)
	if l == 0 {
		return &chg.Fact{}, sql.ErrNoRows
	}

	i, err := intnCrypt(l)

	return &res[i], err
}

func GetEvents() (*[]chg.Event, error) {
	var res []chg.Event

	err := storage.Get().Find(&res).Error
	if err != nil {
		return &res, err
	}

	l := len(res)
	if l == 0 {
		return &res, sql.ErrNoRows
	}

	return &res, err
}

func GetCommand(cmd string) (*chg.Command, error) {
	var res []chg.Command

	err := storage.Get().Find(&res, "cmd = ?", cmd).Error
	if err != nil {
		return &chg.Command{}, err
	}

	l := len(res)
	if l == 0 {
		return &chg.Command{}, sql.ErrNoRows
	}

	return &res[0], err
}

func GetHours(hour, country string) (string, error) {
	var b bytes.Buffer
	args := strings.Split(hour, ":")
	if len(args) != 2 {
		return "", errors.New("invalid time format. Please use HH:MM format")
	}

	h, err := strconv.Atoi(args[0])
	if err != nil {
		return "", errors.New("invalid hour format")
	}
	m, err := strconv.Atoi(args[1])
	if err != nil {
		return "", errors.New("invalid minute format")
	}

	countryCase := wordCase(country)
	timeZoneInfo, ok := chg.TimeZones[countryCase]
	if !ok {
		return "", errors.New("unknown country")
	}

	loc, err := time.LoadLocation(timeZoneInfo.Timezone)
	if err != nil {
		return "", errors.New("unable to load timezone")
	}

	now := time.Now().UTC()
	inTime := time.Date(now.Year(), now.Month(), now.Day(), h, m, 0, 0, loc)
	originTime := inTime.In(loc)

	tzones := make([]string, 0, len(chg.TimeZones))
	for key := range chg.TimeZones {
		tzones = append(tzones, key)
	}
	sort.Strings(tzones)

	b.WriteString(fmt.Sprintf("🕒 %s **%s**: `%s` hrs\n", timeZoneInfo.Flag, countryCase, inTime.Format("15:04")))
	for _, tz := range tzones {
		if tz == countryCase {
			continue
		}
		loc, err := time.LoadLocation(chg.TimeZones[tz].Timezone)
		if err != nil {
			continue
		}
		lTime := originTime.In(loc)
		b.WriteString(fmt.Sprintf("🕒 %s **%s**: `%s` hrs\n", chg.TimeZones[tz].Flag, tz, lTime.Format("15:04")))
	}

	return b.String(), nil
}
