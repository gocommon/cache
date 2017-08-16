package cache

// import (
// 	"fmt"
// 	"reflect"
// 	"testing"
// )

// type TestUser struct {
// 	ID   int64
// 	Name string
// }

// func getTestUserAllTag() string {
// 	return "getTestUserAllTag"
// }

// func getTestUserInfoTag(id int64) string {
// 	return fmt.Sprintf("getTestUserInfoTag:%d", id)
// }

// func getTestUserInfoKey(id int64) string {
// 	return fmt.Sprintf("testuserinfo:%d", id)
// }

// func TestCache(t *testing.T) {
// 	var err error

// 	testUser := &TestUser{1, "weisd"}

// 	c := NewCache()

// 	tags := []string{
// 		getTestUserAllTag(),
// 		getTestUserInfoTag(testUser.ID),
// 	}

// 	key := getTestUserInfoKey(testUser.ID)

// 	err = c.Tags(tags).Set(key, testUser)
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	t.Log("tagid", c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// ret := &testUser
// 	var ret *TestUser
// 	has, err := c.Tags(tags).Get(key, &ret)
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	if !has || !reflect.DeepEqual(ret, testUser) {
// 		t.Log(has, ret, testUser)
// 		t.Fail()
// 		return
// 	}

// 	t.Log("set.get pass", ret, c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// flush
// 	c.Flush([]string{getTestUserInfoTag(testUser.ID)})

// 	var ret1 *TestUser
// 	t.Log("flush get before ", ret1)
// 	has, err = c.Tags(tags).Get(key, &ret1)
// 	if err != nil {
// 		t.Fatal(err)

// 	}

// 	t.Log(has, ret1)

// 	if has || ret1 != nil {
// 		t.Log("get flush fail", has, ret1, ret1 != nil)
// 		t.Fail()
// 		return
// 	}

// 	t.Log("get flush pass", ret1, ret1 == nil, c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// test set nil
// 	retnil := &TestUser{}
// 	retnil = nil
// 	// set nil
// 	err = c.Tags(tags).Set(key, retnil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("tagid", c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// ret2 := &TestUser{}

// 	var ret2 *TestUser

// 	has, err = c.Tags(tags).Get(key, &ret2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if err != nil || ret2 != nil {
// 		t.Log(err, ret2, has)
// 		t.Fail()
// 		return
// 	}

// 	t.Log("tags.Get nil pass", ret2, c.TagID(getTestUserInfoTag(testUser.ID)))

// }
