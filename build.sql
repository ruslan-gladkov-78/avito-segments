CREATE TABLE Segments2 (ID SERIAL PRIMARY KEY, Slug VARCHAR(255));
CREATE TABLE Users2 (ID INT, SegmentsIDs INT[])