const REQUIRED_KEYS = ['name', 'url', 'extent', 'geometryType', 'description']
const debug = require('debug');
const resLog = debug("insert:response");

exports.objToArray = obj => (arr, key) => {
  arr.push(obj[key]);
  return arr;
}

exports.send = (res, error, message) => {
  const response = {
    error: error,
    message: message
  }
  resLog(response);
  res.send(response);
}
/**
 * Throwing this into its own function so we can modify it easily later and potentially re-use it.
 * 
 * Chances are we're going to do filtering in `root` publish step not the consumption trunk step
 */
exports.isValidLayer = (layer) => {
  //Explicitly allow empty strings
  return REQUIRED_KEYS.every(key => layer[key] || layer[key] === "");
}