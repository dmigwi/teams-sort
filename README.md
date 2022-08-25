## Team Matches Sort

The program scraps data from the provided web url, sorts the date from the recent to the earliest and returns the set number of records.

### Code Installation
```shell
    $ git install github.com/dmigwi/teams-sort

    $ teams-sort --help
        Usage of teams-sort:
        -limit int
                Number of records to be returned starting with the most recent. (default 10)
        -path string
                Database path where sqlite db is to be stored. Defaults to the current path.
        -url string
                Its the web url for data to be analysed, If not set program stops.
```

### Code Testing
```
    $ teams-sort --url=https://www.football-data.co.uk/mmz4281/1920/E0.csv --limit=7 
        Teams Matches sort running, please wait...
        Scrapping data from (https://www.football-data.co.uk/mmz4281/1920/E0.csv)...
        Sorted 380 records returned from the data source
        Setting up the database connection...
        Inserting the data collected into the database...
        Fetching 7 records starting with the most recent...
        Data:  [
        {
            "id": "5ba5f519-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Arsenal",
            "away_team": "Watford",
            "fthg": 3,
            "ftag": 2
        },
        {
            "id": "5ba62923-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Southampton",
            "away_team": "Sheffield United",
            "fthg": 3,
            "ftag": 1
        },
        {
            "id": "5ba65937-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Newcastle",
            "away_team": "Liverpool",
            "fthg": 1,
            "ftag": 3
        },
        {
            "id": "5ba6a15c-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Man City",
            "away_team": "Norwich",
            "fthg": 5,
            "ftag": 0
        },
        {
            "id": "5ba6d123-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Leicester",
            "away_team": "Man United",
            "fthg": 0,
            "ftag": 2
        },
        {
            "id": "5ba7022d-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Everton",
            "away_team": "Bournemouth",
            "fthg": 1,
            "ftag": 3
        },
        {
            "id": "5ba741df-24cc-11ed-86b3-7085c2820f73",
            "div": "E0",
            "date_time": "2020-07-26T16:00:00Z",
            "home_team": "Crystal Palace",
            "away_team": "Tottenham",
            "fthg": 1,
            "ftag": 1
        }
    ]
```