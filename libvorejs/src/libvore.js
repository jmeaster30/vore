import wasm from 'libvorejs'

export class LibVoreJS {
  search(source, text) {
    return wasm.voreSearch(source, text);
  }
}
