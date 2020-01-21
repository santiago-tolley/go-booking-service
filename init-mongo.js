db = db.getSiblingDB('clients-service')
db.createUser(
    {
        user: "clients-service",
        pwd: "clients-service",
        roles: [
            {
                role: "readWrite",
                db: "clients-service"
            }
        ]
    }
)
var collection = db.collection('user');
collection.createIndex(
    { user : 1 }, {unique : true}, function(err, result) {
    console.log(result);
    callback(result);
});

db = db.getSiblingDB('test-clients-service')
db.createUser(
    {
        user: "test-clients-service",
        pwd: "test-clients-service",
        roles: [
            {
                role: "readWrite",
                db: "test-clients-service"
            }
        ]
    }
)
var collection = db.collection('user');
collection.createIndex(
    { user : 1 }, {unique : true}, function(err, result) {
    console.log(result);
    callback(result);
});

db = db.getSiblingDB('rooms-service')
db.createUser(
    {
        user: "rooms-service",
        pwd: "rooms-service",
        roles: [
            {
                role: "readWrite",
                db: "rooms-service"
            }
        ]
    }
)
var collection = db.collection('rooms');
collection.createIndex(
    { room_id : 1, date: 1}, {unique : true}, function(err, result) {
    console.log(result);
    callback(result);
});

db = db.getSiblingDB('test-rooms-service')
db.createUser(
    {
        user: "test-rooms-service",
        pwd: "test-rooms-service",
        roles: [
            {
                role: "readWrite",
                db: "test-rooms-service"
            }
        ]
    }
)
var collection = db.collection('rooms');
collection.createIndex(
    { room_id : 1, date: 1}, {unique : true}, function(err, result) {
    console.log(result);
    callback(result);
});