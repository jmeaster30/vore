let libvore = new window.LibVoreJS.LibVoreJS();

console.log(libvore)

const doVoreSearch = () => {
  const sourceCode = document.getElementById("sourceCode").value;
  const searchText = document.getElementById("searchText").value;
  libvore.search(sourceCode, searchText)
    .then(value => {
      console.log("RES: ", value);
      document.getElementById("results").innerHTML = JSON.stringify(value, null, 4);
    });
}