import wasm from 'libvorejs';

export function search(source, text) {
  return wasm.voreSearch(source, text);
}
