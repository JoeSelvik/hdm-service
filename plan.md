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
* Create Contenders
* Create Posts
* Update Contender Variable data
* Create Bracket
* Create first round matchups
* serve endpoints


# manual
* End round
* Create round



# START
* config and main files

# Next work
* in posts.go, write GetHDMPostsByContenderUsername to be used in matchups.go's CreateFirstRoundMatchups
* clean up todo's in CreateFirstRoundMatchups
* in matchups.go write '(c *Contender) CreateMatchup' and CreateMatchupsTable
* write GetHDMMatchup - follow posts or contenders verisons
* How to keep result score hidden until round ends?
* write EndFirstRoundMatchups()
* FB speak method in main instead of listing all contenders
* clean up main, only main calls GetDBHandle?
* error handling
* add an int totalposts to contenders table for convenience and exploratory work

# HDM Qs
* where to get and store a Post's permalink_url?
* move create table methods to own sql file?
* should I be using UNIQUE in my tables? does this show up with .schema?
* INSERT with new, INSERT with existing, vs INSERT OR REPLACE INTO
* should post.PostedDate be a time.Time instead of string
