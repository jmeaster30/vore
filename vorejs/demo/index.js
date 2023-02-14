let libvore = new window.LibVoreJS.LibVoreJS();

console.log(libvore)

const doVoreSearch = () => {
  const sourceCode = document.getElementById("sourceCode")
  const searchText = document.getElementById("searchText")
  const results = libvore.search(sourceCode, searchText)
  document.getElementById("results").innerHTML = JSON.stringify(results, null, 4)
}