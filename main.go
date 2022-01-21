package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	URL        = "https://quasarzone.com/bbs/qb_saleinfo"
	USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/99.0.3538.77 Safari/537.36"
)

var lastItem = 0

type collectorOptions struct {
	bot     *tgbotapi.BotAPI
	chanId  int64
	include []string
	exclude []string
}

func OsPanic(err error) {
	if err != nil {
		zlog.Err(err).Msg("Fail to open file")
		panic(err)
	}
}

func collector(options *collectorOptions) *colly.Collector {
	zlog.Info().Msg("Start scraping")
	c := colly.NewCollector(
		colly.UserAgent(USER_AGENT),
	)

	c.OnHTML("div.market-info-list", func(e *colly.HTMLElement) {
		if category := e.ChildText("span.category"); !strings.HasPrefix(category, "PC") {
			return
		}

		titleLink := e.ChildAttr("a", "href")
		splits := strings.Split(titleLink, "/")
		itemNum, err := strconv.Atoi(splits[len(splits)-1])
		if err != nil {
			fmt.Println("Parse error")
		}

		if itemNum <= lastItem {
			return
		}
		lastItem = itemNum
		zlog.Info().Int("LastItemNum", lastItem).Msg("Last Item checked")

		titleName := e.ChildText("span.ellipsis-with-reply-cnt")
		cost := e.ChildText("span.text-orange")

		if !TitleContains(titleName, options.include, options.exclude) {
			return
		}

		fmsg := fmt.Sprintf("New Deals!\nTitle: %s\nCost: %s\nLink: %s%s", titleName, cost, URL, titleLink)
		msg := tgbotapi.NewMessage(options.chanId, fmsg)
		_, err = options.bot.Send(msg)
		if err != nil {
			zlog.Warn().Msg("Send to tg client failed")
		}

		zlog.Info().Int("ItemNum", itemNum).Str("Title", titleName).Str("Cost", cost).Msg("New Deals")
	})

	err := c.Visit(URL)
	if err != nil {
		zlog.Warn().Msg("URL not responding")
	}
	return c
}

func main() {
	zlog.Logger = zlog.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zlog.Info().Msg("Hello world!")
	r := rand.New(rand.NewSource(time.Now().Unix()))

	dat, err := os.ReadFile("env.toml")
	if err != nil {
		zlog.Err(err).Msg("There is no env file in this directory")
		dat, err = os.ReadFile("/run/secrets/env.toml")
		OsPanic(err)
	}

	var env Env
	err = toml.Unmarshal([]byte(string(dat)), &env)
	OsPanic(err)

	options := collectorOptions{
		bot:     InitBot(env.TelApi),
		chanId:  env.TelChan,
		include: env.FilterInclude,
		exclude: env.FilterExclude,
	}

	zlog.Info().Msgf("Including keywords.. %v", env.FilterInclude)
	zlog.Info().Msgf("Excluding keywords.. %v", env.FilterExclude)

	for {
		c := collector(&options)
		c.Wait()
		dur := time.Duration(r.Intn(100)+10) * 120 * 10 * time.Millisecond
		zlog.Info().Float64("Duration", dur.Seconds()).Msg("Waiting...")
		time.Sleep(dur)
	}
}
