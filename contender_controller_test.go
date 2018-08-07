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

// testContenders returns a slice of Contender pointers and its associated slice of Resource interfaces.
//
// Create the struct you'd like to test with, then convert them to Resource interfaces.
func testContenders() ([]*Contender, []Resource) {
	contenders := []*Contender{
		{
			Name: "Matt Anderson",
			FbId: 666,
		},
		{
			Name: "George Burrows",
			FbId: 777,
		},
	}
	contenderResources := make([]Resource, len(contenders))
	for i, v := range contenders {
		contenderResources[i] = Resource(v)
	}

	return contenders, contenderResources
}

// testPosts returns a slice of Post pointers and its associated slice of Resource interfaces.
//
// Create the struct you'd like to test with, then convert them to Resource interfaces.
func testPosts() ([]*Post, []Resource) {
	posts := []*Post{
		{
			AuthorFbId: 666,
			FbId:       "666_666",
			Likes:      []int{777},
		},
	}
	postResources := make([]Resource, len(posts))
	for i, v := range posts {
		postResources[i] = Resource(v)
	}

	return posts, postResources
}

type fakeFacebookHandle struct{}

func (fh *fakeFacebookHandle) PullContendersFromFb() ([]*Contender, *ApplicationError) {
	contenders := []*Contender{
		{
			Name: "Joe Selvik",
			FbId: 1234},
		{
			Name: "TJ Gordon",
			FbId: 9876,
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
	// Create a ContenderController with the test db
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatalf("Could not open test db: %s\n", err)
	}
	cc := &ContenderController{db: db}

	originalContenders, contenderResources := testContenders()

	// Create the ContenderResources
	cids, aerr := cc.Create(contenderResources)
	if aerr != nil {
		t.Fatalf("Should not error when creating updatedContender: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	if len(cids) != len(originalContenders) {
		t.Fatal("Should only get back a single id")
	}

	// Read all originalContenders
	resources, aerr := cc.ReadCollection()
	if aerr != nil {
		t.Fatalf("Unable to read collection of originalContenders: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	var readContender *Contender
	for _, c := range resources {
		if c.(*Contender).Name == originalContenders[0].Name {
			readContender = c.(*Contender)
			break
		}
	}
	if readContender.Name != originalContenders[0].Name {
		t.Fatal("Unable to find updatedContender in ReadCollection")
	}

	// Update readContender, convert to slice of Resources
	sliceOfPostIds := []string{"111_222", "333_444"}
	readContender.PostsUsed = sliceOfPostIds
	readResource := Resource(readContender)
	aerr = cc.Update([]Resource{readResource})
	if aerr != nil {
		t.Fatalf("Error when updating resource: %s\n%s\n", aerr.Msg, aerr.Err)
	}

	// Read the updated updatedContender
	resource, aerr := cc.Read(originalContenders[0].FbId)
	if aerr != nil {
		t.Fatalf("Error when reading resource: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	updatedContender := resource.(*Contender)
	if updatedContender.Name != originalContenders[0].Name {
		t.Fatal("Read did not find updated updatedContender")
	}
	if updatedContender.PostsUsed == nil {
		t.Fatal("updatedContender did not have any used posts")
	}

	// Destroy the originalContender
	aerr = cc.Destroy([]int{originalContenders[0].FbId})
	if aerr != nil {
		t.Fatalf("Error when destroying resource: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	_, aerr = cc.Read(originalContenders[0].FbId)
	if aerr == nil {
		t.Fatalf("Should find error when reading non-existent resource: %s\n%s\n", aerr.Msg, aerr.Err)
	}
}

// TestContenderController_PopulateContendersTable tests that contenders from FB get created in a db.
func TestContenderController_PopulateContendersTable(t *testing.T) {
	// Create a ContenderController with the test db and a fake facebook handle
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatalf("Could not open test db: %s\n", err)
	}
	fh := fakeFacebookHandle{}
	cc := &ContenderController{fh: &fh, db: db}

	aerr := cc.PopulateContendersTable()
	if aerr != nil {
		t.Fatalf("Error when calling PopulateContendersTable: %s\n%s\n", aerr, aerr.Err)
	}

	contenders, aerr := cc.ReadCollection()
	if len(contenders) < 1 {
		t.Error("Should have found some contenders in DB")
	}
}

func TestContenderController_UpdateContendersVariableDependentData(t *testing.T) {
	// Create a ContenderController with the test db and a fake facebook handle
	config := NewConfig()
	db, err := models.OpenDB(config.DbTestPath)
	if err != nil {
		t.Fatalf("Could not open test db: %s\n", err)
	}
	cc := &ContenderController{db: db}
	pc := &PostController{db: db}

	originalContenders, contenderResources := testContenders()
	originalPosts, postResources := testPosts()

	// Create the ContenderResource and PostResource
	cids, aerr := cc.Create(contenderResources)
	if aerr != nil {
		t.Fatalf("Should not error when creating contender: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	if len(cids) != len(originalContenders) {
		t.Fatal("Should only get back a single id")
	}

	pids, aerr := pc.Create(postResources)
	if aerr != nil {
		t.Fatal("Should not error when creating post")
	}
	if len(pids) != 1 {
		t.Fatal("Should only get back a single id")
	}

	// Update the Contender's variable dependent data
	aerr = cc.UpdateContendersVariableDependentData(pc)
	if aerr != nil {
		t.Fatalf("Should not error when calling UpdateContendersVariableDependentData: %s\n%s\n", aerr.Msg, aerr.Err)
	}

	// Read the updated updatedContender
	resource, aerr := cc.Read(originalContenders[0].FbId)
	if aerr != nil {
		t.Fatalf("Error when reading resource: %s\n%s\n", aerr.Msg, aerr.Err)
	}
	updatedContender := resource.(*Contender)
	if updatedContender.Name != originalContenders[0].Name {
		t.Fatal("Read did not find updated updatedContender")
	}
	if len(updatedContender.Posts) == len(originalPosts) {
		t.Fatal("updatedContender did not have any posts")
	}
}

// TestContenderController_stringPostsToSlicePostIds tests converting strings into a slice of ints.
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
				t.Logf("val: %d, type: %s", val, reflect.TypeOf(val))
				t.Logf("outSlice: %d, type: %s", tt.outSlice, reflect.TypeOf(tt.outSlice))
				t.Errorf("Valid '%s' did not match slice of ints.", tt.inString)
			}
		}
	}
}

// TestContenderController_slicePostIdsToStringPosts tests converting a slice of ints into a specific string.
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
