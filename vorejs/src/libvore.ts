import wasm from './main.go'

export class LibVoreJS {
  search(source: string, text: string) {
    return wasm.then((wasm: any) => {
      console.log(wasm);
      const search = wasm.instance.exports.voreSearch as CallableFunction;
      return search(source, text);
    });
  }
}



// export const voreSearch = async (source: string, text: string) => {
//   const module = await WebAssembly.instantiateStreaming(fetch('../build/libvore.wasm'), {});
//   console.log(module);
//   const search = module.instance.exports.voreSearch as CallableFunction;
//   return search(source, text);
// };