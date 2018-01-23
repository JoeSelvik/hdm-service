// Simple UTs to test for basic functionality, could be cleaner.

package main

import (
	"errors"
	"github.com/JoeSelvik/hdm-service/models"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
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

func (fh *fakeFacebookHandle) PullPostsFromFb() ([]*Post, *ApplicationError) {
	posts := []*Post{
		{
			FbId: "111_222",
		},
	}
	return posts, nil
}

func setup() error {
	log.Println("contender_controller_test setup")
	config := NewConfig()

	// Get db setup commands
	// todo: quick and dirty, split has extra "" entry
	file, err := ioutil.ReadFile(config.DbSetupScript)
	if err != nil {
		log.Printf("Could not open test db script: %s\n", err)
		return err
	}
	cmd := strings.Split(string(file), ";")

	// Open db
	db, err := models.NewDB(config.DbTestPath)
	if err != nil {
		log.Printf("Could not open test db: %s\n", err)
		return err
	}
	log.Printf("Using new test db: %s", config.DbTestPath)

	// Execute each command in db setup script
	for _, commands := range cmd {
		if commands != "" {
			result, err := db.Exec(commands)
			if err != nil {
				log.Printf("Could not execute cmd: %s\n%s\n", err, cmd)
				log.Printf("Result: %+v\n", result)
				return err
			}
		}
	}

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

// Create, ReadCollection, Update, Read, Destroy, Read contender.
func TestContenderController_Create(t *testing.T) {
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatal("Could not open test db")
	}
	cc := &ContenderController{db: db}

	// Create a Contender struct and convert it to a Resource interface
	contenders := []*Contender{
		{
			Name: "matt anderson",
			FbId: 666},
	}
	contenderResources := make([]Resource, len(contenders))
	for i, v := range contenders {
		contenderResources[i] = Resource(v)
	}

	cids, aerr := cc.Create(contenderResources)
	if aerr != nil {
		t.Fatal("Should not error when creating contender")
	}
	if len(cids) != 1 {
		t.Fatal("Should only get back a single id")
	}

	// Read all contenders
	resources, aerr := cc.ReadCollection()
	if aerr != nil {
		t.Fatal("Unable to read collection")
	}
	var lookup *Contender
	for _, c := range resources {
		if c.(*Contender).Name == "matt anderson" {
			lookup = c.(*Contender)
			break
		}
	}
	if lookup.Name != "matt anderson" {
		t.Fatal("Unable to find contender in ReadCollection")
	}

	// Update contender, convert to slice of Resources
	lookup.PostsUsed = []int{1, 2, 3}
	lookupResource := Resource(lookup)
	aerr = cc.Update([]Resource{lookupResource})
	if aerr != nil {
		t.Fatalf("Error when updating resource: %s\n", aerr)
	}

	// Read the updated contender
	resource, aerr := cc.Read(666)
	if aerr != nil {
		t.Fatalf("Error when reading resource: %s\n", aerr)
	}
	contender := resource.(*Contender)
	if contender.Name != "matt anderson" {
		t.Fatal("Read did not find updated contender")
	}

	// Destroy the contender
	aerr = cc.Destroy([]int{666})
	if aerr != nil {
		t.Fatalf("Error when destroying resource: %s\n", aerr)
	}
	_, aerr = cc.Read(666)
	if aerr == nil {
		t.Fatalf("Should find error when reading non-existent resource: %s\n", aerr)
	}
}

// TestContenderController_PopulateContendersTable tests that contenders from FB get created in a db.
func TestContenderController_PopulateContendersTable(t *testing.T) {
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
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
	if len(contenders) < 1 {
		t.Error("Should have found some contenders in DB")
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
		val, err := stringOfIntsToSliceOfInts(tt.inString)
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
		result := sliceOfIntsToString(tt.inSlice)
		if result != tt.outString {
			t.Errorf("Slice '%d' returned incorrect '%s'", tt.inSlice, tt.outString)
		}
	}
}
