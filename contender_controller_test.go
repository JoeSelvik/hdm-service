package main

import (
	"bytes"
	"errors"
	"github.com/JoeSelvik/hdm-service/models"
	"log"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

type fakeFacebookHandle struct{}

func (fh *fakeFacebookHandle) PullContendersFromFb() ([]*Contender, *ApplicationError) {
	contenders := []*Contender{
		{
			Name: "Joe Selvik",
			FbId: 1234},
		{
			Name: "TJ Gordon",
			FbId: 6666,
		},
	}
	return contenders, nil
}

func (fh *fakeFacebookHandle) PullPostsFromFb(startDate time.Time) ([]Post, *ApplicationError) {
	posts := []Post{
		{
			Id: "1234",
		},
	}
	return posts, nil
}

func setup() error {
	log.Println("contender_controller_test setup")
	config := NewConfig()
	log.Println(config.DbTestPath)
	log.Println(config.DbSetupScript)

	cmd := exec.Command("sqlite3", "test.db", "<", "create_tables.sql")
	//cmd := exec.Command("sqlite3 test.db", config.DbTestPath, "<", config.DbSetupScript)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	log.Printf("cmd: %+v\n", cmd)
	log.Println("env:", cmd.Env)
	log.Println("dir:", cmd.Dir)
	log.Println("path:", cmd.Path)
	log.Println("args:", cmd.Args)

	err := cmd.Run()
	if err != nil {
		log.Printf("Error executing query. Command Output: %+v\n: %+v, %v\n", out.String(), stderr.String(), err)
		log.Printf("out: %+v\n", out.String())
		log.Printf("stderr: %+v\n", stderr.String())
		log.Printf("err: %v\n", err)
		return err
	}

	//db, err := models.NewDB(config.DbTestPath)
	//if err != nil {
	//	log.Printf("Could not open test db: %s\n", err)
	//	return err
	//}
	//fh := fakeFacebookHandle{}
	//cc := &ContenderController{config: config, db: db, fh: &fh}

	return nil
}

func teardown() error {
	log.Println("contender_controller_test teardown")
	return nil
}

func TestMain(m *testing.M) {
	err := setup()
	if err != nil {
		log.Fatal("Setup failed")
	}

	retCode := m.Run()
	//if retCode == 0 {
	//	teardown()
	//}

	// Always call teardown()
	err = teardown()
	if err != nil {
		log.Fatal("Teardown failed")
	}

	os.Exit(retCode)
}

// TestContenderController_PopulateContendersTable tests that contenders from FB get created in a db.
func TestContenderController_PopulateContendersTable(t *testing.T) {
	config := NewConfig()

	db, err := models.NewDB(config.DbTestPath)
	if err != nil {
		t.Fatal("Could not open test db")
	}
	fh := fakeFacebookHandle{}
	cc := &ContenderController{fh: &fh, db: db}

	aerr := cc.PopulateContendersTable()
	if aerr != nil {
		t.Fatalf("Error when calling PopulateContendersTable: %s\n%s", aerr, aerr.Err)
	}

	contenders, aerr := cc.ReadCollection()
	if len(contenders) != 2 {
		t.Error("Should have found two contenders in DB")
	}
}

func TestContenderController_stringPostsToSlicePostIds(t *testing.T) {
	var emptySlice []int
	var tests = []struct {
		inString string
		outSlice []int
		outError error
	}{
		{"1, 2, 3", []int{1, 2, 3}, nil},
		{"1000", []int{1000}, nil},
		{"", emptySlice, nil},
		{"1 2 3", nil, errors.New("oops")},
		{"1,", nil, errors.New("oops")},
		{" ", nil, errors.New("oops")},
	}

	for _, tt := range tests {
		val, err := stringPostsToSlicePostIds(tt.inString)
		if tt.outError != nil { // If expecting error...
			if err == nil {
				t.Errorf("Invalid '%s' string ids should generate an error.", tt.inString)
			}
		} else { // If expecting results...
			if err != nil {
				t.Errorf("Valid '%s' string ids should not generate an error.", tt.inString)
			}
			if !reflect.DeepEqual(val, tt.outSlice) {
				t.Logf("val: %s, type: %s", val, reflect.TypeOf(val))
				t.Logf("outSlice: %s, type: %s", tt.outSlice, reflect.TypeOf(tt.outSlice))
				t.Errorf("Valid '%s' did not match slice of ints.", tt.inString)
			}
		}
	}
}

func TestContenderController_slicePostIdsToStringPosts(t *testing.T) {
	var tests = []struct {
		inSlice   []int
		outString string
	}{
		{[]int{1, 2, 3}, "1, 2, 3"},
		{[]int{100}, "100"},
		{[]int{}, ""},
	}

	for _, tt := range tests {
		result := slicePostIdsToStringPosts(tt.inSlice)
		if result != tt.outString {
			t.Errorf("Slice '%d' returned incorrect '%s'", tt.inSlice, tt.outString)
		}
	}
}
