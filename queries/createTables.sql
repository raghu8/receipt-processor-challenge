`
    CREATE TABLE items (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    shortDescription VARCHAR(255),
    price DOUBLE
);

CREATE TABLE transaction (
    ID INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid VARCHAR(255),
    retailer VARCHAR(255),
    purchaseDate TIMESTAMP,
    purchaseTime TIMESTAMP,
    items INTEGER,
    FOREIGN KEY (items) REFERENCES items(id)
);`