const pg = require('pg');
const opts = {
    max: 3,
    host: process.env["PGHOST"] || 'localhost',
    user: process.env["PGUSER"],
    password: process.env["PGPASSWORD"],
    database: process.env["PGDATABASE"] || 'postgres'
};
const pool = new pg.Pool(opts);
console.log(JSON.stringify(opts));

module.exports = pool;
