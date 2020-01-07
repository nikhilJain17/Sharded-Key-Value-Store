# leader election stuff


# leader stuff




# follower stuff
from flask import Flask, request
app = Flask(__name__)

@app.route('/')
def hello():
    return "hello world"

@app.route('/get', methods=["POST"])
def get():
    print("[get]: " + str(request.args) + ", " + str(request.get_data()) + ", " + str(request.get_json()))
    return str(request)

@app.route('/put')
def put(k, v):
    return "put"


# make basic kv store, with static leader and 1 static follower
# then make static leader and many static followers (add in replication and rollback and shit)
# then add leader election and have everything be dynamic

# class Node():
#     node_type = ["leader", "candidate", "follower"]