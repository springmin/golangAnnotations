package event

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/MarcGrol/golangAnnotations/generator/event/eventAnnotation"
	"github.com/MarcGrol/golangAnnotations/model"
	"github.com/stretchr/testify/assert"
)

func cleanup() {
	os.Remove("./testData/$aggregates.go")
	os.Remove("./testData/$wrappers.go")
	os.Remove("./store/$testDataEventStore.go")
}

func TestGenerateForEvents(t *testing.T) {
	cleanup()
	defer cleanup()

	s := []model.Struct{
		{
			PackageName: "testData",
			DocLines:    []string{`//@Event(aggregate = "Test")`},
			Name:        "MyStruct",
			Fields: []model.Field{
				{Name: "StringField", TypeName: "string", IsPointer: false, IsSlice: false},
				{Name: "IntField", TypeName: "int", IsPointer: false, IsSlice: false},
				{Name: "StructField", TypeName: "MyStruct", IsPointer: true, IsSlice: false},
				{Name: "SliceField", TypeName: "MyStruct", IsPointer: false, IsSlice: true},
			},
		},
	}
	err := Generate("testData", model.ParsedSources{Structs: s})
	assert.Nil(t, err)

	// check that generated files exisst
	_, err = os.Stat("./testData/$aggregates.go")
	assert.NoError(t, err)

	data, err := ioutil.ReadFile("./testData/$aggregates.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "type TestAggregate interface {")
	assert.Contains(t, string(data), "ApplyMyStruct(c context.Context, event MyStruct)")
	assert.Contains(t, string(data), "func ApplyTestEvent(c context.Context, envelope events.Envelope, aggregateRoot TestAggregate) error {")
	assert.Contains(t, string(data), "func ApplyTestEvents(c context.Context, envelopes []events.Envelope, aggregateRoot TestAggregate) error {")
	assert.Contains(t, string(data), "func UnWrapTestEvent(envelope *events.Envelope) (events.Event, error) {")
	assert.Contains(t, string(data), "func UnWrapTestEvents(envelopes []events.Envelope) ([]events.Event, error) {")

	// check that generate code has 4 helper functions for MyStruct
	data, err = ioutil.ReadFile("./testData/$wrappers.go")
	assert.NoError(t, err)
	assert.Contains(t, string(data), "func (s *MyStruct) Wrap(sessionUID string) (*events.Envelope,error) {")
	assert.Contains(t, string(data), "func IsMyStruct(envelope *events.Envelope) bool {")
	assert.Contains(t, string(data), "func GetIfIsMyStruct(envelope *events.Envelope) (*MyStruct, bool) {")
	assert.Contains(t, string(data), "func UnWrapMyStruct(envelope *events.Envelope) (*MyStruct,error) {")

	_, err = os.Stat("./testData/$wrappers.go")
	assert.NoError(t, err)

	os.Remove("./testData/$aggregates.go")
	os.Remove("./testData/$wrappers.go")
	os.Remove("./testData/$wrappers_test.go")
	os.Remove("./repo/$storeEvents.go")

}

func TestIsEvent(t *testing.T) {
	eventAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.True(t, IsEvent(s))
}

func TestGetAggregateName(t *testing.T) {
	eventAnnotation.Register()
	s := model.Struct{
		DocLines: []string{
			`//@Event( aggregate = "person")`},
	}
	assert.Equal(t, "person", GetAggregateName(s))
}
