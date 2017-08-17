package cache

// import (
// 	"fmt"
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

// 	c := New()

// 	key := getTestUserInfoKey(testUser.ID)

// 	tags := []string{
// 		getTestUserAllTag(),
// 		getTestUserInfoTag(testUser.ID),
// 	}

// 	// start := time.Now()

// 	// for i := 0; i < 1000; i++ {

// 	err = c.Tags(tags...).Set(key, testUser)
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	t.Log("tagid.", c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// ret := &testUser
// 	var ret *TestUser
// 	has, err := c.Tags(tags...).Get(key, &ret)
// 	if err != nil {
// 		t.Fatal(err)
// 		return
// 	}

// 	t.Log("set.get pass", ret, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	c.Tags(getTestUserInfoTag(testUser.ID)).Flush()

// 	var ret1 *TestUser
// 	has, err = c.Tags(tags...).Get(key, &ret1)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	if has {
// 		t.Log("get flush fail", has, ret1)
// 		t.Fail()
// 		return
// 	}

// 	t.Log("get flush pass", ret1, ret1 == nil, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// test set nil
// 	retnil := &TestUser{}
// 	retnil = nil
// 	// set nil
// 	err = c.Tags(tags...).Set(key, retnil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("tagid", c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// flush

// 	var ret2 *TestUser

// 	has, err = c.Tags(tags...).Get(key, &ret2)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(has, ret1)

// 	if err != nil || ret2 != nil {
// 		t.Log(err, ret2, has)
// 		t.Fail()
// 		return
// 	}

// 	t.Log("tags.Get nil pass", ret2, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// }
