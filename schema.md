# Contenders
* Id int                          --> INTEGER
* Created/UpdatedTime *time.Time  --> DATETIME
* FbGroupId int                   --> INTEGER
* Name string                     --> TEXT
* TotalPosts []int                --> BLOB, vd
* AvgLikesPerPost int             --> INTEGER, vi
* TotalLikesReceived int          --> INTEGER, vd
* TotalLikesGiven int             --> INTEGER, vd
* PostsUsed []int                 --> BLOB, vd


# Posts
* Id int                          --> INTEGER
* Created/UpdatedTime *time.Time  --> DATETIME
* FbGroupId int                   --> INTEGER
* PostedDate *time.Time           --> DATETIME
* Author string                   --> TEXT
* TotalLikes int                  --> INTEGER


# Brackets
* Id int                          --> INTEGER
* Created/UpdatedTime *time.Time  --> DATETIME
* FbGroupId int                   --> INTEGER
* StartTime *time.Time            --> DATETIME
* EndTime *time.Time              --> DATETIME
* teams [][]string                --> BLOB
* results []interface{}           --> BLOB


# Matchups
* Id int                          --> INTEGER
* Created/UpdatedTime *time.Time  --> DATETIME
* BracketId int                   --> INTEGER
* Name string                     --> TEXT  // determins round?
* ContenderAId int                --> INTEGER
* APostIds []int                  --> BLOB
* ContenderBId int                --> INTEGER
* BPostIds []int                  --> BLOB
* InProgress bool                 --> BOOL?, vi
* AVotes     int                  --> INTEGER, vd
* BVotes     int                  --> INTEGER, vd
