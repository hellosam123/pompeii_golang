export {};

export function tableSearch(): void {
	// Declare variables
	const input = document.getElementById("myInput") as HTMLInputElement;
	const vocabListForm = document.getElementById("vocab-list-form") as HTMLInputElement;
	const table = document.getElementById("myTable") as HTMLTableElement;
	const tr = table.getElementsByTagName("tr");
	let vocabListFilter: string = vocabListForm.value.toUpperCase();
	let filter: string = input.value.toUpperCase();
	// Loop through all table rows, and hide those who don't match the search query
	for (let i = 1; i < tr.length; i++) {
		const td = tr[i].getElementsByTagName("td");
		if (td) {
			let rowContainsFilter: boolean = false;
			let rowMatchesGroup: boolean = false;
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
			} else {
				tr[i].style.display = "none";
			}
		}
	}
}

tableSearch();
