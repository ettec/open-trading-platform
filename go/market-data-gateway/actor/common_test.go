package actor

import "testing"

type testActor struct {
	actorImpl
}

func (s *testActor) readInputChannels() (chan<- bool, error) {
	select {
	case d := <-s.closeChan:
		return d, nil
	}

}


func Test_actorImpl_Close(t *testing.T) {

	ta := testActor{}
	ta.actorImpl = newActorImpl("test", ta.readInputChannels)
	ta.Start()

	done:= make(chan bool)
	ta.Close(done)
	<-done

}