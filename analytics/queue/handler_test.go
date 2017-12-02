package queue_test

// A push handler meant for testing the queue package.
type TestPushHandler struct {
    delegate     chan interface{}
}

// Implement the queue.PushHandler interface.
func (t *TestPushHandler) AfterPush(obj interface{}) {
    t.delegate <- obj
}
