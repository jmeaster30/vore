const doVoreSearch = () => {
  const sourceCode = document.getElementById("sourceCode").value;
  const searchText = document.getElementById("searchText").value;
  libvorejs.search(sourceCode, searchText)
    .then(value => {
      console.log(value);
      document.getElementById("results").innerHTML = JSON.stringify(value, null, 4);
    })
    .catch(err => {
      console.log("ERROR ERROR ERROR");
      document.getElementById("results").innerHTML = JSON.stringify(err, null, 4);
    });
}