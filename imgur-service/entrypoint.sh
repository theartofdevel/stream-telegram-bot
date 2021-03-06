#!/bin/sh

export IS_IS_DEBUG=$(cat /s/IS_IS_DEBUG)
export IS_IS_DEVELOPMENT=$(cat /s/IS_IS_DEVELOPMENT)
export IS_RABBIT_HOST=$(cat /s/IS_RABBIT_HOST)
export IS_RABBIT_PORT=$(cat /s/IS_RABBIT_PORT)
export IS_RABBIT_USERNAME=$(cat /s/IS_RABBIT_USERNAME)
export IS_RABBIT_PASSWORD=$(cat /s/IS_RABBIT_PASSWORD)
export IS_RABBIT_CONSUMER_QUEUE=$(cat /s/IS_RABBIT_CONSUMER_QUEUE)
export IS_RABBIT_CONSUMER_MBS=$(cat /s/IS_RABBIT_CONSUMER_MBS)
export IS_RABBIT_PRODUCER_QUEUE=$(cat /s/IS_RABBIT_PRODUCER_QUEUE)
export IS_IMGUR_ACCESS_TOKEN=$(cat /s/IS_IMGUR_ACCESS_TOKEN)
export IS_IMGUR_CLIENT_SECRET=$(cat /s/IS_IMGUR_CLIENT_SECRET)
export IS_IMGUR_CLIENT_ID=$(cat /s/IS_IMGUR_CLIENT_ID)
export IS_IMGUR_URL=$(cat /s/IS_IMGUR_URL)
export IS_EVENTWORKS=$(cat /s/IS_EVENTWORKS)
export IS_LOGLEVEL=$(cat /s/IS_LOGLEVEL)

/app
