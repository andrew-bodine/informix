package queue

type PushHandler interface {

    // Right after an object is pushed into a Queuer, this hook will trigger.
    AfterPush(interface{})
}
