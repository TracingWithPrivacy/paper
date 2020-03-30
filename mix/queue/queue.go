package queue

// ToDo: message-log for additions/removals for persistency.

import (
	"io"
	"sync"
)

// timeSlice implements a queue for one moment in time. Messages are randomized on addition.
type timeSlice struct {
	messages []interface{} // The messages in the queue.
	mutex    *sync.Mutex   // Mutex for messages access.
	random   io.Reader     // Random source.
}

// newTimeSlice creates a new timeSlice.
func newTimeSlice(size uint64, randomSource io.Reader) *timeSlice {
	return &timeSlice{
		messages: make([]interface{}, 0, uint64ToInt(size)), // We only allocate memory.
		mutex:    new(sync.Mutex),
		random:   randomSource,
	}
}

// replace returns the current messages of the timeslice and replaces the queue.
func (slice *timeSlice) replace() []interface{} {
	slice.mutex.Lock()
	defer slice.mutex.Unlock()
	r := slice.messages
	slice.messages = make([]interface{}, 0, (len(r)*11)/10) // Keep ~10% spare entries.
	return r
}

// add a message to a timeSlice.
func (slice *timeSlice) add(message interface{}) {
	var aR, bR int
	slice.mutex.Lock()
	defer slice.mutex.Unlock()
	messageCount := uint64(len(slice.messages))
	// Short queues cannot be randomized.
	if messageCount < 3 {
		slice.messages = append(slice.messages, message)
		return
	}
	// Create two unequal random positions from the message slice.
	aR = uint64ToInt(randomUint64Max(slice.random, messageCount))
sampleLoop:
	for {
		if bR = uint64ToInt(randomUint64Max(slice.random, messageCount)); bR != aR {
			break sampleLoop
		}
	}
	buf := make([]interface{}, 3)
	buf[0] = slice.messages[aR]
	buf[1] = slice.messages[bR]
	buf[2] = message

	index := randomPermutation(slice.random, 3)
	slice.messages[aR] = buf[index[0]]
	slice.messages[bR] = buf[index[1]]
	slice.messages = append(slice.messages, buf[index[2]])
}

// RandomQueue implements a time-sliced, randomized message queue.
type RandomQueue struct {
	slices        []*timeSlice  // Queues per time-slice.
	mutex         *sync.RWMutex // Mutex for slices access.
	readerMutex   *sync.Mutex   // Mutex to control readers.
	windowSize    uint64        // windowSize in seconds.
	rounds        uint64        // The number of rounds in the queue (that is, the number of sub-queues).
	granularity   uint64        // The granularity in seconds for slice distribution
	random        io.Reader     // Random source.
	readPos       uint64        // Last read timeslice.
	readerStarted bool          // True after first read from queue.
}

// New creates a new RandomQueue that covers a window of windowSize seconds with the given granularity in seconds.
// MessagesPerWindow is the total amount of messages expected per windowSize.
func New(windowSize, granularity uint64, messagesPerWindow uint64) (*RandomQueue, error) {
	random, err := randomSourceHKDF(globalRandomSource)
	if err != nil {
		return nil, err
	}
	r := &RandomQueue{
		mutex:       new(sync.RWMutex),
		readerMutex: new(sync.Mutex),
		windowSize:  windowSize,
		rounds:      windowSize / granularity,
		granularity: granularity,
		random:      random,
	}
	rRounds := uint64ToInt(r.rounds)
	r.slices = make([]*timeSlice, rRounds)
	initialTimeSliceLength := messagesPerWindow / r.rounds
	for i := 0; i < rRounds; i++ {
		randomSourceSlice, _ := randomSourceHKDF(r.random)
		r.slices[i] = newTimeSlice(initialTimeSliceLength, randomSourceSlice)
	}
	return r, nil
}

func (queue *RandomQueue) sendTimeToQueue(sendTime uint64) uint64 {
	return (sendTime / queue.granularity) % queue.rounds
}

// Add a message into the RandomQueue for sendTime. If sendTime is 0 or <currentTime,
// a random sendTime will be generated. unixTimeNow can be the current time, or 0 if we fetch that locally.
func (queue *RandomQueue) Add(sendTime uint64, message interface{}, unixTimeNow uint64) {
	if unixTimeNow == 0 {
		unixTimeNow = unixTimeSource()
	}
	if sendTime < unixTimeNow {
		sendTime = unixTimeNow + randomUint64(queue.random)
	}
	sliceNum := queue.sendTimeToQueue(sendTime)
	queue.mutex.RLock()
	slice := queue.slices[sliceNum]
	queue.mutex.RUnlock()
	slice.add(message)
}

// GetSendQueue returns the list of messages that should be sent now. If unixTimeNow is 0, the current time will be used.
func (queue *RandomQueue) GetSendQueue(unixTimeNow uint64) []interface{} {
	// ToDo: Requires real testing.
	queue.readerMutex.Lock()
	defer queue.readerMutex.Unlock()

	if unixTimeNow == 0 {
		unixTimeNow = unixTimeSource()
	}
	sliceNum := unixTimeNow / queue.granularity
	if sliceNum <= queue.readPos && queue.readerStarted {
		// Going back in time.
		return nil
	}
	queue.readerStarted = true
	queue.readPos++
	queue.mutex.Lock()
	defer queue.mutex.Unlock()
	return queue.slices[queue.readPos%queue.rounds].replace()
}
