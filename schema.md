# Contenders
* FbId int                        --> INTEGER
* FbGroupId int                   --> INTEGER
* Name string                     --> TEXT
* TotalPosts []int                --> BLOB, vd
* AvgLikesPerPost int             --> INTEGER, vi
* TotalLikesReceived int          --> INTEGER, vd
* TotalLikesGiven int             --> INTEGER, vd
* PostsUsed []int                 --> BLOB, vd
* CreatedAt *time.Time            --> DATETIME
* UpdatedAt *time.Time            --> DATETIME


# Posts
* FbId int                        --> INTEGER
* CreatedAt/UpdatedAt *time.Time  --> DATETIME
* FbGroupId int                   --> INTEGER
* PostedDate *time.Time           --> DATETIME
* Author string                   --> TEXT
* TotalLikes int                  --> INTEGER


# Brackets
* Id int                          --> INTEGER
* CreatedAt/UpdatedAt *time.Time  --> DATETIME
* FbGroupId int                   --> INTEGER
* StartTime *time.Time            --> DATETIME
* EndTime *time.Time              --> DATETIME
* teams [][]string                --> BLOB
* results []interface{}           --> BLOB


# Matchups
* Id int                          --> INTEGER
* CreatedAt/UpdatedAt *time.Time  --> DATETIME
* BracketId int                   --> INTEGER
* Name string                     --> TEXT  // determins round?
* ContenderAId int                --> INTEGER
* APostIds []int                  --> BLOB
* ContenderBId int                --> INTEGER
* BPostIds []int                  --> BLOB
* InProgress bool                 --> BOOL?, vi
* AVotes     int                  --> INTEGER, vd
* BVotes     int                  --> INTEGER, vd
