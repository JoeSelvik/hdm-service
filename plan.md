# Contenders
* CreateContendersTable()
* CreateContenders()

* GetContenders()
* GetContenderByName()
* PullContendersFromFB()

* UpdateContendersVDData()
* UpdateContenderByName()


# Posts
* CreatePostsTable()
* CreatePosts()

* GetPosts()
* PullPostsFromFB()


# Brackets
* CreateBracketsTable()

* CreateBracket()
* CreateTeams()
* CreateResults()

* GetBracketById()
* GetResultsByRound()
* GetTeams()

* UpdateResults()


# Matchup
* CreateMatchupsTable()

* CreateFirstRoundMatchups(), SecondRound(), ...
* CreateMatchup()

* GetMatchupsByRound()
* GetMatchupByName()

* UpdateMatchupProgress()
* UpdateMatchupVoters()


# Handler interfaces
* GetBracketById()
  * GET /bracket/1

* GetMatchupByName()
  *GET /bracket/1?matchup='firstRound_m1'

* PostVote()
  * POST /matchup?match='firstRound_m1',vote='contenderA'

* GetScore()
  * GET /matchup/1

# readme
* sqlite3 data.db < create_tables.sql


# main
* read config file

* CreateContenders()
  * PullContendersFromFb()
* Create Posts
  * PullPostsFromFb()

* UpdateContendersVariableDependentData()
  * pc.GetPosts()
* UpdateContendersIndependantData()
  * cc.GetContenders()

* CreateBracket()
  * CreateTeams()
  * CreateResults()
  
* CreateFirstRoundMatchups()

* serve endpoints
  * http.HandleFunc("/bracket/", bracketViewHandler)
  * http.HandleFunc("/matchup/", matchupViewHandler)
  * http.ListenAndServe(":8080", nil)


# manual
* End round
* Create round



# START
* finish create post, inserting likes and serializing
* PullPostsFromFb and Create and ReadCollection
* UpdateContendersVariableDependentData

# HDM Qs
* where to get and store a Post's permalink_url?
* should I be using UNIQUE in my tables? does this show up with .schema?

# golang
* optional arguments, like make() length and capacity arguments
