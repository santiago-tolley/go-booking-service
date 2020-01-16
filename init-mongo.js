use clients-service
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