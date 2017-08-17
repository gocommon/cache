package cache

// import (
// 	"fmt"
// 	"testing"
// 	"time"
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

// 	// testUser := &TestUser{1, "weisd"}

// 	c := NewCache()

// 	lockey := "ssss"

// 	l1 := c.Locker(lockey)

// 	err = l1.Lock()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	go func() {
// 		l2 := c.Locker(lockey)
// 		for {
// 			err := l2.Lock()
// 			t.Log("l2", err)
// 			if err == nil {
// 				t.Log("l2 lock ok")
// 				l2.Unlock()
// 				return
// 			}

// 			time.Sleep(300 * time.Millisecond)
// 			t.Log("l2 lock again")
// 		}
// 	}()

// 	time.Sleep(1 * time.Second)

// 	l1.Unlock()
// 	t.Log("unlock l1")

// 	time.Sleep(2 * time.Second)

// 	// 	// start := time.Now()

// 	// 	// for i := 0; i < 1000; i++ {

// 	// 	err = c.Tags(tags...).Set(key, testUser)
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 		return
// 	// 	}

// 	// 	// t.Log("tagid", c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// 	// ret := &testUser
// 	// 	var ret *TestUser
// 	// 	has, err := c.Tags(tags...).Get(key, &ret)
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 		return
// 	// 	}

// 	// 	// t.Log("set.get pass", ret, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// 	c.Tags(getTestUserInfoTag(testUser.ID)).Flush()

// 	// 	var ret1 *TestUser
// 	// 	has, err = c.Tags(tags...).Get(key, &ret1)
// 	// 	if err != nil {
// 	// 		t.Fatal(err)

// 	// key := getTestUserInfoKey(testUser.ID)

// 	// err = c.Tags(tags).Set(key, testUser)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// 	return
// 	// }

// 	// t.Log("tagid", c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// 	// t.Log("get flush pass", ret1, ret1 == nil, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// 	// test set nil
// 	// 	retnil := &TestUser{}
// 	// 	retnil = nil
// 	// 	// set nil
// 	// 	err = c.Tags(tags...).Set(key, retnil)
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 	}

// 	// 	// t.Log("tagid", c.Tags().TagID(getTestUserInfoTag(testUser.ID)))

// 	// // flush
// 	// c.Flush([]string{getTestUserInfoTag(testUser.ID)})

// 	// var ret1 *TestUser
// 	// t.Log("flush get before ", ret1)
// 	// has, err = c.Tags(tags).Get(key, &ret1)
// 	// if err != nil {
// 	// 	t.Fatal(err)

// 	// 	has, err = c.Tags(tags...).Get(key, &ret2)
// 	// 	if err != nil {
// 	// 		t.Fatal(err)
// 	// 	}

// 	// t.Log(has, ret1)

// 	// 	t.Log("tags.Get nil pass", ret2, c.Tags().TagID(getTestUserInfoTag(testUser.ID)))
// 	// 	// }

// 	// t.Log("get flush pass", ret1, ret1 == nil, c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// // test set nil
// 	// retnil := &TestUser{}
// 	// retnil = nil
// 	// // set nil
// 	// err = c.Tags(tags).Set(key, retnil)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// t.Log("tagid", c.TagID(getTestUserInfoTag(testUser.ID)))

// 	// // ret2 := &TestUser{}

// 	// var ret2 *TestUser

// 	// has, err = c.Tags(tags).Get(key, &ret2)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }

// 	// if err != nil || ret2 != nil {
// 	// 	t.Log(err, ret2, has)
// 	// 	t.Fail()
// 	// 	return
// 	// }

// 	// t.Log("tags.Get nil pass", ret2, c.TagID(getTestUserInfoTag(testUser.ID)))

// }
