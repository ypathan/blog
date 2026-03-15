function handleViewBlog(element){
	id = element.dataset.id
	window.location.href = "/blog/"+id
}

let allpages = document.querySelectorAll('.page')
let currentpage = 1
let totalpages = allpages.length

function showpage(pagenum){
	// hide all pages first incase a previous page is showing
	allpages.forEach(page => page.style.display = 'none')

	// then show the page needed
	allpages.forEach(page =>{
		let p = `page-${pagenum}`
		if (page.id == p ){
			console.log("found it")
			page.style.display = 'flex'
		}
	})
}


function changePage(direction){
	next = currentpage + direction

	//check if page exists
	if (next < 1 || next > totalpages ){
		return
	}

	showpage(next)
	currentpage = next
}

showpage(1)