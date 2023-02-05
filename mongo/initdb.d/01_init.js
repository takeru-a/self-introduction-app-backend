var user = {
    user: "superuser",
    pwd: "password123",
    roles: [
      {
        role: "dbAdmin",
        db: "sample"
      }
    ]
  };
  
db.createUser(user);
db.createCollection("sample_collection");