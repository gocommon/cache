package cache

import (
	"fmt"
	"reflect"
	"testing"
)

type TestUser struct {
	ID   int64
	Name string
}

func getTestUserAllTag() string {
	return "getTestUserAllTag"
}

func getTestUserInfoTag(id int64) string {
	return fmt.Sprintf("getTestUserInfoTag:%d", id)
}

func getTestUserInfoKey(id int64) string {
	return fmt.Sprintf("testuserinfo:%d", id)
}

func TestCache(t *testing.T) {
	var err error

	testUser := &TestUser{1, "weisd"}

	c := NewCache()

	tags := []string{
		getTestUserAllTag(),
		getTestUserInfoTag(testUser.ID),
	}

	key := getTestUserInfoKey(testUser.ID)

	err = c.Tags(tags).Set(key, testUser)
	if err != nil {
		t.Fatal(err)
	}

	var ret *TestUser
	err = c.Tags(tags).Get(key, &ret)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil || !reflect.DeepEqual(ret, testUser) {
		t.Fail()
	}

	t.Log("set.get pass", ret)

	// flush
	c.Flush([]string{getTestUserInfoTag(testUser.ID)})

	var ret1 *TestUser
	err = c.Tags(tags).Get(key, &ret1)
	if err != nil {
		if err != ErrNil {
			t.Fatal(err)
		}

	}

	if err != ErrNil || ret1 != nil {
		t.Fail()
	}
	t.Log("get flush pass", ret1, ret1 == nil)

	// test set nil
	var retnil *TestUser
	retnil = nil
	// set nil
	err = c.Tags(tags).Set(key, retnil)
	if err != nil {
		t.Fatal(err)
	}

	var ret2 *TestUser

	err = c.Tags(tags).Get(key, &ret2)
	if err != nil {
		if err != ErrNil {
			t.Fatal(err)
		}
	}

	if err != nil || ret2 != nil {
		t.Fail()
	}

	t.Log("tags.Get nil pass", ret2)

}
