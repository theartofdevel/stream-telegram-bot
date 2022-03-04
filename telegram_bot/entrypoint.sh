#!/bin/sh

export ST_BOT_IS_DEBUG=$(cat /s/ST_BOT_IS_DEBUG)
export ST_BOT_IS_DEVELOPMENT=$(cat /s/ST_BOT_IS_DEVELOPMENT)
export ST_BOT_TELEGRAM_TOKEN=$(cat /s/ST_BOT_TELEGRAM_TOKEN)
export ST_BOT_RABBIT_HOST=$(cat /s/ST_BOT_RABBIT_HOST)
export ST_BOT_RABBIT_PORT=$(cat /s/ST_BOT_RABBIT_PORT)
export ST_BOT_RABBIT_USERNAME=$(cat /s/ST_BOT_RABBIT_USERNAME)
export ST_BOT_RABBIT_PASSWORD=$(cat /s/ST_BOT_RABBIT_PASSWORD)
export ST_BOT_RABBIT_CONSUMER_YOUTUBE=$(cat /s/ST_BOT_RABBIT_CONSUMER_YOUTUBE)
export ST_BOT_RABBIT_CONSUMER_IMGUR=$(cat /s/ST_BOT_RABBIT_CONSUMER_IMGUR)
export ST_BOT_RABBIT_CONSUMER_MBS=$(cat /s/ST_BOT_RABBIT_CONSUMER_MBS)
export ST_BOT_RABBIT_PRODUCER_YOUTUBE=$(cat /s/ST_BOT_RABBIT_PRODUCER_YOUTUBE)
export ST_BOT_RABBIT_PRODUCER_IMGUR=$(cat /s/ST_BOT_RABBIT_PRODUCER_IMGUR)
export ST_BOT_EVENT_WORKERS_YT=$(cat /s/ST_BOT_EVENT_WORKERS_YT)
export ST_BOT_EVENT_WORKERS_IMGUR=$(cat /s/ST_BOT_EVENT_WORKERS_IMGUR)
export ST_BOT_LOG_LEVEL=$(cat /s/ST_BOT_LOG_LEVEL)

/app