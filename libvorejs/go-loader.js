const {readFileSync, unlinkSync} = require("fs");
const {join} = require("path");
const {execFile} = require("child_process");
const which = require('which');

const proxyBuilder = (wasmFile) => `
require('../lib/wasm_exec.js');
const g = global || window || self

if (!g.__libvore__) {
  g.__libvore__ = {};
}

const bridge = g.__libvore__;

function sleep() {
  return new Promise(requestAnimationFrame);
}

function _base64ToArrayBuffer(base64) {
  var binary_string = window.atob(base64);
  var len = binary_string.length;
  var bytes = new Uint8Array(len);
  for (var i = 0; i < len; i++) {
      bytes[i] = binary_string.charCodeAt(i);
  }
  return bytes.buffer;
}

function loadwasm(bytes) {
  let ready = false;

  async function init() {
    const go = new g.Go();
    let result = await WebAssembly.instantiate(_base64ToArrayBuffer(bytes), go.importObject);
    go.run(result.instance);
    ready = true;
  }

  init();

  let proxy = new Proxy(
    {},
    {
      get: (a, key) => {
        return (...args) => {
          return new Promise(async (resolve, reject) => {
            while (!ready) {
              await sleep();
            }
  
            if (!(key in bridge)) {
              reject(\`There is nothing defined with the name "\${key.toString()}"\`);
              return;
            }
  
            if (typeof bridge[key] !== 'function') {
              resolve(bridge[key]);
              return;
            }

            bridge[key].apply(undefined, [...args, resolve, reject]);
          });
        };
      }
    }
  );

  return proxy;
}

export default loadwasm('${wasmFile}');
`;

module.exports = async function loader(contents) {
  const cb = this.async();

  const opts = {
    env: {
      GOPATH: process.env.GOPATH,
      GOROOT: process.env.GOROOT,
      GOCACHE: join(__dirname, "./.gocache"),
      GOOS: "js",
      GOARCH: "wasm"
    }
  };

  const outFile = `${this.resourcePath}.wasm`;
  const args = ["build", "-o", outFile, this.resourcePath];

  execFile(which.sync('go'), args, opts, (err, stdout, stderr) => {
    console.log(stdout)
    console.log(stderr)
    console.log("compiled")
    if (err) {
      cb(err);
      return;
    }

    let out = readFileSync(outFile).toString('base64');
    unlinkSync(outFile);

    cb(
      null,
      proxyBuilder(out)
    );
  });
}
