package cache

// func TestCache(t *testing.T) {
// 	var err error

// 	c := NewCache()

// 	tags := []string{
// 		"userall",
// 		"user1",
// 	}
// 	start := time.Now()
// 	err = c.Tags(tags).Set("user1", "info1")
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log("set use", time.Now().Sub(start))

// 	start = time.Now()

// 	var ret string
// 	err = c.Tags(tags).Get("user1", &ret)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	t.Log(ret)
// 	t.Log("get use", time.Now().Sub(start))
// }
