package utils

import "github.com/rs/zerolog/log"

func LogFatalIfError(err error) {
	if err != nil {
		log.Fatal().Err(err).Msg(err.Error())
	}
}

func LogErrorIfError(err error) {
	if err != nil {
		log.Error().Err(err).Msg(err.Error())
	}
}
