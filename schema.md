Working notes on the data schema

## Contenders
* FbId int                        --> INTEGER
* FbGroupId int                   --> INTEGER
* Name string                     --> TEXT
* Posts []string                  --> BLOB, vd
* AvgLikesPerPost float64         --> INTEGER, vi
* TotalLikesReceived int          --> INTEGER, vd
* TotalLikesGiven int             --> INTEGER, vd
* PostsUsed []string              --> BLOB, vd
* CreatedAt time.Time             --> DATETIME
* UpdatedAt time.Time             --> DATETIME


## Posts
* FbId string                     --> TEXT
* FbGroupId int                   --> INTEGER
* PostedDate time.Time            --> DATETIME
* AuthorFbId int                  --> INTEGER
* Likes []int                     --> BLOB
* CreatedAt time.Time             --> DATETIME
* UpdatedAt time.Time             --> DATETIME


## Brackets
* Id int                          --> INTEGER
* FbGroupId int                   --> INTEGER
* StartTime time.Time             --> DATETIME
* EndTime time.Time               --> DATETIME
* teams [][]string                --> BLOB
* results []interface{}           --> BLOB
* CreatedAt time.Time             --> DATETIME
* UpdatedAt time.Time             --> DATETIME


## Matchups
* Id int                          --> INTEGER
* BracketId int                   --> INTEGER
* Name string                     --> TEXT  // determines round?
* ContenderAId int                --> INTEGER
* APostIds []int                  --> BLOB
* ContenderBId int                --> INTEGER
* BPostIds []int                  --> BLOB
* InProgress bool                 --> BOOL?, vi
* AVotes     int                  --> INTEGER, vd
* BVotes     int                  --> INTEGER, vd
* CreatedAt time.Time             --> DATETIME
* UpdatedAt time.Time             --> DATETIME
