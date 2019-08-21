package main

import (
	"gotest.tools/assert"
	"io/ioutil"
	"reflect"
	"testing"
)

const SourceFile = "testdata/SourceTestFile"
const DestFile = "testdata/DestTestFile"
const FullCopyFile = "testdata/fullcopy"
const Offset40CopyFile = "testdata/offset40"
const Offset1CopyFile = "testdata/offset1"

func TestFullCopy(t *testing.T) {
	bytes, err := CopyFile(DestFile, SourceFile, 0, 0)
	assert.Equal(t, bytes, int64(120))
	assert.NilError(t, err)
	read, _ := ioutil.ReadFile(DestFile)
	shouldBe, _ := ioutil.ReadFile(FullCopyFile)
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))

}

func TestOffset1Copy(t *testing.T) {
	bytes, err := CopyFile(DestFile, SourceFile, 0, 1)
	assert.Equal(t, bytes, int64(119))
	assert.NilError(t, err)
	read, _ := ioutil.ReadFile(DestFile)
	shouldBe, _ := ioutil.ReadFile(Offset1CopyFile)
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))
}

func TestOffset40Copy(t *testing.T) {
	bytes, err := CopyFile(DestFile, SourceFile, 0, 40)
	assert.Equal(t, bytes, int64(80))
	assert.NilError(t, err)
	read, _ := ioutil.ReadFile(DestFile)
	shouldBe, _ := ioutil.ReadFile(Offset40CopyFile)
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))
}

func TestLimit80Copy(t *testing.T) {
	bytes, err := CopyFile(DestFile, SourceFile, 80, 0)
	assert.Equal(t, bytes, int64(80))
	assert.NilError(t, err)
	read, _ := ioutil.ReadFile(DestFile)
	shouldBe, _ := ioutil.ReadFile(Offset40CopyFile)
	assert.Equal(t, true, reflect.DeepEqual(shouldBe, read))
}

func TestOffsetLimitCopy(t *testing.T) {
	_, err := CopyFile(DestFile, SourceFile, 80, 80)
	assert.Error(t, err, "EOF")
}
