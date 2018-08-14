package main

import (
	"database/sql"
	"encoding/json"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"net/http"
	"sort"
	"time"
)

type JSBracket struct {
	Teams   [][]interface{} `json:"teams"`
	Results []interface{}   `json:"results"`
}

type Bracket struct {
	Id        int
	Teams     []TeamPair       `json:"teams"`
	Results   SixtyFourResults `json:"results"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// teams := make([][]interface{}, 32)
// teams[0] = []interface{}{"hank", nil}
type TeamPair struct {
	ContenderAName string
	ContenderBName string
}

// results := make([]interface{}, 3)
// .. firstRound := make([][]interface{}, 4)
// .. firstRound[0] = []interface{}{1, 0, "firstRound-m0"}
// results[0] = firstRound
// ...
type SixtyFourResults struct {
	FirstRound   [][]interface{}
	SecondRound  [][]interface{}
	SweetSixteen [][]interface{}
	EliteEight   [][]interface{}
	FinalFour    [][]interface{}
	Championship [][]interface{}
}

func (b *Bracket) DBTableName() string {
	return "brackets"
}

func (b *Bracket) Path() string {
	return "/brackets/"
}

// Returns the TeamPair struct in the format of an arrary of names
func (t *TeamPair) serialize() []interface{} {
	return []interface{}{t.ContenderAName, t.ContenderBName}
}

// Return a JSBracket from Bracket. When Marshalled, JSBracket is properly formatted for jQueryBracket
func (b *Bracket) serialize() *JSBracket {
	var bracketJS JSBracket

	// Teams needs to be a list of arrays
	// [["joe","matt"], ["jim","mike"], ...]
	teamJS := make([][]interface{}, 32)
	for i := 0; i < len(b.Teams); i++ {
		t := TeamPair{b.Teams[i].ContenderAName, b.Teams[i].ContenderBName}
		teamJS[i] = t.serialize()
	}
	bracketJS.Teams = teamJS

	// Results is a multi-dimension list
	// first list contains a list of main and a list of consolation results
	// winner list contains a list for each round
	// each round contains a list of each game
	// [
	// 	  [ // main
	// 	    [[1,0,"g1"],[null,null,"g2"],[null,null,"g3"],[null,null,"g4"]],  // round 1
	// 		[[null,null,"g5"],[null,null,"g6"]],  // round 2
	// 		[[null,null,"g7"],[null,null,"g8"]]  // round 3
	// 	  ]
	// ] // no consolation
	resultJS := make([]interface{}, 6)
	resultJS[0] = b.Results.FirstRound
	resultJS[1] = b.Results.SecondRound
	resultJS[2] = b.Results.SweetSixteen
	resultJS[3] = b.Results.EliteEight
	resultJS[4] = b.Results.FinalFour
	resultJS[5] = b.Results.Championship
	bracketJS.Results = resultJS

	return &bracketJS
}

func (b *Bracket) CreateBracket(tx *sql.Tx) (int64, error) {
	q := `
	INSERT INTO brackets (
		Id,
		Teams,
		Results,
		CreatedAt,
		UpdatedAt
	) values (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
	`
	teams, err := json.Marshal(b.Teams)
	if err != nil {
		return 0, err
	}

	results, err := json.Marshal(b.Results)
	if err != nil {
		return 0, err
	}

	res, err := tx.Exec(q, b.Id, teams, results)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// CreateBracketTable creates the brackets table if it does not exist
func CreateBracketsTable(db *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS brackets(
		Id INT NOT NULL,
		Teams BLOB,
		Results BLOB,
		CreatedAt DATETIME,
		UpdatedAt DATETIME
	);
	`
	_, err := db.Exec(q)
	if err != nil {
		log.Println("Failed to CREATE brackets table")
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("Failed to BEGIN txn:", err)
		return err
	}
	defer tx.Rollback()

	bracket, _ := GenerateInitialBracket()
	_, _ = bracket.CreateBracket(tx)

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		log.Println("Failed to COMMIT txn:", err)
		return err
	}

	return nil
}

func GenerateInitialBracket() (*Bracket, error) {
	teams, _ := CreateInitialTeams()
	results, _ := CreateInitialResults()
	bracket := Bracket{1, *teams, *results, time.Now(), time.Now()}
	return &bracket, nil
}

func CreateInitialTeams() (*[]TeamPair, error) {
	log.Println("Creating initial teams")

	//contenders, _ := GetHDMContenders(db)
	var contenders map[int]*models.Contender

	// seedContenders
	sortedContenders := make(models.ContenderSlice, 0, len(contenders))
	for _, c := range contenders {
		sortedContenders = append(sortedContenders, c)
	}
	sort.Sort(sortedContenders)

	// pairContenders
	teams := make([]TeamPair, 32)

	// East - 1, top left
	teams[0] = TeamPair{sortedContenders[0].Name, ""}
	teams[1] = TeamPair{sortedContenders[24].Name, sortedContenders[36].Name}
	teams[2] = TeamPair{sortedContenders[20].Name, sortedContenders[44].Name}
	teams[3] = TeamPair{sortedContenders[4].Name, ""}
	teams[4] = TeamPair{sortedContenders[8].Name, ""}
	teams[5] = TeamPair{sortedContenders[16].Name, sortedContenders[40].Name}
	teams[6] = TeamPair{sortedContenders[28].Name, sortedContenders[32].Name}
	teams[7] = TeamPair{sortedContenders[12].Name, ""}

	// West - 2, bottom left
	teams[8] = TeamPair{sortedContenders[2].Name, ""}
	teams[9] = TeamPair{sortedContenders[26].Name, sortedContenders[38].Name}
	teams[10] = TeamPair{sortedContenders[22].Name, "46th seed"} // 46?
	teams[11] = TeamPair{sortedContenders[6].Name, ""}
	teams[12] = TeamPair{sortedContenders[10].Name, ""}
	teams[13] = TeamPair{sortedContenders[18].Name, sortedContenders[42].Name}
	teams[14] = TeamPair{sortedContenders[30].Name, sortedContenders[34].Name}
	teams[15] = TeamPair{sortedContenders[14].Name, ""}

	// Midwest - 3, top right
	teams[16] = TeamPair{sortedContenders[1].Name, ""}
	teams[17] = TeamPair{sortedContenders[25].Name, sortedContenders[37].Name}
	teams[18] = TeamPair{sortedContenders[21].Name, sortedContenders[45].Name}
	teams[19] = TeamPair{sortedContenders[5].Name, ""}
	teams[20] = TeamPair{sortedContenders[9].Name, ""}
	teams[21] = TeamPair{sortedContenders[17].Name, sortedContenders[41].Name}
	teams[22] = TeamPair{sortedContenders[29].Name, sortedContenders[33].Name}
	teams[23] = TeamPair{sortedContenders[13].Name, ""}

	// South - 4, bottom right
	teams[24] = TeamPair{sortedContenders[3].Name, ""}
	teams[25] = TeamPair{sortedContenders[27].Name, sortedContenders[39].Name}
	teams[26] = TeamPair{sortedContenders[23].Name, "47th seed"} // 47?
	teams[27] = TeamPair{sortedContenders[7].Name, ""}
	teams[28] = TeamPair{sortedContenders[11].Name, ""}
	teams[29] = TeamPair{sortedContenders[19].Name, sortedContenders[43].Name}
	teams[30] = TeamPair{sortedContenders[31].Name, sortedContenders[35].Name}
	teams[31] = TeamPair{sortedContenders[15].Name, ""}

	// for i, c := range sortedContenders {
	// 	fmt.Println(i, c.Name, c.AvgLikesPerPost, c.TotalLikesReceived, c.TotalLikesGiven)
	// }
	// fmt.Println(teams)

	return &teams, nil
}

func CreateInitialResults() (*SixtyFourResults, error) {
	log.Println("Creating initial results")

	// firstRound
	firstRound := make([][]interface{}, 32)
	firstRound[0] = []interface{}{1, 0, "firstRound_0"}
	firstRound[1] = []interface{}{3, 4, "firstRound_1"}
	firstRound[2] = []interface{}{5, 69, "firstRound_2"}
	firstRound[3] = []interface{}{1, 0, "firstRound_3"}
	firstRound[4] = []interface{}{nil, nil, "firstRound_4"}
	firstRound[5] = []interface{}{nil, nil, "firstRound_5"}
	firstRound[6] = []interface{}{nil, nil, "firstRound_6"}
	firstRound[7] = []interface{}{nil, nil, "firstRound_7"}

	firstRound[8] = []interface{}{nil, nil, "firstRound_8"}
	firstRound[9] = []interface{}{nil, nil, "firstRound_9"}
	firstRound[10] = []interface{}{nil, nil, "firstRound_10"}
	firstRound[11] = []interface{}{nil, nil, "firstRound_11"}
	firstRound[12] = []interface{}{nil, nil, "firstRound_12"}
	firstRound[13] = []interface{}{nil, nil, "firstRound_13"}
	firstRound[14] = []interface{}{nil, nil, "firstRound_14"}
	firstRound[15] = []interface{}{nil, nil, "firstRound_15"}

	firstRound[16] = []interface{}{nil, nil, "firstRound_16"}
	firstRound[17] = []interface{}{nil, nil, "firstRound_17"}
	firstRound[18] = []interface{}{nil, nil, "firstRound_18"}
	firstRound[19] = []interface{}{nil, nil, "firstRound_19"}
	firstRound[20] = []interface{}{nil, nil, "firstRound_20"}
	firstRound[21] = []interface{}{nil, nil, "firstRound_21"}
	firstRound[22] = []interface{}{nil, nil, "firstRound_22"}
	firstRound[23] = []interface{}{nil, nil, "firstRound_23"}

	firstRound[24] = []interface{}{nil, nil, "firstRound_24"}
	firstRound[25] = []interface{}{nil, nil, "firstRound_25"}
	firstRound[26] = []interface{}{nil, nil, "firstRound_26"}
	firstRound[27] = []interface{}{nil, nil, "firstRound_27"}
	firstRound[28] = []interface{}{nil, nil, "firstRound_28"}
	firstRound[29] = []interface{}{nil, nil, "firstRound_29"}
	firstRound[30] = []interface{}{nil, nil, "firstRound_30"}
	firstRound[31] = []interface{}{nil, nil, "firstRound_31"}

	// secondRound
	secondRound := make([][]interface{}, 16)
	secondRound[0] = []interface{}{nil, nil, "secondRound_0"}
	secondRound[1] = []interface{}{nil, nil, "secondRound_1"}
	secondRound[2] = []interface{}{nil, nil, "secondRound_2"}
	secondRound[3] = []interface{}{nil, nil, "secondRound_3"}
	secondRound[4] = []interface{}{nil, nil, "secondRound_4"}
	secondRound[5] = []interface{}{nil, nil, "secondRound_5"}
	secondRound[6] = []interface{}{nil, nil, "secondRound_6"}
	secondRound[7] = []interface{}{nil, nil, "secondRound_7"}

	secondRound[8] = []interface{}{nil, nil, "secondRound_8"}
	secondRound[9] = []interface{}{nil, nil, "secondRound_9"}
	secondRound[10] = []interface{}{nil, nil, "secondRound_10"}
	secondRound[11] = []interface{}{nil, nil, "secondRound_11"}
	secondRound[12] = []interface{}{nil, nil, "secondRound_12"}
	secondRound[13] = []interface{}{nil, nil, "secondRound_13"}
	secondRound[14] = []interface{}{nil, nil, "secondRound_14"}
	secondRound[15] = []interface{}{nil, nil, "secondRound_15"}

	// sweetSixteen
	sweetSixteen := make([][]interface{}, 8)
	sweetSixteen[0] = []interface{}{nil, nil, "sweetSixteen_0"}
	sweetSixteen[1] = []interface{}{nil, nil, "sweetSixteen_1"}
	sweetSixteen[2] = []interface{}{nil, nil, "sweetSixteen_2"}
	sweetSixteen[3] = []interface{}{nil, nil, "sweetSixteen_3"}
	sweetSixteen[4] = []interface{}{nil, nil, "sweetSixteen_4"}
	sweetSixteen[5] = []interface{}{nil, nil, "sweetSixteen_5"}
	sweetSixteen[6] = []interface{}{nil, nil, "sweetSixteen_6"}
	sweetSixteen[7] = []interface{}{nil, nil, "secondRound_7"}

	// eliteEight
	eliteEight := make([][]interface{}, 4)
	eliteEight[0] = []interface{}{nil, nil, "eliteEight_0"}
	eliteEight[1] = []interface{}{nil, nil, "eliteEight_1"}
	eliteEight[2] = []interface{}{nil, nil, "eliteEight_2"}
	eliteEight[3] = []interface{}{nil, nil, "eliteEight_3"}

	// finalFour
	finalFour := make([][]interface{}, 2)
	finalFour[0] = []interface{}{nil, nil, "finalFour_0"}
	finalFour[1] = []interface{}{nil, nil, "finalFour_1"}

	// championship
	championship := make([][]interface{}, 1)
	championship[0] = []interface{}{nil, nil, "championship_0"}

	results := SixtyFourResults{}
	results.FirstRound = firstRound
	results.SecondRound = secondRound
	results.SweetSixteen = sweetSixteen
	results.EliteEight = eliteEight
	results.FinalFour = finalFour
	results.Championship = championship

	return &results, nil
}

func (b *Bracket) UpdateResults() {

}

func GetHDMBracket(db *sql.DB, x int) (*Bracket, error) {
	var id int
	var strTeams string
	var strResults string
	var createdAt time.Time
	var updatedAt time.Time

	err := db.QueryRow("SELECT * FROM brackets WHERE Id=?", x).Scan(&id, &strTeams, &strResults, &createdAt, &updatedAt)
	if err != nil {
		log.Println("Failed to scan bracket from Brackets table:", err)
		return nil, err
	}

	teams := []TeamPair{}
	json.Unmarshal([]byte(strTeams), &teams)

	results := SixtyFourResults{}
	json.Unmarshal([]byte(strResults), &results)

	bracket := Bracket{
		Id:        id,
		Teams:     teams,
		Results:   results,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	return &bracket, nil
}

// SampleBracket makes a Bracket with 8 contenders
func SampleBracket() *Bracket {
	// Teams
	teams := make([]TeamPair, 4)
	teams[0] = TeamPair{"joe", "matt"}
	teams[1] = TeamPair{"jim", "mike"}
	teams[2] = TeamPair{"tim", "amy"}
	teams[3] = TeamPair{"tim", "amy"}

	// Each round of results
	firstRound := make([][]interface{}, 4)
	firstRound[0] = []interface{}{1, 0, "g1"}
	firstRound[1] = []interface{}{nil, nil, "g2"}
	firstRound[2] = []interface{}{nil, nil, "g3"}
	firstRound[3] = []interface{}{nil, nil, "g4"}

	secondRound := make([][]interface{}, 2)
	secondRound[0] = []interface{}{nil, nil, "g5"}
	secondRound[1] = []interface{}{nil, nil, "g6"}

	sweetSixteen := make([][]interface{}, 2)
	sweetSixteen[0] = []interface{}{nil, nil, "g7"}
	sweetSixteen[1] = []interface{}{nil, nil, "g8"}

	// Total results
	results := SixtyFourResults{}
	results.FirstRound = firstRound
	results.SecondRound = secondRound
	results.SweetSixteen = sweetSixteen

	bracket := Bracket{666, teams, results, time.Now(), time.Now()}
	return &bracket
}

// FullBracketExample contains an example of raw interfaces formatted for jQueryBracket
func FullBracketExample() Bracket {
	teams := make([][]interface{}, 32)
	teams[0] = []interface{}{"hank", nil}
	teams[1] = []interface{}{"joe", "matt"}
	teams[2] = []interface{}{"tj", "cody"}
	teams[3] = []interface{}{"nil", nil}
	teams[4] = []interface{}{"nil", nil}
	teams[5] = []interface{}{"george", "jim"}
	teams[6] = []interface{}{"ted", "tim"}
	teams[7] = []interface{}{"nil", nil}

	teams[8] = []interface{}{"nil", nil}
	teams[9] = []interface{}{"ted", "tim"}
	teams[10] = []interface{}{"ted", "tim"}
	teams[11] = []interface{}{"nil", nil}
	teams[12] = []interface{}{"nil", nil}
	teams[13] = []interface{}{"ted", "tim"}
	teams[14] = []interface{}{"ted", "tim"}
	teams[15] = []interface{}{"nil", nil}

	teams[16] = []interface{}{"nil", nil}
	teams[17] = []interface{}{"ted", "tim"}
	teams[18] = []interface{}{"ted", "tim"}
	teams[19] = []interface{}{"nil", nil}
	teams[20] = []interface{}{"nil", nil}
	teams[21] = []interface{}{"ted", "tim"}
	teams[22] = []interface{}{"ted", "tim"}
	teams[23] = []interface{}{"nil", nil}

	teams[24] = []interface{}{"nil", nil}
	teams[25] = []interface{}{"ted", "tim"}
	teams[26] = []interface{}{"ted", "tim"}
	teams[27] = []interface{}{"nil", nil}
	teams[28] = []interface{}{"nil", nil}
	teams[29] = []interface{}{"ted", "tim"}
	teams[30] = []interface{}{"ted", "tim"}
	teams[31] = []interface{}{"nil", nil}

	firstRound := make([][]interface{}, 32)
	firstRound[0] = []interface{}{1, 0, "firstRound-m0"}
	firstRound[1] = []interface{}{4, 9, "firstRound-m1"}
	firstRound[2] = []interface{}{nil, nil, "firstRound-m2"}
	firstRound[3] = []interface{}{nil, nil, "firstRound-m3"}
	firstRound[4] = []interface{}{nil, nil, "firstRound-m4"}
	firstRound[5] = []interface{}{nil, nil, "firstRound-m5"}
	firstRound[6] = []interface{}{nil, nil, "firstRound-m6"}
	firstRound[7] = []interface{}{nil, nil, "firstRound-m7"}
	firstRound[8] = []interface{}{nil, nil, "firstRound-m8"}
	firstRound[9] = []interface{}{nil, nil, "firstRound-m9"}
	firstRound[10] = []interface{}{nil, nil, "firstRound-m10"}
	firstRound[11] = []interface{}{nil, nil, "firstRound-m11"}
	firstRound[12] = []interface{}{nil, nil, "firstRound-m12"}
	firstRound[13] = []interface{}{nil, nil, "firstRound-m13"}
	firstRound[14] = []interface{}{nil, nil, "firstRound-m14"}
	firstRound[15] = []interface{}{nil, nil, "firstRound-m15"}
	firstRound[16] = []interface{}{nil, nil, "firstRound-m16"}
	firstRound[17] = []interface{}{nil, nil, "firstRound-m17"}
	firstRound[18] = []interface{}{nil, nil, "firstRound-m18"}
	firstRound[19] = []interface{}{nil, nil, "firstRound-m19"}
	firstRound[20] = []interface{}{nil, nil, "firstRound-m20"}
	firstRound[21] = []interface{}{nil, nil, "firstRound-m21"}
	firstRound[22] = []interface{}{nil, nil, "firstRound-m22"}
	firstRound[23] = []interface{}{nil, nil, "firstRound-m23"}
	firstRound[24] = []interface{}{nil, nil, "firstRound-m24"}
	firstRound[25] = []interface{}{nil, nil, "firstRound-m25"}
	firstRound[26] = []interface{}{nil, nil, "firstRound-m26"}
	firstRound[27] = []interface{}{nil, nil, "firstRound-m27"}
	firstRound[28] = []interface{}{nil, nil, "firstRound-m28"}
	firstRound[29] = []interface{}{nil, nil, "firstRound-m29"}
	firstRound[30] = []interface{}{nil, nil, "firstRound-m30"}
	firstRound[31] = []interface{}{nil, nil, "firstRound-m31"}

	secondRound := make([][]interface{}, 16)
	secondRound[0] = []interface{}{nil, nil, "secondRound-m0"}
	secondRound[1] = []interface{}{nil, nil, "secondRound-m1"}
	secondRound[2] = []interface{}{nil, nil, "secondRound-m2"}
	secondRound[3] = []interface{}{nil, nil, "secondRound-m3"}
	secondRound[4] = []interface{}{nil, nil, "secondRound-m4"}
	secondRound[5] = []interface{}{nil, nil, "secondRound-m5"}
	secondRound[6] = []interface{}{nil, nil, "secondRound-m6"}
	secondRound[7] = []interface{}{nil, nil, "secondRound-m7"}
	secondRound[8] = []interface{}{nil, nil, "secondRound-m8"}
	secondRound[9] = []interface{}{nil, nil, "secondRound-m9"}
	secondRound[10] = []interface{}{nil, nil, "secondRound-m10"}
	secondRound[11] = []interface{}{nil, nil, "secondRound-m11"}
	secondRound[12] = []interface{}{nil, nil, "secondRound-m12"}
	secondRound[13] = []interface{}{nil, nil, "secondRound-m13"}
	secondRound[14] = []interface{}{nil, nil, "secondRound-m14"}
	secondRound[15] = []interface{}{nil, nil, "secondRound-m15"}

	sweetSixteen := make([][]interface{}, 8)
	sweetSixteen[0] = []interface{}{nil, nil, "sweetSixteen-m0"}
	sweetSixteen[1] = []interface{}{nil, nil, "sweetSixteen-m1"}
	sweetSixteen[2] = []interface{}{nil, nil, "sweetSixteen-m2"}
	sweetSixteen[3] = []interface{}{nil, nil, "sweetSixteen-m3"}
	sweetSixteen[4] = []interface{}{nil, nil, "sweetSixteen-m4"}
	sweetSixteen[5] = []interface{}{nil, nil, "sweetSixteen-m5"}
	sweetSixteen[6] = []interface{}{nil, nil, "sweetSixteen-m6"}
	sweetSixteen[7] = []interface{}{nil, nil, "sweetSixteen-m7"}

	eliteEight := make([][]interface{}, 4)
	eliteEight[0] = []interface{}{nil, nil, "eliteEight-m0"}
	eliteEight[1] = []interface{}{nil, nil, "eliteEight-m1"}
	eliteEight[2] = []interface{}{nil, nil, "eliteEight-m2"}
	eliteEight[3] = []interface{}{nil, nil, "eliteEight-m3"}

	finalFour := make([][]interface{}, 2)
	finalFour[0] = []interface{}{nil, nil, "finalFour-m0"}
	finalFour[1] = []interface{}{nil, nil, "finalFour-m1"}

	championship := make([][]interface{}, 1)
	championship[0] = []interface{}{nil, nil, "championship"}

	results := make([]interface{}, 6)
	results[0] = firstRound
	results[1] = secondRound
	results[2] = sweetSixteen
	results[3] = eliteEight
	results[4] = finalFour
	results[5] = championship

	// return Bracket{666, teams, results, time.Now(), time.Now()}
	return Bracket{}
}

// example on how to serve a bracket from hdm db
func bracketDataHandler(w http.ResponseWriter, r *http.Request) {
	bracket, _ := GenerateInitialBracket()

	bracketJS := bracket.serialize()

	// bundle up JSBracket for transport!
	js, err := json.Marshal(bracketJS)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
