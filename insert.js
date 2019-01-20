const pool = require('./db');
const debug = require("debug")('insert');
const {objToArray, isValidLayer, send} = require('./util');
const MAX_STRING_LENGTH = 5000;

const mapLayerToSql = (row, index) => {
    const numKeys = 11;
    let i = (index * numKeys) + 1;
    const ex = row.extent;
    const srid = ex.spatialReference.latestWkid || ex.spatialReference.wkid;
    const sql = `(ST_Transform(ST_MakeEnvelope($${i++},  $${i++}, $${i++}, $${i++}, $${i++}), 4326), MD5($${i++})::uuid, $${i++}::uuid, $${i++}, $${i++}, $${i++}, $${i++}, to_tsvector($${i++}))`
    return {
        sql: sql,
        values: [ex.xmin, ex.ymin, ex.xmax, ex.ymax, srid, row.url, row.parent_id, row.url, row.name, row.geometryType, row.description.substr(0, MAX_STRING_LENGTH), row.description]
    }

}

const buildInsert = (layers, parentId) => {
    const rowsWithData = layers.map(row => {row['parent_id'] = parentId; return row}).map(mapLayerToSql);
    const values = rowsWithData.reduce((acc, row) => acc.concat(row.values), []);
    const sql = `INSERT INTO layers(extent, md5, parent_id, url, name, geometry, description, desc_vector) VALUES ${rowsWithData.map(row => row.sql).join(",")} ON CONFLICT DO NOTHING`
    return {
        values: values,
        sql: sql
    }
}

const buildDelete = (layers, parentId) => {
    return {
        sql: `DELETE FROM layers WHERE parent_id=$1::uuid AND md5 NOT IN (${layers.map((_, i) => `MD5($${i + 2})::uuid`).join(",")}) RETURNING url`,
        values: [parentId, ...layers.map(layer => layer.url)]
    }
}

module.exports = async (req, res) => {
  try {
    const banned = await pool.query("SELECT url from banned_servers WHERE md5=MD5($1)::uuid", [req.body.server]);
    if (banned.rowCount > 0) {
      send(res, 1, "Layers will not be added because it is in the banned_servers list");
      return;
    }
    let layers = Object.keys(req.body.layers).reduce(objToArray(req.body.layers), []);
    layers = layers.filter(isValidLayer);
    if(layers.length > 0) {
      let result = await pool.query("SELECT * FROM servers WHERE md5=MD5($1)::uuid AND current_timestamp < last_updated + '1 second'::interval", [req.body.server])
      if(result.rows.length == 0) {
        debug("before insert into servers");
        result = await pool.query("INSERT INTO servers(md5, url) VALUES (MD5($1)::uuid, $1) ON CONFLICT (md5) DO UPDATE set last_updated=current_timestamp RETURNING md5", [req.body.server]);
        debug(`updated md5: ${JSON.stringify(result.rows[0].md5)}`);
        const parentId = result.rows[0].md5;
        const deleteQuery = buildDelete(layers, result.rows[0].md5);
        result = await pool.query(deleteQuery.sql, deleteQuery.values);
        debug(`successfully deleted '${result.rowCount}' old layers`);
        result = await pool.query("SELECT url FROM layers WHERE parent_id=$1::uuid", [parentId]);
        const leftovers = result.rows.reduce((prev, curr) => {prev[curr.url] = true; return prev}, {});
        layers.forEach((l, i) => {
          //Delete layers that already exist in table
          if(leftovers[l.url]) {
            layers.splice(i, 1);
          }
        });
        if(layers.length > 0) {
          const insertWithData = buildInsert(layers, parentId);
          debug("sql: %s", insertWithData.sql);
          result = await pool.query(insertWithData.sql, insertWithData.values);
          debug("inserted successfully");
        } else {
          debug("no new layers for server: " + req.body.server);
        }
        send(res, 0, "success!");
      } else {
        send(res, 1, "no valid layers sent");
      }
    } else {
      send(res, 1, "server has already been updated recently")
    }
  } catch (e) {
    console.log("fuck");
    send(res, 1, e);
  }

}