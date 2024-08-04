
CREATE TABLE jobs (
    uuid UUID PRIMARY KEY,
    name TEXT NOT NULL,
    state TEXT NOT NULL,
    status TEXT NOT NULL,
    description TEXT
);

