function tableSearch() {
	// Declare variables
	var input, group, filter, groupFilter, table, tr, td, i, j, txtValue, groupValue, rowContainsFilter, rowMatchesGroup;
	input = document.getElementById("myInput");
	vocabListForm = document.getElementById("vocab-list-form");
	vocabListFilter = vocabListForm.value.toUpperCase();
	filter = input.value.toUpperCase();
	table = document.getElementById("myTable");
	tr = table.getElementsByTagName("tr");

	// Loop through all table rows, and hide those who don't match the search query
	for (i = 1; i < tr.length; i++) {
		td = tr[i].getElementsByTagName("td");
		if (td) {
			rowContainsFilter = false;
			rowMatchesGroup = false;
			for (j = 0; j < 2; j++) {
				txtValue = td[j].textContent || td[j].innerText;
				if (txtValue.toUpperCase().indexOf(filter) > -1) {
					rowContainsFilter = true;
					break;
				}
			}

			groupValue = td[2].textContent || td[2].innerText;
			if (vocabListFilter === "" || groupValue.toUpperCase().indexOf(vocabListFilter) > -1) {
				rowMatchesGroup = true;
			}

			if (rowContainsFilter && rowMatchesGroup) {
				tr[i].style.display = "";
			} else {
				tr[i].style.display = "none";
			}
		}
	}
}

tableSearch();
