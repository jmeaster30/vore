const webpack = require("webpack");
const {readFileSync, unlinkSync} = require("fs");
const {basename, join} = require("path");
const {execFile} = require("child_process");

const proxyBuilder = (filename) => `
let ready = false;

const g = self || window || global

if (!g.__libvore__) {
  g.__libvore__ = {};
}

const bridge = g.__libvore__;

async function init() {
  const go = new Go();
  let result = await WebAssembly.instantiateStreaming(fetch("${filename}"), go.importObject);
  go.run(result.instance);
  ready = true;
}

function sleep() {
  return new Promise(requestAnimationFrame);
}

init();

let proxy = new Proxy(
  {},
  {
    get: (_, key) => {
      return (...args) => {
        return new Promise(async (resolve, reject) => {
          let run = () => {
            let cb = (err, ...msg) => (err ? reject(err) : resolve(...msg));
            bridge[key].apply(undefined, [...args, cb]);
          };

          while (!ready) {
            await sleep();
          }

          if (!(key in bridge)) {
            reject(\`There is nothing defined with the name "$\{key\}"\`);
            return;
          }

          if (typeof bridge[key] !== 'function') {
            resolve(bridge[key]);
            return;
          }

          run();
        });
      };
    }
  }
);
  
export default proxy;`;

module.exports = function loader(contents) {
  const cb = this.async();

  console.log(cb)

  const opts = {
    env: {
      GOPATH: process.env.GOPATH,
      GOROOT: process.env.GOROOT,
      GOCACHE: join(__dirname, "./.gocache"),
      GOOS: "js",
      GOARCH: "wasm"
    }
  };

  console.log(opts);

  //const goBin = `${opts.env.GOROOT}/bin/go`;
  const outFile = `${this.resourcePath}.wasm`;
  const args = ["build", "-o", outFile, this.resourcePath];

  execFile("go", args, opts, (err) => {
    console.log("compiled")
    if (err) {
      cb(err);
      return;
    }

    let out = readFileSync(outFile);
    unlinkSync(outFile);
    const emittedFilename = basename(this.resourcePath, ".go") + ".wasm";
    this.emitFile(emittedFilename, out, null);

    console.log(`emittedFilename: '${emittedFilename}'`)

    cb(
      null,
      [
        "require('!",
        join(__dirname, "..", "lib", "wasm_exec.js"),
        "');",
        proxyBuilder(emittedFilename)
      ].join("")
    );
  });
}
