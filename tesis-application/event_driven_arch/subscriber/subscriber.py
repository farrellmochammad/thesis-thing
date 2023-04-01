import redis
import threading

class Listener(threading.Thread):
    def __init__(self, r, channels):
        threading.Thread.__init__(self)
        self.redis = r
        self.pubsub = self.redis.pubsub()
        self.pubsub.subscribe(channels)
    
    def work(self, item):
        print(item['channel'], ":", item['data'])
    
    def run(self):
        for item in self.pubsub.listen():
            self.work(item)



if __name__ == "__main__":
    r = redis.StrictRedis(host='0.0.0.0', port=6379, db=0)
    client = Listener(r, ['Musik-Klasik'])
    client.start()