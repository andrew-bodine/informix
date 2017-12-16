package queue_test

// A push handler meant for testing the queue package.
type testPushHandler struct {
    delegate     chan interface{}
}

// Implement the queue.PushHandler interface.
func (t *testPushHandler) AfterPush(objs []interface{}) {
    t.delegate <- objs[0]
}
