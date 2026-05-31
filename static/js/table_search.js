export function tableSearch() {
    // Declare variables
    const input = document.getElementById("myInput");
    const vocabListForm = document.getElementById("vocab-list-form");
    const table = document.getElementById("myTable");
    const tr = table.getElementsByTagName("tr");
    let vocabListFilter = vocabListForm.value.toUpperCase();
    let filter = input.value.toUpperCase();
    // Loop through all table rows, and hide those who don't match the search query
    for (let i = 1; i < tr.length; i++) {
        const td = tr[i].getElementsByTagName("td");
        if (td) {
            let rowContainsFilter = false;
            let rowMatchesGroup = false;
            for (let j = 0; j < 2; j++) {
                const txtValue = td[j].textContent || td[j].innerText;
                if (txtValue.toUpperCase().includes(filter)) {
                    rowContainsFilter = true;
                    break;
                }
            }
            const groupValue = td[2].textContent || td[2].innerText;
            if (vocabListFilter === "" || groupValue.toUpperCase().indexOf(vocabListFilter) > -1) {
                rowMatchesGroup = true;
            }
            if (rowContainsFilter && rowMatchesGroup) {
                tr[i].style.display = "";
            }
            else {
                tr[i].style.display = "none";
            }
        }
    }
}
tableSearch();
