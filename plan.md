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
* cc.Read where to scan values into
* Create ContenderController
* copy scooby application error system?
* clean Sprintf vs log.Println
* address todo's

# HDM Qs
* should facebook_controller panic or return error? conform to resource paradigm?
* add an int totalposts to contenders table for convenience and exploratory work?
* where to get and store a Post's permalink_url?
* move create table methods to own sql file?
* should I be using UNIQUE in my tables? does this show up with .schema?
* INSERT with new, INSERT with existing, vs INSERT OR REPLACE INTO
* should post.PostedDate be a time.Time instead of string
