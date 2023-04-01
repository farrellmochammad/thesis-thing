import redis
class TestPub():
    def __init__(self, host='0.0.0.0', port=6379, db=0, password='eYVX7EwVmmxKPCDmwMtyKVge8oLd2t81'):
        self.queue = redis.StrictRedis(host=host, port=port, db=db)

    def pub(self, name, value):
        self.queue.publish(name, value)
        
channel = "Musik-Klasik"
message ="Canon in D major"
TestPub().pub(channel, message)