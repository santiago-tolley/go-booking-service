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
db.users.createIndex({"user": 1}, {unique: true});

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
db.users.createIndex({"user": 1}, {unique: true});

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
db.users.createIndex({"room_id": 1, "date": 1}, {unique: true});

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
db.users.createIndex({"room_id": 1, "date": 1}, {unique: true});
