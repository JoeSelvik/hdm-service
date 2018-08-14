Working notes as I develop

## Contenders
* CreateContendersTable()
* CreateContenders()

* GetContenders()
* GetContenderByName()
* PullContendersFromFB()

* UpdateContendersVDData()
* UpdateContenderByName()


## Posts
* CreatePostsTable()
* CreatePosts()

* GetPosts()
* PullPostsFromFB()


## Brackets
* CreateBracketsTable()

* CreateBracket()
* CreateTeams()
* CreateResults()

* GetBracketById()
* GetResultsByRound()
* GetTeams()

* UpdateResults()


## Matchup
* CreateMatchupsTable()

* CreateFirstRoundMatchups(), SecondRound(), ...
* CreateMatchup()

* GetMatchupsByRound()
* GetMatchupByName()

* UpdateMatchupProgress()
* UpdateMatchupVoters()


## Handler interfaces
* GetBracketById()
  * GET /bracket/1

* GetMatchupByName()
  *GET /bracket/1?matchup='firstRound_m1'

* PostVote()
  * POST /matchup?match='firstRound_m1',vote='contenderA'

* GetScore()
  * GET /matchup/1


## On startup
* read config file

* CreateContenders()
  * PullContendersFromFb()
* Create Posts
  * PullPostsFromFb()

* UpdateContendersVariableDependentData()

* CreateBracket()
  * CreateTeams()
  * CreateResults()
  
* CreateFirstRoundMatchups()

* serve endpoints
  * http.HandleFunc("/bracket/", bracketViewHandler)
  * http.HandleFunc("/matchup/", matchupViewHandler)
  * http.ListenAndServe(":8080", nil)


## Planned Manual Operations
* End round
* Create round



## Work on next
* If project is worked on again, it will need a new way to pull contenders and posts from facebook
* re-evaulate http://www.aropupu.fi/bracket/ for front end use
* update bracket and matchup endpoints
