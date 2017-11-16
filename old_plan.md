# Contenders
* GetFBContenders()
* GetHDMContenders()
* CreateContenderTable(all)
* UpdateContenderTable(all)
* UpdateContender(tx *sql.Tx) (int64, error)
* UpdateHDMContenderDependentData()


# Posts
* createPostsTable(startDate)
* updateTable()


# Brackets
* CreateBracket() Bracket
* CreateInitialTeams()
* CreateInitialResults()
* serialize
* deserialize


# Macthup
* GetMatchup(matchName string) []Match
* CreateFirstRoundMatches(teams [][]string)
* StopFirstFound() teams [][]string


# main func
* CreateContenderTable
* CreatePostTable
* UpdateHDMContenderDependentData

## manual
* CreateBracket()
* CreateFirstRoundMatches(teams [][]string)
* StopFirstFound() teams [][]string
