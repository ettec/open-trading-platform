import { ClientReadableStream, Error, Status } from 'grpc-web';

/**
 * Wraps a grpc stream with reconnect behaviour
 */
export default class Stream<T> {
    connect: () => ClientReadableStream<T>
    listener: (t: T) => void
    nextResubscribeInterval: number
    readonly initialResubscribeInterval = 1000
    readonly maxResubscribeInterval = 30000
    streamName: string
  
    constructor(connect: () => ClientReadableStream<T>, listener: (t: T) => void, streamName: string) {
  
      this.subscribe = this.subscribe.bind(this);
  
      this.connect = connect
      this.listener = listener
      this.nextResubscribeInterval = this.initialResubscribeInterval
      this.streamName = streamName
      this.subscribe()
    }
  
    private subscribe() {
      console.log("subscribing to stream:" + this.streamName)
      let stream = this.connect()
  
      stream.on('data', (data: T) => {
        this.nextResubscribeInterval = this.initialResubscribeInterval
        this.listener(data)
      });
      stream.on('status', (status: Status) => {
        console.log(this.streamName + " status:" + status.details)
      });
      stream.on('error', (err: Error) => {
        stream.cancel()
        let nextInterval = this.getNextResubscribeInterval()
        console.log("error occurred on  " + this.streamName + ", resubscribing in " + this.nextResubscribeInterval + "ms.  Error:", err);
        setTimeout(this.subscribe, nextInterval)
      });
      stream.on('end', () => {
        stream.cancel()
        let nextInterval = this.getNextResubscribeInterval()
        console.log(this.streamName + " end signal received, resubscribing in " + this.nextResubscribeInterval + "ms");
        setTimeout(this.subscribe, nextInterval)
      });
    }
  
    private getNextResubscribeInterval(): number {
      this.nextResubscribeInterval = this.nextResubscribeInterval * 2
      if (this.nextResubscribeInterval > this.maxResubscribeInterval) {
        this.nextResubscribeInterval = this.maxResubscribeInterval
      }
      return this.nextResubscribeInterval
    }
  }