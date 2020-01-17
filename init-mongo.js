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
db = db.getSiblingDB('test-clients-service')
db.createUser(
    {
        user: "clients-service",
        pwd: "clients-service",
        roles: [
            {
                role: "readWrite",
                db: "test-clients-service"
            }
        ]
    }
)
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