import wasm from 'libvorejs'

export class LibVoreJS {
  search(source, text) {
    console.log(source);
    console.log(text);
    console.log(wasm.instance);
    return wasm.voreSearch(source, text);
  }
}
