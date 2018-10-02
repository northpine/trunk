const pool = require('./db');
const debug = require('debug')('get');
const {send} = require("./util");
const getFeatureCollectionSql = `SELECT jsonb_build_object(
  'type', 'FeatureCollection',
  'features', jsonb_agg(features.feature)
) AS collection
FROM (
  SELECT jsonb_build_object(
          'type', 'Feature',
          'id', md5,
          'geometry', ST_AsGeoJSON(extent)::jsonb,
          'properties', to_jsonb(inputs) - 'md5' - 'extent' - 'parent_id'
  ) AS feature
  FROM (
          SELECT l.*, s.url AS base_url
          FROM layers l, servers s
          WHERE parent_id = s.md5
          AND (description like $1 OR name like $1 OR l.url like $1)
          AND ST_MakeEnvelope($2, $3, $4, $5, 4326) && extent
  ) inputs
) features;`


module.exports = async (req, res) => {
  try {
    debug("query: %o", req.query);
    const {phrase, xmin, ymin, xmax, ymax} = req.query;
    const result = await pool.query(getFeatureCollectionSql, [`%${phrase}%`, xmin, ymin, xmax, ymax]);
    //We're aggregating the entire table into a single geojson collection. Therefore we only need the first row
    const featureCollection = result.rows[0].collection;
    send(res, 0, featureCollection);
  } catch (e) {
    debug(e);
    send(res, 1, e);
  }
}