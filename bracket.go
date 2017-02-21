package main

import (
	"database/sql"
)

type Bracket struct {
	Teams   [][]interface{} `json:"teams"`
	Results []interface{}   `json:"results"`
}

// CreateBracketTable creates the brackets table if it does not exist
func CreateBracketTable(db *sql.DB) {
	sql_table := `
	CREATE TABLE IF NOT EXISTS brackets(
		Id TEXT NOT NULL,
		Teams ??,
		Results ??
	);
	`

	_, err := db.Exec(sql_table)
	if err != nil {
		panic(err)
	}
}

func (b *Bracket) DBTableName() string {
	return "brackets"
}

// Path will supply the URL extension for a Bracket resource
func (b *Bracket) Path() string {
	return "/brackets/"
}

func (b *Bracket) UpdateResults() {

}

func fullBracket() Bracket {
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

	return Bracket{teams, results}
}
