const {arrayToObj, objToArray} = require('./util')

test('every key in obj is included in array', () => {
  const obj = {
    "a": "b",
    "c": "d"
  }
  const arr = Object.keys(obj).reduce(objToArray(obj), []);
  expect(arr).toEqual(['b', 'd'])
});