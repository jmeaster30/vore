import wasm from 'libvorejs'

export class LibVoreJS {
  async search(source, text) {
    console.log(source);
    console.log(text);
    return wasm.voreSearch(source, text);
  }
}
