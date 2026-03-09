const express = require("express");

const app = express();
app.use(express.json());

let users = [];

app.get("/users", (req, res) => {
    res.json(users);
});

app.post("/users", (req, res) => {
    const user = req.body;
    users.push(user);
    res.json(user);
});

app.listen(3001, () => {
    console.log("User Service running on port 3001");
});
