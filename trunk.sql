CREATE EXTENSION IF NOT EXISTS postgis;
DROP TABLE IF EXISTS layers;
DROP TABLE IF EXISTS servers;
CREATE TABLE servers(
    url VARCHAR(1024),
    md5 uuid PRIMARY KEY NOT NULL,
    last_updated TIMESTAMP DEFAULT current_timestamp
);

CREATE TABLE layers(
    md5 UUID PRIMARY KEY NOT NULL,
    parent_id UUID,
    url VARCHAR(1024) NOT NULL,
    geometry VARCHAR(32),
    extent geometry(POLYGON, 4326), --4326 is WGS84 spatial reference. We'll need to make sure we always transform before inputting into table http://spatialreference.org/ref/epsg/wgs-84/
    name VARCHAR(512),
    description VARCHAR(10000),
    desc_vector tsvector,
    FOREIGN KEY (parent_id) REFERENCES servers(md5)
);

CREATE INDEX layer_extent_index ON layers USING GIST(extent);