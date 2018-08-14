# HerpDerpMadness Service

A fun project to hone my Go and backend development skills. Developed to use with a long-standing Facebook group of
high school friends and their history silly posts.

Using Facebook's [Graph API](https://developers.facebook.com/docs/graph-api/), pull a Facebook group's members
and feed. Gather various data from each and create a "Bracket" similar to college basketball's March Madness tournament.
"Contenders" are paired against each other per-round. Their "Matchup" displays five random "Posts" from each contender and anyone can vote
for which contender had the best posts.

For the front end I originally planned to use [jquery-bracket](http://www.aropupu.fi/bracket/). However, as the
project dragged on I started to learn Vue. If continued I may switch to this [vue-tournament-bracket](https://github.com/antonionoca/vue-tournament-bracket)
repo.

Unfortunately, after the Cambridge Analytica scandal Facebook removed many of the Graph API endpoints I used to gather
data from. Instead of finding a new way (if even possible) I've decided to retire the project and move on. It is a well
thought out effort though and I currently use it as a reference for other work.

## Running

1. Install Dependencies
2. Create Configuration
3. Setup Database
4. Run

### Install Dependencies
Just kidding, there aren't any yet.

### Create Configuration
hdm-service can be configured via environment variables or a config file. The default is a `config.json`
file in the project's root directory. I'm lazy and committed the one I use while developing. *This is only okay because
there is no sensitive data in it yet.* The access token is long expired and this specific group id is for a
secret group.

* fb_access_token - Grab one for development with the [Graph API Explorer](https://developers.facebook.com/tools/explorer/)
* fb_group_id - Find this in the url for the group you'd like to use
* the db related fields are currently only supported for you project's root
* start_time and end_time are for which posts to gather from the group feed

```
{
  "fb_access_token": "EAAKweAs72uUBAJa2oNfSEG45MV8cYMxhBwakvbQI2z4MwlLWFqLM4BmYkY1HlygPOmjb17ddIoevcfEe6M3J6kpfFIqbg2TGlwdpOCZBZC7BHtTW6Sl5DgIWRICo79gw5vAcigR7NV9Xtqaa0WFmpA1PB4I9aIvPhTSf0xhm08Htm58zxUEFSFR2nF3yMZD",
  "fb_group_id": 208678979226870,
  "db_path": "/Users/joeselvik/projects/go/src/github.com/JoeSelvik/hdm-service/data.db",
  "db_setup_script": "/Users/joeselvik/projects/go/src/github.com/JoeSelvik/hdm-service/create_tables.sql",
  "db_test_path": "/Users/joeselvik/projects/go/src/github.com/JoeSelvik/hdm-service/test.db",
  "start_time": "2017-01-01T00:00:00+0000",
  "end_time": "2017-12-31T00:00:00+0000"
}
```

### Setup Database
So far only two tables are needed, if project is completed this will become more complex.

```
$ sqlite3 data.db < create_tables.sql
```

### Run
```
$ hdm-service
```

## Tests
```
$ go test
```
