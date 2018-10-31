const pool = require('./db');
const {send} = require('./util');

const deleteServer = async (server) => {
  console.log(server)
  //Layers table has a ON DELETE CASCADE so we don't need to worry about the layers table
  await pool.query("DELETE FROM servers WHERE md5=MD5($1)::uuid", [server]);
  await pool.query("INSERT INTO banned_servers(md5, url) VALUES (MD5($1)::uuid, $1) ON CONFLICT DO NOTHING", [server]);
  console.log(`Successfully added '${server}' to banned_servers list`)
}


module.exports = async (req, res) => {
  switch(req.method) {
    case "DELETE":
      try {
        deleteServer(req.body['server']);
        send(res, 0, "Successfully deleted and blacklisted server");
      } catch(err) {
        send(res, 1, err);
      }
      break;
  }
}