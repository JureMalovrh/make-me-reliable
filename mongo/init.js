db = db.getSiblingDB(_getEnv("DATABASE_NAME"));

db.createUser(
    {
        user: _getEnv("DATABASE_USERNAME"),
        pwd: _getEnv("DATABASE_PASSWORD"),

        roles: [
            {
                role: "readWrite",
                db: "reliable-api"
            }
        ]
    }
);
