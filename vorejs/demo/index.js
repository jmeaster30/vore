let libvore = new window.LibVoreJS.LibVoreJS();

console.log(libvore)

const doVoreSearch = () => {
  const sourceCode = document.getElementById("sourceCode").value;
  const searchText = document.getElementById("searchText").value;
  const results = libvore.search(sourceCode, searchText);
  document.getElementById("results").innerHTML = JSON.stringify(results, null, 4);
}