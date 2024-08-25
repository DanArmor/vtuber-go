# vtuber-go

- [About](#about)
- [Running instance](#running-instance)
- [How to selfhost](#how-to-selfhost)
    - [Deploy](#deploy)
    - [Configuration](#configuration)
    - [How to populate data](#how-to-populate-data)

## About

Telegram bot that reminds users about streams of selected vtubers.

![message example](/backend/internal/resources/photos/files/1.jpg)

It has user interface implemented via Telegram Mini App mechanism. You can access it by pressing 'Menu' button in the lower left corner.

Select the streamers whose broadcasts you want to be notified about. You can use filters by:
- *name*;
- *company*;
- *wave*;
- *all/selected/not selected property*.

Also in "Show" tab you can hide tags of vtubers and avatars.

![screen example](/backend/internal/resources/photos/files/2.jpg)

In "Settings" tab you can select the shift from GMT of your local time to get proper "start time" value. If you didn't select any - it will format start time as GMT +0.

![screen example](/backend/internal/resources/photos/files/3.jpg)

Also you can click on streamer in list and got streamer card with links to their accounts on youtube/twitch/twitter and etc.
![screeb example](/backend/internal/resources/photos/files/4.jpg)

To get help message send `/start` command to the bot.

Can be used as learning example of Telegram Mini App

Based on example from [vkruglikov/react-telegram-web-app](https://github.com/vkruglikov/react-telegram-web-app)

Built with Golang + Ent + Gin / Typescript + React + Ant Design.

## Running instance

[Link to running instance](https://t.me/vtuber_go_bot).

Sends a notification 30 minutes before the stream.

## How to selfhost

### Deploy

You need 3 things:
- [Holodex API key](https://holodex.net/);
- [Telegram Bot (its token)](https://t.me/botfather);
- Server to run on with the access to external users to connect to it (to access web interface. Also you want some domain to got certifciate to run web interface over https).

When you setuped your server, you need to clone the repo to it, enter `.docker` dir, change .env and config.yaml according to your needs and run:
```
docker compose up -d
```
Done - it will build all the images and run the application.

### Configuration

Clarification of some config.yaml values:

```
base_path - the path after hostname and port to bind backend to.
expiration_hours - expiration time in hours for jwt token
admin_token - token for administration endpoints
time_notify_after - time threshold in minutes after which the notification should be sent
time_step - time step in minutes for check streams time. If you set it to 2 - it will check time of streams every 2 minutes.
```

Also don't forget to change backend url in `frontend/src/config.ts` accordingly to your setup.

Also you need to create empty `secrets.ts` file in `frontend/src` directory.

### How to populate data

You can use Holodex API to get data about vtubers. Example:

```shell
curl --request GET \
--header "X-APIKEY: <YOUR_API_KEY>" "https://holodex.net/api/v2/channels?org=Hololive&limit=50" > out.txt
```

And use data from out.txt to send it into Vtuber-Go (`/api/admin/vtubers`) with curl or [Postman](https://www.postman.com/) in format:

header: `ADMIN-TOKEN: <YOUR_ADMIN_TOKEN>`

body:

```json
{
    "vtubers": <output of the /channels endpoint of Holodex API>
}
```

At the moment you can use only vtubers, that are available with Holodex, cause we use Holodex API to check available streams.